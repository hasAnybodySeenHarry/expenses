package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"harry2an.com/expenses/internal/validator"
)

type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  password  `json:"-"`
	Activated bool      `json:"activated"`
	CreatedAt time.Time `json:"-"`
	Version   uuid.UUID `json:"version"`
}

type password struct {
	text *string
	hash []byte
}

func (p *password) Matches(pwd string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(pwd))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}

	return true, nil
}

func (p *password) Create(pwd string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	p.text = &pwd
	p.hash = hash
	return nil
}

func ValidateUser(v *validator.Validator, u *User) {
	v.Check(u.Name != "", "name", "is blank")
	v.Check(len(u.Name) <= 100, "name", "must not exceed 100 chars")
	v.Check(u.Email != "", "email", "is blank")

	if u.Password.text != nil {
		ValidatePassword(v, *u.Password.text)
	}

	if u.Password.hash == nil {
		panic("hash is not set")
	}
}

func ValidatePassword(v *validator.Validator, pwd string) {
	v.Check(pwd != "", "password", "is blank")
	v.Check(len(pwd) >= 8, "password", "must be at least 8 chars long")
	v.Check(len(pwd) <= 72, "password", "must not exceed 72 chars")
}

type UserModel struct {
	db *sql.DB
}

func (m UserModel) Insert(user *User) error {
	stmt := `
		INSERT INTO users (name, email, password, activated)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version
	`

	args := []interface{}{user.ID, user.Email, user.Password.hash, user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	err := m.db.QueryRowContext(ctx, stmt, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Version,
	)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (m UserModel) GetForToken(token string, scope string) (*User, error) {
	stmt := `
		SELECT users.id, users.name, users.email, users.created_at, users.password, users.activated, users.version
		FROM users
		INNER JOIN tokens
		ON users.id = tokens.user_id
		WHERE tokens.hash = $1 AND tokens.scope = $2 AND tokens.expiry > $3
	`

	hash := sha256.Sum256([]byte(token))
	args := []interface{}{hash[:], scope, time.Now()}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.db.QueryRowContext(ctx, stmt, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.Password.hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNoRecord
		default:
			return nil, err
		}
	}

	return &user, nil
}
