package request

type UserRequest struct {
	Login    string `json:"login" binding:"required,min=3,max=100"`
	Password string `json:"password" binding:"required,min=6,max=100"`
	Name     string `json:"name" binding:"max=100"`
	LastName string `json:"last_name" binding:"max=100"`
	Email    string `json:"email" binding:"required,email"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty" binding:"omitempty,max=100"`
	LastName *string `json:"last_name,omitempty" binding:"omitempty,max=100"`
	Email    *string `json:"email,omitempty" binding:"omitempty,email"`
}
