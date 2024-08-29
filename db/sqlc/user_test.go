package db

import (
	"context"
	"github.com/sherifzaher/clone-simplebank/util"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func randomUser(t *testing.T) (User, string) {
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       "Sherif Zaher",
		Email:          util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)

	return user, password
}

func TestCreateUser(t *testing.T) {
	randomUser(t)
}

func TestGetUser(t *testing.T) {
	user, password := randomUser(t)

	gotUser, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, gotUser)

	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.WithinDuration(t, user.CreatedAt, gotUser.CreatedAt, time.Second)

	passwordIsMatch := util.VerifyPassword(gotUser.HashedPassword, password)
	require.NoError(t, passwordIsMatch)

	passwordIsNotMatch := util.VerifyPassword(gotUser.HashedPassword, "secret1")
	require.Error(t, passwordIsNotMatch)
}
