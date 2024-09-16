package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	db "github.com/sherifzaher/clone-simplebank/db/sqlc"
	"github.com/sherifzaher/clone-simplebank/token"
	"github.com/sherifzaher/clone-simplebank/util"
	"github.com/stretchr/testify/require"
	"net/http"
	"os"
	"testing"
	"time"
)

func newTestStore(t *testing.T, store db.Store) *Server {
	config, err := util.LoadConfig("../")
	require.NoError(t, err)

	server, err := NewServer(config, store)
	require.NoError(t, err)
	require.NotEmpty(t, server)

	return server
}

func createAndSetAuthToken(t *testing.T, request *http.Request, tokenMaker token.Maker, username string) {
	if len(username) == 0 {
		return
	}

	authToken, _, err := tokenMaker.CreateToken(username, time.Minute)
	require.NoError(t, err)

	authorizationHeader := fmt.Sprintf("%s %s", authorizationTypeBearer, authToken)
	request.Header.Set(authorizationHeaderKey, authorizationHeader)
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
