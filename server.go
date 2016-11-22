package main

import (
	"log"
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
		log.Fatalln("Could not load Swagger file. Error: " + err.Error())
	}

	port := ":8001"
	router := echo.New()
	router.Pre(middleware.RemoveTrailingSlash())
	router.Use(middleware.CORS())

	// Use the `secure` routing group to require authentication
	secure := router.Group("", echo.WrapMiddleware(wso2jwt.ValidateJWT))

	router.GET("/", echo.WrapHandler(http.HandlerFunc(hateoas.RootResponse)))
	router.GET("/health", echo.WrapHandler(http.HandlerFunc(health.Check)))

	secure.GET("/prompt/:address", controllers.GetPrompt)
	secure.GET("/project/:address", controllers.GetProjectInfo)

	secure.GET("/command", controllers.CommandInfo)
	secure.POST("/command", controllers.Command)
	secure.GET("/confirmed", controllers.CommandWithConfirmInfo)
	secure.POST("/confirmed", controllers.CommandWithConfirm)

	server := http.Server{
		Addr:           port,
		MaxHeaderBytes: 1024 * 10,
	}

	router.StartServer(&server)
}
