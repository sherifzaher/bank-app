package util

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPassword(t *testing.T) {
	password := "secret"

	hashedPassword, err := HashPassword(password)
	require.NoError(t, err)

	err = VerifyPassword(hashedPassword, password)
	require.NoError(t, err)

	err = VerifyPassword(hashedPassword, "secret1")
	require.Error(t, err)
}
