package main

import (
	"fmt"
	"kiraform/src/infras/configs"
	"kiraform/src/interfaces/rest/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	// define core entities
	CONFIG := configs.Environment()
	DB := configs.Connection(CONFIG)
	e := echo.New()

	// calling main route
	routes.Routes(e, DB)

	// run applications
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", CONFIG.APP_PORT)))
}
