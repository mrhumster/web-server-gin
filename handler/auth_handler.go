package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mrhumster/web-server-gin/dto/request"
	"github.com/mrhumster/web-server-gin/dto/response"
	"github.com/mrhumster/web-server-gin/models"
	"github.com/mrhumster/web-server-gin/service"
)

type AuthHandler struct {
	UserService  *service.UserService
	TokenService *service.TokenService
	JwtSecret    string
}

func NewAuthHandler(userService *service.UserService, tokenService *service.TokenService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		UserService:  userService,
		TokenService: tokenService,
		JwtSecret:    jwtSecret,
	}
}

func (a *AuthHandler) Login(c *gin.Context) {
	var (
		req request.LoginRequest
		t   string
		u   *models.User
		err error
	)
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if u, err = a.UserService.ValidateUser(c, req.Email, req.Password); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse(err.Error()))
		return
	}

	if t, err = a.TokenService.GenerateToken(u); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("failed to generate token"))
	}

	var user response.UserResponse
	user.FillInTheModel(u)
	c.JSON(http.StatusOK, response.LoginResponse{
		Token: t,
		User:  user,
	})
}

func (a *AuthHandler) Logout(c *gin.Context) {
	// TODO: Реализовать логику logout (добавление токена в blacklist и т.д.)
	c.JSON(http.StatusOK, response.SuccessResponse("Logged out successfully"))
}
