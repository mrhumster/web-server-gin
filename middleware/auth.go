package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
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

		userRole, ok := claims["role"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse("invalid role in token"))
			return
		}
		c.Set("userID", uint64(userID))
		c.Set("userEmail", userEmail)
		c.Set("role", userRole)
		c.Next()
	}
}

func Authorize(obj string, act string, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDinterface, exists := c.Get("userID")
		userID := fmt.Sprintf("%v", userIDinterface)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse("User hasn't logged in yet"))
			return
		}

		resourceID := c.Param("id")

		fullResource := obj
		if resourceID != "" {
			fullResource = fmt.Sprintf("%s/%s", obj, resourceID)
		}

		log.Printf("🚩 Authorize debug! SUB: %s;  OBJ: %s; ACT: %s", userID, fullResource, act)

		if ok, _ := enforcer.Enforce(userID, fullResource, act); !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, response.ErrorResponse("Access denied"))
			return
		}

		c.Next()
	}
}
