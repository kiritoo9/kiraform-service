package masterdi

import (
	masterrepo "kiraform/src/applications/repos/masters"
	storerepo "kiraform/src/applications/repos/stores"
	masterusecase "kiraform/src/applications/usecases/masters"

	"gorm.io/gorm"
)

type FormEntryDependencies struct {
	DB         *gorm.DB
	UC         masterusecase.FormEntryUsecase
	UCcampaign masterusecase.CampaignUsecase
}

func NewFormEntryDependencies(DB *gorm.DB) *FormEntryDependencies {
	// load repositories
	formEntryRepo := masterrepo.NewFormEntryRepository(DB)
	campaignRepo := masterrepo.NewCampaignRepository(DB)
	workspaceRepo := masterrepo.NewWorkspaceRepository(DB)
	storeRepo := storerepo.NewStoreRepository(DB)

	// load usecase
	UC := masterusecase.NewFormEntryUsecase(formEntryRepo)
	UCcampaign := masterusecase.NewCampaignUsecase(campaignRepo, workspaceRepo, storeRepo)
	return &FormEntryDependencies{
		DB:         DB,
		UC:         UC,
		UCcampaign: UCcampaign,
	}
}
