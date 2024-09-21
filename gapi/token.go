package gapi

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/sherifzaher/clone-simplebank/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) RefreshToken(ctx context.Context, body *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	var refreshToken = body.GetRefreshToken()

	refreshPayload, err := server.tokenMaker.VerifyToken(refreshToken)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "token is not valid or expired")
	}

	session, err := server.store.GetSession(ctx, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			status.Errorf(codes.Internal, "session is not found: %s", err)
		}
		return nil, status.Errorf(codes.Internal, "%s", err)
	}

	if session.IsBlocked {
		err := errors.New("session is blocked")
		return nil, status.Errorf(codes.DeadlineExceeded, "%s", err)
	}

	if session.Username != refreshPayload.Username {
		err := errors.New("incorrect session user")
		return nil, status.Errorf(codes.Canceled, "%s", err)
	}

	if session.RefreshToken != refreshToken {
		err := errors.New("mismatched session token")
		return nil, status.Errorf(codes.Canceled, "%s", err)
	}

	if time.Now().After(session.ExpiresAt) {
		err := errors.New("expired session")
		return nil, status.Errorf(codes.DeadlineExceeded, "%s", err)
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(refreshPayload.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed during create access-token: %s", err)
	}

	responsePayload := &pb.RefreshTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: timestamppb.New(accessPayload.ExpiredAt),
	}

	return responsePayload, nil
}
