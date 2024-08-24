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

	if !isAuthorizedUser(u, debt) {
		app.forbidden(w, r)
		return
	}

	if !isTransactionAllowed(u, debt, t.Amount) {
		app.forbidden(w, r)
		return
	}

	amount, err := app.models.Transactions.Insert(t)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	t.Borrower.ID = debt.Borrower.ID
	t.Lender.ID = debt.Lender.ID

	// Fetch counterparty name
	var counterpartyID int64
	if u.ID == t.Borrower.ID {
		counterpartyID = t.Lender.ID
	} else {
		counterpartyID = t.Borrower.ID
	}

	counterparty, err := app.models.Users.GetUsernameByID(counterpartyID)
	if err != nil {
		if errors.Is(err, data.ErrNoRecord) {
			app.accepted(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	if u.ID == t.Borrower.ID {
		t.Borrower.Name = u.Name
		t.Lender.Name = counterparty.Name
	} else {
		t.Borrower.Name = counterparty.Name
		t.Lender.Name = u.Name
	}

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

func isAuthorizedUser(user *data.User, debt *data.Debt) bool {
	return user.ID == debt.Lender.ID || user.ID == debt.Borrower.ID
}

func isTransactionAllowed(user *data.User, debt *data.Debt, amount float64) bool {
	if user.ID == debt.Borrower.ID {
		return amount > 0
	}
	if user.ID == debt.Lender.ID {
		return amount < 0
	}
	return false
}
