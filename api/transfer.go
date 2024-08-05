package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	db "github.com/RenanWinter/bank/db/sqlc"
)

type transferRequest struct {
	FromAccountId int64   `json:"from_account_id" binding:"required,min=1"`
	ToAccountId   int64   `json:"to_account_id" binding:"required,min=1"`
	Amount        float64 `json:"amount" binding:"required,gt=0"`
	Currency      string  `json:"currency" binding:"required,currency"` //currency is a custo validation. see server.go for registration
}
type transferResponse struct {
	Transfer db.Transfer `json:"transfer"`
	Movement db.Movement `json:"movement"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var request transferRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		handleError(ctx, err, gin.H{})
		return
	}

	user, valid := server.validAccount(ctx, request.FromAccountId)
	if !valid {
		return
	}

	loggedUser := getLoggedUser(ctx)
	if user.ID != loggedUser.ID {
		err := fmt.Errorf("Origin account does not belong to the authenticated user")
		handleError(ctx, err, gin.H{"code": http.StatusUnauthorized})
		return
	}

	_, valid = server.validAccount(ctx, request.ToAccountId)
	if !valid {
		return
	}

	arg := db.TransferParams{
		FromAccountID: request.FromAccountId,
		ToAccountID:   request.ToAccountId,
		Amount:        request.Amount,
	}

	result, err := server.store.Transfer(ctx, arg)

	if err != nil {
		handleError(ctx, err, gin.H{})
		return
	}

	data := transferResponse{
		Transfer: result.Transfer,
		Movement: result.FromMovement,
	}
	ctx.JSON(http.StatusOK, data)
}

func (server *Server) validAccount(ctx *gin.Context, accountID int64) (*db.User, bool) {
	account, err := server.store.GetAccount(ctx, accountID)
	if err != nil {
		handleError(ctx, err, gin.H{})
		return nil, false
	}

	user, err := server.store.GetUserById(ctx, account.OwnerID)

	if err != nil {
		err = fmt.Errorf("error getting user from account: %v", account.Name)
		handleError(ctx, err, gin.H{})
		return nil, false
	}

	return &user, true
}
