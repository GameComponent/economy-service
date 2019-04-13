package accountrepository

import (
	"context"
	"database/sql"
)

// AccountRepository struct
type AccountRepository struct {
	db *sql.DB
}

// Account struct
type Account struct {
	ID    string
	Email string
	Hash  string
}

// NewAccountRepository constructor
func NewAccountRepository(db *sql.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

// Get an account
func (r *AccountRepository) Get(ctx context.Context, email string) *Account {
	account := Account{}

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, email, password FROM account WHERE email = $1`,
		email,
	).Scan(
		&account.ID,
		&account.Email,
		&account.Hash,
	)

	if err != nil {
		return nil
	}

	return &account
}

// Create an account
func (r *AccountRepository) Create(ctx context.Context, email string, password string) *Account {
	var id string

	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO account(email, password) VALUES($1, $2) RETURNING id`,
		email,
		password,
	).Scan(&id)

	if err != nil {
		return nil
	}

	account := Account{
		ID:    id,
		Email: email,
		Hash:  password,
	}

	return &account
}
