package authroute

import (
	"errors"
	authdi "kiraform/src/applications/dependencies/auths"
	authschema "kiraform/src/interfaces/rest/schemas/auths"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	DB           *gorm.DB
	Validator    *validator.Validate
	Dependencies authdi.Dependencies
}

func NewHandler(DB *gorm.DB, validate *validator.Validate, dependencies authdi.Dependencies) *Handler {
	return &Handler{
		DB:           DB,
		Validator:    validate,
		Dependencies: dependencies,
	}
}

func NewHTTP(g *echo.Group, DB *gorm.DB) {
	validate := validator.New()
	h := NewHandler(DB, validate, *authdi.NewDependencies(DB))

	// define endpoints
	g.POST("/login", h.Login)
	g.POST("/register", h.Register)
}

// @Security BearerAuth
// @Summary      Login
// @Description  User login
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        loginPayload  body      authschema.LoginPayload   true  "Login credentials"
// @Success      200  {string} string "Login success"
// @Failure      400  {string} string "Login failure"
// @Router       /api/login [post]
func (h *Handler) Login(c echo.Context) error {
	var body authschema.LoginPayload

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	signedToken, err := h.Dependencies.UC.Login(body)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		} else {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	return c.JSON(http.StatusOK, map[string]any{
		"access_token":  signedToken,
		"refresh_token": signedToken,
	})
}

// @Security BearerAuth
// @Summary      Registration
// @Description  You can regist new user here
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        registerPayload  body      authschema.RegisterPayload   true  "Register credentials"
// @Success      200  {string} string "Registration success"
// @Failure      400  {string} string "Registration failure"
// @Router       /api/register [post]
func (h *Handler) Register(c echo.Context) error {
	var body authschema.RegisterPayload

	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	msg, err := h.Dependencies.UC.Register(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]any{
		"message": msg,
	})
}
