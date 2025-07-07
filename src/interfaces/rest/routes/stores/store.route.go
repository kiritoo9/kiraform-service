package storeroute

import (
	"errors"
	storedi "kiraform/src/applications/dependencies/stores"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	storeschema "kiraform/src/interfaces/rest/schemas/stores"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type StoreHandler struct {
	DB           *gorm.DB
	Validator    *validator.Validate
	Dependencies storedi.StoreDependencies
}

func NewStoreHandler(DB *gorm.DB, validator *validator.Validate, dependencies storedi.StoreDependencies) *StoreHandler {
	return &StoreHandler{
		DB:           DB,
		Validator:    validator,
		Dependencies: dependencies,
	}
}

func NewStoreHTTP(g *echo.Group, DB *gorm.DB) {
	validator := validator.New()
	h := NewStoreHandler(DB, validator, *storedi.NewStoreDependencies(DB))

	// define store routes
	s := g.Group("/store")
	s.GET("", h.FindStore)
	s.PUT("", h.UpdateStore)

	// define store product category routes
	// spc := s.Group("/product_categories")
	// spc.GET("", h.FindStoreProductCategories)
	// spc.GET("/:id", h.FindStoreProductCategory)
	// spc.POST("", h.CreateStoreProductCategory)
	// spc.PUT("", h.UpdateStoreProductCategory)
	// spc.DELETE("/:id", h.DeleteStoreProductCategory)

	// define store product routes
	// sp := s.Group("/products")
	// sp.GET("", h.FindStoreProducts)
	// sp.GET("/:id", h.FindStoreProduct)
	// sp.POST("", h.CreateStoreProduct)
	// sp.PUT("", h.UpdateStoreProduct)
	// sp.DELETE("/:id", h.DeleteStoreProduct)
}

// @Security BearerAuth
// @Summary      Store Profile
// @Description  Get store profile for logged user
// @Tags         Store
// @Accept  	 json
// @Produce  	 json
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/store [get]
func (h *StoreHandler) FindStore(c echo.Context) error {
	// prepare usable data
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}
	userID, _ := c.Get("user_id").(string)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("your token is not valid"))
	}

	// get data from usecase
	data, err := h.Dependencies.UC.FindStore(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Code = http.StatusNotFound
		}
		return echo.NewHTTPError(response.Code, err.Error())
	}

	// send response
	response.Code = http.StatusOK
	response.Message = "Request success"
	response.Data = data
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      Update Store Profile
// @Description  Update store profile for logged user
// @Tags         Store
// @Accept  	 json
// @Produce  	 json
// @Param        storePayload  body      storeschema.StorePayload   true  "store payload"
// @Success      204  {object} commonschema.ResponseHTTP "Data updated"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/store [put]
func (h *StoreHandler) UpdateStore(c echo.Context) error {
	// prepare usable data
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}
	userID, _ := c.Get("user_id").(string)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("your token is not valid"))
	}
	var body storeschema.StorePayload

	// validate body
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send to usecase for store-update process
	err := h.Dependencies.UC.UpdateStore(userID, body)
	if err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	// send response
	response.Code = http.StatusNoContent
	response.Message = "Data updated!"
	return c.JSON(response.Code, response)
}
