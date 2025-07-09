package masterroute

import (
	"fmt"
	masterdi "kiraform/src/applications/dependencies/masters"
	"kiraform/src/interfaces/rest/middlewares"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	masterschema "kiraform/src/interfaces/rest/schemas/masters"
	"kiraform/src/utils"
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
// @Param 		 product_id query string false "Product ID"
// @Param        formEntryPayload  body      []masterschema.FormEntryPayload   true  "form entry payload"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/form_entries/{campaign_id} [post]
func (h *FormEntryHandler) EntryForm(c echo.Context) error {
	// get parameters
	campaignID := c.Param("campaign_id")
	userID, _ := c.Get("user_id").(string)
	var body []masterschema.FormEntryPayload
	productID := c.QueryParam("product_id")

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
	err := h.Dependencies.UC.EntryForm(campaignID, &userID, body, &productID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// send success response
	return c.JSON(http.StatusCreated, commonschema.ResponseHTTP{
		Code:    http.StatusCreated,
		Message: "Your data is successfully submitted",
	})
}

// @Security BearerAuth
// @Summary      History of your entries
// @Description  Get the list of history of your entries
// @Tags         Transaction - Form Entry
// @Accept  	 json
// @Produce  	 json
// @Param 		 page query int true "Page of list data"
// @Param 		 limit query int true "Limitting data you want to get"
// @Param 		 search query string false "Find your data with keywords"
// @Param 		 orderBy query string false "Ordering data" example(created_at:desc)
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/form_entries/history [get]
func (h *FormEntryHandler) GetHistory(c echo.Context) error {
	// get parameters
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Your account is not auhtorized yet")
	}
	params := utils.QParams(c)

	// get history data
	data, err := h.Dependencies.UC.GetHistory(userID, params)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send response
	return c.JSON(http.StatusOK, data)

}

// @Security BearerAuth
// @Summary      Detail entry
// @Description  Get detail of your curren entry
// @Tags         Transaction - Form Entry
// @Accept  	 json
// @Produce  	 json
// @Param 		 id path string false "ID of your history"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/form_entries/history/{id} [get]
func (h *FormEntryHandler) GetDetailHistory(c echo.Context) error {
	// get parameters
	ID := c.Param("id")
	userID, ok := c.Get("user_id").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusBadRequest, "Your account is not auhtorized yet")
	}

	// get detail of history
	data, err := h.Dependencies.UC.GetDetailHistory(userID, ID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// send response
	return c.JSON(http.StatusOK, commonschema.ResponseHTTP{
		Code:    http.StatusOK,
		Message: "Request success",
		Data:    data,
	})
}
