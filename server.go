package main

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/hateoas"
	"github.com/byuoitav/telnet-microservice/controllers"
	"github.com/byuoitav/wso2jwt"
	"github.com/jessemillar/health"
	"github.com/labstack/echo"
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
	router.Use(middleware.CORS())

	router.GET("/", hateoas.RootResponse)
	router.GET("/health", health.Check)

	router.GET("/prompt/:address", controllers.GETPrompt, wso2jwt.ValidateJWT())
	router.GET("/project/:address", controllers.GETProjectInfo, wso2jwt.ValidateJWT())

	router.GET("/command", controllers.CommandInfo, wso2jwt.ValidateJWT())
	router.POST("/command", controllers.Command, wso2jwt.ValidateJWT())
	router.GET("/confirmed", controllers.CommandWithConfirmInfo, wso2jwt.ValidateJWT())
	router.POST("/confirmed", controllers.CommandWithConfirm, wso2jwt.ValidateJWT())

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
