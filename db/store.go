package db

import (
	"context"
	"database/sql"
	"fmt"
	db "github.com/sherifzaher/clone-simplebank/db/sqlc"
)

type Store interface {
	db.Querier
	// TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type SQLStore struct {
	*db.Queries
	db *sql.DB
}

func NewStore(dbValue *sql.DB) Store {
	return &SQLStore{
		db:      dbValue,
		Queries: db.New(dbValue),
	}
}

func (s *SQLStore) withTx(ctx context.Context, fn func(queries *db.Queries) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := db.New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, eb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
