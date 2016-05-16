package helpers

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/ziutek/telnet"
)

func SendCommand(c echo.Context) (error, error) {
	var connection *telnet.Conn

	request := Request{}
	c.Bind(&request)

	request, err := CheckRequest(request)
	if err != nil {
		return nil, c.JSON(http.StatusBadRequest, "Error: "+err.Error())
	}

	connection, err = telnet.Dial("tcp", request.Address+":"+request.Port)
	if err != nil {
		return nil, c.JSON(http.StatusBadRequest, "Error contacting host: "+err.Error())
	}

	defer connection.Close()
	connection.SetUnixWriteMode(true) // Convert any '\n' (LF) to '\r\n' (CR LF)

	_, err = connection.Write([]byte(request.Command + "\n"))
	if err != nil {
		return nil, c.JSON(http.StatusBadRequest, "Error: "+err.Error())
	}

	connection.SkipUntil(request.Prompt)
	output, err := connection.ReadUntil(request.Prompt)
	if err != nil {
		return nil, c.JSON(http.StatusBadRequest, "Error: "+err.Error())
	}

	output = output[:len(output)-len(request.Prompt)] // Trim the prompt off the output

	switch request.Command {
	case "xget ~.LocalInfo.vtpage":
		return c.JSON(http.StatusOK, "Use the /project endpoint for this command"), nil
	case "iptable":
		iptable, err := GetIPTable(output)
		if err != nil {
			return nil, c.JSON(http.StatusBadRequest, "Error: "+err.Error())
		}

		return c.JSON(http.StatusOK, iptable), nil
	default:
		response := strings.Replace(string(output), "\r\n", "", -1) // Remove line returns that are added automatically
		response = strings.TrimSpace(response)                      // Kill any leading or trailing spaces

		return c.JSON(http.StatusOK, response), nil
	}
}

func SendCommandWithConfirm(c echo.Context) (error, error) {
	var connection *telnet.Conn

	request := Request{}
	c.Bind(&request)

	request, err := CheckRequest(request)
	if err != nil {
		return nil, c.JSON(http.StatusBadRequest, "Error: "+err.Error())
	}

	connection, err = telnet.Dial("tcp", request.Address+":"+request.Port)
	if err != nil {
		return nil, c.JSON(http.StatusInternalServerError, "Error: "+err.Error())
	}

	connection.SetUnixWriteMode(true) // Convert any '\n' (LF) to '\r\n' (CR LF)

	_, err = connection.Write([]byte(request.Command + "\n"))
	if err != nil {
		return nil, c.JSON(http.StatusInternalServerError, "Error: "+err.Error())
	}

	time.Sleep(1000 * time.Millisecond) // Wait for the prompt to appear

	connection.Write([]byte("y")) // Send the confirmation

	return c.JSON(http.StatusOK, "Success"), nil
}
