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
	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())

	router.Get("/", hateoas.RootResponse)
	router.Get("/health", health.Check)

	router.Get("/prompt/:address", controllers.GetPrompt)
	router.Get("/project/:address", controllers.GetProjectInfo)

	router.Get("/command", controllers.CommandInfo)
	router.Post("/command", controllers.Command)
	router.Get("/confirmed", controllers.CommandWithConfirmInfo)
	router.Post("/confirmed", controllers.CommandWithConfirm)

	fmt.Printf("The Telnet Microservice is listening on %s\n", port)
	server := fasthttp.New(port)
	server.ReadBufferSize = 1024 * 10 // Needed to interface properly with WSO2
	router.Run(server)
}
