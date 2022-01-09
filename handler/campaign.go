package handler

import (
	"bwa-startup/campaign"
	"bwa-startup/helper"
	"bwa-startup/user"
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
