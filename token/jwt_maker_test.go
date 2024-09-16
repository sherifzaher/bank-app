package token

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJWTMaker(t *testing.T) {
	maker := JWTMaker{secretKey: "test"}
	token, payload, err := maker.CreateToken("sherif", time.Minute*15)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.NotEmpty(t, payload)
	fmt.Println(token)

	isValid, err := maker.VerifyToken(token)
	require.NoError(t, err)
	fmt.Printf("After decoding %s", isValid)
}
