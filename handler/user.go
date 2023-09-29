package handler

import (
	"net/http"

	"github.com/gin-api/service"
	"github.com/gin-gonic/gin"
)

type usersHandler struct {
	userHandler service.UserService
}

func NewUserHandler(userHandler service.UserService) usersHandler {
	return usersHandler{userHandler: userHandler}
}

func (h usersHandler) Login(c *gin.Context) {

	auth := &service.AuthInput{}

	if err := c.Bind(auth); err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	token, err := h.userHandler.Login(*auth)

	if err != nil {
		c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": *token,
	})
}

func (h usersHandler) RegisterUser(c *gin.Context) {

	register := &service.RegisterUser{}
	if err := c.Bind(register); err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	message, err := h.userHandler.RegisterUser(*register)
	if err != nil {
		c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, map[string]string{"message ": *message})
}
