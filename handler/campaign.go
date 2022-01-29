package handler

import (
	"bwa-startup/campaign"
	"bwa-startup/helper"
	"bwa-startup/user"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campaignHandler struct {
	service campaign.Service
}

func NewCampaignHandler(service campaign.Service) *campaignHandler {
	return &campaignHandler{service}
}

func (h *campaignHandler) GetCampaigns(c *gin.Context) {
	UserID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(UserID)
	if err != nil {
		response := helper.APIResponse("Error to get campaigns", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", campaign.FormatCampaigns(campaigns))
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) GetCampaign(c *gin.Context) {
	var input campaign.GetCampaignDetailInput
	err := c.ShouldBindUri(&input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	campaignDetail, err := h.service.GetCampaignByID(input)
	if err != nil {
		response := helper.APIResponse("Failed to get detail of campaign", http.StatusInternalServerError, "error", err.Error())
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if campaignDetail.ID == 0 {
		response := helper.APIResponse("Campaign Not Found", http.StatusNotFound, "success", nil)
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := helper.APIResponse("Success get campaign detail", http.StatusOK, "success", campaign.FormatCampaignDetail(campaignDetail))
	c.JSON(http.StatusOK, response)

}

func (h *campaignHandler) CreateCampaign(c *gin.Context) {
	var input campaign.CreateCampaignInput
	response := helper.APIResponse

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		c.JSON(http.StatusBadRequest, response("Bad Request!", http.StatusUnprocessableEntity, "error", errorMessage))
		return
	}

	input.User = c.MustGet("currentUser").(user.User)

	newCampaign, err := h.service.CreateCampaign(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, response("Failed Create Campaign!", http.StatusBadRequest, "error", err.Error()))
		return
	}

	c.JSON(http.StatusOK, response("Success Create Campaign!", http.StatusOK, "success", campaign.FormatCampaign(newCampaign)))
}

func (h *campaignHandler) UpdatedCampaign(c *gin.Context) {
	var inputID campaign.GetCampaignDetailInput
	response := helper.APIResponse

	err := c.ShouldBindUri(&inputID)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	var inputData campaign.CreateCampaignInput
	err = c.ShouldBindJSON(&inputData)
	if err != nil {
		response := helper.APIResponse("Failed to update campaign", http.StatusBadRequest, "error", err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}

	inputData.User = c.MustGet("currentUser").(user.User)
	updatedCampaign, err := h.service.UpdateCampaign(inputID, inputData)
	if err != nil {
		c.JSON(http.StatusBadRequest, response("Failed Update Campaign!", http.StatusBadRequest, "error", err.Error()))
		return
	}
	c.JSON(http.StatusOK, response("Success Update Campaign!", http.StatusOK, "success", campaign.FormatCampaign(updatedCampaign)))
}

func (h *campaignHandler) UploadImage(c *gin.Context) {
	var input campaign.CreateCampaignImageInput
	response := helper.APIResponse

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		c.JSON(http.StatusBadRequest, response("Bad Request!", http.StatusUnprocessableEntity, "error", errorMessage))
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		data := gin.H{"is_uploaded": false, "errors": err.Error()}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	UserID := currentUser.ID
	input.User = currentUser
	path := fmt.Sprintf("images/campaign-user-%d-%s", UserID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false, "errors": err.Error()}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.service.SaveCampaignImage(input, path)
	if err != nil {
		data := gin.H{"is_uploaded": false, "errors": err.Error()}
		response := helper.APIResponse("Failed to upload campaign image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	c.JSON(http.StatusOK, response("Campaign image successfuly uploaded", http.StatusOK, "success", data))

}
