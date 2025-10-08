package service

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/google/uuid"
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
	if db == nil {
		t.Fatal("Test DB is nil")
	}
	defer testutils.CleanTestDatabase()

	repo := repository.NewUserRepository(db)
	if repo == nil {
		t.Fatal("UserRepository is nil")
	}

	service := NewUserService(repo)
	if service == nil {
		t.Fatal("UserService is nil")
	}

	ctx := context.Background()

	t.Run("Create user", func(t *testing.T) {

		login := fmt.Sprintf("testuser-%s", uuid.New().String())
		password := "password"
		name := "name"
		lastName := "lastname"
		email := fmt.Sprintf("%s@test.local", login)

		user := models.User{
			Login:        &login,
			PasswordHash: &password,
			Name:         &name,
			LastName:     &lastName,
			Email:        &email,
		}

		if user.Login == nil || user.PasswordHash == nil || user.Email == nil {
			t.Fatal("User fields contain nil pointers")
		}

		id, err := service.CreateUser(ctx, user)
		if err != nil {
			t.Fatalf("Failed to create user: %v", err)
		}
		assert.NoError(t, err)
		assert.NotZero(t, id)

		_, err = service.CreateUser(ctx, user)
		log.Printf("🚩 %v", err.Error())
		assert.Error(t, err)

		changedEmail := "changed@test.local"
		upd := request.UpdateUserRequest{
			Email: &changedEmail,
		}

		_, err = service.UpdateUser(ctx, id, upd)
		assert.NoError(t, err)

		u, err := service.ReadUser(ctx, id)
		assert.NoError(t, err)
		if u != nil && u.Email != nil {
			assert.Equal(t, "changed@test.local", *u.Email)
		} else {
			t.Fatal("ReadUser returned nil user or nil email")
		}

		err = service.DeleteUser(ctx, id)
		assert.NoError(t, err)

		_, _, err = service.ReadUserList(ctx, 100, 1)
		assert.NoError(t, err)
	})
}

func TestUserService_Validate(t *testing.T) {
	db := testutils.GetTestDB()
	defer testutils.CleanTestDatabase()
	repo := repository.NewUserRepository(db)
	service := NewUserService(repo)
	ctx := context.Background()
	t.Run("Password validate success", func(t *testing.T) {

		user := models.User{
			Login:        strPtr("testuser"),
			PasswordHash: strPtr("password"),
			Name:         strPtr("name"),
			LastName:     strPtr("lastname"),
			Email:        strPtr("email@test.local"),
		}
		id, err := service.CreateUser(ctx, user)
		u, err := service.ReadUser(ctx, id)
		u, err = service.GetUserByEmail(ctx, "email@test.local")
		assert.NoError(t, err)
		assert.Equal(t, id, u.ID)
	})
}
