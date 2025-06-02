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
	masterroute.NewFormEntryHTTP(publicApi, DB)

	// re-define /api for authorized endpoint
	// then regist middleware
	privateApi := e.Group("/api")
	privateApi.Use(middlewares.VerifyToken)

	// regist all secure routes
	masterroute.NewFormHTTP(privateApi, DB)
	masterroute.NewWorkspaceHTTP(privateApi, DB)
	masterroute.NewCampaignHTTP(privateApi, DB)
}
