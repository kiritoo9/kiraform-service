package masterroute

import (
	"fmt"
	masterdi "kiraform/src/applications/dependencies/masters"
	"kiraform/src/interfaces/rest/middlewares"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	masterschema "kiraform/src/interfaces/rest/schemas/masters"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type FormEntryHandler struct {
	DB           *gorm.DB
	Validator    *validator.Validate
	Dependencies masterdi.FormEntryDependencies
}

func NewFormEntryHandler(DB *gorm.DB, validator *validator.Validate, dependencies masterdi.FormEntryDependencies) *FormEntryHandler {
	return &FormEntryHandler{
		DB:           DB,
		Validator:    validator,
		Dependencies: dependencies,
	}
}

func NewFormEntryHTTP(g *echo.Group, DB *gorm.DB) {
	validator := validator.New()
	h := NewFormEntryHandler(DB, validator, *masterdi.NewFormEntryDependencies(DB))

	// define [unauthrozid] endpoints
	fe := g.Group("/form_entries")
	fe.GET("/:campaign_key", h.PreviewForm)
	fe.POST("/:campaign_id", h.EntryForm)

	// define [authorized] endpointes
	// pfe = private_form_entries
	pfe := g.Group("/form_entries") // re-define
	pfe.Use(middlewares.VerifyToken)

	pfe.GET("/history", h.GetHistory)
	pfe.GET("/history/:id", h.GetDetailHistory)
}

// @Summary      Preview Form
// @Description  Get detail form by key
// @Tags         Transaction - Form Entry
// @Accept  	 json
// @Produce  	 json
// @Param 		 campaign_key path string true "campaign key"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/form_entries/{campaign_key} [get]
func (h *FormEntryHandler) PreviewForm(c echo.Context) error {
	// get parameters
	campaignKey := c.Param("campaign_key")

	// send to usecase to validate and get detail of campaign
	isPublish := true // find only published campaign
	data, err := h.Dependencies.UCcampaign.FindCampaignByKey(campaignKey, &isPublish)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, nil)
	}

	// get detail form by this campaign
	forms, err := h.Dependencies.UCcampaign.FindFormsByCampaign(data.ID.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// get attributes each form
	// handle if forms is null, because there is possibility for it
	for i, v := range forms {
		attr, err := h.Dependencies.UCcampaign.FindFormAttributes(v.ID.String())
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
		forms[i].Attributes = attr
	}

	// prepare response
	data.Forms = forms
	response := commonschema.ResponseHTTP{
		Code:    http.StatusOK,
		Message: "Request success",
		Data:    data,
	}
	return c.JSON(response.Code, response)
}

// @Summary      Form Entry
// @Description  Submit user value for this form
// @Tags         Transaction - Form Entry
// @Accept  	 json
// @Produce  	 json
// @Param 		 campaign_id path string true "Campaign ID"
// @Param        formEntryPayload  body      []masterschema.FormEntryPayload   true  "form entry payload"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/form_entries/{campaign_id} [post]
func (h *FormEntryHandler) EntryForm(c echo.Context) error {
	// get parameters
	campaignID := c.Param("campaign_id")
	userID, _ := c.Get("user_id").(string)
	var body []masterschema.FormEntryPayload

	// validate body
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body payload")
	}

	for i, v := range body {
		if err := h.Validator.Struct(v); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid item at index %d: %s", i, err.Error()))
		}
	}

	// send to usecase for business process
	err := h.Dependencies.UC.EntryForm(campaignID, &userID, body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// send success response
	return c.JSON(http.StatusCreated, commonschema.ResponseHTTP{
		Code:    http.StatusCreated,
		Message: "Your data is successfully submitted",
	})
}

func (h *FormEntryHandler) GetHistory(c echo.Context) error {
	return echo.NewHTTPError(http.StatusBadRequest, nil)
}

func (h *FormEntryHandler) GetDetailHistory(c echo.Context) error {
	return echo.NewHTTPError(http.StatusBadRequest, nil)
}
