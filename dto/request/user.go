package request

type UserRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=20"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email" binding:"required,email"`
}
