package main

import (
	"fmt"

	"github.com/byuoitav/hateoas"
	"github.com/byuoitav/telnet-microservice/controllers"
	"github.com/jessemillar/health"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"
)

func main() {
	err := hateoas.Load("https://raw.githubusercontent.com/byuoitav/telnet-microservice/master/swagger.yml")
	if err != nil {
		fmt.Printf("Could not load swagger.yaml file. Error: %s", err.Error())
		panic(err)
	}

	port := ":8001"
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())

	e.Get("/", hateoas.RootResponse)
	e.Get("/health", health.Check)

	e.Get("/prompt/:address", controllers.GetPrompt)
	e.Get("/project/:address", controllers.GetProjectInfo)

	e.Post("/command", controllers.Command)
	e.Post("/command/confirm", controllers.CommandWithConfirm)

	fmt.Printf("The Telnet Microservice is listening on %s\n", port)
	e.Run(fasthttp.New(port))
}
