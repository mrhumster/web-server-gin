package repository

import (
	"context"
	"testing"

	"github.com/mrhumster/web-server-gin/dto/request"
	"github.com/mrhumster/web-server-gin/models"
	"github.com/mrhumster/web-server-gin/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func strPtr(s string) *string {
	return &s
}

func TestUserRepository_CreateRepo(t *testing.T) {
	db := testutils.GetTestDB()
	defer db.Exec("DELETE FROM users")
	repo := NewUserRepository(db)
	ctx := context.Background()
	t.Run("Create and Read user", func(t *testing.T) {
		user := models.User{
			Name:         strPtr("Name"),
			LastName:     strPtr("Lastname"),
			Login:        strPtr("testuser"),
			PasswordHash: strPtr("#########"),
			Email:        strPtr("testuser@test.local"),
		}
		id, err := repo.CreateUser(ctx, user)
		assert.NoError(t, err)
		assert.Greater(t, id, uint(0))
		foundUser, err := repo.ReadUserByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, user.Login, foundUser.Login)
		assert.Equal(t, user.Email, foundUser.Email)
		foundUser, err = repo.ReadUserByID(ctx, id+1)
		assert.Error(t, err)
		user = models.User{
			Name:         strPtr("Name"),
			LastName:     strPtr("Lastname"),
			Login:        nil,
			PasswordHash: strPtr("#########"),
			Email:        strPtr("testuser@test.local"),
		}
		_, err = repo.CreateUser(ctx, user)
		assert.Error(t, err)
	})
}

func TestUserRepository_ReadList(t *testing.T) {
	db := testutils.GetTestDB()
	defer db.Exec("DELETE FROM users")
	repo := NewUserRepository(db)
	ctx := context.Background()
	t.Run("Read userlist", func(t *testing.T) {
		user := models.User{
			Name:         strPtr("Name"),
			LastName:     strPtr("Lastname"),
			Login:        strPtr("testuser1"),
			PasswordHash: strPtr("#########"),
			Email:        strPtr("testuser@test.local"),
		}
		repo.CreateUser(ctx, user)
		user = models.User{
			Name:         strPtr("Name"),
			LastName:     strPtr("Lastname"),
			Login:        strPtr("testuser2"),
			PasswordHash: strPtr("#########"),
			Email:        strPtr("testuser2@test.local"),
		}
		repo.CreateUser(ctx, user)
		users, _, err := repo.ReadUserList(ctx, 2, 1)
		assert.NoError(t, err)
		assert.NotEmpty(t, users)
		assert.Equal(t, 2, len(users))
	})

}

func TestUserRepository_UpdateRepo(t *testing.T) {
	db := testutils.GetTestDB()
	defer db.Exec("DELETE FROM users")
	repo := NewUserRepository(db)
	ctx := context.Background()
	t.Run("Update user", func(t *testing.T) {
		user := models.User{
			Name:         strPtr("Name"),
			LastName:     strPtr("Lastname"),
			Login:        strPtr("testuser1"),
			PasswordHash: strPtr("#########"),
			Email:        strPtr("testuser@test.local"),
		}
		id, _ := repo.CreateUser(ctx, user)
		updateReq := request.UpdateUserRequest{
			Name: strPtr("ChangedName"),
		}
		id, err := repo.UpdateUser(ctx, id, updateReq)
		assert.NoError(t, err)
		updatedUser, err := repo.ReadUserByID(ctx, id)
		assert.Equal(t, *updatedUser.Name, "ChangedName")
		assert.Equal(t, *updatedUser.Login, "testuser1")
		err = repo.DeleteUserByID(ctx, id)
		assert.NoError(t, err)
	})
}
