package main

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"harry2an.com/expenses/internal/core"
	"harry2an.com/expenses/internal/data"
	"harry2an.com/expenses/internal/validator"
)

func (app *application) showDebts(w http.ResponseWriter, r *http.Request) {
	user := app.getUser(r)

	debts, err := app.models.Debts.GetForUserByCategories(user.ID)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"debts": debts}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) createDebt(w http.ResponseWriter, r *http.Request) {
	var input struct {
		LenderID *int64  `json:"lender_id"`
		Total    float64 `json:"total"`
		Category string  `json:"category"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequest(w, r, err)
		return
	}

	u := app.getUser(r)
	v := validator.New()

	if input.LenderID == nil || *input.LenderID == 0 {
		v.AddError("lender's id", "must be provided")
		app.failedValidation(w, r, v.Errors)
		return
	}

	debt := &data.Debt{
		Lender: core.Entity{
			ID: *input.LenderID,
		},
		Borrower: core.Entity{
			ID:   u.ID,
			Name: u.Name,
		},
		Total:    input.Total,
		Category: input.Category,
	}

	if data.ValidateDebt(v, debt); !v.Validate() {
		app.failedValidation(w, r, v.Errors)
		return
	}

	err = app.models.Debts.Insert(debt)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateCategory):
			v.AddError("category", "already exists")
			app.failedValidation(w, r, v.Errors)
		case errors.Is(err, data.ErrUserForeignKey):
			app.badRequest(w, r, errors.New("invalid lender's id"))
		default:
			app.serverError(w, r, err)
		}
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

	debt.Lender.Name = lender.Name

	err = app.notifiers.Debts.Send(debt)
	if err != nil {
		app.log(r, err)
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"debt": debt}, nil)
	if err != nil {
		app.serverError(w, r, err)
	}
}

func (app *application) deleteDebt(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFound(w, r)
		return
	}

	u := app.getUser(r)
	ver := r.Header.Get("X-Expected-Version")

	if ver == "" {
		app.badRequest(w, r, errors.New("requires version to modify the resource"))
		return
	}

	debt, err := app.models.Debts.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoRecord):
			app.notFound(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	v, err := uuid.Parse(ver)
	if err != nil {
		app.badRequest(w, r, errors.New("invalid version"))
		return
	}

	if debt.Version != v {
		app.conflict(w, r)
		return
	}

	if debt.Lender.ID != u.ID {
		app.forbidden(w, r)
		return
	}

	err = app.models.Debts.DeleteByID(id, v)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrWriteConflict):
			app.conflict(w, r)
		default:
			app.serverError(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusNoContent,
		envelope{"message": "successfully deleted the debt"}, nil,
	)
	if err != nil {
		app.serverError(w, r, err)
	}
}
