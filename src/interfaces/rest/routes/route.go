package routes

import (
	"kiraform/src/interfaces/rest/middlewares"
	authroute "kiraform/src/interfaces/rest/routes/auths"
	masterroute "kiraform/src/interfaces/rest/routes/masters"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Routes(e *echo.Echo, DB *gorm.DB) {
	// unauthorized endpoint
	publicApi := e.Group("/api")
	authroute.NewAuthHTTP(publicApi, DB)

	// authorized endpoint
	privateApi := e.Group("/api")

	// regist middlewares
	privateApi.Use(middlewares.VerifyToken)

	// regist all secure routes
	masterroute.NewFormHTTP(privateApi, DB)
	masterroute.NewWorkspaceHTTP(privateApi, DB)
	masterroute.NewCampaignHTTP(privateApi, DB)
}
