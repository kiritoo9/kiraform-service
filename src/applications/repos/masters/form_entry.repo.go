package masterrepo

import (
	"kiraform/src/applications/models"
	masterschema "kiraform/src/interfaces/rest/schemas/masters"

	"gorm.io/gorm"
)

type FormEntryRepository interface {
	EntryForm(formEntry models.FormEntries, formDetailEntries []models.FormDetailEntries) error
	GetHistory(userID string) ([]masterschema.FormEntrySchema, error)
	GetDetailHistory(userID string, ID string) (*masterschema.FormEntrySchema, error)
}

type FormEntryQuery struct {
	DB *gorm.DB
}

func NewFormEntryRepository(DB *gorm.DB) *FormEntryQuery {
	return &FormEntryQuery{DB: DB}
}

func (q *FormEntryQuery) EntryForm(formEntry models.FormEntries, formDetailEntries []models.FormDetailEntries) error {
	err := q.DB.Transaction(func(tx *gorm.DB) error {
		// insert form entry header
		if err := tx.Create(&formEntry).Error; err != nil {
			return err
		}

		// insert form detail entries
		if err := tx.Create(&formDetailEntries).Error; err != nil {
			return err
		}
		return nil // commit query
	})
	if err != nil {
		return err
	}
	return nil
}

func (q *FormEntryQuery) GetHistory(userID string) ([]masterschema.FormEntrySchema, error) {
	return nil, nil
}

func (q *FormEntryQuery) GetDetailHistory(userID string, ID string) (*masterschema.FormEntrySchema, error) {
	return nil, nil
}
