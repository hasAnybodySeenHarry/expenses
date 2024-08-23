package main

import (
	"errors"
	"net/http"

	"harry2an.com/expenses/internal/data"
	"harry2an.com/expenses/internal/validator"
)

func (app *application) createTransaction(w http.ResponseWriter, r *http.Request) {
	var input struct {
		DebtID      *int64  `json:"debt_id"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	v := validator.New()

	if input.DebtID == nil || *input.DebtID == 0 {
		v.AddError("debt's id", "must be provided")
		app.failedValidation(w, r, v.Errors)
		return
	}

	t := &data.Transaction{
		DebtID:      *input.DebtID,
		Amount:      input.Amount,
		Description: input.Description,
	}

	if data.ValidateTransaction(v, t); !v.Validate() {
		app.failedValidation(w, r, v.Errors)
		return
	}

	u := app.getUser(r)
	debt, err := app.models.Debts.GetByID(t.DebtID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			app.notFound(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	if u.ID != debt.Lender.ID && u.ID != debt.Borrower.ID {
		app.forbidden(w, r)
		return
	}

	if u.ID == debt.Borrower.ID && t.Amount <= 0 {
		app.forbidden(w, r)
		return
	}

	if u.ID == debt.Lender.ID && t.Amount >= 0 {
		app.forbidden(w, r)
		return
	}

	t.Borrower.ID = debt.Borrower.ID
	t.Borrower.Name = u.Name

	amount, err := app.models.Transactions.Insert(t)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	lender, err := app.models.Users.GetUsernameByID(debt.Lender.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			app.accepted(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	t.Lender.ID = lender.ID
	t.Lender.Name = lender.Name

	err = app.notifiers.Transactions.Send(t)
	if err != nil {
		app.log(r, err)
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{
		"transaction": t,
		"new_total":   amount,
	}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) showTransactions(w http.ResponseWriter, r *http.Request) {
	var f data.Filters

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	q := r.URL.Query()
	v := validator.New()

	f.Page = app.readInt(q, "page", 1, v)
	f.PageSize = app.readInt(q, "page_size", 10, v)

	if data.ValidateFilters(v, &f); !v.Validate() {
		app.failedValidation(w, r, v.Errors)
		return
	}

	transactions, meta, err := app.models.Transactions.GetAll(id, f)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"data": transactions, "metadata": meta}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}
