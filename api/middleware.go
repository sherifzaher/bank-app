package api

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sherifzaher/clone-simplebank/token"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	abort := func(ctx *gin.Context, err error) {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		ctx.Abort()
	}

	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			abort(ctx, errors.New("authorization header is not provided"))
			return
		}

		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			abort(ctx, errors.New("invalid authorization header format"))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			abort(ctx, fmt.Errorf("unsupported authorization type %s", authorizationType))
			return
		}

		accessToken := fields[1]
		payload, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			abort(ctx, err)
		}
		ctx.Set(authorizationPayloadKey, payload)
		ctx.Next()
	}
}
