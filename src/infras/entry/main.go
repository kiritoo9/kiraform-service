package main

import (
	"fmt"
	"kiraform/src/infras/configs"
	"kiraform/src/interfaces/rest/routes"
	"strings"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "kiraform/docs"
)

func main() {
	// define core entities
	CONFIG := configs.Environment()
	DB := configs.Connection(CONFIG)
	e := echo.New()

	// load swagger only for development environment
	if strings.ToLower(CONFIG.APP_ENV) == "dev" {
		// @securityDefinitions.apikey BearerAuth
		// @in header
		// @name Authorization
		// @description Type "Bearer " followed by a space and JWT token.
		e.GET("/docs/*", echoSwagger.WrapHandler)
	}

	// calling main route
	routes.Routes(e, DB)

	// run applications
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", CONFIG.APP_PORT)))
}
