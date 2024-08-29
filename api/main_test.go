package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/sherifzaher/clone-simplebank/db/sqlc"
	"github.com/sherifzaher/clone-simplebank/util"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func newTestStore(t *testing.T, store db.Store) *Server {
	config, err := util.LoadConfig("../")
	require.NoError(t, err)

	server, err := NewServer(config, store)
	require.NoError(t, err)
	require.NotEmpty(t, server)

	return server
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
