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
	FindCountFormSubmissionByCampaign(campaignID string) (int64, error)
	FindSummaryEntriesByDate(workspaceID string, campaignID string) ([]masterschema.CampaignFormEntryChart, error)
	FindFormEntries(workspaceID string, campaignID string, params *commonschema.QueryParams) ([]masterschema.FormEntryList, error)
	FindCountFormEntries(workspaceID string, campaignID string, params *commonschema.QueryParams) (int64, error)
	CheckAllowedUserForCampaign(workspaceID string, campaignID string, userID string) (*masterschema.CampaignSchema, error)
	FindFormEntry(ID string) (*masterschema.FormEntrySchema, error)
	FindDetailFormEntry(formEntryID string) ([]masterschema.FormDetailEntrySchema, error)
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
	st := q.DB.Model(&models.Campaigns{}).Where("deleted = ? AND workspace_id::TEXT = ?", false, workspaceID).Select("id", "workspace_id", "title", "key", "slug", "description", "is_publish", "created_at")

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

func (q *CampaignQuery) FindCountFormSubmissionByCampaign(campaignID string) (int64, error) {
	var count int64
	if err := q.DB.Model(&models.FormEntries{}).Where("deleted = ? AND campaign_id = ?", false, campaignID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (q *CampaignQuery) FindSummaryEntriesByDate(workspaceID string, campaignID string) ([]masterschema.CampaignFormEntryChart, error) {
	var data []masterschema.CampaignFormEntryChart
	if err := q.DB.Raw("SELECT COUNT(1) AS total, TO_CHAR(form_entries.created_at, 'YYYY-MM-DD') AS date FROM form_entries JOIN campaigns ON campaigns.id = form_entries.campaign_id WHERE campaigns.deleted = ? AND form_entries.deleted = ? AND campaigns.workspace_id = ? AND campaigns.id = ? AND form_entries.created_at >= CURRENT_DATE - INTERVAL '60 days' GROUP BY TO_CHAR(form_entries.created_at, 'YYYY-MM-DD') ORDER BY TO_CHAR(form_entries.created_at, 'YYYY-MM-DD') ASC", false, false, workspaceID, campaignID).Scan(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (q *CampaignQuery) FindFormEntries(workspaceID string, campaignID string, params *commonschema.QueryParams) ([]masterschema.FormEntryList, error) {
	var formEntries []masterschema.FormEntryList

	// define offset
	offset := 0
	if params.Limit > 0 && params.Page > 0 {
		offset = params.Limit * (params.Page - 1)
	}

	// define query
	query := `
		SELECT 
			form_entries.id, 
			form_entries.user_id, 
			form_entries.campaign_id, 
			form_entries.status, 
			form_entries.remark, 
			form_entries.created_at::TEXT AS created_at, 
			campaigns.title AS campaign_title, 
			campaigns.key AS campaign_key, 
			campaigns.slug AS campaign_slug, 
			campaigns.workspace_id, 
			users.fullname AS user_name, 
			users.email as user_email 
		FROM form_entries
		JOIN campaigns ON campaigns.id = form_entries.campaign_id 
		JOIN workspaces ON workspaces.id = campaigns.workspace_id
		LEFT JOIN users ON users.id = form_entries.user_id 
		WHERE 
			form_entries.deleted = ? 
			AND campaigns.deleted = ? 
			AND workspaces.deleted = ? 
			AND campaigns.workspace_id = ? 
			AND campaigns.id = ?
	`
	args := []any{false, false, false, workspaceID, campaignID}

	// add search condition
	if params.Search != "" {
		query += " AND LOWER(users.fullname) LIKE ? "
		args = append(args, "%"+strings.ToLower(params.Search)+"%")
	}

	// add limit:offset
	query += " ORDER BY form_entries.created_at ASC LIMIT ? OFFSET ? "
	args = append(args, params.Limit, offset)

	// perform to get the data
	if err := q.DB.Raw(query, args...).Scan(&formEntries).Error; err != nil {
		return nil, err
	}
	return formEntries, nil
}

func (q *CampaignQuery) FindCountFormEntries(workspaceID string, campaignID string, params *commonschema.QueryParams) (int64, error) {
	var count int64

	// define query
	query := `
		SELECT 
			COUNT(1)
		FROM form_entries
		JOIN campaigns ON campaigns.id = form_entries.campaign_id 
		JOIN workspaces ON workspaces.id = campaigns.workspace_id
		LEFT JOIN users ON users.id = form_entries.user_id 
		WHERE 
			form_entries.deleted = ? 
			AND campaigns.deleted = ? 
			AND workspaces.deleted = ? 
			AND campaigns.workspace_id = ? 
			AND campaigns.id = ?
	`
	args := []any{false, false, false, workspaceID, campaignID}

	// add search condition
	if params.Search != "" {
		query += " AND LOWER(users.fullname) LIKE ? "
		args = append(args, "%"+strings.ToLower(params.Search)+"%")
	}

	// perform to get data
	if err := q.DB.Raw(query, args...).Scan(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (q *CampaignQuery) CheckAllowedUserForCampaign(workspaceID string, campaignID string, userID string) (*masterschema.CampaignSchema, error) {
	var campaign *masterschema.CampaignSchema

	query := `
		SELECT
			campaigns.id
		FROM campaigns
		JOIN workspaces ON workspaces.id = campaigns.workspace_id
		WHERE
			campaigns.deleted = ?
			AND workspaces.deleted = ?
			AND workspaces.id = ?
			AND campaigns.id = ?
			AND workspaces.id IN (
				SELECT 
					workspace_users.workspace_id
				FROM workspace_users
				WHERE
					workspace_users.workspace_id = workspaces.id
					AND workspace_users.user_id = ?
					AND workspace_users.deleted = ?
					AND workspace_users.status IN ('S3','S5')
			)	
	`
	args := []any{false, false, workspaceID, campaignID, userID, false}

	// perform to query
	if err := q.DB.Raw(query, args...).Scan(&campaign).Error; err != nil {
		return nil, err
	}
	return campaign, nil
}

func (q *CampaignQuery) FindFormEntry(ID string) (*masterschema.FormEntrySchema, error) {
	var formEntry masterschema.FormEntrySchema

	st := q.DB.Model(&models.FormEntries{}).
		Where("form_entries.deleted = ? AND form_entries.id = ?", false, ID).
		Select("form_entries.*", "campaigns.title AS campaign_title", "campaigns.description AS campaign_description", "users.fullname AS user_name", "users.email AS user_email").
		Joins("LEFT JOIN users ON users.id = form_entries.user_id").
		Joins("JOIN campaigns ON campaigns.id = form_entries.campaign_id")

	if err := st.First(&formEntry).Error; err != nil {
		return nil, err
	}

	return &formEntry, nil
}

func (q *CampaignQuery) FindDetailFormEntry(formEntryID string) ([]masterschema.FormDetailEntrySchema, error) {
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
