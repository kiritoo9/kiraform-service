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

type WorkspaceHandler struct {
	DB           *gorm.DB
	Validator    *validator.Validate
	Dependencies masterdi.WorkspaceDependencies
}

func NewWorkspaceHandler(DB *gorm.DB, validator *validator.Validate, dependencies masterdi.WorkspaceDependencies) *WorkspaceHandler {
	return &WorkspaceHandler{
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
	w.GET("", h.FindWorkspaces)
	w.GET("/:id", h.FindWorkspace)
	w.POST("", h.CreateWorkspace)
	w.PUT("/:id", h.UpdateWorkspace)
	w.DELETE("/:id", h.DeleteWokspace)
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
func (h *WorkspaceHandler) FindWorkspaces(c echo.Context) error {
	// define response
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// perform to get data
	params := utils.QParams(c)
	list, err := h.Dependencies.UC.FindWorkspaces(params)
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
// @Summary      Detail Workspace
// @Description  Get detail data of workspace you choose
// @Tags         Master - Workspace
// @Accept  	 json
// @Produce  	 json
// @Param 		 id path string true "ID of your data"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/workspaces/{id} [get]
func (h *WorkspaceHandler) FindWorkspace(c echo.Context) error {
	// define response
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// get data by id
	id := c.Param("id")
	data, err := h.Dependencies.UC.FindWorkspaceByID(id)
	if err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusNotFound, response)
	}

	// send response
	response.Code = http.StatusOK
	response.Message = "Request success"
	response.Data = data
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      Create Workspace
// @Description  Create new workspace data
// @Tags         Master - Workspace
// @Accept  	 json
// @Produce  	 json
// @Param        workspacePayload  body      masterschema.WorkspacePayload   true  "Workspace payload"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/workspaces [post]
func (h *WorkspaceHandler) CreateWorkspace(c echo.Context) error {
	var body masterschema.WorkspacePayload

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body")
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// call usecase for busines validation
	err := h.Dependencies.UC.CreateWorkspace(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send success response
	response := commonschema.ResponseHTTP{
		Code:    http.StatusCreated,
		Message: "Data is successfully created",
	}
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      Update Workspace
// @Description  Update existing workspace data
// @Tags         Master - Workspace
// @Accept  	 json
// @Produce  	 json
// @Param 		 id path string true "ID of your data"
// @Param        workspacePayload  body      masterschema.WorkspacePayload   true  "Workspace payload"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/workspaces/{id} [put]
func (h *WorkspaceHandler) UpdateWorkspace(c echo.Context) error {
	ID := c.Param("id")
	var body masterschema.WorkspacePayload

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body")
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// call usecase for business validation
	err := h.Dependencies.UC.UpdateWorkspace(ID, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send success response
	return c.JSON(http.StatusNoContent, nil)
}

// @Security BearerAuth
// @Summary      Delete Workspace
// @Description  Delete existing workspace data
// @Tags         Master - Workspace
// @Accept  	 json
// @Produce  	 json
// @Param 		 id path string true "ID of your data"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/workspaces/{id} [delete]
func (h *WorkspaceHandler) DeleteWokspace(c echo.Context) error {
	ID := c.Param("id")

	// call usecase for business validation
	err := h.Dependencies.UC.DeleteWorkspace(ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send success response
	return c.JSON(http.StatusNoContent, nil)
}
