package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mrhumster/web-server-gin/dto/request"
	"github.com/mrhumster/web-server-gin/dto/response"
	"github.com/mrhumster/web-server-gin/models"
	"github.com/mrhumster/web-server-gin/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user request.UserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errors := make(map[string]string)
			for _, fieldError := range validationErrors {
				errors[fieldError.Field()] = getErrorMessage(fieldError)
			}
			c.JSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password can't be empty"})
		return
	}

	id, err := h.service.CreateUser(c, models.User{
		Login:    user.Login,
		Password: user.Password,
		Email:    user.Email,
	})
	if err != nil {
		log.Printf("⚠️ CreateUser error: %v, type: %T", err, err)
		switch err {
		case service.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, gin.H{
				"error": "user already exists",
			})

		default:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *UserHandler) ReadUser(c *gin.Context) {
	strId := c.Param("id")
	id, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.ReadUser(c, uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response.UserResponse{
		ID:        user.ID,
		Login:     user.Login,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
}

func (h *UserHandler) Update(c *gin.Context) {
	strId := c.Param("id")
	id, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = h.service.UpdateUser(c, uint(id), user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *UserHandler) Delete(c *gin.Context) {
	strId := c.Param("id")
	id, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.service.DeleteUser(c, uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, id)
}

func (h *UserHandler) ReadUsers(c *gin.Context) {
	page := int64(1)
	limit := int64(10)
	pageStr := c.Query("page")

	if pageStr != "" {
		page, _ = strconv.ParseInt(pageStr, 10, 64)
	}

	if page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "page is incorrect"})
		return
	}

	limitStr := c.Query("limit")
	if limitStr != "" {
		limit, _ = strconv.ParseInt(limitStr, 10, 64)
	}
	if limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit is incorrect"})
		return
	}

	users, total, err := h.service.ReadUserList(c, limit, page)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	var usersReponse []response.UserResponse
	for _, user := range users {
		usersReponse = append(usersReponse, response.UserResponse{
			ID:        user.ID,
			Login:     user.Login,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		})
	}
	c.JSON(http.StatusOK, response.UsersListReponse{
		Users: usersReponse,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}
