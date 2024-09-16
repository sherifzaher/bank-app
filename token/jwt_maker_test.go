package token

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestJWTMaker(t *testing.T) {
	maker := JWTMaker{secretKey: "test"}
	token, err := maker.CreateToken("sherif", time.Minute*15)
	require.NoError(t, err)
	require.NotEmpty(t, token)
	fmt.Println(token)

	isValid, err := maker.VerifyToken(token)
	require.NoError(t, err)
	fmt.Printf("After decoding %s", isValid)
}
