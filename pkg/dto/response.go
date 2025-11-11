package dto

func ErrorResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"error": message,
	}
}
