package storeroute

import (
	storedi "kiraform/src/applications/dependencies/stores"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	"kiraform/src/utils"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type StorePublicHandler struct {
	DB           *gorm.DB
	Validator    *validator.Validate
	Dependencies storedi.StoreDependencies
}

func NewStorePublicHandler(DB *gorm.DB, validator *validator.Validate, dependencies storedi.StoreDependencies) *StorePublicHandler {
	return &StorePublicHandler{
		DB:           DB,
		Validator:    validator,
		Dependencies: dependencies,
	}
}

func NewStorePublicHTTP(g *echo.Group, DB *gorm.DB) {
	validator := validator.New()
	h := NewStorePublicHandler(DB, validator, *storedi.NewStoreDependencies(DB))

	s := g.Group("/storepub")
	s.GET("/:key", h.FindStore)
	s.GET("/categories/:key", h.FindStoreCategories)
	s.GET("/products/:key", h.FindStoreProducts)
}

// @Security BearerAuth
// @Summary      Store Informations
// @Description  Get store information by key
// @Tags         Public - Stores
// @Accept  	 json
// @Produce  	 json
// @Param 		 key path string true "key of store you want to get"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/storepub/{key} [get]
func (h *StorePublicHandler) FindStore(c echo.Context) error {
	// prepare usable data
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}
	key := c.Param("key")

	// get data from usecase
	data, err := h.Dependencies.UC.FindStoreByKey(c, key)
	if err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	// send response
	response.Code = http.StatusOK
	response.Message = "Request success"
	response.Data = data
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      Store Categories
// @Description  Get store categories
// @Tags         Public - Stores
// @Accept  	 json
// @Produce  	 json
// @Param 		 key path string true "key of store you want to get"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/storepub/categories/{key} [get]
func (h *StorePublicHandler) FindStoreCategories(c echo.Context) error {
	// prepare usable data
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}
	key := c.Param("key")

	// get data from usecase
	data, err := h.Dependencies.UC.FindStoreCategoriesByKey(key)
	if err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	// send response
	response.Code = http.StatusOK
	response.Message = "Request success"
	response.Data = data
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      Store Products
// @Description  Get store products
// @Tags         Public - Stores
// @Accept  	 json
// @Produce  	 json
// @Param 		 key path string true "key of store you want to get"
// @Param 		 page query int true "Page of list data"
// @Param 		 limit query int true "Limitting data you want to get"
// @Param 		 search query string false "Find your data with keywords"
// @Param 		 category_id query string false "category of product you want to get"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/storepub/products/{key} [get]
func (h *StorePublicHandler) FindStoreProducts(c echo.Context) error {
	// get parameters
	params := utils.QParams(c)
	key := c.Param("key")
	category_id := c.QueryParam("category_id")
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// get list of product by store key
	list, err := h.Dependencies.UC.FindStoreProductsByStoreKey(c, key, params, &category_id)
	if err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	// send response
	response.Code = http.StatusOK
	response.Message = "Request success"
	response.Data = list
	return c.JSON(response.Code, response)
}
