package controllers

import (
	"github.com/byuoitav/telnet-microservice/helpers"
	"github.com/labstack/echo"
)

func Command(c echo.Context) error {
	response, err := helpers.SendCommand(c)
	if err != nil {
		return err
	}

	return response
}
