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

func (m DebtModel) GetForUserByCategories(userID int64) ([]*Debt, error) {
	stmt := `
		SELECT 
			debts.id, debts.category, debts.total_amount, debts.created_at, debts.version,
			lender.id AS lender_id, lender.name AS lender_name,
			borrower.id AS borrower_id, borrower.name AS borrower_name
		FROM debts
		INNER JOIN users lender ON debts.lender_id = lender.id
		INNER JOIN users borrower ON debts.borrower_id = borrower.id
		WHERE debts.borrower_id = $1
		ORDER BY debts.id DESC;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, stmt, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	debts := make([]*Debt, 0)

	for rows.Next() {
		var debt Debt

		err := rows.Scan(
			&debt.ID, &debt.Category, &debt.Total, &debt.CreatedAt, &debt.Version,
			&debt.Lender.ID, &debt.Lender.Name,
			&debt.Borrower.ID, &debt.Borrower.Name,
		)
		if err != nil {
			return nil, err
		}
		debts = append(debts, &debt)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return debts, nil
}

func (m DebtModel) Insert(debt *Debt) error {
	stmt := `
		INSERT INTO debts (lender_id, borrower_id, category, total_amount)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version
	`

	args := []interface{}{debt.Lender.ID, debt.Borrower.ID, debt.Category, debt.Total}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, stmt, args...).Scan(&debt.ID, &debt.CreatedAt, &debt.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "debts_category_key"`:
			return ErrDuplicateCategory
		case err.Error() == `pq: insert or update on table "debts" violates foreign key constraint "debts_lender_id_fkey"`:
			return ErrUserForeignKey
		default:
			return err
		}
	}

	return nil
}

func (m DebtModel) DeleteByID(debtID int64, version uuid.UUID) error {
	stmt := `
		DELETE FROM debts
		WHERE id = $1 AND version = $2
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.db.ExecContext(ctx, stmt, debtID, version)
	if err != nil {
		return err
	}

	n, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if n != 1 {
		return ErrWriteConflict
	}

	return nil
}
