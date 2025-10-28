package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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
	Domain       string
}

func NewAuthHandler(userService *service.UserService, tokenService *service.TokenService, jwtSecret, domain string) *AuthHandler {
	return &AuthHandler{
		UserService:  userService,
		TokenService: tokenService,
		JwtSecret:    jwtSecret,
		Domain:       domain,
	}
}

func (a *AuthHandler) Login(c *gin.Context) {
	var (
		req request.LoginRequest
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

	tokenPair, err := a.TokenService.GenerateToken(u)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			response.ErrorResponse(
				fmt.Sprintf("generate token: %v", err.Error()),
			),
		)
	}

	c.SetCookie(
		"refresh_token",
		tokenPair.RefreshToken,
		int(a.TokenService.GetRefreshExpiry().Seconds()),
		"/api/refresh",
		a.Domain,
		false,
		true,
	)

	var user response.UserResponse
	user.FillInTheModel(u)
	c.JSON(http.StatusOK, response.LoginResponse{
		AccessToken: tokenPair.AccessToken,
		ExpiresIn:   tokenPair.ExpiresIn,
		TokenType:   tokenPair.TokenType,
	})
}

func (a *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse("refresh token required"))
		return
	}

	claims, err := a.TokenService.ValidateRefreshToken(refreshToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse("invalid refresh token"))
		return
	}
	userID, err := strconv.ParseUint(claims.UserID, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse("invalid user id in claim"))
	}
	u, err := a.UserService.ReadUser(c, uint(userID))
	if err != nil || *u.TokenVersion != claims.TokenVersion {
		log.Printf("⚠️ AUTH HANDLER: User Token ver %v; claims token ver %v", *u.TokenVersion, claims.TokenVersion)
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse("token revoke"))
		return
	}
	tokenPair, err := a.TokenService.GenerateToken(u)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, response.ErrorResponse("failed to generate token"))
		return
	}

	c.SetCookie(
		"refresh_token",
		tokenPair.RefreshToken,
		int(a.TokenService.GetRefreshExpiry().Seconds()),
		"/api/refresh",
		a.Domain,
		false,
		true,
	)

	var user response.UserResponse
	user.FillInTheModel(u)
	c.JSON(http.StatusOK, response.LoginResponse{
		AccessToken: tokenPair.AccessToken,
		ExpiresIn:   tokenPair.ExpiresIn,
		TokenType:   tokenPair.TokenType,
	})

}

func (a *AuthHandler) Logout(c *gin.Context) {
	c.SetCookie("refresh_token", "", -1, "/api/refresh", "", false, true)
	c.JSON(http.StatusOK, response.SuccessResponse("Logged out successfully"))
}

func (a *AuthHandler) LogoutAll(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	id, err := strconv.ParseUint(userID.(string), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, response.ErrorResponse("imvalid user id in claim"))
	}
	err = a.UserService.UpdateTokenVersion(c, uint64(id), generateNewTokenVersion())
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("failed to logout"))
		return
	}

	c.SetCookie("refresh_token", "", -1, "/api/refresh", "", false, true)

	c.JSON(http.StatusOK, response.ErrorResponse("logged out from all devices"))
}

func generateNewTokenVersion() string {
	return "v" + time.Now().Format("20060102150405")
}
