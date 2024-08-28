package db

import (
	"context"
	"testing"

	"github.com/sherifzaher/clone-simplebank/util"
	"github.com/stretchr/testify/require"
)

func randomAccount() (Account, error) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Currency: util.RandomCurrency(),
		Balance:  util.RandomInt(1, 1000),
	}

	return testQueries.CreateAccount(context.Background(), arg)
}

func TestCreateAccount(t *testing.T) {
	acc, err := randomAccount()
	require.NoError(t, err)
	require.NotEmpty(t, acc)
	require.True(t, acc.Balance >= 1)
}

func TestGetAccount(t *testing.T) {
	acc, err := randomAccount()
	require.NoError(t, err)
	require.NotEmpty(t, acc)
	require.True(t, acc.Balance >= 1)

	account, err := testQueries.GetAccount(context.Background(), acc.ID)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, account, acc)
}
