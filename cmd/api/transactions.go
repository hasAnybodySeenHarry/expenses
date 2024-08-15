package main

import (
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

	_ = app.getUser(r)
	// to be continued
}
