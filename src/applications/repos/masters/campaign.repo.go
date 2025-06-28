package masterrepo

import (
	"fmt"
	"kiraform/src/applications/models"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	masterschema "kiraform/src/interfaces/rest/schemas/masters"
	"strings"

	"gorm.io/gorm"
)

type CampaignRepository interface {
	FindCampaigns(workspaceID string, params *commonschema.QueryParams) ([]masterschema.CampaignSchema, error)
	FindCountCampaign(workspaceID string, params *commonschema.QueryParams) (int64, error)
	FindCampaignByID(workspaceID string, ID string) (*masterschema.CampaignSchema, error)
	FindCampaignByKey(key string, isPublish *bool) (*masterschema.CampaignSchema, error)
	FindFormsByCampaign(campaignID string) ([]masterschema.CampaignFormSchema, error)
	FindFormAttributes(campaignFormID string) ([]masterschema.CampaignFormAttributeSchemas, error)
	CreateCampaign(campaign models.Campaigns, campaignForms []models.CampaignForms, campaignFormAttributes []models.CampaignFormAttributes) error
	UpdateCampaign(ID string, campaign models.Campaigns) error
	UpdateEntireCampaign(ID string, campaign models.Campaigns, campaignFormActions map[string][]models.CampaignForms, campaignFormAttributesCreate []models.CampaignFormAttributes) error
	FindCampaignSeos(campaignID string, params *commonschema.QueryParams) ([]masterschema.CampaignSeoSchema, error)
	FindCountCampaignSeo(campaignID string, params *commonschema.QueryParams) (int64, error)
	FindCampaignSeoByID(campaignID string, ID string) (*masterschema.CampaignSeoSchema, error)
	CreateCampaignSeo(body models.CampaignSeos) error
	UpdateCampaignSeo(campaignID string, ID string, body models.CampaignSeos) error
	CreateFormAttribute(formAttribute models.CampaignFormAttributes) error
	UpdateFormAttribute(formAttribute models.CampaignFormAttributes, ID string) error
}

type CampaignQuery struct {
	DB *gorm.DB
}

func NewCampaignRepository(DB *gorm.DB) *CampaignQuery {
	return &CampaignQuery{
		DB: DB,
	}
}

func (q *CampaignQuery) FindCampaigns(workspaceID string, params *commonschema.QueryParams) ([]masterschema.CampaignSchema, error) {
	var campaigns []masterschema.CampaignSchema

	// define offset
	offset := 0
	if params.Limit > 0 && params.Page > 0 {
		offset = params.Limit * (params.Page - 1)
	}

	// define statemetns
	st := q.DB.Model(&models.Campaigns{}).Where("deleted = ? AND workspace_id::TEXT = ?", false, workspaceID).Select("id", "workspace_id", "title", "description", "is_publish", "created_at")

	// add search condition
	if params.Search != "" {
		st = st.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(params.Search)+"%")
	}

	// add orderby
	if params.OrderBy != "" {
		st = st.Order(params.OrderBy)
	}

	// add limit:offset
	st = st.Limit(params.Limit).Offset(offset)

	// perform to get the data
	if err := st.Find(&campaigns).Error; err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (q *CampaignQuery) FindCountCampaign(workspaceID string, params *commonschema.QueryParams) (int64, error) {
	var count int64

	// prepare condition
	st := q.DB.Model(&models.Campaigns{}).Where("deleted = ? AND workspace_id::TEXT = ?", false, workspaceID)
	if params.Search != "" {
		st = st.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(params.Search)+"%")
	}

	// perform to get data
	if err := st.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (q *CampaignQuery) FindCampaignByID(workspaceID string, ID string) (*masterschema.CampaignSchema, error) {
	var campaign masterschema.CampaignSchema

	// perform to query
	st := q.DB.Model(&models.Campaigns{}).Where("deleted = ? AND workspace_id = ? and id = ?", false, workspaceID, ID)
	if err := st.First(&campaign).Error; err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (q *CampaignQuery) FindCampaignByKey(key string, isPublish *bool) (*masterschema.CampaignSchema, error) {
	var campaign masterschema.CampaignSchema

	// perform to query
	st := q.DB.Model(&models.Campaigns{}).Where("deleted = ? AND key = ?", false, key)
	if isPublish != nil {
		st = st.Where("is_publish = ?", isPublish)
	}

	if err := st.First(&campaign).Error; err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (q *CampaignQuery) FindFormsByCampaign(campaignID string) ([]masterschema.CampaignFormSchema, error) {
	var campaignForms []masterschema.CampaignFormSchema

	// perform to query
	st := q.DB.Model(&models.CampaignForms{}).
		Joins("JOIN forms ON forms.id = campaign_forms.form_id").
		Select("campaign_forms.*", "forms.name AS form_name", "forms.code AS form_code").
		Where("campaign_forms.deleted = ? AND campaign_forms.campaign_id = ?", false, campaignID)
	st = st.Order("campaign_forms.created_at ASC")
	if err := st.Find(&campaignForms).Error; err != nil {
		return nil, err
	}
	return campaignForms, nil
}

func (q *CampaignQuery) FindFormAttributes(campaignFormID string) ([]masterschema.CampaignFormAttributeSchemas, error) {
	var campaignFormAttributes []masterschema.CampaignFormAttributeSchemas

	// perform query
	st := q.DB.Model(&models.CampaignFormAttributes{}).Where("deleted = ? AND campaign_form_id = ?", false, campaignFormID)
	st = st.Order("created_at ASC")
	if err := st.Find(&campaignFormAttributes).Error; err != nil {
		return nil, err
	}
	return campaignFormAttributes, nil
}

func (q *CampaignQuery) CreateCampaign(campaign models.Campaigns, campaignForms []models.CampaignForms, campaignFormAttributes []models.CampaignFormAttributes) error {
	// insert all data using transaction [commit:rollback]
	// to prevent error coming
	err := q.DB.Transaction(func(tx *gorm.DB) error {
		// insert campaign header
		if err := tx.Create(&campaign).Error; err != nil {
			return err
		}

		// insert campaign forms
		if len(campaignForms) > 0 {
			if err := tx.Create(&campaignForms).Error; err != nil {
				return err
			}
		}

		//  insert campaign form attributes
		if len(campaignFormAttributes) > 0 {
			if err := tx.Create(&campaignFormAttributes).Error; err != nil {
				return err
			}
		}

		return nil // flag as commit
	})
	if err != nil {
		return err
	}

	// tell usecase if everything is OK
	return nil
}

func (q *CampaignQuery) UpdateCampaign(ID string, campaign models.Campaigns) error {
	if err := q.DB.Where("deleted = ? AND id = ?", false, ID).Updates(&campaign).Error; err != nil {
		return err
	}
	return nil
}

func (q *CampaignQuery) UpdateEntireCampaign(ID string, campaign models.Campaigns, campaignFormActions map[string][]models.CampaignForms, campaignFormAttributesCreate []models.CampaignFormAttributes) error {
	err := q.DB.Transaction(func(tx *gorm.DB) error {
		// update campaign
		if err := tx.Where("deleted = ? AND id = ?", false, ID).Updates(campaign).Error; err != nil {
			return err
		}

		// action create new campaign form
		if len(campaignFormActions["create"]) > 0 {
			if c, ok := campaignFormActions["create"]; ok {
				for _, cv := range c {
					if err := tx.Model(&models.CampaignForms{}).Create(&cv).Error; err != nil {
						return err
					}
				}
			}
		}

		// action update campaign form
		if len(campaignFormActions["update"]) > 0 {
			if u, ok := campaignFormActions["update"]; ok {
				for _, uv := range u {
					if err := tx.Model(&models.CampaignForms{}).Where("id = ?", uv.ID).Updates(&uv).Error; err != nil {
						return err
					}
				}
			}
		}

		// action delete campaign form
		if len(campaignFormActions["delete"]) > 0 {
			if d, ok := campaignFormActions["delete"]; ok {
				for _, dv := range d {
					if err := tx.Model(&models.CampaignForms{}).Where("id = ?", dv.ID).Updates(&dv).Error; err != nil {
						fmt.Println(err)
						return err
					}
				}
			}
		}

		// action create attributes for new form
		if len(campaignFormAttributesCreate) > 0 {
			if err := tx.Model(&models.CampaignFormAttributes{}).Create(&campaignFormAttributesCreate).Error; err != nil {
				return err
			}
		}

		// commit transaction
		return nil
	})
	if err != nil {
		return nil
	}

	// return success response
	return nil
}

func (q *CampaignQuery) FindCampaignSeos(campaignID string, params *commonschema.QueryParams) ([]masterschema.CampaignSeoSchema, error) {
	var campaigns []masterschema.CampaignSeoSchema

	// define offset
	offset := 0
	if params.Limit > 0 && params.Page > 0 {
		offset = params.Limit * (params.Page - 1)
	}

	// define statemetns
	st := q.DB.Model(&models.CampaignSeos{}).Where("deleted = ? AND campaign_id::TEXT = ?", false, campaignID)

	// add search condition
	if params.Search != "" {
		st = st.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(params.Search)+"%")
	}

	// add orderby
	if params.OrderBy != "" {
		st = st.Order(params.OrderBy)
	}

	// add limit:offset
	st = st.Limit(params.Limit).Offset(offset)

	// perform to get the data
	if err := st.Find(&campaigns).Error; err != nil {
		return nil, err
	}
	return campaigns, nil
}

func (q *CampaignQuery) FindCountCampaignSeo(campaignID string, params *commonschema.QueryParams) (int64, error) {
	var count int64

	// prepare condition
	st := q.DB.Model(&models.CampaignSeos{}).Where("deleted = ? AND campaign_id::TEXT = ?", false, campaignID)
	if params.Search != "" {
		st = st.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(params.Search)+"%")
	}

	// perform to get data
	if err := st.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (q *CampaignQuery) FindCampaignSeoByID(campaignID string, ID string) (*masterschema.CampaignSeoSchema, error) {
	var campaign masterschema.CampaignSeoSchema

	// perform to query
	st := q.DB.Model(&models.CampaignSeos{}).Where("deleted = ? AND campaign_id = ? and id = ?", false, campaignID, ID)
	if err := st.First(&campaign).Error; err != nil {
		return nil, err
	}
	return &campaign, nil
}

func (q *CampaignQuery) CreateCampaignSeo(campaignSeo models.CampaignSeos) error {
	if err := q.DB.Create(&campaignSeo).Error; err != nil {
		return err
	}
	return nil
}

func (q *CampaignQuery) UpdateCampaignSeo(campaignID string, ID string, campaignSeo models.CampaignSeos) error {
	if err := q.DB.Where("deleted = ? AND campaign_id = ? AND id = ?", false, campaignID, ID).Updates(&campaignSeo).Error; err != nil {
		return err
	}
	return nil
}

func (q *CampaignQuery) CreateFormAttribute(formAttribute models.CampaignFormAttributes) error {
	if err := q.DB.Model(&models.CampaignFormAttributes{}).Create(&formAttribute).Error; err != nil {
		return err
	}
	return nil
}

func (q *CampaignQuery) UpdateFormAttribute(formAttribute models.CampaignFormAttributes, ID string) error {
	if err := q.DB.Model(&models.CampaignFormAttributes{}).Where("id = ?", ID).Updates(&formAttribute).Error; err != nil {
		return err
	}
	return nil
}

func (q *CampaignQuery) FindCampaignFormAttributes(campaignFormID string) ([]models.CampaignFormAttributes, error) {
	var campaignFormAttributes []models.CampaignFormAttributes
	if err := q.DB.Model(&models.CampaignFormAttributes{}).Where("deleted = ? AND campaign_form_id = ?", false, campaignFormID).Find(&campaignFormAttributes).Error; err != nil {
		return nil, err
	}
	return campaignFormAttributes, nil
}
