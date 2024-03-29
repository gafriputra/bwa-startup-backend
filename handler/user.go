package handler

import (
	"bwa-startup/auth"
	"bwa-startup/helper"
	"bwa-startup/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput
	response := helper.APIResponse

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		c.JSON(http.StatusBadRequest, response("Bad Request!", http.StatusUnprocessableEntity, "error", errorMessage))
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, response("Register account failed!", http.StatusBadRequest, "error", err.Error()))
		return
	}

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response("Register account failed!", http.StatusBadRequest, "error", err.Error()))
		return
	}

	formatter := user.FormatUser(newUser, token)
	c.JSON(http.StatusOK, response("Register Success", http.StatusCreated, "success", formatter))
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginInput
	response := helper.APIResponse

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		c.JSON(http.StatusBadRequest, response("Bad Request!", http.StatusUnprocessableEntity, "error", errorMessage))
		return
	}

	loggedUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		c.JSON(http.StatusBadRequest, response("Login failed!", http.StatusBadRequest, "error", errorMessage))
		return
	}

	token, err := h.authService.GenerateToken(loggedUser.ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response("Login failed!", http.StatusBadRequest, "error", err.Error()))
		return
	}
	formatter := user.FormatUser(loggedUser, token)
	c.JSON(http.StatusOK, response("Login Success", http.StatusCreated, "success", formatter))
}

func (h *userHandler) CheckEmailAvailability(c *gin.Context) {
	var input user.CheckEmailInput
	response := helper.APIResponse

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		c.JSON(http.StatusBadRequest, response("Bad Request!", http.StatusUnprocessableEntity, "error", errorMessage))
		return
	}

	IsEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"errors": "Server Error"}
		c.JSON(http.StatusBadRequest, response("Server Error", http.StatusBadRequest, "error", errorMessage))
		return
	}

	data := gin.H{
		"is_available": IsEmailAvailable,
	}

	message := "Email has been registered"
	if IsEmailAvailable {
		message = "Email is available"
	}
	c.JSON(http.StatusOK, response(message, http.StatusCreated, "success", data))
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false, "errors": err.Error()}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	UserID := currentUser.ID
	path := fmt.Sprintf("images/%d-%s", UserID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false, "errors": err.Error()}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(UserID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false, "errors": err.Error()}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar successfuly uploaded", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}

func (h *userHandler) FetchUser(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(user.User)
	response := helper.APIResponse("Avatar successfuly uploaded", http.StatusOK, "success", user.FormatUser(currentUser, ""))
	c.JSON(http.StatusOK, response)
}
