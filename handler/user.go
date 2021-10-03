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

func (h *userHandler) Login(c *gin.Context){
	var input user.LoginInput
	response := helper.APIResponse

	err := c.ShouldBindJSON(&input)
	if err != nil{
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors" : errors}

		c.JSON(http.StatusBadRequest, response("Bad Request!",http.StatusUnprocessableEntity,"error", errorMessage))
		return
	}

	loggedUser, err := h.userService.Login(input)
	if err != nil{
		errorMessage := gin.H{"errors" : err.Error()}
		c.JSON(http.StatusBadRequest, response("Login failed!",http.StatusBadRequest,"error", errorMessage))
		return
	}

	formatter :=  user.FormatUser(loggedUser, "tokenya")
	c.JSON(http.StatusOK, response("Login Success", http.StatusCreated, "success", formatter))
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context){
	var input user.CheckEmailInput
	response := helper.APIResponse

	err := c.ShouldBindJSON(&input)
	if err != nil{
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors" : errors}

		c.JSON(http.StatusBadRequest, response("Bad Request!",http.StatusUnprocessableEntity,"error", errorMessage))
		return
	}

	IsEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil{
		errorMessage := gin.H{"errors" : "Server Error"}
		c.JSON(http.StatusBadRequest, response("Server Error",http.StatusBadRequest,"error", errorMessage))
		return
	}

	data := gin.H{
		"is_available" : IsEmailAvailable,
	}

	message := "Email has been registered"
	if IsEmailAvailable {
		message = "Email is available"
	}
	c.JSON(http.StatusOK, response(message, http.StatusCreated, "success", data))
}