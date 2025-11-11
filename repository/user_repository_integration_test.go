package repository

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/mrhumster/web-server-gin/internal/delivery/http/dto/request"
	"github.com/mrhumster/web-server-gin/models"
	"github.com/mrhumster/web-server-gin/tests/testutils"
	"github.com/stretchr/testify/assert"
)

func strPtr(s string) *string {
	return &s
}

func generateUniqueUser() models.User {
	timestamp := time.Now().UnixNano()
	uniqueLogin := fmt.Sprintf("testuser_%d", timestamp)
	uniqueEmail := fmt.Sprintf("testuser_%d@test.local", timestamp)

	return models.User{
		Name:         strPtr("Name"),
		LastName:     strPtr("Lastname"),
		Login:        &uniqueLogin,
		PasswordHash: strPtr("#########"),
		Email:        &uniqueEmail,
	}
}

func TestUserRepository_CreateRepo(t *testing.T) {
	db := testutils.GetTestDB()
	defer testutils.CleanTestDatabase()
	repo := NewUserRepository(db)
	ctx := context.Background()

	t.Run("Create and Read user", func(t *testing.T) {
		user := generateUniqueUser()

		id, err := repo.CreateUser(ctx, user)
		assert.NoError(t, err)
		assert.Greater(t, id, uint(0))

		foundUser, err := repo.ReadUserByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, user.Login, foundUser.Login)
		assert.Equal(t, user.Email, foundUser.Email)

		foundUser, err = repo.ReadUserByID(ctx, id+1)
		assert.Error(t, err)
	})

	t.Run("Create user with nil login", func(t *testing.T) {
		user := generateUniqueUser()
		user.Login = nil

		_, err := repo.CreateUser(ctx, user)
		assert.Error(t, err)
	})

	t.Run("Create user with duplicate login", func(t *testing.T) {
		user := generateUniqueUser()

		id, err := repo.CreateUser(ctx, user)
		assert.NoError(t, err)
		assert.Greater(t, id, uint(0))

		user2 := user
		user2.Email = strPtr("different@test.local")
		_, err = repo.CreateUser(ctx, user2)
		log.Printf("⚠️ EBANINA %s", err.Error())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicate")
	})
}

func TestUserRepository_ReadList(t *testing.T) {
	db := testutils.GetTestDB()
	defer testutils.CleanTestDatabase()
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
	defer testutils.CleanTestDatabase()
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

func TestUserRepository_ReadByEmail(t *testing.T) {
	db := testutils.GetTestDB()
	defer testutils.CleanTestDatabase()
	repo := NewUserRepository(db)
	ctx := context.Background()
	t.Run("Create user and get user by email", func(t *testing.T) {
		user := models.User{
			Login:        strPtr("billy"),
			Email:        strPtr("billy@test.local"),
			PasswordHash: strPtr("***1234***"),
		}
		id, err := repo.CreateUser(ctx, user)
		assert.NotEmpty(t, id)
		assert.NoError(t, err)
		u, _ := repo.ReadUserByID(ctx, id)
		u, err = repo.ReadUserByEmail(ctx, "billy@test.local")
		assert.NoError(t, err)
		assert.Equal(t, id, u.ID)
		assert.Equal(t, user.Email, u.Email)
	})
}

func TestUserRepository_UserExist(t *testing.T) {
	db := testutils.GetTestDB()
	defer testutils.CleanTestDatabase()
	repo := NewUserRepository(db)
	ctx := context.Background()
	t.Run("Create user and check that user exisr", func(t *testing.T) {
		user := models.User{
			Login:        strPtr("billy"),
			Email:        strPtr("billy@test.local"),
			PasswordHash: strPtr("***1234***"),
		}
		id, err := repo.CreateUser(ctx, user)
		if err != nil {
			t.Error(err.Error())
		}
		exists := repo.Exists(ctx, id)
		assert.True(t, exists)
	})
}
