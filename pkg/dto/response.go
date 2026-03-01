package dto

func ErrorResponse(message string) map[string]any {
	return map[string]any{
		"error": message,
	}
}
