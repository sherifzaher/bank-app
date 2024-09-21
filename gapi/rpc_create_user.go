package gapi

import (
	"context"
	db "github.com/sherifzaher/clone-simplebank/db/sqlc"
	"github.com/sherifzaher/clone-simplebank/pb"
	"github.com/sherifzaher/clone-simplebank/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password: %s", err)
	}

	user, err := server.store.CreateUser(ctx, db.CreateUserParams{
		Username:       req.GetUsername(),
		Email:          req.GetEmail(),
		FullName:       req.GetFullName(),
		HashedPassword: hashedPassword,
	})

	if err != nil {
		return nil, status.Errorf(codes.AlreadyExists, "username already exists: %s", err)
	}

	rsp := &pb.CreateUserResponse{
		User: convertUser(user),
	}
	return rsp, nil
}
