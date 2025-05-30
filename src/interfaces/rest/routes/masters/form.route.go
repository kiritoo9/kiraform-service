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

type FormHandler struct {
	DB           *gorm.DB
	Validator    *validator.Validate
	Dependencies masterdi.FormDependencies
}

func NewFormHandler(DB *gorm.DB, validate *validator.Validate, dependencies masterdi.FormDependencies) *FormHandler {
	return &FormHandler{
		DB:           DB,
		Validator:    validate,
		Dependencies: dependencies,
	}
}

func NewFormHTTP(g *echo.Group, DB *gorm.DB) {
	validator := validator.New()
	h := NewFormHandler(DB, validator, *masterdi.NewFormDependencies(DB))

	// define endpoints
	w := g.Group("/forms")
	w.GET("", h.FindForms)
	w.GET("/:id", h.FindForm)
}

// @Security BearerAuth
// @Summary      List Forms
// @Description  Get the list of forms you created
// @Tags         Master - Forms
// @Accept  	 json
// @Produce  	 json
// @Param 		 page query int true "Page of list data"
// @Param 		 limit query int true "Limitting data you want to get"
// @Param 		 search query string false "Find your data with keywords"
// @Param 		 orderBy query string false "Ordering data" example(created_at:desc)
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/forms [get]
func (h *FormHandler) FindForms(c echo.Context) error {
	// define response
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// perform to get data
	params := utils.QParams(c)
	list, err := h.Dependencies.UC.FindForms(params)
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
// @Summary      Detail Form
// @Description  Get detail data of form you choose
// @Tags         Master - Forms
// @Accept  	 json
// @Produce  	 json
// @Param 		 id path string true "ID of your data"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/forms/{id} [get]
func (h *FormHandler) FindForm(c echo.Context) error {
	// define response
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// get data by id
	id := c.Param("id")
	data, err := h.Dependencies.UC.FindFormByID(id)
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
