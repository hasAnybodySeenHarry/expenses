package data

import (
	"database/sql"
	"errors"
)

var (
	ErrNoRecord       = errors.New("no record found")
	ErrDuplicateEmail = errors.New("duplicate email")
)

type Models struct {
	Users        UserModel
	Debts        DebtModel
	Transactions TransactionModel
}

func New(db *sql.DB) Models {
	return Models{
		Users:        UserModel{db: db},
		Debts:        DebtModel{db: db},
		Transactions: TransactionModel{db: db},
	}
}
