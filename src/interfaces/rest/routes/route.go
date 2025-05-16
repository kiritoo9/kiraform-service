package routes

import (
	authroute "kiraform/src/interfaces/rest/routes/auths"
	masterroute "kiraform/src/interfaces/rest/routes/masters"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Routes(e *echo.Echo, DB *gorm.DB) {
	api := e.Group("/api")
	authroute.NewAuthHTTP(api, DB)
	masterroute.NewWorkspaceHTTP(api, DB)
}
