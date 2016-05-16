package helpers

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/ziutek/telnet"
)

func SendCommand(c echo.Context) (error, error) {
	var connection *telnet.Conn

	req := Request{}
	c.Bind(&req)

	req, err := CheckRequest(req)
	if err != nil {
		return nil, err
	}

	connection, err = telnet.Dial("tcp", req.Address+":"+req.Port)
	if err != nil {
		return nil, errors.New("Error contacting host: " + err.Error())
	}

	defer connection.Close()
	connection.SetUnixWriteMode(true) // Convert any '\n' (LF) to '\r\n' (CR LF)

	_, err = connection.Write([]byte(req.Command + "\n"))
	if err != nil {
		return nil, err
	}

	switch req.Command {
	case "xget ~.LocalInfo.vtpage":
		resp, err := GetProjectInfo(req, connection)
		if err != nil {
			return nil, errors.New("Error contacting host: " + err.Error())
		}

		return c.JSON(http.StatusOK, strings.TrimSpace(string(resp))), nil
	case "iptable":
		ipTable, err := GetIPTable(req.Prompt)
		if err != nil {
			return nil, errors.New("Error: " + err.Error())
		}

		return c.JSON(http.StatusOK, ipTable), nil
	default:
		connection.SkipUntil(req.Prompt)
		output, err := connection.ReadUntil(req.Prompt)
		if err != nil {
			return nil, err
		}

		output = output[:len(output)-len(req.Prompt)]               // Trim the prompt off the output
		response := strings.Replace(string(output), "\r\n", "", -1) // Remove line returns that are added automatically
		response = strings.Trim(response, " ")                      // Kill any trailing spaces

		return c.JSON(http.StatusOK, response), nil
	}
}
