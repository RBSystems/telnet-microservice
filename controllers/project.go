package controllers

import (
	"net/http"

	"github.com/byuoitav/telnet-microservice/helpers"
	"github.com/labstack/echo"
)

func GetProjectInfo(c echo.Context) error {
	output, err := helpers.GetProjectInfo(c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, output)
}
