package service

import (
	"context"
	"testing"

	"github.com/mrhumster/web-server-gin/internal/domain/models"
	"github.com/mrhumster/web-server-gin/internal/repository"
	"github.com/mrhumster/web-server-gin/tests/testutils"
	"github.com/stretchr/testify/assert"
)

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
	enforcer := testutils.GetEnforcer(db)
	permissionService := NewPermissionService(enforcer)
	service := NewUserService(repo, permissionService)
	if service == nil {
		t.Fatal("UserService is nil")
	}

}

func TestUserService_Validate(t *testing.T) {
	db := testutils.GetTestDB()
	defer testutils.CleanTestDatabase()
	repo := repository.NewUserRepository(db)
	enforcer := testutils.GetEnforcer(db)
	permissionService := NewPermissionService(enforcer)
	service := NewUserService(repo, permissionService)
	ctx := context.Background()
	t.Run("Password validate success", func(t *testing.T) {

		user := models.User{
			PasswordHash: "password",
			Email:        "email@test.local",
		}
		id, err := service.CreateUser(ctx, user)
		assert.NoError(t, err)
		_, err = service.ReadUser(ctx, *id)
		assert.NoError(t, err)
		_, err = service.GetUserByEmail(ctx, "email@test.local")
		assert.NoError(t, err)
	})
}
