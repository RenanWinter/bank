package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"

	db "github.com/RenanWinter/bank/db/sqlc"
	"github.com/RenanWinter/bank/util/config"
	"github.com/RenanWinter/bank/util/token"
)

// Server serve the HTTP requests
type Server struct {
	config     config.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(store db.Store, cfg config.Config) (*Server, error) {

	tokenMaker, err := token.NewPasetoMaker(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		config:     cfg,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/signup", server.signUp)
	router.POST("/signin", server.signIn)

	authRoutes := router.Group("/").Use(authMiddleware(server))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/users", server.listUsers)
	authRoutes.POST("/transfer", server.createTransfer)
	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// Tryies to automatically handle the error response
//
// It's possible to override de code with the key code on detail as int
func handleError(ctx *gin.Context, err error, detail map[string]any) {
	code, detail := errorResponseBody(err, detail)
	ctx.JSON(code, detail)
}

func errorResponseBody(err error, detail map[string]any) (int, gin.H) {
	if err == nil {
		if detail["message"] != nil {
			err = fmt.Errorf(detail["message"].(string))
		} else {
			err = fmt.Errorf("unknown error")
		}
	}

	code := http.StatusInternalServerError
	message := err.Error()
	debugDetail := err.Error()

	if pqErr, ok := err.(*pq.Error); ok {
		errorName := pqErr.Code.Name()
		switch errorName {
		case "foreign_key_violation":
			code = http.StatusForbidden
			message = "This request breaks some relation constraint"
		case "unique_violation":
			code = http.StatusForbidden
			message = "This request is trying to create a duplicated resource"
		default:
			code = http.StatusInternalServerError
		}
	}

	if err == sql.ErrNoRows {
		code = http.StatusNotFound
		message = "No record found"
	}

	// Check for validationErrors
	if ve, ok := err.(validator.ValidationErrors); ok {
		code = http.StatusBadRequest
		message = "Your request is invalid"
		errors := make(map[string]string)
		for _, field := range ve {
			errors[field.Field()] = field.Tag()
		}
		detail["validation_errors"] = errors
	}

	response := gin.H{
		"message": message,
	}

	if config.Env.Debug == true {
		response["error"] = err
		response["debug"] = debugDetail
	}

	detailCode := detail["code"]
	if detailCode != nil {
		switch detailCode.(type) {
		case int:
			code = detailCode.(int)
			delete(detail, "code")
		}
	}

	return code, mergeResponseMap(response, detail)
}

// Alias to handleError with code 400
func badRequestError(ctx *gin.Context, err error, detail map[string]any) {
	detail["code"] = http.StatusBadRequest
	handleError(ctx, err, detail)
}

// Alias to handleError with code 403
func forbiddenRequestError(ctx *gin.Context, err error, detail map[string]any) {
	detail["code"] = http.StatusForbidden
	handleError(ctx, err, detail)
}

// Alias to handleError with code 401
func unauthorizedRequestError(ctx *gin.Context, err error, detail map[string]any) {
	detail["code"] = http.StatusUnauthorized
	handleError(ctx, err, detail)
}

// Alias to handleError with code 404
func notFoundRequestError(ctx *gin.Context, err error, detail map[string]any) {
	detail["code"] = http.StatusNotFound
	handleError(ctx, err, detail)
}

func mergeResponseMap(map1, map2 map[string]any) map[string]any {
	for key, value := range map2 {
		map1[key] = value
	}

	return map1
}

func getLoggedUser(c *gin.Context) db.User {
	return c.MustGet(authenticatedUser).(db.User)
}
