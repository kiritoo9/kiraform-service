package storeroute

import (
	storedi "kiraform/src/applications/dependencies/stores"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
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
	// s.PUT("", h.UpdateStore)

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

func (h *StoreHandler) FindStore(c echo.Context) error {
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// get data from usecase
	data, err := h.Dependencies.UC.FindStore()
	if err != nil {
		return echo.NewHTTPError(response.Code, err)
	}

	// send response
	response.Code = http.StatusOK
	response.Message = "Request success"
	response.Data = data
	return c.JSON(response.Code, response)
}
