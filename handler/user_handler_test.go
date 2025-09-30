package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mrhumster/web-server-gin/models"
	"github.com/mrhumster/web-server-gin/repository"
	"github.com/mrhumster/web-server-gin/service"
	"github.com/mrhumster/web-server-gin/tests/testutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupTest() (*UserHandler, *gorm.DB) {
	db := testutils.GetTestDB()
	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userHandler := NewUserHandler(userService)
	return userHandler, db
}

func createUserRequest(router *gin.Engine, login, password string) *httptest.ResponseRecorder {
	user := models.User{
		Login:    login,
		Password: password,
	}

	userJSON, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func TestUserService_Success(t *testing.T) {
	handler, db := setupTest()
	defer db.Exec("DELETE FROM users")

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users", handler.CreateUser)

	user := models.User{
		Login:    "testuser",
		Password: "testuser",
	}

	userJSON, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var response map[string]interface{}

	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response, "id")
	assert.NotZero(t, response["id"])
}

func TestUserService_DiplucateLogin(t *testing.T) {
	handler, db := setupTest()
	defer db.Exec("DELETE FROM users")
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users", handler.CreateUser)
	resp1 := createUserRequest(router, "testuser", "password1")
	assert.Equal(t, http.StatusCreated, resp1.Code)

	resp2 := createUserRequest(router, "testuser", "password2")
	assert.Equal(t, http.StatusConflict, resp2.Code)

	var response map[string]interface{}
	json.Unmarshal(resp2.Body.Bytes(), &response)

	assert.Contains(t, response, "error")
}

func TestUserService_EmptyPassword(t *testing.T) {
	handler, db := setupTest()
	defer db.Exec("DELETE FROM users")
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/users", handler.CreateUser)
	resp1 := createUserRequest(router, "testuser", "")
	assert.Equal(t, http.StatusBadRequest, resp1.Code)

	var response map[string]interface{}
	json.Unmarshal(resp1.Body.Bytes(), &response)
	assert.Contains(t, response, "error")
	assert.Contains(t, response["error"], "password can't be empty")
}
