package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mrhumster/web-server-gin/dto/response"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse("Authorization header required"))
			c.Abort()
			return
		}
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse("Authorization header format must be Bearer {token}"))
			c.Abort()
			return
		}
		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse("Invalid token"))
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse("invalid token claims"))
			c.Abort()
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse("invalid user id in token"))
			c.Abort()
			return
		}

		userEmail, ok := claims["email"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse("invalid email in token"))
			c.Abort()
			return
		}
		c.Set("userID", uint64(userID))
		c.Set("userEmail", userEmail)
		c.Next()
	}
}
