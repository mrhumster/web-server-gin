package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mrhumster/web-server-gin/internal/delivery/http/dto/response"
	"github.com/mrhumster/web-server-gin/internal/service"
	"github.com/mrhumster/web-server-gin/pkg/auth"
	"github.com/mrhumster/web-server-gin/pkg/dto"
)

func Authorize(client auth.PermissionClient, obj, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID := c.MustGet("user").(uuid.UUID)
		resourceID := c.Param("id")

		fullResource := obj
		if resourceID != "" {
			fullResource = fmt.Sprintf("%s/%s", obj, resourceID)
		}

		ok, err := client.CheckPermission(c.Request.Context(), userUUID.String(), fullResource, act)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse("⚠️ Authorize middleware error"))
			return
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse("Access denied"))
		}
		c.Next()
	}
}

func AuthMiddleware(tokenService *service.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c.Request)
		claims, err := tokenService.ValidateAccessToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse("invalid token claims"))
			c.Abort()
			return
		}
		userUUID, err := uuid.Parse(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.ErrorResponse("error parse user id in auth middleware"))
		}
		c.Set("user", userUUID)
		c.Set("claims", claims)
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
