package masterroute

import (
	masterdi "kiraform/src/applications/dependencies/masters"
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
	w := g.Group("/campaigns")
	w.GET("/:workspace_id", h.FindCampaigns)
	w.GET("/:workspace_id/:id", h.FindCampaign)
	w.POST("/:workspace_id", h.CreateCampaign)
	w.PUT("/:workspace_id/:id", h.UpdateCampaign)
	w.DELETE("/:workspace_id/:id", h.DeleteCampaign)
}

// @Security BearerAuth
// @Summary      List Campaigns
// @Description  Get the list of campaigns you created
// @Tags         Master - Campaign
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
// @Summary      Detaiol Campaigns
// @Description  Get detail of campaigns you created
// @Tags         Master - Campaign
// @Accept  	 json
// @Produce  	 json
// @Param 		 workspace_id path string true "Workspace ID"
// @Param 		 id path string true "ID of your data"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/campaigns/{workspace_id}/{id} [get]
func (h *CampaignHandler) FindCampaign(c echo.Context) error {
	// get paramters
	workspaceID := c.Param("workspace_id")
	ID := c.Param("id")

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
// @Summary      Create Campaign
// @Description  Create new campaign data
// @Tags         Master - Campaign
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

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body payload")
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send to usecase for insert logic
	err := h.Dependencies.UC.CreateCampaign(workspaceID, body)
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
// @Tags         Master - Campaign
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

	// check for valid body
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body payload")
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send to usecase for update logic
	err := h.Dependencies.UC.UpdateCampaign(workspaceID, ID, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusNoContent, nil)
}

// @Security BearerAuth
// @Summary      Delete Campaign
// @Description  Delete existing campaign data
// @Tags         Master - Campaign
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

	// send to usecase to do delete process
	if err := h.Dependencies.UC.DeleteCampaign(workspaceID, ID); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusNoContent, nil)
}
