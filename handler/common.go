package handler

import "github.com/go-playground/validator/v10"

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
