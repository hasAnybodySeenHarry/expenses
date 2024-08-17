package data

import (
	"database/sql"
	"errors"
)

var (
	ErrNoRecord          = errors.New("no record found")
	ErrDuplicateEmail    = errors.New("duplicate email")
	ErrDuplicateCategory = errors.New("duplicate category")
	ErrUserForeignKey    = errors.New("user foreign key error")
	ErrWriteConflict     = errors.New("stale data")
)

type Models struct {
	Users        UserModel
	Debts        DebtModel
	Transactions TransactionModel
	Tokens       TokenModel
}

func New(db *sql.DB) Models {
	return Models{
		Users:        UserModel{db: db},
		Debts:        DebtModel{db: db},
		Transactions: TransactionModel{db: db},
		Tokens:       TokenModel{db: db},
	}
}
