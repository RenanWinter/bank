package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	db "github.com/RenanWinter/bank/db/sqlc"
)

var (
	errAccountNotFromUser = errors.New("account does not belong to the authenticated user")
)

type createAccountRequest struct {
	Name          string  `json:"name" binding:"required"`
	OwnerID       int64   `json:"owner_id" binding:"required"`
	AccountTypeID int64   `json:"account_type_id" binding:"required"`
	Balance       float64 `json:"balance" binding:"required"`
}

func (server *Server) createAccount(ctx *gin.Context) {

	var request createAccountRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		handleError(ctx, err, gin.H{})
		return
	}

	arg := db.CreateAccountParams{
		Name:          request.Name,
		OwnerID:       request.OwnerID,
		AccountTypeID: request.AccountTypeID,
		Balance:       request.Balance,
	}

	user := getLoggedUser(ctx)

	if user.ID != arg.OwnerID {
		handleError(ctx, errAccountNotFromUser, gin.H{"code": http.StatusForbidden})
		return
	}

	account, err := server.store.CreateAccount(ctx, arg)

	if err != nil {
		handleError(ctx, err, gin.H{})
		return
	}

	ctx.JSON(http.StatusCreated, account)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccount(ctx *gin.Context) {
	var request getAccountRequest
	if err := ctx.ShouldBindUri(&request); err != nil {
		badRequestError(ctx, err, gin.H{})
		return
	}

	account, err := server.store.GetAccount(ctx, request.ID)
	if err != nil {
		handleError(ctx, err, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, account)
}
