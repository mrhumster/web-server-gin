package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mrhumster/web-server-gin/dto/request"
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
	user := request.UserRequest{
		Login:    login,
		Password: password,
		Email:    fmt.Sprintf("%s@test.local", login),
	}

	userJSON, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	return resp
}

func createRouter(userHandler *UserHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users", userHandler.CreateUser)
	router.GET("/users/:id", userHandler.ReadUser)
	router.DELETE("/users/:id", userHandler.Delete)
	router.GET("/users", userHandler.ReadUsers)
	return router
}

func TestUserHandler_Success(t *testing.T) {
	handler, db := setupTest()
	defer db.Exec("DELETE FROM users")

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users", handler.CreateUser)

	user := request.UserRequest{
		Login:    "testuser",
		Password: "testuser",
		Email:    "testuser@test.local",
	}

	userJSON, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusCreated, resp.Code)

	var response map[string]any

	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response, "id")
	assert.NotZero(t, response["id"])
}

func TestUserHandler_DiplucateLogin(t *testing.T) {
	handler, db := setupTest()
	defer db.Exec("DELETE FROM users")
	router := createRouter(handler)
	resp1 := createUserRequest(router, "testuser", "password1")
	assert.Equal(t, http.StatusCreated, resp1.Code)

	resp2 := createUserRequest(router, "testuser", "password2")
	assert.Equal(t, http.StatusConflict, resp2.Code)

	var response map[string]any
	json.Unmarshal(resp2.Body.Bytes(), &response)

}

func TestUserHandler_EmptyPassword(t *testing.T) {
	handler, db := setupTest()
	defer db.Exec("DELETE FROM users")
	router := createRouter(handler)
	resp1 := createUserRequest(router, "testuser", "")

	assert.Equal(t, http.StatusBadRequest, resp1.Code)

	var response map[string]any
	json.Unmarshal(resp1.Body.Bytes(), &response)
	assert.Contains(t, response, "errors")
	assert.Equal(t, response["errors"], map[string]interface{}{"Password": "This field is required"})
}

func TestUserHandler_InvalidDate(t *testing.T) {
	handler, db := setupTest()
	defer db.Exec("DELETE FROM users")
	router := createRouter(handler)
	userJSON := `{"login": 123456, "password": 123456}`
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer([]byte(userJSON)))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]any
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response, "error")
	log.Printf("Invalid Data error: %s", response["error"])
}

func TestUserHandler_ReadUser_InvalidID(t *testing.T) {
	handler, db := setupTest()
	defer db.Exec("DELETE FROM users")
	router := createRouter(handler)
	resp := createUserRequest(router, "testuser", "password")

	var repsonse1 map[string]any
	json.Unmarshal(resp.Body.Bytes(), &repsonse1)
	userID := repsonse1["id"]
	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%.0f", userID), nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var response map[string]any
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response, "id")
	assert.Equal(t, response["id"], userID)

	req, _ = http.NewRequest("GET", "/users/-1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestUserHandler_DeleteUser(t *testing.T) {
	handler, db := setupTest()
	defer db.Exec("DELETE FROM users")
	router := createRouter(handler)
	resp := createUserRequest(router, "testuser", "password")

	var response map[string]any
	json.Unmarshal(resp.Body.Bytes(), &response)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/users/%.0f", response["id"]), nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNoContent, resp.Code)

	req, _ = http.NewRequest("DELETE", "/users/-1", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestUserHandler_ReadUsers(t *testing.T) {
	page := float64(1)
	limit := float64(5)
	total := 50
	handler, db := setupTest()
	defer db.Exec("DELETE FROM users")
	router := createRouter(handler)
	for i := range total {
		createUserRequest(router, fmt.Sprintf("testuser%d", i), "password")
	}
	req, _ := http.NewRequest("GET", "/users", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	var respMap1 map[string]any
	json.Unmarshal(resp.Body.Bytes(), &respMap1)
	req, _ = http.NewRequest("GET", fmt.Sprintf("/users?page=%.0f&limit=%.0f", page, limit), nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	var respMap map[string]any
	json.Unmarshal(resp.Body.Bytes(), &respMap)
	assert.Contains(t, respMap, "total")
	assert.Contains(t, respMap, "page")
	assert.Contains(t, respMap, "limit")
	assert.Contains(t, respMap, "users")
	assert.Equal(t, float64(total), respMap["total"])
	assert.Equal(t, page, respMap["page"])
	users, ok := respMap["users"].([]interface{})
	assert.True(t, ok, "users should be an array")
	assert.Equal(t, int(limit), len(users))
}

func TestUserHandler_ReadUsers_QueryValidate(t *testing.T) {
	handler, db := setupTest()
	defer db.Exec("DELETE FROM users")
	router := createRouter(handler)
	req, _ := http.NewRequest("GET", "/users?page=-1&limit=s", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	req, _ = http.NewRequest("GET", "/users?page=1&limit=-10", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	req, _ = http.NewRequest("GET", "/users", nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestUserHandler_UserResponseIncludeEmail(t *testing.T) {
	handler, db := setupTest()
	defer db.Exec("DELETE FROM users")
	router := createRouter(handler)
	resp := createUserRequest(router, "testuser1", "testuser1")

	var respMap map[string]any
	json.Unmarshal(resp.Body.Bytes(), &respMap)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/users/%.0f", respMap["id"]), nil)
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	respMap = map[string]any{}
	json.Unmarshal(resp.Body.Bytes(), &respMap)

	assert.Contains(t, respMap, "email")
}

func TestUserHandler_EmailIsUniq(t *testing.T) {
	return
}
