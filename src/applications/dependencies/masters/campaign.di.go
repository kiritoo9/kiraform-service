package masterdi

import (
	masterrepo "kiraform/src/applications/repos/masters"
	masterusecase "kiraform/src/applications/usecases/masters"

	"gorm.io/gorm"
)

type CampaignDependencies struct {
	DB *gorm.DB
	UC masterusecase.CampaignUsecase
}

func NewCampaignDependencies(DB *gorm.DB) *CampaignDependencies {
	campaignRepo := masterrepo.NewCampaignRepository(DB)
	UC := masterusecase.NewCampaignUsecase(campaignRepo)
	return &CampaignDependencies{
		DB: DB,
		UC: UC,
	}
}
