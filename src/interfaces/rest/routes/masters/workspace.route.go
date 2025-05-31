package masterroute

import (
	masterdi "kiraform/src/applications/dependencies/masters"
	"kiraform/src/applications/helpers"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	masterschema "kiraform/src/interfaces/rest/schemas/masters"
	"kiraform/src/utils"
	"net/http"
	"strings"

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

	// workspace endpoints
	w := g.Group("/workspaces")
	w.GET("", h.FindWorkspaces)
	w.GET("/:id", h.FindWorkspace)
	w.POST("", h.CreateWorkspace)
	w.PUT("/:id", h.UpdateWorkspace)
	w.DELETE("/:id", h.DeleteWokspace)

	// workspace user endpoints
	wu := w.Group("/users")
	wu.GET("/:workspace_id", h.FindWorkspaceUsers)
	wu.GET("/:workspace_id/:id", h.FindWorkspaceUser)
	wu.POST("/:workspace_id", h.CreateWorkspaceUser)
	wu.PUT("/:workspace_id/:id", h.UpdateWorkspaceUser)
	wu.DELETE("/:workspace_id/:id", h.DeleteWokspaceUser)
}

// @Security BearerAuth
// @Summary      List Workspaces
// @Description  Get the list of workspaces you created
// @Tags         Master - Workspaces
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

	// get logged data
	var userID *string
	roleName := c.Get("role_name").(string)
	if strings.ToLower(roleName) != "admin" {
		val, ok := c.Get("user_id").(string)
		if !ok {
			return c.JSON(response.Code, "missing user id")
		}
		userID = &val
	}

	// perform to get data
	params := utils.QParams(c)
	list, err := h.Dependencies.UC.FindWorkspaces(userID, params)
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
// @Tags         Master - Workspaces
// @Accept  	 json
// @Produce  	 json
// @Param 		 id path string true "ID of your data"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/workspaces/{id} [get]
func (h *WorkspaceHandler) FindWorkspace(c echo.Context) error {
	// define data
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}
	ID := c.Param("id")

	// check allowed access
	err := helpers.CheckAllowedWorkspace(c, ID, h.Dependencies.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// get data by id
	data, err := h.Dependencies.UC.FindWorkspaceByID(ID)
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
// @Tags         Master - Workspaces
// @Accept  	 json
// @Produce  	 json
// @Param        workspacePayload  body      masterschema.WorkspacePayload   true  "Workspace payload"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/workspaces [post]
func (h *WorkspaceHandler) CreateWorkspace(c echo.Context) error {
	// define data
	var body masterschema.WorkspacePayload
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing user id")
	}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body")
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// call usecase for busines validation
	err := h.Dependencies.UC.CreateWorkspace(userID, body)
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
// @Tags         Master - Workspaces
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
// @Tags         Master - Workspaces
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

// @Security BearerAuth
// @Summary      List Workspace Users
// @Description  Get the list of users available in this workspace
// @Tags         Master - Workspace Users
// @Accept  	 json
// @Produce  	 json
// @Param 		 workspace_id path string true "Workspace ID"
// @Param 		 page query int true "Page of list data"
// @Param 		 limit query int true "Limitting data you want to get"
// @Param 		 search query string false "Find your data with keywords"
// @Param 		 orderBy query string false "Ordering data" example(created_at:desc)
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/workspaces/users/{workspace_id} [get]
func (h *WorkspaceHandler) FindWorkspaceUsers(c echo.Context) error {
	// define data
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}
	workspaceID := c.Param("workspace_id")

	// check allowed access
	err := helpers.CheckAllowedWorkspace(c, workspaceID, h.Dependencies.DB)
	if err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	// perform to get data
	params := utils.QParams(c)
	list, err := h.Dependencies.UC.FindWorkspaceUsers(workspaceID, params)
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
// @Summary      Detail Workspace User
// @Description  Get detail data of workspace user you choose
// @Tags         Master - Workspace Users
// @Accept  	 json
// @Produce  	 json
// @Param 		 workspace_id path string true "Workspace ID"
// @Param 		 id path string true "ID of your data"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/workspaces/users/{workspace_id}/{id} [get]
func (h *WorkspaceHandler) FindWorkspaceUser(c echo.Context) error {
	// define data
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}
	workspaceID := c.Param("workspace_id")
	ID := c.Param("id")

	// check allowed access
	err := helpers.CheckAllowedWorkspace(c, workspaceID, h.Dependencies.DB)
	if err != nil {
		response.Message = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	// get data by id
	data, err := h.Dependencies.UC.FindWorkspaceUserByID(workspaceID, ID)
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
// @Summary      Create Workspace User
// @Description  Create new user for this workspace
// @Tags         Master - Workspace Users
// @Accept  	 json
// @Produce  	 json
// @Param 		 workspace_id path string true "Workspace ID"
// @Param        workspaceUserPayload  body      masterschema.WorkspaceUserPayload   true  "Workspace user payload"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/workspaces/users/{workspace_id} [post]
func (h *WorkspaceHandler) CreateWorkspaceUser(c echo.Context) error {
	// define data
	var body masterschema.WorkspaceUserPayload
	workspaceID := c.Param("workspace_id")

	// check allowed access
	err := helpers.CheckAllowedWorkspace(c, workspaceID, h.Dependencies.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body")
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// call usecase for busines validation
	err = h.Dependencies.UC.CreateWorkspaceUser(workspaceID, body)
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
// @Summary      Update Workspace User
// @Description  Update existing user in this workspace
// @Tags         Master - Workspace Users
// @Accept  	 json
// @Produce  	 json
// @Param 		 workspace_id path string true "Workspace ID"
// @Param 		 id path string true "ID of your data"
// @Param        workspaceUserUpdatePayload  body      masterschema.WorkspaceUserUpdatePayload   true  "Workspace user payload"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/workspaces/users/{workspace_id}/{id} [put]
func (h *WorkspaceHandler) UpdateWorkspaceUser(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	ID := c.Param("id")
	var body masterschema.WorkspaceUserUpdatePayload

	// check allowed access
	err := helpers.CheckAllowedWorkspace(c, workspaceID, h.Dependencies.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body")
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// call usecase for business validation
	err = h.Dependencies.UC.UpdateWorkspaceUser(workspaceID, ID, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send success response
	return c.JSON(http.StatusNoContent, nil)
}

// @Security BearerAuth
// @Summary      Delete Workspace User
// @Description  Delete existing user in this workspace
// @Tags         Master - Workspace Users
// @Accept  	 json
// @Produce  	 json
// @Param 		 workspace_id path string true "Workspace ID"
// @Param 		 id path string true "ID of your data"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/workspaces/users/{workspace_id}/{id} [delete]
func (h *WorkspaceHandler) DeleteWokspaceUser(c echo.Context) error {
	workspaceID := c.Param("workspace_id")
	ID := c.Param("id")

	// check allowed access
	err := helpers.CheckAllowedWorkspace(c, workspaceID, h.Dependencies.DB)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// call usecase for business validation
	err = h.Dependencies.UC.DeleteWorkspaceUser(workspaceID, ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send success response
	return c.JSON(http.StatusNoContent, nil)
}
