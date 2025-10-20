package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mrhumster/web-server-gin/dto/request"
	"github.com/mrhumster/web-server-gin/dto/response"
	"github.com/mrhumster/web-server-gin/models"
	"github.com/mrhumster/web-server-gin/service"
)

type AuthHandler struct {
	UserService *service.UserService
	JwtSecret   string
}

func NewAuthHandler(userService *service.UserService, jwtSecret string) *AuthHandler {
	return &AuthHandler{
		UserService: userService,
		JwtSecret:   jwtSecret,
	}
}

func (a *AuthHandler) Login(c *gin.Context) {
	var (
		req request.LoginRequest
		t   *jwt.Token
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
	t = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(u.ID),
		"role":    &u.Role,
		"email":   &u.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := t.SignedString([]byte(a.JwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("could not generate token"))
		return
	}
	var user response.UserResponse
	user.FillInTheModel(u)
	c.JSON(http.StatusOK, response.LoginResponse{
		Token:   tokenString,
		Expires: time.Now().Add(time.Hour * 24),
		User:    user,
	})
}

func (a *AuthHandler) Logout(c *gin.Context) {
	// TODO: Реализовать логику logout (добавление токена в blacklist и т.д.)
	c.JSON(http.StatusOK, response.SuccessResponse("Logged out successfully"))
}
