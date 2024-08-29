package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type GetAccountParams struct {
	ID int64 `uri:"id" json:"id"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var accountID GetAccountParams
	if err := ctx.ShouldBindQuery(&accountID); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	account, err := server.store.GetAccount(ctx, accountID.ID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, errorResponse(err))
	}

	ctx.JSON(http.StatusOK, account)
}
