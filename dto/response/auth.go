package response

import (
	"fmt"
	"time"
)

type LoginResponse struct {
	Token   string       `json:"token"`
	Expires time.Time    `json:"expires"`
	User    UserResponse `json:"user"`
}

func (l *LoginResponse) GetTokenAsBearerHeader() string {
	return fmt.Sprintf("Bearer %s", l.Token)
}
