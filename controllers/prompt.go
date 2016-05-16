package controllers

import (
	"fmt"
	"net/http"

	"github.com/byuoitav/ftp-microservice/helpers"
	"github.com/labstack/echo"
	"github.com/ziutek/telnet"
)

func getPromptHandler(c echo.Context) error {
	request := helpers.Request{}
	c.Bind(request)

	if len(request.Port) < 1 {
		request.Port = "41795"
	}

	connection := *telnet.Conn

	if request.Port == "" {
		connection, err := telnet.Dial("tcp", request.IPAddress+":41795")
	} else {
		connection, err := telnet.Dial("tcp", request.IPAddress+":"+request.Port)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error with contacting host: "+err.Error())
	}

	defer connection.Close()
	connection.SetUnixWriteMode(true) // Convert any '\n' (LF) to '\r\n' (CR LF) This is apparently very important

	prompt, err := helpers.GetPrompt(request, connection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error with contacting host %s", err.Error())
		return
	}

	request.Prompt = prompt

	return c.JSON(http.StatusOK, prompt)
}
