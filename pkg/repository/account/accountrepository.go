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

// GetByEmail gets an account by email
func (r *AccountRepository) GetByEmail(ctx context.Context, email string) (*v1.Account, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT account.id, account.email, account.password, account_permission.permission
			FROM account
			LEFT JOIN account_permission ON (account.id = account_permission.account_id)
			WHERE email = $1
		`,
		email,
	)

	if err != nil {
		return nil, err
	}

	type row struct {
		AccountID    string
		AccountEmail string
		AccountHash  string
		Permission   string
	}

	accountPermissions := []string{}

	var res row
	for rows.Next() {
		err = rows.Scan(
			&res.AccountID,
			&res.AccountEmail,
			&res.AccountHash,
			&res.Permission,
		)

		if err != nil {
			return nil, err
		}

		accountPermissions = append(accountPermissions, res.Permission)
	}

	account := &v1.Account{
		Id:          res.AccountID,
		Email:       res.AccountEmail,
		Hash:        res.AccountHash,
		Permissions: accountPermissions,
	}

	return account, nil
}

// Get gets an account
func (r *AccountRepository) Get(ctx context.Context, accountID string) (*v1.Account, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT
				account.id,
				account.email,
				account.password,
				account_permission.permission
			FROM account
			LEFT JOIN account_permission ON (account.id = account_permission.account_id)
			WHERE account.id = $1
		`,
		accountID,
	)

	if err != nil {
		return nil, err
	}

	type row struct {
		AccountID    string
		AccountEmail string
		AccountHash  string
		Permission   string
	}

	accountPermissions := []string{}

	var res row
	for rows.Next() {
		err = rows.Scan(
			&res.AccountID,
			&res.AccountEmail,
			&res.AccountHash,
			&res.Permission,
		)

		if err != nil {
			return nil, err
		}

		accountPermissions = append(accountPermissions, res.Permission)
	}

	account := &v1.Account{
		Id:          res.AccountID,
		Email:       res.AccountEmail,
		Hash:        res.AccountHash,
		Permissions: accountPermissions,
	}

	return account, nil
}

// List all accounts
func (r *AccountRepository) List(ctx context.Context, limit int32, offset int32) ([]*v1.Account, int32, error) {
	// Query configs from the database
	rows, err := r.db.QueryContext(
		ctx,
		`
			SELECT
				id,
				email,
				(SELECT COUNT(*) FROM account) AS total_size
			FROM account
			LIMIT $1
			OFFSET $2
		`,
		limit,
		offset,
	)

	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	// Unwrap rows into configs
	accounts := []*v1.Account{}
	totalSize := int32(0)

	for rows.Next() {
		account := v1.Account{}

		err := rows.Scan(
			&account.Id,
			&account.Email,
			&totalSize,
		)
		if err != nil {
			return nil, 0, err
		}

		accounts = append(accounts, &account)
	}

	return accounts, totalSize, nil
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

// AssignPermission assigns a permission to an account
func (r *AccountRepository) AssignPermission(ctx context.Context, accountID string, permission string) (*v1.Account, error) {
	_, err := r.db.ExecContext(
		ctx,
		`
			INSERT INTO account_permission(account_id, permission)
			VALUES($1, $2)
		`,
		accountID,
		permission,
	)

	if err != nil {
		return nil, err
	}

	return r.Get(ctx, accountID)
}

// RevokePermission rvokes a permission from an account
func (r *AccountRepository) RevokePermission(ctx context.Context, accountID string, permission string) (*v1.Account, error) {
	_, err := r.db.ExecContext(
		ctx,
		`
			DELETE FROM account_permission
			WHERE account_id = $1
			AND permission = $2
		`,
		accountID,
		permission,
	)

	if err != nil {
		return nil, err
	}

	return r.Get(ctx, accountID)
}
