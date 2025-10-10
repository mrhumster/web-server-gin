package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func getErrorMessage(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	case "gte":
		return "Value is too small"
	case "lte":
		return "Value is too large"
	default:
		return fieldError.Error()
	}
}

func GetUserIDFromContext(c *gin.Context) (uint64, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("user ID not found in context")
	}
	id, ok := userID.(uint64)
	if !ok {
		return 0, errors.New("invalid user ID type in context")
	}

	return id, nil
}

func GetUserEmailFromContext(c *gin.Context) (string, error) {
	email, exists := c.Get("userEmail")
	if !exists {
		return "", errors.New("user email not found in context")
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", errors.New("invalid user email type in context")
	}

	return emailStr, nil
}
