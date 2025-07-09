package masterroute

import (
	"errors"
	masterdi "kiraform/src/applications/dependencies/masters"
	"kiraform/src/applications/helpers"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	masterschema "kiraform/src/interfaces/rest/schemas/masters"
	"kiraform/src/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CampaignHandler struct {
	DB           *gorm.DB
	Validator    *validator.Validate
	Dependencies masterdi.CampaignDependencies
}

func NewCampaignHandler(DB *gorm.DB, validator *validator.Validate, dependencies masterdi.CampaignDependencies) *CampaignHandler {
	return &CampaignHandler{
		DB:           DB,
		Validator:    validator,
		Dependencies: dependencies,
	}
}

func NewCampaignHTTP(g *echo.Group, DB *gorm.DB) {
	validator := validator.New()
	h := NewCampaignHandler(DB, validator, *masterdi.NewCampaignDependencies(DB))

	// define endpoints
	c := g.Group("/campaigns")
	c.GET("/:workspace_id", h.FindCampaigns)
	c.GET("/dashboard/:workspace_id", h.CampaignDashboard)
	c.GET("/detail/:workspace_id/:id", h.FindCampaign)
	c.POST("/:workspace_id", h.CreateCampaign)
	c.PUT("/:workspace_id/:id", h.UpdateCampaign)
	c.DELETE("/:workspace_id/:id", h.DeleteCampaign)

	// for analytic pages
	a := c.Group("/analytics")
	a.GET("/dashboard/:workspace_id/:campaign_id", h.DashboardAnalytics)
	a.GET("/form_entries/:workspace_id/:campaign_id", h.FindFormEntries)
	a.GET("/form_entries/:workspace_id/:campaign_id/:id", h.FindDetailFormEntry)

	// for campaign seos
	s := c.Group("/seos")
	s.GET("/:campaign_id", h.FindCampaignSeos)
	s.GET("/:campaign_id/:id", h.FindCampaignSeo)
	s.POST("/:campaign_id", h.CreateCampaignSeo)
	s.PUT("/:campaign_id/:id", h.UpdateCampaignSeo)
	s.DELETE("/:campaign_id/:id", h.DeleteCampaignSeo)
}

// @Security BearerAuth
// @Summary      List Campaigns
// @Description  Get the list of campaigns you created
// @Tags         Master - Campaigns
// @Accept  	 json
// @Produce  	 json
// @Param 		 workspace_id path string true "Workspace ID"
// @Param 		 page query int true "Page of list data"
// @Param 		 limit query int true "Limitting data you want to get"
// @Param 		 search query string false "Find your data with keywords"
// @Param 		 orderBy query string false "Ordering data" example(created_at:desc)
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/campaigns/{workspace_id} [get]
func (h *CampaignHandler) FindCampaigns(c echo.Context) error {
	// because campaign is detail of workspace
	// so this endpoint require workspace_id to get data
	// do not get all data direclty
	workspaceID := c.Param("workspace_id")
	params := utils.QParams(c)
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// send to usecase to get data
	list, err := h.Dependencies.UC.FindCampaigns(workspaceID, params)
	if err != nil {
		response.Message = err.Error()
		return c.JSON(response.Code, response)
	}

	// send response
	response.Code = http.StatusOK
	response.Message = "Request success"
	response.Data = list
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      Detail Campaigns
// @Description  Get detail of campaigns you created
// @Tags         Master - Campaigns
// @Accept  	 json
// @Produce  	 json
// @Param 		 workspace_id path string true "Workspace ID"
// @Param 		 id path string true "ID of your data"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/campaigns/detail/{workspace_id}/{id} [get]
func (h *CampaignHandler) FindCampaign(c echo.Context) error {
	// get parameters
	workspaceID := c.Param("workspace_id")
	ID := c.Param("id")

	// check allowed user
	err := helpers.CheckAllowedCampaign(c, workspaceID, ID, h.Dependencies.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// get existing data
	campaign, err := h.Dependencies.UC.FindCampaign(workspaceID, ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Data is not found")
	}

	// get detail form by this campaign
	forms, err := h.Dependencies.UC.FindFormsByCampaign(campaign.ID.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// get attributes each form
	// handle if forms is null, because there is possibility for it
	for i, v := range forms {
		attr, err := h.Dependencies.UC.FindFormAttributes(v.ID.String())
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		forms[i].Attributes = attr
	}

	// prepare response
	campaign.Forms = forms
	response := commonschema.ResponseHTTP{
		Code:    http.StatusOK,
		Message: "Request success",
		Data:    campaign,
	}
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      Campaign Dashboard
// @Description  Get summary data campaign for dashboard
// @Tags         Master - Campaigns
// @Accept  	 json
// @Produce  	 json
// @Param 		 workspace_id path string true "Workspace ID"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/campaigns/dashboard/{workspace_id} [get]
func (h *CampaignHandler) CampaignDashboard(c echo.Context) error {
	// get parameters
	workspaceID := c.Param("workspace_id")

	// check allowed user
	err := helpers.CheckAllowedWorkspace(c, workspaceID, h.Dependencies.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// get existing data
	dashboard, err := h.Dependencies.UC.CampaignDashboard(workspaceID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Data is not found")
	}

	// send response
	response := commonschema.ResponseHTTP{
		Code:    http.StatusOK,
		Message: "Request success",
		Data:    dashboard,
	}
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      Create Campaign
// @Description  Create new campaign data
// @Tags         Master - Campaigns
// @Accept  	 json
// @Produce  	 json
// @Param 		 workspace_id path string true "Workspace ID"
// @Param        campaignPayload  body      masterschema.CampaignPayload   true  "campaign payload"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/campaigns/{workspace_id} [post]
func (h *CampaignHandler) CreateCampaign(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	var body masterschema.CampaignPayload

	// check allowed user
	err := helpers.CheckAllowedWorkspace(c, workspaceID, h.Dependencies.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body payload")
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send to usecase for insert logic
	err = h.Dependencies.UC.CreateCampaign(workspaceID, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send success response
	return c.JSON(http.StatusCreated, commonschema.ResponseHTTP{
		Code:    http.StatusCreated,
		Message: "Data is successfully created",
	})
}

// @Security BearerAuth
// @Summary      Update Campaign
// @Description  Update existing campaign data
// @Tags         Master - Campaigns
// @Accept  	 json
// @Produce  	 json
// @Param 		 workspace_id path string true "Workspace ID"
// @Param 		 id path string true "ID of your data"
// @Param        campaignPayload  body      masterschema.CampaignPayload   true  "Campaign payload"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/campaigns/{workspace_id}/{id} [put]
func (h *CampaignHandler) UpdateCampaign(c echo.Context) error {
	// get payload and parameters
	workspaceID := c.Param("workspace_id")
	ID := c.Param("id")
	var body masterschema.CampaignPayload

	// check allowed user
	err := helpers.CheckAllowedCampaign(c, workspaceID, ID, h.Dependencies.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// check for valid body
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body payload")
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send to usecase for update logic
	err = h.Dependencies.UC.UpdateCampaign(workspaceID, ID, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusNoContent, nil)
}

// @Security BearerAuth
// @Summary      Delete Campaign
// @Description  Delete existing campaign data
// @Tags         Master - Campaigns
// @Accept  	 json
// @Produce  	 json
// @Param 		 workspace_id path string true "Workspace ID"
// @Param 		 id path string true "ID of your data"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/campaigns/{workspace_id}/{id} [delete]
func (h *CampaignHandler) DeleteCampaign(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	ID := c.Param("id")

	// check allowed user
	err := helpers.CheckAllowedCampaign(c, workspaceID, ID, h.Dependencies.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send to usecase to do delete process
	if err := h.Dependencies.UC.DeleteCampaign(workspaceID, ID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusNoContent, nil)
}

// @Security BearerAuth
// @Summary      List SEO Campaign
// @Description  Get the list of seo campaign you created
// @Tags         Master - Campaign SEO
// @Accept  	 json
// @Produce  	 json
// @Param 		 campaign_id path string true "Campaign ID"
// @Param 		 page query int true "Page of list data"
// @Param 		 limit query int true "Limitting data you want to get"
// @Param 		 search query string false "Find your data with keywords"
// @Param 		 orderBy query string false "Ordering data" example(created_at:desc)
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/campaigns/seos/{campaign_id} [get]
func (h *CampaignHandler) FindCampaignSeos(c echo.Context) error {
	campaignID := c.Param("campaign_id")
	params := utils.QParams(c)
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// send to usecase to get data
	list, err := h.Dependencies.UC.FindCampaignSeos(campaignID, params)
	if err != nil {
		response.Message = err.Error()
		return c.JSON(response.Code, response)
	}

	// send response
	response.Code = http.StatusOK
	response.Message = "Request success"
	response.Data = list
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      Detail SEO Campaign
// @Description  Get detail of seo campaign you created
// @Tags         Master - Campaign SEO
// @Accept  	 json
// @Produce  	 json
// @Param 		 campaign_id path string true "Campaign ID"
// @Param 		 id path string true "ID of your data"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/campaigns/seos/{campaign_id}/{id} [get]
func (h *CampaignHandler) FindCampaignSeo(c echo.Context) error {
	// get parameters
	campaignID := c.Param("campaign_id")
	ID := c.Param("id")

	// get existing data
	campaignSeo, err := h.Dependencies.UC.FindCampaignSeoByID(campaignID, ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Data is not found")
	}

	// send response
	response := commonschema.ResponseHTTP{
		Code:    http.StatusOK,
		Message: "Request success",
		Data:    campaignSeo,
	}
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      Create SEO Campaign
// @Description  Create new seo campaign data
// @Tags         Master - Campaign SEO
// @Accept  	 json
// @Produce  	 json
// @Param 		 campaign_id path string true "Campaign ID"
// @Param        campaignSeoPayload  body      masterschema.CampaignSeoPayload   true  "SEO campaign payload"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/campaigns/seos/{campaign_id} [post]
func (h *CampaignHandler) CreateCampaignSeo(c echo.Context) error {
	campaignID := c.Param("campaign_id")
	var body masterschema.CampaignSeoPayload

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body payload")
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send to usecase for insert logic
	err := h.Dependencies.UC.CreateCampaignSeo(campaignID, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send success response
	return c.JSON(http.StatusCreated, commonschema.ResponseHTTP{
		Code:    http.StatusCreated,
		Message: "Data is successfully created",
	})
}

// @Security BearerAuth
// @Summary      Update SEO Campaign
// @Description  Update existing seo campaign data
// @Tags         Master - Campaign SEO
// @Accept  	 json
// @Produce  	 json
// @Param 		 campaign_id path string true "Campaign ID"
// @Param 		 id path string true "ID of your data"
// @Param        campaignSeoPayload  body      masterschema.CampaignSeoPayload   true  "SEO campaign payload"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/campaigns/seos/{campaign_id}/{id} [put]
func (h *CampaignHandler) UpdateCampaignSeo(c echo.Context) error {
	// get payload and parameters
	campaignID := c.Param("campaign_id")
	ID := c.Param("id")
	var body masterschema.CampaignSeoPayload

	// check for valid body
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body payload")
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send to usecase for update logic
	err := h.Dependencies.UC.UpdateCampaignSeo(campaignID, ID, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusNoContent, nil)
}

// @Security BearerAuth
// @Summary      Delete SEO Campaign
// @Description  Delete existing seo campaign data
// @Tags         Master - Campaign SEO
// @Accept  	 json
// @Produce  	 json
// @Param 		 campaign_id path string true "Campaign ID"
// @Param 		 id path string true "ID of your data"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/campaigns/seos/{campaign_id}/{id} [delete]
func (h *CampaignHandler) DeleteCampaignSeo(c echo.Context) error {
	// get parameters
	campaignID := c.Param("campaign_id")
	ID := c.Param("id")

	// send to usecase to do delete process
	if err := h.Dependencies.UC.DeleteCampaignSeo(campaignID, ID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusNoContent, nil)
}

// @Security BearerAuth
// @Summary      Form Entries Graphic
// @Description  Get line-graph for form entries for last 30days
// @Tags         Master - Campaign Analytics
// @Accept  	 json
// @Produce  	 json
// @Param 		 workspace_id path string true "Workspace ID"
// @Param 		 campaign_id path string true "Campaign ID"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/campaigns/analytics/dashboard/{workspace_id}/{campaign_id} [get]
func (h *CampaignHandler) DashboardAnalytics(c echo.Context) error {
	// get parameters
	workspaceID := c.Param("workspace_id")
	campaignID := c.Param("campaign_id")

	// check allowed user
	err := helpers.CheckAllowedCampaign(c, workspaceID, campaignID, h.Dependencies.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// get data form entries by workspace and campaign group by date
	// get only for last 60 days (for this version)
	data, err := h.Dependencies.UC.FindSummaryEntriesByDate(workspaceID, campaignID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send response
	response := commonschema.ResponseHTTP{
		Code:    http.StatusOK,
		Message: "Request success",
		Data:    data,
	}
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      Form Entries List
// @Description  List of form entries based on workspace and campaign
// @Tags         Master - Campaign Analytics
// @Accept  	 json
// @Produce  	 json
// @Param 		 workspace_id path string true "Workspace ID"
// @Param 		 campaign_id path string true "Campaign ID"
// @Param 		 page query int true "Page of list data"
// @Param 		 limit query int true "Limitting data you want to get"
// @Param 		 search query string false "Find your data with keywords"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/campaigns/analytics/form_entries/{workspace_id}/{campaign_id} [get]
func (h *CampaignHandler) FindFormEntries(c echo.Context) error {
	// get parameters
	workspaceID := c.Param("workspace_id")
	campaignID := c.Param("campaign_id")
	params := utils.QParams(c)
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// check allowed user
	err := helpers.CheckAllowedCampaign(c, workspaceID, campaignID, h.Dependencies.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send to usecase to get data
	list, err := h.Dependencies.UC.FindFormEntries(workspaceID, campaignID, params)
	if err != nil {
		response.Message = err.Error()
		return c.JSON(response.Code, response)
	}

	// send response
	response.Code = http.StatusOK
	response.Message = "Request success"
	response.Data = list
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      Form Entry detail
// @Description  Detail of user entry for each campaign
// @Tags         Master - Campaign Analytics
// @Accept  	 json
// @Produce  	 json
// @Param 		 workspace_id path string true "Workspace ID"
// @Param 		 campaign_id path string true "Campaign ID"
// @Param 		 id path string true "ID"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/campaigns/analytics/form_entries/{workspace_id}/{campaign_id}/{id} [get]
func (h *CampaignHandler) FindDetailFormEntry(c echo.Context) error {
	// get parameters
	workspaceID := c.Param("workspace_id")
	campaignID := c.Param("campaign_id")
	ID := c.Param("id")
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// check allowed user
	err := helpers.CheckAllowedCampaign(c, workspaceID, campaignID, h.Dependencies.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send to usecase to get data
	data, err := h.Dependencies.UC.FindFormEntry(c, ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	// send response
	response.Code = http.StatusOK
	response.Message = "Request success"
	response.Data = data
	return c.JSON(response.Code, response)
}
