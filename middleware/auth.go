package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/mrhumster/web-server-gin/dto/response"
	"github.com/mrhumster/web-server-gin/service"
)

func AuthMiddleware(tokenService *service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {

		token := extractToken(c.Request)
		claims, err := tokenService.ValidateToken(token)
		log.Printf("⚠️ AuthMiddleware: CLAIM ERROR %v %v", err, claims)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse("invalid token claims"))
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("claims", claims)
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

func extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return ""
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}
