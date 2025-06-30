package masterdi

import (
	masterrepo "kiraform/src/applications/repos/masters"
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

	// load usecase
	UC := masterusecase.NewFormEntryUsecase(formEntryRepo)
	UCcampaign := masterusecase.NewCampaignUsecase(campaignRepo, workspaceRepo)
	return &FormEntryDependencies{
		DB:         DB,
		UC:         UC,
		UCcampaign: UCcampaign,
	}
}
