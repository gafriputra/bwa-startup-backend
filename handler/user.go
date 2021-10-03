package handler

import (
	"bwa-startup/helper"
	"bwa-startup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context){
	var input user.RegisterUserInput
	response := helper.APIResponse
	responseBadRequest := response("Bad Request!",http.StatusBadRequest,"error",nil)

	err := c.ShouldBindJSON(&input)
	if err != nil{
		c.JSON(http.StatusBadRequest, responseBadRequest)
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil{
		c.JSON(http.StatusBadRequest, responseBadRequest)
	}

	// token, err := h.JSONToken(c)
	formatter :=  user.FormatUser(newUser, "tokenya")
	c.JSON(http.StatusOK, response("Register Success", http.StatusCreated, "success", formatter))
}