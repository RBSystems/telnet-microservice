package controllers

import (
	"net/http"

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

func CommandInfo(c echo.Context) error {
	return c.JSON(http.StatusOK, "Send a POST request to the /command endpoint with a body including at least Address and Command tokens to send a telnet command to the specified machine")
}

func CommandWithConfirm(c echo.Context) error {
	response, err := helpers.SendCommandWithConfirm(c)
	if err != nil {
		return err
	}

	return response
}

func CommandWithConfirmInfo(c echo.Context) error {
	return c.JSON(http.StatusOK, "Send a POST request to the /confirmed endpoint with a body including at least Address and Command tokens to send a confirmed (command followed by 'y') telnet command to the specified machine")
}
