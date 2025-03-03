// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package db

import (
	"context"

	"github.com/google/uuid"
)

type Querier interface {
	CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error)
	CreateCredential(ctx context.Context, arg CreateCredentialParams) (Credential, error)
	CreateMovement(ctx context.Context, arg CreateMovementParams) (Movement, error)
	CreateTransfer(ctx context.Context, arg CreateTransferParams) (Transfer, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetAccount(ctx context.Context, id int64) (Account, error)
	GetAccountForUpdate(ctx context.Context, id int64) (Account, error)
	GetAccountMovements(ctx context.Context, arg GetAccountMovementsParams) ([]Movement, error)
	GetMovement(ctx context.Context, id int64) (Movement, error)
	GetTransfer(ctx context.Context, id int64) (Transfer, error)
	GetUserAccounts(ctx context.Context, ownerID int64) ([]Account, error)
	GetUserActiveCredential(ctx context.Context, userID int64) (Credential, error)
	GetUserByEmail(ctx context.Context, email string) (User, error)
	GetUserById(ctx context.Context, id int64) (User, error)
	GetUserByUUID(ctx context.Context, argUuid uuid.UUID) (User, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetUserCredentials(ctx context.Context, userID int64) ([]Credential, error)
	ListUsers(ctx context.Context, arg ListUsersParams) ([]User, error)
	RemoveAccount(ctx context.Context, id int64) error
	RemoveMovement(ctx context.Context, id int64) error
	RemoveTransfer(ctx context.Context, id int64) error
	RemoveUser(ctx context.Context, id int64) error
	RemoveUserCredential(ctx context.Context, userID int64) error
	UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error)
	UpdateAccountBalance(ctx context.Context, arg UpdateAccountBalanceParams) (Account, error)
	UpdateMovement(ctx context.Context, arg UpdateMovementParams) (Movement, error)
	UpdateUser(ctx context.Context, arg UpdateUserParams) (User, error)
}

var _ Querier = (*Queries)(nil)
