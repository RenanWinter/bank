package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	db "github.com/RenanWinter/bank/db/sqlc"
	"github.com/RenanWinter/bank/util/config"
	"github.com/RenanWinter/bank/util/cript"
)

type signUpRequest struct {
	Name     string `json:"name" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (server *Server) signUp(ctx *gin.Context) {
	var request signUpRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		handleError(ctx, err, gin.H{})
		return
	}

	arg := db.CreateUserParams{
		Email:    request.Email,
		Username: request.Email,
		Name:     request.Name,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		handleError(ctx, err, gin.H{})
		return
	}
	password, err := cript.HashPassword(request.Password, bcrypt.DefaultCost)
	if err != nil {
		handleError(ctx, err, gin.H{})
		return
	}
	_, err = server.store.CreateCredential(ctx, db.CreateCredentialParams{
		UserID:   user.ID,
		Password: password,
	})

	if err != nil {
		handleError(ctx, err, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

type signInRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type signInResponse struct {
	AccessToken string  `json:"access_token"`
	User        db.User `json:"user"`
}

func (server *Server) signIn(ctx *gin.Context) {
	var request signInRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		handleError(ctx, err, gin.H{})
		return
	}

	user, err := server.store.GetUserByUsername(ctx, request.Username)
	if err != nil {
		unauthorizedRequestError(ctx, err, gin.H{"message": "invalid credentials"})
		return
	}

	credential, err := server.store.GetUserActiveCredential(ctx, user.ID)
	if err != nil {
		unauthorizedRequestError(ctx, err, gin.H{"message": "invalid credentials"})
		return
	}

	err = cript.CheckPassword(request.Password, credential.Password)
	if err != nil {
		unauthorizedRequestError(ctx, err, gin.H{"message": "invalid credentials"})
		return
	}

	token, err := server.tokenMaker.CreateToken(user.Uuid.String(), config.Env.TokenDuration)
	if err != nil {
		handleError(ctx, err, gin.H{})
		return
	}

	response := signInResponse{
		AccessToken: token,
		User:        user,
	}

	ctx.JSON(http.StatusOK, response)
}

type listUsersRequest struct {
	Page  int32 `form:"page" binding:"required,min=1"`
	Limit int32 `form:"limit"`
}

func (server *Server) listUsers(ctx *gin.Context) {

	var request listUsersRequest
	if err := ctx.ShouldBindQuery(&request); err != nil {
		handleError(ctx, err, gin.H{})
		return
	}

	arg := db.ListUsersParams{
		Limit:  request.Limit,
		Offset: (request.Page - 1) * request.Limit,
	}

	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		handleError(ctx, err, gin.H{})
		return
	}

	ctx.JSON(http.StatusOK, users)
}
