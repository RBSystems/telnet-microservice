package helpers

import (
	"errors"
	"strings"

	"github.com/labstack/echo"
	"github.com/ziutek/telnet"
)

func GetProjectInfo(c echo.Context) (Project, error) {
	var connection *telnet.Conn

	request := Request{}
	c.Bind(&request)

	request.Address = c.Param("address")
	request.Port = "41795" // Creston uses this port instead of :23 for some reason

	prompt, err := GetPrompt(request.Address)
	if err != nil {
		return Project{}, errors.New("Error contacting host: " + err.Error())
	}

	request.Prompt = prompt

	connection, err = telnet.Dial("tcp", request.Address+":"+request.Port)
	if err != nil {
		return Project{}, err
	}

	defer connection.Close()
	connection.SetUnixWriteMode(true) // Convert any '\n' (LF) to '\r\n' (CR LF)

	connection.Write([]byte("udir \\romdisk\\user\\display\\\n\n"))
	connection.SkipUntil(request.Prompt)
	output, err := connection.ReadUntil(request.Prompt)
	if err != nil {
		return Project{}, err
	}

	if strings.Contains(string(output), ".vtpage") {
		response := strings.Split(string(output), "\n")

		for i := range response {
			if strings.Contains(response[i], "LocalInfo.vtpage") {
				date := strings.Split(strings.TrimSpace(response[i]), " ")

				project := Project{
					Date: date[2] + " " + date[3],
				}

				return project, nil
			}
		}
	}

	return Project{}, errors.New("File ~.LocalInfo.vtpage does not exist")
}
