package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
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

type TransferTxParams struct {
	FromAccountId int64 `json:"from_account_id"`
	ToAccountId   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToAccount   Account  `json:"to_account"`
	ToEntry     Entry    `json:"to_entry"`
}

func (s *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := s.withTx(ctx, func(queries *Queries) error {
		var err error
		result.Transfer, err = queries.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountId,
			ToAccountID:   arg.ToAccountId,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountId,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = queries.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountId,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//account1, err := queries.GetAccountForUpdate(ctx, arg.FromAccountId)
		//if err != nil {
		//	return err
		//}
		//result.FromAccount, err = queries.UpdateAccount(ctx, UpdateAccountParams{
		//	ID:      account1.ID,
		//	Balance: account1.Balance - arg.Amount,
		//})
		//if err != nil {
		//	return err
		//}
		//
		//account2, err := queries.GetAccountForUpdate(ctx, arg.ToAccountId)
		//if err != nil {
		//	return err
		//}
		//result.ToAccount, err = queries.UpdateAccount(ctx, UpdateAccountParams{
		//	ID:      account2.ID,
		//	Balance: account2.Balance + arg.Amount,
		//})
		//if err != nil {
		//	return err
		//}
		if arg.FromAccountId < arg.ToAccountId {
			result.FromAccount, result.ToAccount, err = addMoney(arg.FromAccountId, arg.ToAccountId, -arg.Amount, arg.Amount, ctx, queries)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(arg.ToAccountId, arg.FromAccountId, arg.Amount, -arg.Amount, ctx, queries)
		}
		return nil
	})
	return result, err
}

func addMoney(
	accountID1, accountID2,
	balance1, balance2 int64,
	ctx context.Context,
	queries *Queries,
) (account1, account2 Account, err error) {
	account1, err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: balance1,
	})
	if err != nil {
		return account1, account2, err
	}

	account2, err = queries.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: balance2,
	})

	return account1, account2, err
}
