package storeroute

import (
	"errors"
	storedi "kiraform/src/applications/dependencies/stores"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	storeschema "kiraform/src/interfaces/rest/schemas/stores"
	"kiraform/src/utils"
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
	spc := s.Group("/product_categories")
	spc.GET("", h.FindStoreProductCategories)
	spc.GET("/:id", h.FindStoreProductCategory)
	spc.POST("", h.CreateStoreProductCategory)
	spc.PUT("/:id", h.UpdateStoreProductCategory)
	spc.DELETE("/:id", h.DeleteStoreProductCategory)

	// define store product routes
	sp := s.Group("/products")
	sp.GET("", h.FindStoreProducts)
	sp.GET("/form_entries/:id", h.FindStoreProductFormEntries)
	sp.GET("/:id", h.FindStoreProduct)
	sp.POST("", h.CreateStoreProduct)
	sp.PUT("/:id", h.UpdateStoreProduct)
	sp.DELETE("/:id", h.DeleteStoreProduct)
}

// @Security BearerAuth
// @Summary      Store Profile
// @Description  Get store profile for logged user
// @Tags         Store - Profile
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
	data, err := h.Dependencies.UC.FindStore(c, userID)
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
// @Summary      Update Store Profile
// @Description  Update store profile for logged user
// @Tags         Store - Profile
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

// @Security BearerAuth
// @Summary      List Product Categories
// @Description  Get the list of product categories based on logged user
// @Tags         Store - Product Categories
// @Accept  	 json
// @Produce  	 json
// @Param 		 page query int true "Page of list data"
// @Param 		 limit query int true "Limitting data you want to get"
// @Param 		 search query string false "Find your data with keywords"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/store/product_categories [get]
func (h *StoreHandler) FindStoreProductCategories(c echo.Context) error {
	// get parameters
	params := utils.QParams(c)
	userID, _ := c.Get("user_id").(string)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("your token is not valid"))
	}
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// get list of product categories
	list, err := h.Dependencies.UC.FindStoreProductCategories(userID, params)
	if err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	// send response
	response.Code = http.StatusOK
	response.Message = "Request success"
	response.Data = list
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      Detail Product Category
// @Description  Get detail of product category
// @Tags         Store - Product Categories
// @Accept  	 json
// @Produce  	 json
// @Param 		 id path string true "ID of your data"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/store/product_categories/{id} [get]
func (h *StoreHandler) FindStoreProductCategory(c echo.Context) error {
	// get parameters
	ID := c.Param("id")
	userID, _ := c.Get("user_id").(string)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("your token is not valid"))
	}
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// get detail data
	data, err := h.Dependencies.UC.FindStoreProductCategory(userID, ID)
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
// @Summary      Create Product Category
// @Description  Create new product category for logged user store
// @Tags         Store - Product Categories
// @Accept  	 json
// @Produce  	 json
// @Param        productCategoryPayload  body      storeschema.ProductCategoryPayload   true  "product category payload"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/store/product_categories [post]
func (h *StoreHandler) CreateStoreProductCategory(c echo.Context) error {
	// get parameters
	userID, _ := c.Get("user_id").(string)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("your token is not valid"))
	}
	var body storeschema.ProductCategoryPayload
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// validate body
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	// perform to create data
	err := h.Dependencies.UC.CreateStoreProductCategory(userID, body)
	if err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	// send response
	response.Code = http.StatusCreated
	response.Message = "Data created"
	response.Data = body
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      Update Product Category
// @Description  Update existing product category for logged user store
// @Tags         Store - Product Categories
// @Accept  	 json
// @Produce  	 json
// @Param 		 id path string true "ID of your data"
// @Param        productCategoryPayload  body      storeschema.ProductCategoryPayload   true  "product category payload"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/store/product_categories/{id} [put]
func (h *StoreHandler) UpdateStoreProductCategory(c echo.Context) error {
	// get parameters
	ID := c.Param("id")
	userID, _ := c.Get("user_id").(string)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("your token is not valid"))
	}
	var body storeschema.ProductCategoryPayload
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// validate body
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	// perform to create data
	err := h.Dependencies.UC.UpdateStoreProductCategory(userID, ID, body)
	if err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	// send response
	response.Code = http.StatusNoContent
	response.Message = "Data updated"
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      Delete Product Category
// @Description  Delete existing product category for logged user store data
// @Tags         Store - Product Categories
// @Accept  	 json
// @Produce  	 json
// @Param 		 id path string true "ID of your data"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/store/product_categories/{id} [delete]
func (h *StoreHandler) DeleteStoreProductCategory(c echo.Context) error {
	// get parameters
	ID := c.Param("id")
	userID, _ := c.Get("user_id").(string)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("your token is not valid"))
	}
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// perform to create data
	err := h.Dependencies.UC.DeleteStoreProductCategory(userID, ID)
	if err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	// send response
	response.Code = http.StatusNoContent
	response.Message = "Data deleted"
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      List Products
// @Description  Get the list of products based on logged user
// @Tags         Store - Products
// @Accept  	 json
// @Produce  	 json
// @Param 		 page query int true "Page of list data"
// @Param 		 limit query int true "Limitting data you want to get"
// @Param 		 search query string false "Find your data with keywords"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/store/products [get]
func (h *StoreHandler) FindStoreProducts(c echo.Context) error {
	// prepare usable data
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}
	userID, _ := c.Get("user_id").(string)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("your token is not valid"))
	}
	params := utils.QParams(c)

	// get data from usecase
	data, err := h.Dependencies.UC.FindStoreProducts(c, userID, params)
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
// @Summary      List Form Entries
// @Description  Get the list form entries of product
// @Tags         Store - Products
// @Accept  	 json
// @Produce  	 json
// @Param 		 page query int true "Page of list data"
// @Param 		 limit query int true "Limitting data you want to get"
// @Param 		 search query string false "Find your data with keywords"
// @Param 		 id path string true "ID of your data"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/store/products/form_entries/{id} [get]
func (h *StoreHandler) FindStoreProductFormEntries(c echo.Context) error {
	// prepare usable data
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}
	ID := c.Param("id")
	userID, _ := c.Get("user_id").(string)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("your token is not valid"))
	}
	params := utils.QParams(c)

	// get data from usecase
	data, err := h.Dependencies.UC.FindStoreProductFormEntries(userID, ID, params)
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
// @Summary      Detail Product
// @Description  Get detail of product
// @Tags         Store - Products
// @Accept  	 json
// @Produce  	 json
// @Param 		 id path string true "ID of your data"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/store/products/{id} [get]
func (h *StoreHandler) FindStoreProduct(c echo.Context) error {
	// prepare usable data
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}
	ID := c.Param("id")
	userID, _ := c.Get("user_id").(string)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("your token is not valid"))
	}

	// get data from usecase
	data, err := h.Dependencies.UC.FindStoreProduct(c, userID, ID)
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
// @Summary      Create Product
// @Description  Create new product for logged user store
// @Tags         Store - Products
// @Accept  	 json
// @Produce  	 json
// @Param        productPayload  body      storeschema.ProductPayload   true  "product payload"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/store/products [post]
func (h *StoreHandler) CreateStoreProduct(c echo.Context) error {
	// get parameters
	userID, _ := c.Get("user_id").(string)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("your token is not valid"))
	}
	var body storeschema.ProductPayload
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// validate body
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	// perform to create data
	err := h.Dependencies.UC.CreateStoreProduct(userID, body)
	if err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	// send response
	response.Code = http.StatusCreated
	response.Message = "Data created"
	response.Data = body
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      Update Product
// @Description  Update existing product for logged user store
// @Tags         Store - Products
// @Accept  	 json
// @Produce  	 json
// @Param 		 id path string true "ID of your data"
// @Param        productPayload  body      storeschema.ProductPayload   true  "product payload"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/store/products/{id} [put]
func (h *StoreHandler) UpdateStoreProduct(c echo.Context) error {
	// get parameters
	ID := c.Param("id")
	userID, _ := c.Get("user_id").(string)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("your token is not valid"))
	}
	var body storeschema.ProductPayload
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// validate body
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	// perform to create data
	err := h.Dependencies.UC.UpdateStoreProduct(userID, ID, body)
	if err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	// send response
	response.Code = http.StatusNoContent
	response.Message = "Data updated"
	return c.JSON(response.Code, response)
}

// @Security BearerAuth
// @Summary      Delete Product
// @Description  Delete existing product for logged user store data
// @Tags         Store - Products
// @Accept  	 json
// @Produce  	 json
// @Param 		 id path string true "ID of your data"
// @Success      200  {object} commonschema.ResponseHTTP "Request success"
// @Failure      400  {object} commonschema.ResponseHTTP "Request failure"
// @Router       /api/store/products/{id} [delete]
func (h *StoreHandler) DeleteStoreProduct(c echo.Context) error {
	// get parameters
	ID := c.Param("id")
	userID, _ := c.Get("user_id").(string)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, errors.New("your token is not valid"))
	}
	response := commonschema.ResponseHTTP{Code: http.StatusBadRequest}

	// perform to create data
	err := h.Dependencies.UC.DeleteStoreProduct(userID, ID)
	if err != nil {
		return echo.NewHTTPError(response.Code, err.Error())
	}

	// send response
	response.Code = http.StatusNoContent
	response.Message = "Data deleted"
	return c.JSON(response.Code, response)
}
