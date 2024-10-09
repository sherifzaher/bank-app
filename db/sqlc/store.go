package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParams) (VerifyEmailTxResult, error)
}

type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(dbValue *sql.DB) Store {
	return &SQLStore{
		db:      dbValue,
		Queries: New(dbValue),
	}
}

func (s *SQLStore) withTx(ctx context.Context, fn func(queries *Queries) error) error {
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	db := New(tx)
	err = fn(db)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, eb err: %v", err, rbErr)
		}
		return err
	}
	return tx.Commit()
}
