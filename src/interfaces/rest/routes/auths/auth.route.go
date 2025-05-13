package authroute

import (
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

	g.POST("/login", h.Login)
}

func (h *Handler) Login(c echo.Context) error {
	var body authschema.LoginPayload

	// bind body
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid JSON")
	}

	// validate struct
	if err := h.Validator.Struct(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// call usecase
	err := h.Dependencies.UC.Login(body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, body)
}
