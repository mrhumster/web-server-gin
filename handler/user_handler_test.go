package handler_test

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
	"github.com/mrhumster/web-server-gin/dto/response"
	"github.com/mrhumster/web-server-gin/routes"
	"github.com/mrhumster/web-server-gin/tests/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func strPtr(s string) *string {
	return &s
}

func setupTest() (*gin.Engine, *gorm.DB) {
	db := testutils.GetTestDB()
	router := routes.SetupRoutes(db, "debug")
	return router, db
}

func createUserRequest(router *gin.Engine, login, password, email string) *httptest.ResponseRecorder {
	user := request.UserRequest{
		Login:    login,
		Password: password,
		Email:    email,
		Name:     login,
		LastName: login,
	}

	userJSON, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(userJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	var body map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &body)
	log.Printf("➕ createUserRequest")
	for v, k := range body {
		log.Printf("%v: %v", v, k)
	}
	return resp
}

func LoginAndGetToken(router *gin.Engine, email, password string) string {
	body := request.LoginRequest{Email: email, Password: password}
	bodyJson, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(bodyJson))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	var response response.LoginResponse
	json.Unmarshal(resp.Body.Bytes(), &response)
	return response.Token
}

func TestUserHandler_Success(t *testing.T) {
	router, _ := setupTest()
	defer testutils.CleanTestDatabase()

	user := request.UserRequest{
		Login:    "testuser9",
		Password: "testuser",
		Email:    "testuser9@test.local",
	}

	userJSON, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer(userJSON))
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
	router, _ := setupTest()
	defer testutils.CleanTestDatabase()
	resp1 := createUserRequest(router, "testuser", "password1", "testuser@test.local")
	assert.Equal(t, http.StatusCreated, resp1.Code)

	resp2 := createUserRequest(router, "testuser", "password2", "testuser@test.local")
	assert.Equal(t, http.StatusConflict, resp2.Code)

	var response map[string]any
	json.Unmarshal(resp2.Body.Bytes(), &response)

}

func TestUserHandler_EmptyPassword(t *testing.T) {
	router, _ := setupTest()
	defer testutils.CleanTestDatabase()
	resp1 := createUserRequest(router, "testuser", "", "testuser@test.local")

	assert.Equal(t, http.StatusBadRequest, resp1.Code)

	var response map[string]any
	json.Unmarshal(resp1.Body.Bytes(), &response)
	assert.Contains(t, response, "errors")
	assert.Equal(t, response["errors"], map[string]interface{}{"Password": "This field is required"})
}

func TestUserHandler_InvalidDate(t *testing.T) {
	router, _ := setupTest()
	defer testutils.CleanTestDatabase()
	userJSON := `{"login": 123456, "password": 123456}`
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBuffer([]byte(userJSON)))
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var response map[string]any
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response, "error")
	log.Printf("Invalid Data error: %s", response["error"])
}

func TestUserHandler_ReadUser_InvalidID(t *testing.T) {
	router, _ := setupTest()
	defer testutils.CleanTestDatabase()
	resp := createUserRequest(router, "testuser", "password", "testuser@test.local")
	loginResponse, _ := AuthByLogin(router, "testuser@test.local", "password")
	var repsonse1 map[string]any
	json.Unmarshal(resp.Body.Bytes(), &repsonse1)
	userID := repsonse1["id"]
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/users/%.0f", userID), nil)
	req.Header.Set("Authorization", loginResponse.GetTokenAsBearerHeader())
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	var response map[string]any
	json.Unmarshal(resp.Body.Bytes(), &response)
	assert.Contains(t, response, "id")
	assert.Equal(t, response["id"], userID)
	assert.Contains(t, response, "email")
	assert.NotEmpty(t, response["email"])

	req, _ = http.NewRequest("GET", "/api/users/-1", nil)
	req.Header.Set("Authorization", loginResponse.GetTokenAsBearerHeader())
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestUserHandler_DeleteUser(t *testing.T) {
	router, _ := setupTest()
	defer testutils.CleanTestDatabase()
	resp := createUserRequest(router, "testuser", "password", "testuser@test.local")
	loginReponse, _ := AuthByLogin(router, "testuser@test.local", "password")
	var response map[string]any
	json.Unmarshal(resp.Body.Bytes(), &response)

	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/users/%.0f", response["id"]), nil)
	req.Header.Set("Authorization", loginReponse.GetTokenAsBearerHeader())
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusNoContent, resp.Code)

	req, _ = http.NewRequest("DELETE", "/api/users/-1", nil)
	req.Header.Add("Authorization", loginReponse.GetTokenAsBearerHeader())
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestUserHandler_ReadUsers(t *testing.T) {
	page := float64(1)
	limit := float64(5)
	total := 10
	router, _ := setupTest()
	defer testutils.CleanTestDatabase()
	for i := range total {
		createUserRequest(router, fmt.Sprintf("testuser%d", i), "password", fmt.Sprintf("testuser%d@test.local", i))
	}

	req, _ := http.NewRequest("GET", "/api/users", nil)
	loginResponse, _ := AuthByLogin(router, "testuser1@test.local", "password")
	req.Header.Set("Authorization", loginResponse.GetTokenAsBearerHeader())

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	var respMap1 map[string]any
	json.Unmarshal(resp.Body.Bytes(), &respMap1)
	req, _ = http.NewRequest("GET", fmt.Sprintf("/api/users?page=%.0f&limit=%.0f", page, limit), nil)
	req.Header.Set("Authorization", loginResponse.GetTokenAsBearerHeader())
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
	for _, u := range users {

		userMap, ok := u.(map[string]interface{})
		require.True(t, ok, "each user should be an object")

		require.NotEmpty(t, userMap["email"])
		require.NotEmpty(t, userMap["login"])
	}
}

func TestUserHandler_ReadUsers_QueryValidate(t *testing.T) {
	router, _ := setupTest()
	defer testutils.CleanTestDatabase()
	createUserRequest(router, "testuser", "password", "testuser@test.local")
	loginResponse, _ := AuthByLogin(router, "testuser@test.local", "password")
	req, _ := http.NewRequest("GET", "/api/users?page=-1&limit=s", nil)
	req.Header.Set("Authorization", loginResponse.GetTokenAsBearerHeader())
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
	req, _ = http.NewRequest("GET", "/api/users?page=1&limit=-10", nil)
	req.Header.Set("Authorization", loginResponse.GetTokenAsBearerHeader())
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	req, _ = http.NewRequest("GET", "/api/users", nil)
	req.Header.Set("Authorization", loginResponse.GetTokenAsBearerHeader())
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestUserHandler_UserResponseIncludeEmail(t *testing.T) {
	router, _ := setupTest()
	defer testutils.CleanTestDatabase()
	resp := createUserRequest(router, "testuser1", "password", "testuser1@test.local")
	loginResponse, _ := AuthByLogin(router, "testuser1@test.local", "password")
	var respMap map[string]any
	json.Unmarshal(resp.Body.Bytes(), &respMap)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/users/%.0f", respMap["id"]), nil)
	req.Header.Set("Authorization", loginResponse.GetTokenAsBearerHeader())
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	json.Unmarshal(resp.Body.Bytes(), &respMap)

	assert.Contains(t, respMap, "email")
}

func TestUserHandler_EmailIsUniq(t *testing.T) {
	router, _ := setupTest()
	defer testutils.CleanTestDatabase()
	createUserRequest(router, "testuser1", "password", "testuser1@test.local")
	resp := createUserRequest(router, "testuser2", "password", "testuser1@test.local")
	assert.Equal(t, http.StatusConflict, resp.Code)
	var respMap map[string]any
	json.Unmarshal(resp.Body.Bytes(), &respMap)
	assert.Contains(t, respMap, "error")
	assert.Equal(t, respMap["error"], "user already exists")

}

func TestUserHandler_Response_IncludeAllFields(t *testing.T) {
	router, _ := setupTest()
	defer testutils.CleanTestDatabase()
	resp := createUserRequest(router, "testuser1", "password", "testuser1@test.local")
	loginResponse, _ := AuthByLogin(router, "testuser1@test.local", "password")
	var body map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &body)
	userId := body["id"]
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/users/%.0f", userId), nil)
	req.Header.Set("Authorization", loginResponse.GetTokenAsBearerHeader())
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	json.Unmarshal(resp.Body.Bytes(), &body)
	assert.Contains(t, body, "name")
	assert.NotEmpty(t, body["name"])
	assert.Contains(t, body, "last_name")
	assert.NotEmpty(t, body["last_name"])
}

func TestUserHandler_UpdateUser(t *testing.T) {
	router, _ := setupTest()
	defer testutils.CleanTestDatabase()
	resp := createUserRequest(router, "testuser1", "password", "testuser1@test.local")
	loginResponse, _ := AuthByLogin(router, "testuser1@test.local", "password")
	var body, updatedBody map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &body)
	var userUpdate request.UpdateUserRequest
	userUpdate.Name = strPtr("Larry")
	userUpdate.LastName = strPtr("Coat")
	userJson, _ := json.Marshal(userUpdate)
	req, _ := http.NewRequest(
		"PATCH",
		fmt.Sprintf("/api/users/%.0f", body["id"]),
		bytes.NewBuffer(userJson),
	)
	req.Header.Set("Authorization", loginResponse.GetTokenAsBearerHeader())
	resp = httptest.NewRecorder()
	router.ServeHTTP(resp, req)
	json.Unmarshal(resp.Body.Bytes(), &updatedBody)
	v, exist := updatedBody["error"]
	if exist {
		log.Printf("⚠️ DEBUG: ERROR %v", v)
		assert.True(t, false)
		return
	}
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, updatedBody, "name")
	assert.Equal(t, updatedBody["name"], "Larry")
	assert.Contains(t, updatedBody, "last_name")
	assert.Equal(t, updatedBody["last_name"], "Coat")
	assert.Contains(t, updatedBody, "login")
	assert.Equal(t, updatedBody["login"], "testuser1")
	assert.Contains(t, updatedBody, "email")
	assert.Equal(t, updatedBody["email"], "testuser1@test.local")
}
