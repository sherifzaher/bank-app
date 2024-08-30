package db

import (
	"context"
	"testing"

	"github.com/sherifzaher/clone-simplebank/util"
	"github.com/stretchr/testify/require"
)

func randomAccount(t *testing.T) Account {
	user, _ := randomUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Currency: util.RandomCurrency(),
		Balance:  util.RandomInt(1, 1000),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.True(t, account.Balance >= 1)

	return account
}

func TestCreateAccount(t *testing.T) {
	randomAccount(t)
}

func TestGetAccount(t *testing.T) {
	acc := randomAccount(t)

	account, err := testQueries.GetAccount(context.Background(), GetAccountParams{
		Owner:    acc.Owner,
		Currency: acc.Currency,
		ID:       acc.ID,
	})
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, account, acc)
}

func TestDeleteAccount(t *testing.T) {
	acc := randomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), acc.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), GetAccountParams{
		Owner:    acc.Owner,
		Currency: acc.Currency,
		ID:       acc.ID,
	})
	require.Error(t, err)
	require.Empty(t, account)
}

func TestUpdateAccount(t *testing.T) {
	account1 := randomAccount(t)

	arg := UpdateAccountParams{
		ID:      account1.ID,
		Balance: util.RandomMoney(),
	}

	account2, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.ID, account2.ID)
	require.Equal(t, account2.Balance, arg.Balance)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Owner, account2.Owner)
}
