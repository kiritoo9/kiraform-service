package authroute

import (
	"errors"
	authdi "kiraform/src/applications/dependencies/auths"
	authschema "kiraform/src/interfaces/rest/schemas/auths"
	commonschema "kiraform/src/interfaces/rest/schemas/commons"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AuthHandler struct {
	DB           *gorm.DB
	Validator    *validator.Validate
	Dependencies authdi.AuthDependencies
}

func NewAuthHandler(DB *gorm.DB, validator *validator.Validate, dependencies authdi.AuthDependencies) *AuthHandler {
	return &AuthHandler{
		DB:           DB,
		Validator:    validator,
		Dependencies: dependencies,
	}
}

func NewAuthHTTP(g *echo.Group, DB *gorm.DB) {
	validator := validator.New()
	h := NewAuthHandler(DB, validator, *authdi.NewAuthDependencies(DB))

	// define endpoints
	g.POST("/login", h.Login)
	g.POST("/register", h.Register)
}

// @Summary      Login
// @Description  User login
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        loginPayload  body      authschema.LoginPayload   true  "Login credentials"
// @Success      200  {object} commonschema.ResponseHTTP "Login success"
// @Failure      400  {object} commonschema.ResponseHTTP "Login failure"
// @Router       /api/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var body authschema.LoginPayload

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// call usecase for busines validation
	signedToken, err := h.Dependencies.UC.Login(body)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	// send response
	response := commonschema.ResponseHTTP{
		Code:    http.StatusOK,
		Message: "Login success",
		Data: map[string]any{
			"access_token":  signedToken,
			"refresh_token": signedToken,
		},
	}
	return c.JSON(response.Code, response)
}

// @Summary      Registration
// @Description  You can regist new user here
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        registerPayload  body      authschema.RegisterPayload   true  "Register credentials"
// @Success      200  {object} commonschema.ResponseHTTP "Registration success"
// @Failure      400  {object} commonschema.ResponseHTTP "Registration failure"
// @Router       /api/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var body authschema.RegisterPayload

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body")
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// call usecase for business validation
	msg, err := h.Dependencies.UC.Register(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// send response
	response := commonschema.ResponseHTTP{
		Code:    http.StatusCreated,
		Message: *msg,
	}
	return c.JSON(response.Code, response)
}
