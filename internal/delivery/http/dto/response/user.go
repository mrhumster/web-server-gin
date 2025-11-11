package response

import (
	"time"

	"github.com/mrhumster/web-server-gin/internal/domain/models"
)

type UserResponse struct {
	ID        uint      `json:"id"`
	Login     string    `json:"login"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UsersListReponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total"`
	Page  int64          `json:"page"`
	Limit int64          `json:"limit"`
}

func (u *UserResponse) FillInTheModel(m *models.User) {
	u.ID = m.ID
	u.Login = *m.Login
	u.Email = *m.Email
	u.CreatedAt = m.CreatedAt
	u.UpdatedAt = m.UpdatedAt
	u.Name = *m.Name
	u.LastName = *m.LastName
}
