package main

import (
	"errors"
	"net/http"

	"harry2an.com/expenses/internal/data"
	"harry2an.com/expenses/internal/validator"
)

func (app *application) createTransactionHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		DebtID      int64   `json:"debt_id"`
		Amount      float64 `json:"amount"`
		Description string  `json:"description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	t := &data.Transaction{
		DebtID:      input.DebtID,
		Amount:      input.Amount,
		Description: input.Description,
	}

	v := validator.New()
	if data.ValidateTransaction(v, t); !v.Validate() {
		app.failedValidation(w, r, v.Errors)
		return
	}

	u := app.getUser(r)
	debt, err := app.models.Debts.GetByID(input.DebtID)
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

	amount, err := app.models.Transactions.Insert(t)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{
		"transaction": t,
		"new_total":   amount,
	}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}
