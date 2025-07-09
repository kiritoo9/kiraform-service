package masterrepo

import (
	"kiraform/src/applications/models"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	masterschema "kiraform/src/interfaces/rest/schemas/masters"
	"strings"

	"gorm.io/gorm"
)

type FormEntryRepository interface {
	EntryForm(formEntry models.FormEntries, formDetailEntries []models.FormDetailEntries) error
	FindFormEntries(userID string, params *commonschema.QueryParams) ([]masterschema.FormEntrySchema, error)
	FindCountFormEntry(userID string, params *commonschema.QueryParams) (int64, error)
	FindFormEntry(userID string, ID string) (*masterschema.FormEntrySchema, error)
	FindDetailFormEntry(formEntryID string) ([]masterschema.FormDetailEntrySchema, error)
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

func (q *FormEntryQuery) FindFormEntries(userID string, params *commonschema.QueryParams) ([]masterschema.FormEntrySchema, error) {
	var formEntries []masterschema.FormEntrySchema

	// calculating offset
	offset := 0
	if params.Limit > 0 && params.Page > 0 {
		offset = params.Limit * (params.Page - 1)
	}

	// define statments
	st := q.DB.Model(&models.FormEntries{}).Where("form_entries.deleted = ? AND form_entries.user_id = ?", false, userID).
		Select("form_entries.*", "users.fullname AS user_name", "users.email AS user_email", "campaigns.title AS campaign_title", "campaigns.description AS campaign_description").
		Joins("JOIN users ON users.id = form_entries.user_id").
		Joins("JOIN campaigns ON campaigns.ID = form_entries.campaign_id")

	// add search condition
	if params.Search != "" {
		st = st.Where("LOWER(campaigns.title) LIKE ?", "%"+strings.ToLower(params.Search)+"%")
	}

	// add orderby
	if params.OrderBy != "" {
		st = st.Order(params.OrderBy)
	} else {
		st = st.Order("form_entries.created_at DESC")
	}

	// add limit:offset
	st = st.Limit(params.Limit).Offset(offset)

	// perform to get the data
	if err := st.Find(&formEntries).Error; err != nil {
		return nil, err
	}

	// send back
	return formEntries, nil
}

func (q *FormEntryQuery) FindCountFormEntry(userID string, params *commonschema.QueryParams) (int64, error) {
	var count int64

	// preparing conditions
	st := q.DB.Model(&models.FormEntries{}).Where("form_entries.deleted = ? AND form_entries.user_id = ?", false, userID).
		Joins("JOIN campaigns ON campaigns.id = form_entries.campaign_id")

	if params.Search != "" {
		st = st.Where("LOWER(campaigns.title) LIKE ?", "%"+strings.ToLower(params.Search)+"%")
	}

	// perform to get data
	if err := st.Count(&count).Error; err != nil {
		return 0, err
	}

	// send back
	return count, nil
}

func (q *FormEntryQuery) FindFormEntry(userID string, ID string) (*masterschema.FormEntrySchema, error) {
	var formEntry masterschema.FormEntrySchema

	st := q.DB.Model(&models.FormEntries{}).
		Where("form_entries.deleted = ? AND form_entries.user_id = ? AND form_entries.id = ?", false, userID, ID).
		Select("form_entries.*", "campaigns.title AS campaign_title", "campaigns.description AS campaign_description", "users.fullname AS user_name", "users.email AS user_email").
		Joins("JOIN users ON users.id = form_entries.user_id").
		Joins("JOIN campaigns ON campaigns.id = form_entries.campaign_id")

	if err := st.First(&formEntry).Error; err != nil {
		return nil, err
	}

	return &formEntry, nil
}

func (q *FormEntryQuery) FindDetailFormEntry(formEntryID string) ([]masterschema.FormDetailEntrySchema, error) {
	var formDetailEntries []masterschema.FormDetailEntrySchema

	st := q.DB.Model(&models.FormDetailEntries{}).
		Where("form_detail_entries.deleted = ? AND form_detail_entries.form_entry_id = ?", false, formEntryID).
		Select("form_detail_entries.*", "forms.name AS form_name", "forms.code AS form_code", "campaign_forms.title AS campaign_form_title", "campaign_forms.description AS campaign_form_description").
		Joins("JOIN campaign_forms ON campaign_forms.id = form_detail_entries.campaign_form_id").
		Joins("JOIN forms ON forms.id = campaign_forms.form_id")

	if err := st.Find(&formDetailEntries).Error; err != nil {
		return nil, err
	}

	return formDetailEntries, nil
}
