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

	err := c.ShouldBindJSON(&input)
	if err != nil{
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors" : errors}

		c.JSON(http.StatusBadRequest, response("Bad Request!",http.StatusUnprocessableEntity,"error", errorMessage))
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil{
		c.JSON(http.StatusBadRequest, response("Register account failed!",http.StatusBadRequest,"error", err.Error()))
		return
	}

	// token, err := h.JSONToken(c)
	formatter :=  user.FormatUser(newUser, "tokenya")
	c.JSON(http.StatusOK, response("Register Success", http.StatusCreated, "success", formatter))
}