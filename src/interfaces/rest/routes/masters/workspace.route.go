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

type WorksapceHandler struct {
	DB           *gorm.DB
	Validator    *validator.Validate
	Dependencies masterdi.WorkspaceDependencies
}

func NewWorkspaceHandler(DB *gorm.DB, validator *validator.Validate, dependencies masterdi.WorkspaceDependencies) *WorksapceHandler {
	return &WorksapceHandler{
		DB:           DB,
		Validator:    validator,
		Dependencies: dependencies,
	}
}

func NewWorkspaceHTTP(g *echo.Group, DB *gorm.DB) {
	validator := validator.New()
	h := NewWorkspaceHandler(DB, validator, *masterdi.NewWorkspaceDependencies(DB))

	// define endpoints
	w := g.Group("/workspaces")
	w.GET("", h.GetWorkspaces)
}

// @Security BearerAuth
// @Summary      List Workspaces
// @Description  Get the list of workspaces you created
// @Tags         Master - Workspace
// @Accept  	 json
// @Produce  	 json
// @Param 		 page query int true "Page of list data"
// @Param 		 limit query int true "Limitting data you want to get"
// @Param 		 search query string false "Find your data with keywords"
// @Param 		 orderBy query string false "Ordering data" example(created_at:desc)
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/workspaces [get]
func (h *WorksapceHandler) GetWorkspaces(c echo.Context) error {
	// define response
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// perform to get data
	params := utils.QParams(c)
	list, err := h.Dependencies.UC.GetWorkspaces(params)
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
