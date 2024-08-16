package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"harry2an.com/expenses/internal/validator"
)

type Debt struct {
	ID int64 `json:"id"`

	Lender struct {
		ID   int64  `json:"id"`
		Name string `json:"name,omitempty"`
	} `json:"lender"`

	Borrower struct {
		ID   int64  `json:"id"`
		Name string `json:"name,omitempty"`
	} `json:"borrower"`

	Category string  `json:"category"`
	Total    float64 `json:"total"`

	CreatedAt time.Time `json:"created_at"`
	Version   uuid.UUID `json:"-"`
}

func ValidateDebt(v *validator.Validator, d *Debt) {
	v.Check(d.Lender.ID != d.Borrower.ID, "user", "cannot lend to oneself")
	v.Check(d.Category != "", "category", "cannot be empty")
	v.Check(len(d.Category) >= 4 && len(d.Category) <= 60, "category", "must be between 4 and 60 characters")
	v.Check(d.Total >= 0, "debt", "cannot be issued in a negative amount")
}

type DebtModel struct {
	db *sql.DB
}

func (m DebtModel) GetByID(debtID int64) (*Debt, error) {
	stmt := `
		SELECT id, borrower_id, lender_id, version
		FROM debts
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var d Debt
	err := m.db.QueryRowContext(ctx, stmt, debtID).Scan(
		&d.ID,
		&d.Borrower.ID,
		&d.Lender.ID,
		&d.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecord
		default:
			return nil, err
		}
	}

	return &d, nil
}
