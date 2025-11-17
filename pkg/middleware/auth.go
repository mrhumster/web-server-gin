package middleware

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrhumster/web-server-gin/pkg/auth"
	"github.com/mrhumster/web-server-gin/pkg/dto"
)

func Authorize(client *auth.PermissionClient, obj, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse("User hasn't logged in yet"))
			return
		}
		resourceID := c.Param("id")

		fullResource := obj
		if resourceID != "" {
			fullResource = fmt.Sprintf("%s/%s", obj, resourceID)
		}
		fmt.Printf("🐞 Authorize middleware: sub: %s; obj: %s; act: %s\n", userID, fullResource, act)
		ok, err := client.Client.CheckPermission(c.Request.Context(), userID.(string), fullResource, act)
		if err != nil {
			log.Fatal("⚠️ Authorize middleware error: ", err)
			return
		}
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse("Access denied"))
		}
		c.Next()
	}
}
