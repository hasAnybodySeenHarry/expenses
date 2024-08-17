package data

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"

	"harry2an.com/expenses/internal/validator"
)

const (
	ScopeAuthentication = "authentication"
	ScopeActivation     = "activation"
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

func CreateToken(userID int64, scope string, ttl time.Duration) (*Token, error) {
	t := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	t.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b)
	hash := sha256.Sum256([]byte(t.Plaintext))
	t.Hash = hash[:]

	return t, nil
}

type TokenModel struct {
	db *sql.DB
}

func (m TokenModel) New(userID int64, scope string, ttl time.Duration) (*Token, error) {
	token, err := CreateToken(userID, scope, ttl)
	if err != nil {
		return nil, err
	}

	err = m.Insert(token)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (m TokenModel) Insert(token *Token) error {
	stmt := `
		INSERT INTO tokens (hash, user_id, expiry, scope)
		VALUES ($1, $2, $3, $4)
	`

	args := []interface{}{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.db.ExecContext(ctx, stmt, args...)
	return err
}
