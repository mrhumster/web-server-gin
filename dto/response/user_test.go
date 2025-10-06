package response

import (
	"testing"

	"github.com/mrhumster/web-server-gin/models"
	"github.com/stretchr/testify/assert"
)

func strPtr(s string) *string {
	return &s
}

func TestUserResponse_Fill(t *testing.T) {
	user := models.User{
		Login:        strPtr("billy"),
		PasswordHash: strPtr("***123***"),
		Name:         strPtr("Billy"),
		LastName:     strPtr("Coat"),
		Email:        strPtr("billy@test.local"),
	}
	res := UserResponse{}
	res.FillInTheModel(&user)
	assert.Equal(t, *user.LastName, res.LastName)
	assert.Equal(t, *user.Name, res.Name)
	assert.Equal(t, *user.Login, res.Login)
	assert.Equal(t, *user.Email, res.Email)
}
