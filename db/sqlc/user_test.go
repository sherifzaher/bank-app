package db

import (
	"context"
	"database/sql"
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

func TestUpdateUserOnlyFullName(t *testing.T) {
	oldUser, _ := randomUser(t)

	newFullName := util.RandomOwner()
	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.NotEqual(t, oldUser.FullName, newFullName)
	require.Equal(t, newFullName, newUser.FullName)
	require.Equal(t, oldUser.Email, newUser.Email)
	require.Equal(t, oldUser.HashedPassword, newUser.HashedPassword)
}

func TestUpdateUserOnlyEmail(t *testing.T) {
	oldUser, _ := randomUser(t)

	newEmail := util.RandomEmail()
	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.NotEqual(t, oldUser.Email, newUser.Email)
	require.Equal(t, newEmail, newUser.Email)
	require.Equal(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, oldUser.HashedPassword, newUser.HashedPassword)
}

func TestUpdateUserOnlyPassword(t *testing.T) {
	oldUser, _ := randomUser(t)

	newPassword := util.RandomString(6)
	hashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.NotEqual(t, oldUser.HashedPassword, newUser.HashedPassword)
	require.Equal(t, hashedPassword, newUser.HashedPassword)
	require.Equal(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, oldUser.Email, newUser.Email)
}

func TestUpdateUserAllFields(t *testing.T) {
	oldUser, _ := randomUser(t)

	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()
	newPassword := util.RandomString(6)
	hashedPassword, err := util.HashPassword(newPassword)
	require.NoError(t, err)

	newUser, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		Username: oldUser.Username,
		HashedPassword: sql.NullString{
			String: hashedPassword,
			Valid:  true,
		},
		Email: sql.NullString{
			String: newEmail,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: newFullName,
			Valid:  true,
		},
	})
	require.NoError(t, err)
	require.NotEmpty(t, newUser)
	require.NotEqual(t, oldUser.HashedPassword, newUser.HashedPassword)
	require.Equal(t, hashedPassword, newUser.HashedPassword)
	require.NotEqual(t, oldUser.Email, newUser.Email)
	require.Equal(t, newEmail, newUser.Email)
	require.NotEqual(t, oldUser.FullName, newUser.FullName)
	require.Equal(t, newFullName, newUser.FullName)
	require.Equal(t, oldUser.Username, newUser.Username)
}
