package masterroute

import (
	masterdi "kiraform/src/applications/dependencies/masters"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	"kiraform/src/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CampaingHandler struct {
	DB           *gorm.DB
	Validator    *validator.Validate
	Dependencies masterdi.CampaignDependencies
}

func NewCampaignHandler(DB *gorm.DB, validator *validator.Validate, dependencies masterdi.CampaignDependencies) *CampaingHandler {
	return &CampaingHandler{
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
func (h *CampaingHandler) FindCampaigns(c echo.Context) error {
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
