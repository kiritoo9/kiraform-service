package masterdi

import (
	masterrepo "kiraform/src/applications/repos/masters"
	storerepo "kiraform/src/applications/repos/stores"
	masterusecase "kiraform/src/applications/usecases/masters"

	"gorm.io/gorm"
)

type CampaignDependencies struct {
	DB *gorm.DB
	UC masterusecase.CampaignUsecase
}

func NewCampaignDependencies(DB *gorm.DB) *CampaignDependencies {
	campaignRepo := masterrepo.NewCampaignRepository(DB)
	workspaceRepo := masterrepo.NewWorkspaceRepository(DB)
	storeRepo := storerepo.NewStoreRepository(DB)

	UC := masterusecase.NewCampaignUsecase(campaignRepo, workspaceRepo, storeRepo)
	return &CampaignDependencies{
		DB: DB,
		UC: UC,
	}
}
