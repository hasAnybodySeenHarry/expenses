package data

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/google/uuid"
	"harry2an.com/expenses/internal/validator"
)

type Transaction struct {
	ID          int64     `json:"id"`
	DebtID      int64     `json:"-"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Version     uuid.UUID `json:"version"`
}

func ValidateTransaction(v *validator.Validator, t *Transaction) {
	v.Check(t.Amount != 0, "amount", "must be provided")
	v.Check(t.Description != "", "description", "must be provided")
	v.Check(len(t.Description) >= 4 && len(t.Description) <= 200, "description", "must be between 4 and 200 characters")
}

type TransactionModel struct {
	db *sql.DB
}

func (m TransactionModel) Insert(t *Transaction) (float64, error) {
	tStmt := `
		INSERT INTO transactions (debt_id, amount, description, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, version
	`

	tArgs := []interface{}{t.DebtID, t.Amount, t.Description}

	dStmt := `
		UPDATE debts
		SET total_amount = total_amount + $1
		WHERE id = $2
	`

	dArgs := []interface{}{t.Amount, t.DebtID}

	tx, err := m.db.Begin()
	if err != nil {
		return 0, err
	}

	var rollBackErr error

	defer func() {
		if p := recover(); p != nil {
			rollBackErr = tx.Rollback()
			if rollBackErr != nil {
				log.Printf("tx.Rollback failed: %v", rollBackErr)
			}
			panic(p)
		} else if err != nil {
			rollBackErr = tx.Rollback()
			if rollBackErr != nil {
				log.Printf("tx.Rollback failed: %v", rollBackErr)
			}
		} else {
			err = tx.Commit()
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()

	err = tx.QueryRowContext(ctx, tStmt, tArgs...).Scan(&t.ID, &t.Version)
	if err != nil {
		return 0, err
	}

	_, err = tx.ExecContext(ctx, dStmt, dArgs...)
	if err != nil {
		return 0, err
	}

	totalStmt := `
		SELECT total_amount
		FROM debts
		WHERE id = $1
	`

	var newTotal float64
	err = tx.QueryRowContext(ctx, totalStmt, t.DebtID).Scan(&newTotal)
	if err != nil {
		return 0, err
	}

	return newTotal, nil
}

func (m TransactionModel) GetAll(debtID int64, f Filters) ([]Transaction, Metadata, error) {
	stmt := `
		SELECT count(*) OVER(), id, debt_id, amount, description, created_at, version
		FROM transactions
		WHERE debt_id = $1
		ORDER BY id ASC
		LIMIT $2
		OFFSET $3
	`

	args := []interface{}{debtID, f.limit(), f.offset()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, stmt, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	transactions := make([]Transaction, 0, f.limit())
	var total int

	for rows.Next() {
		var t Transaction
		err := rows.Scan(&total, &t.ID, &t.DebtID, &t.Amount, &t.Description, &t.CreatedAt, &t.Version)
		if err != nil {
			return nil, Metadata{}, err
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	meta := getMetadata(f.Page, f.PageSize, total)

	return transactions, meta, nil
}
