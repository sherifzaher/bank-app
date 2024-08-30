package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/sherifzaher/clone-simplebank/db/sqlc"
)

type CreateAccountParams struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req CreateAccountParams
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateAccountParams{
		Owner:    req.Owner,
		Currency: req.Currency,
		Balance:  0,
	}

	account, err := server.store.CreateAccount(ctx, arg)
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		switch pqErr.Code.Name() {
		case "foreign_key_violation", "unique_violation":
			ctx.JSON(http.StatusForbidden, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type GetAccountUriParams struct {
	ID int64 `uri:"id" binding:"required"`
}

type GetAccountQueryParams struct {
	Owner    string `form:"owner" binding:"required"`
	Currency string `form:"currency" binding:"required,currency"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var uriParams GetAccountUriParams
	if err := ctx.ShouldBindUri(&uriParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var queryParams GetAccountQueryParams
	if err := ctx.ShouldBindQuery(&queryParams); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, db.GetAccountParams{
		ID:       uriParams.ID,
		Currency: queryParams.Currency,
		Owner:    queryParams.Owner,
	})
	if err != nil {
		errString, statusCode := GetError(err)
		ctx.JSON(statusCode, errString)
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type ListAccountParams struct {
	Owner      string `form:"owner" binding:"required"`
	PageNumber int32  `form:"page_number" binding:"required,min=1"`
	PageSize   int32  `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listAccounts(ctx *gin.Context) {
	var req ListAccountParams
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	accounts, err := server.store.ListAccounts(ctx, db.ListAccountsParams{
		Limit:  req.PageSize,
		Offset: (req.PageNumber - 1) * req.PageSize,
		Owner:  req.Owner,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, accounts)
}
