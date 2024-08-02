// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: account.sql

package db

import (
	"context"
)

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts (
  name, owner_id, account_type_id, balance
) VALUES (
  $1, $2, $3, $4
) RETURNING id, uuid, name, owner_id, account_type_id, balance, created_at, updated_at, deleted_at
`

type CreateAccountParams struct {
	Name          string  `json:"name"`
	OwnerID       int64   `json:"owner_id"`
	AccountTypeID int64   `json:"account_type_id"`
	Balance       float64 `json:"balance"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createAccount,
		arg.Name,
		arg.OwnerID,
		arg.AccountTypeID,
		arg.Balance,
	)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Name,
		&i.OwnerID,
		&i.AccountTypeID,
		&i.Balance,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getAccount = `-- name: GetAccount :one
SELECT id, uuid, name, owner_id, account_type_id, balance, created_at, updated_at, deleted_at
FROM accounts
WHERE deleted_at is null and id = $1
`

func (q *Queries) GetAccount(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Name,
		&i.OwnerID,
		&i.AccountTypeID,
		&i.Balance,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getAccountForUpdate = `-- name: GetAccountForUpdate :one
SELECT id, uuid, name, owner_id, account_type_id, balance, created_at, updated_at, deleted_at
FROM accounts
WHERE deleted_at is null and id = $1
FOR NO KEY UPDATE
`

func (q *Queries) GetAccountForUpdate(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountForUpdate, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Name,
		&i.OwnerID,
		&i.AccountTypeID,
		&i.Balance,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const getUserAccounts = `-- name: GetUserAccounts :many
SELECT id, uuid, name, owner_id, account_type_id, balance, created_at, updated_at, deleted_at
FROM accounts
WHERE deleted_at is null and owner_id = $1
order by id desc
`

func (q *Queries) GetUserAccounts(ctx context.Context, ownerID int64) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, getUserAccounts, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Uuid,
			&i.Name,
			&i.OwnerID,
			&i.AccountTypeID,
			&i.Balance,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.DeletedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const removeAccount = `-- name: RemoveAccount :exec
UPDATE accounts
SET deleted_at = now()
WHERE deleted_at is null and id = $1
`

func (q *Queries) RemoveAccount(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, removeAccount, id)
	return err
}

const updateAccount = `-- name: UpdateAccount :one
UPDATE accounts
SET name = $2
WHERE deleted_at is null and id = $1
RETURNING id, uuid, name, owner_id, account_type_id, balance, created_at, updated_at, deleted_at
`

type UpdateAccountParams struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, updateAccount, arg.ID, arg.Name)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Name,
		&i.OwnerID,
		&i.AccountTypeID,
		&i.Balance,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const updateAccountBalance = `-- name: UpdateAccountBalance :one
UPDATE accounts
SET balance = balance + $1
WHERE deleted_at is null and id = $2
RETURNING id, uuid, name, owner_id, account_type_id, balance, created_at, updated_at, deleted_at
`

type UpdateAccountBalanceParams struct {
	Amount float64 `json:"amount"`
	ID     int64   `json:"id"`
}

func (q *Queries) UpdateAccountBalance(ctx context.Context, arg UpdateAccountBalanceParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, updateAccountBalance, arg.Amount, arg.ID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Uuid,
		&i.Name,
		&i.OwnerID,
		&i.AccountTypeID,
		&i.Balance,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}
