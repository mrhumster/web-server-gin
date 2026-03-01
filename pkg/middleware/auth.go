package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/mrhumster/web-server-gin/pkg/auth"
	"github.com/mrhumster/web-server-gin/pkg/dto"
)

func Authorize(client auth.PermissionClient, obj, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userUUID := c.MustGet("userID").(uuid.UUID)
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
