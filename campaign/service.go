package campaign

import (
	"errors"
	"fmt"

	"github.com/gosimple/slug"
)

type Service interface {
	GetCampaigns(UserID int) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
	CreateCampaign(input CreateCampaignInput) (Campaign, error)
	UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error)
	SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(UserID int) ([]Campaign, error) {
	if UserID != 0 {
		return s.repository.FindByUserID(UserID)
	}
	return s.repository.FindAll()
}

func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {
	return s.repository.FindByID(input.ID)
}

func (s *service) CreateCampaign(input CreateCampaignInput) (Campaign, error) {
	campaign := Campaign{
		Name:             input.Name,
		ShortDescription: input.ShortDescription,
		Description:      input.Description,
		GoalAmount:       input.GoalAmount,
		Perks:            input.Perks,
		UserID:           input.User.ID,
		Slug:             slug.Make(fmt.Sprintf("%s %d", input.Name, input.User.ID)),
	}

	return s.repository.Save(campaign)
}

func (s *service) UpdateCampaign(inputID GetCampaignDetailInput, inputData CreateCampaignInput) (Campaign, error) {
	campaign, err := s.repository.FindByID(inputID.ID)
	if err != nil {
		return campaign, err
	} else if campaign.ID == 0 {
		return campaign, errors.New("campaign not found")
	} else if campaign.UserID != inputData.User.ID {
		return campaign, errors.New("not an owner of the campaign")
	}

	campaign.Name = inputData.Name
	campaign.ShortDescription = inputData.ShortDescription
	campaign.Description = inputData.Description
	campaign.Perks = inputData.Perks
	campaign.GoalAmount = inputData.GoalAmount

	updatedCampaign, err := s.repository.Update(campaign)
	return updatedCampaign, err

}

func (s *service) SaveCampaignImage(input CreateCampaignImageInput, fileLocation string) (CampaignImage, error) {
	campaignImage := CampaignImage{}
	campaign, err := s.repository.FindByID(input.CampaignID)
	if err != nil {
		return campaignImage, err
	} else if campaign.ID == 0 {
		return campaignImage, errors.New("campaign not found")
	} else if campaign.UserID != input.User.ID {
		return campaignImage, errors.New("not an owner of the campaign")
	}

	campaignImage.CampaignID = input.CampaignID
	campaignImage.IsPrimary = 0
	campaignImage.FileName = fileLocation

	if input.IsPrimary {
		campaignImage.IsPrimary = 1
		_, err := s.repository.MarkAllImagesAsNonPrimary(input.CampaignID)
		if err != nil {
			return CampaignImage{}, err
		}
	}

	newCampaignImage, err := s.repository.CreateImage(campaignImage)
	return newCampaignImage, err

}
