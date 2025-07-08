package main

import (
	"fmt"
	"kiraform/src/infras/configs"
	"kiraform/src/interfaces/rest/routes"
	"log"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "kiraform/docs"
)

func main() {
	// define core entities
	CONFIG := configs.Environment()
	DB := configs.Connection(CONFIG)
	e := echo.New()

	// cors handler
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
	}))

	// serve static file
	cdnPath, err := filepath.Abs("./cdn")
	if err != nil {
		log.Fatal("Failed to resolve CDN path:", err)
	}
	e.Static("/cdn", cdnPath)

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
