package masterrepo

import (
	"kiraform/src/applications/models"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	masterschema "kiraform/src/interfaces/rest/schemas/masters"
	"strings"

	"gorm.io/gorm"
)

type FormRepository interface {
	FindForms(params *commonschema.QueryParams) ([]masterschema.FormSchema, error)
	FindCountForm(params *commonschema.QueryParams) (int64, error)
	FindFormByID(ID string) (*masterschema.FormSchema, error)
}

type FormQuery struct {
	DB *gorm.DB
}

func NewFormRepository(DB *gorm.DB) *FormQuery {
	return &FormQuery{DB: DB}
}

func (q *FormQuery) FindForms(params *commonschema.QueryParams) ([]masterschema.FormSchema, error) {
	var forms []masterschema.FormSchema

	// calculating offset
	offset := 0
	if params.Limit > 0 && params.Page > 0 {
		offset = params.Limit * (params.Page - 1)
	}

	// define statments
	st := q.DB.Model(&models.Forms{}).Where("deleted = ?", false)

	// add search condition
	if params.Search != "" {
		st = st.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(params.Search)+"%")
	}

	// add orderby
	if params.OrderBy != "" {
		st = st.Order(params.OrderBy)
	}

	// add limit:offset
	st = st.Limit(params.Limit).Offset(offset)

	// perform to get the data
	if err := st.Find(&forms).Error; err != nil {
		return nil, err
	}

	// send back
	return forms, nil
}

func (q *FormQuery) FindCountForm(params *commonschema.QueryParams) (int64, error) {
	var count int64

	// preparing conditions
	st := q.DB.Model(&models.Forms{}).Where("deleted = ?", false)
	if params.Search != "" {
		st = st.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(params.Search)+"%")
	}

	// perform to get data
	if err := st.Count(&count).Error; err != nil {
		return 0, err
	}

	// send back
	return count, nil
}

func (q *FormQuery) FindFormByID(ID string) (*masterschema.FormSchema, error) {
	var form masterschema.FormSchema

	// preparing query
	err := q.DB.Model(&models.Forms{}).Where("deleted = ? AND id = ?", false, ID).First(&form).Error
	if err != nil {
		return nil, err
	}
	return &form, nil
}
