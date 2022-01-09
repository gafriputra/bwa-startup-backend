package campaign

type Service interface {
	GetCampaigns(userID int) ([]Campaign, error)
	GetCampaignByID(input GetCampaignDetailInput) (Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetCampaigns(userID int) ([]Campaign, error) {
	if userID != 0 {
		return s.repository.FindByUserID(userID)
	}
	return s.repository.FindAll()
}

func (s *service) GetCampaignByID(input GetCampaignDetailInput) (Campaign, error) {
	return s.repository.FindByID(input.ID)
}
