package accountrepository

import (
	"context"
	"database/sql"
	
	"go.uber.org/zap"
	repository "github.com/GameComponent/economy-service/pkg/repository"
)

// AccountRepository struct
type AccountRepository struct {
	db *sql.DB
	logger *zap.Logger
}


// NewAccountRepository constructor
func NewAccountRepository(db *sql.DB, logger *zap.Logger) repository.AccountRepository {
	return &AccountRepository{
		db: db,
		logger: logger,
	}
}

// Get an account
func (r *AccountRepository) Get(ctx context.Context, email string) *repository.Account {
	account := repository.Account{}

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
func (r *AccountRepository) Create(ctx context.Context, email string, password string) *repository.Account {
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

	account := repository.Account{
		ID:    id,
		Email: email,
		Hash:  password,
	}

	return &account
}
