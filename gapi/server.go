package gapi

import (
	db "github.com/sherifzaher/clone-simplebank/db/sqlc"
	"github.com/sherifzaher/clone-simplebank/pb"
	"github.com/sherifzaher/clone-simplebank/token"
	"github.com/sherifzaher/clone-simplebank/util"
	"github.com/sherifzaher/clone-simplebank/worker"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	pb.UnimplementedSimpleBankServer
	store           db.Store
	config          util.Config
	tokenMaker      token.Maker
	taskDistributor worker.TaskDistributor
}

// NewServer creates a new gRPC server.
func NewServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) (*Server, error) {
	pasetoTokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, err
	}
	server := &Server{
		store:           store,
		config:          config,
		tokenMaker:      pasetoTokenMaker,
		taskDistributor: taskDistributor,
	}

	return server, nil
}
