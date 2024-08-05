package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
	authenticatedUser       = "authenticated_user"
)

var (
	errMissingHeader              = errors.New("Authorization header is not provided")
	errInvalidAuthorizationType   = errors.New("Invalid authorization type")
	errInvalidAuthorizationFormat = errors.New("Invalid authorization header format")
	errUserNotFound               = errors.New("User not found on our records")
	errInvalidSubject             = errors.New("Invalid subject")
)

func authMiddleware(server *Server) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader(authorizationHeaderKey)
		if len(authorization) == 0 {
			_, body := errorResponseBody(errMissingHeader, gin.H{})
			c.AbortWithStatusJSON(http.StatusUnauthorized, body)
			return
		}
		fields := strings.Fields(authorization)
		if len(fields) < 2 {
			_, body := errorResponseBody(errInvalidAuthorizationFormat, gin.H{})
			c.AbortWithStatusJSON(http.StatusUnauthorized, body)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			_, body := errorResponseBody(errInvalidAuthorizationType, gin.H{})
			c.AbortWithStatusJSON(http.StatusUnauthorized, body)
			return
		}

		tokenString := fields[1]
		claims, err := server.tokenMaker.VerifyToken(tokenString)
		if err != nil {
			_, body := errorResponseBody(err, gin.H{})
			c.AbortWithStatusJSON(http.StatusUnauthorized, body)
			return
		}

		c.Set(authorizationPayloadKey, claims)
		uuid, err := uuid.Parse(claims.Sub)
		if err != nil {
			_, body := errorResponseBody(errInvalidSubject, gin.H{})
			c.AbortWithStatusJSON(http.StatusUnauthorized, body)
			return
		}
		user, err := server.store.GetUserByUUID(c, uuid)
		if err != nil {
			_, body := errorResponseBody(errUserNotFound, gin.H{})
			c.AbortWithStatusJSON(http.StatusUnauthorized, body)
			return
		}

		c.Set(authenticatedUser, user)
		c.Next()
	}
}
