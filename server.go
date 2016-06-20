package main

import (
	"fmt"
	"log"

	"github.com/byuoitav/hateoas"
	"github.com/byuoitav/telnet-microservice/controllers"
	"github.com/byuoitav/wso2jwt"
	"github.com/jessemillar/health"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	"github.com/labstack/echo/middleware"
)

func main() {
	err := hateoas.Load("https://raw.githubusercontent.com/byuoitav/telnet-microservice/master/swagger.json")
	if err != nil {
		fmt.Printf("Could not load swagger.json file. Error: %s", err.Error())
		panic(err)
	}

	port := ":8001"
	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())

	router.Get("/", hateoas.RootResponse)
	router.Get("/health", health.Check)

	router.Get("/prompt/:address", controllers.GetPrompt, wso2jwt.ValidateJWT())
	router.Get("/project/:address", controllers.GetProjectInfo, wso2jwt.ValidateJWT())

	router.Get("/command", controllers.CommandInfo, wso2jwt.ValidateJWT())
	router.Post("/command", controllers.Command, wso2jwt.ValidateJWT())
	router.Get("/confirmed", controllers.CommandWithConfirmInfo, wso2jwt.ValidateJWT())
	router.Post("/confirmed", controllers.CommandWithConfirm, wso2jwt.ValidateJWT())

	log.Println("The Telnet Microservice is listening on " + port)
	server := fasthttp.New(port)
	server.ReadBufferSize = 1024 * 10 // Needed to interface properly with WSO2
	router.Run(server)
}
