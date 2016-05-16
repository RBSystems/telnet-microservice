package controllers

import (
	"net/http"

	"github.com/byuoitav/telnet-microservice/helpers"
	"github.com/labstack/echo"
)

func GetPrompt(c echo.Context) error {
	prompt, err := helpers.GetPrompt(c.Param("address"))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, prompt)
}
