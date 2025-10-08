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
	var u models.User
	u.FillInTheRequest(user)
	log.Printf("👷🏻‍♂️ SERVICE. USER CREATE")
	u.Debug()
	id, err := h.service.CreateUser(c, u)
	if err != nil {
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
	requestUserId, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse(err.Error()))
		return
	}

	strId := c.Param("id")
	id, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if requestUserId != id {
		c.JSON(http.StatusForbidden, response.ErrorResponse("you can update only that records, when you owner"))
		c.Abort()
		return
	}

	user, err := h.service.ReadUser(c, uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	var u response.UserResponse
	u.FillInTheModel(user)
	c.JSON(http.StatusOK, u)
}

func (h *UserHandler) Update(c *gin.Context) {
	requestUserId, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Authentication required"))
		return
	}

	strId := c.Param("id")
	id, err := strconv.ParseUint(strId, 10, 64)
	if requestUserId != id {
		c.JSON(http.StatusForbidden, response.ErrorResponse("you can update only that records, when you owner"))
		c.Abort()
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var user request.UpdateUserRequest
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err = h.service.UpdateUser(c, uint(id), user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedUser, err := h.service.ReadUser(c, uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var response response.UserResponse
	response.FillInTheModel(updatedUser)
	c.JSON(http.StatusOK, response)
}

func (h *UserHandler) Delete(c *gin.Context) {
	requestUserId, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Authentication required"))
		return
	}

	strId := c.Param("id")
	id, err := strconv.ParseUint(strId, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if requestUserId != id {
		c.JSON(http.StatusForbidden, response.ErrorResponse("you can update only that records, when you owner"))
		c.Abort()
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
	_, err := GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse("Authentication required"))
		return
	}
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
		var u response.UserResponse
		u.FillInTheModel(&user)
		usersReponse = append(usersReponse, u)
	}
	c.JSON(http.StatusOK, response.UsersListReponse{
		Users: usersReponse,
		Total: total,
		Page:  page,
		Limit: limit,
	})
}
