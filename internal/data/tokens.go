package data

import (
	"database/sql"
	"time"

	"harry2an.com/expenses/internal/validator"
)

const (
	ScopeAuthentication = "authentication"
)

type Token struct {
	Plaintext string    `json:"token"`
	Hash      []byte    `json:"-"`
	UserID    int64     `json:"-"`
	Expiry    time.Time `json:"expired_at"`
	Scope     string    `json:"-"`
}

func ValidateToken(v *validator.Validator, token string) {
	v.Check(token != "", "token", "must not be empty")
	v.Check(len(token) == 26, "token", "must be 26 bytes long")
}

type TokenModel struct {
	db *sql.DB
}
