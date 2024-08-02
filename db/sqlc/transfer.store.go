package db

import (
	"context"
	"database/sql"
)

type TransferParams struct {
	FromAccountID int64   `json:"from_account_id"`
	ToAccountID   int64   `json:"to_account_id"`
	Amount        float64 `json:"amount"`
}

type TransferResult struct {
	Transfer     Transfer `json:"transfer"`
	FromAccount  Account  `json:"from_account"`
	ToAccount    Account  `json:"to_account"`
	FromMovement Movement `json:"from_movement"`
	ToMovement   Movement `json:"to_movement"`
}

func (s *Store) Transfer(ctx context.Context, arg TransferParams) (TransferResult, error) {
	var result TransferResult

	err := s.transaction(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		toAccount, err := q.GetAccount(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		fromAccount, err := q.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		result.FromMovement, err = q.CreateMovement(ctx, CreateMovementParams{
			AccountID:   arg.FromAccountID,
			Amount:      -arg.Amount,
			TransferID:  sql.NullInt64{Int64: result.Transfer.ID, Valid: true},
			Description: "Transfer to account " + toAccount.Name,
			Validated:   false,
		})
		if err != nil {
			return err
		}

		result.ToMovement, err = q.CreateMovement(ctx, CreateMovementParams{
			AccountID:   arg.ToAccountID,
			Amount:      arg.Amount,
			TransferID:  sql.NullInt64{Int64: result.Transfer.ID, Valid: true},
			Description: "Transfer from account " + fromAccount.Name,
			Validated:   false,
		})

		if err != nil {
			return err
		}

		// Update first the account with the lowest ID to avoid deadlocks in the database
		if fromAccount.ID < toAccount.ID {
			result.FromAccount, result.ToAccount, err = moveMoney(ctx, q, fromAccount.ID, -arg.Amount, toAccount.ID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = moveMoney(ctx, q, toAccount.ID, arg.Amount, fromAccount.ID, -arg.Amount)
			if err != nil {
				return err
			}
		}

		return nil
	})

	return result, err
}

func moveMoney(
	ctx context.Context,
	q *Queries,
	account1Id int64,
	amount1 float64,
	account2Id int64,
	amount2 float64,

) (account1 Account, account2 Account, err error) {

	account1, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     account1Id,
		Amount: amount1,
	})

	if err != nil {
		return
	}

	account2, err = q.UpdateAccountBalance(ctx, UpdateAccountBalanceParams{
		ID:     account2Id,
		Amount: amount2,
	})

	return
}
