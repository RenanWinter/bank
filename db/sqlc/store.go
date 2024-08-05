package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	Transfer(ctx context.Context, arg TransferParams) (TransferResult, error) 
}

// Store provides all functions to execute SQL and transactions
type SQLStore struct {
	*Queries
	db *sql.DB
}

var txKey = struct{}{}

// Creates a new Store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// Transact executes a function within a database transaction
func (store *SQLStore) transaction(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %v, failed to rollback: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
