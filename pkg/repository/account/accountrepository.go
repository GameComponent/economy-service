package accountrepository

import (
	"context"
	"database/sql"

	v1 "github.com/GameComponent/economy-service/pkg/api/v1"
	repository "github.com/GameComponent/economy-service/pkg/repository"
	"go.uber.org/zap"
)

// AccountRepository struct
type AccountRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewAccountRepository constructor
func NewAccountRepository(db *sql.DB, logger *zap.Logger) repository.AccountRepository {
	return &AccountRepository{
		db:     db,
		logger: logger,
	}
}

// Get an account
func (r *AccountRepository) Get(ctx context.Context, email string) (*v1.Account, error) {
	account := &v1.Account{}

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, email, password FROM account WHERE email = $1`,
		email,
	).Scan(
		&account.Id,
		&account.Email,
		&account.Hash,
	)

	if err != nil {
		return nil, err
	}

	return account, nil
}

// Create an account
func (r *AccountRepository) Create(ctx context.Context, email string, password string) (*v1.Account, error) {
	var id string

	err := r.db.QueryRowContext(
		ctx,
		`INSERT INTO account(email, password) VALUES($1, $2) RETURNING id`,
		email,
		password,
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	account := &v1.Account{
		Id:    id,
		Email: email,
		Hash:  password,
	}

	return account, nil
}

// Update an account
func (r *AccountRepository) Update(ctx context.Context, email string, password string) (*v1.Account, error) {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE account SET password = $1 WHERE email = $2`,
		password,
		email,
	)

	if err != nil {
		return nil, err
	}

	return r.Get(ctx, email)
}
