package controllers

import (
	"net/http"

	"github.com/byuoitav/telnet-microservice/helpers"
	"github.com/labstack/echo"
)

func Command(c echo.Context) error {
	response, err := helpers.SendCommand(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, response)
}
