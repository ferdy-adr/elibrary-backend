package middleware

import (
	"net/http"
	"strings"

	"github.com/ferdy-adr/elibrary-backend/internal/configs"
	"github.com/ferdy-adr/elibrary-backend/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, model.APIResponse{
				Success: false,
				Message: "Authorization header required",
				Error:   "missing_auth_header",
			})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, model.APIResponse{
				Success: false,
				Message: "Invalid authorization format",
				Error:   "invalid_auth_format",
			})
			c.Abort()
			return
		}

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(configs.Get().JWT.SecretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, model.APIResponse{
				Success: false,
				Message: "Invalid token",
				Error:   "invalid_token",
			})
			c.Abort()
			return
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("user_id", int(claims["user_id"].(float64)))
			c.Set("username", claims["username"].(string))
		}

		c.Next()
	}
}
