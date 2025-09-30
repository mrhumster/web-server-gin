package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/mrhumster/web-server-gin/models"
	"github.com/mrhumster/web-server-gin/service"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password can't be empty"})
		return
	}

	id, err := h.service.CreateUser(c, user)
	if err != nil {
		log.Printf("CreateUser error: %v, type: %T", err, err)
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
	c.JSON(http.StatusOK, &user)
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
	users, err := h.service.ReadUserList(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.JSON(http.StatusOK, users)
}
