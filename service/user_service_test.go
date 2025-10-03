package service

import (
	"context"
	"testing"

	"github.com/mrhumster/web-server-gin/dto/request"
	"github.com/mrhumster/web-server-gin/models"
	"github.com/mrhumster/web-server-gin/repository"
	"github.com/mrhumster/web-server-gin/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func strPtr(s string) *string {
	return &s
}

func TestUserService_Create(t *testing.T) {
	db := testutils.GetTestDB()
	defer db.Exec("DELETE FROM users")
	repo := repository.NewUserRepository(db)
	service := NewUserService(repo)
	ctx := context.Background()
	t.Run("Create user", func(t *testing.T) {
		user := models.User{
			Login:        strPtr("testuser"),
			PasswordHash: strPtr("password"),
			Name:         strPtr("name"),
			LastName:     strPtr("lastname"),
			Email:        strPtr("email@test.local"),
		}
		id, err := service.CreateUser(ctx, user)
		assert.NoError(t, err)
		_, err = service.CreateUser(ctx, user)
		assert.Error(t, err)
		upd := request.UpdateUserRequest{
			Email: strPtr("changed@test.local"),
		}
		_, err = service.UpdateUser(ctx, id, upd)
		u, err := service.ReadUser(ctx, id)
		assert.Equal(t, *u.Email, "changed@test.local")
		assert.NoError(t, err)
		err = service.DeleteUser(ctx, id)
		assert.NoError(t, err)
		_, _, err = service.ReadUserList(ctx, 100, 1)
		assert.NoError(t, err)
	})
}
