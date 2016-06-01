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
		connection.Write([]byte("cd \\romdisk\\user\\display\\\n\n"))
		connection.SkipUntil(request.Prompt)

		connection.Write([]byte("xget ~.LocalInfo.vtpage\n\n")) // Send a second newline so we get the prompt

		connection.SkipUntil("[TPS]", "ERROR")
		response, err := connection.ReadUntil("[END_INFO]", "Panel", "not")
		if err != nil {
			return Project{}, err
		}

		info := strings.Split(strings.TrimSpace(string(response)), "\n")
		project := Project{}

		for i := range info {
			if strings.Contains(info[i], "VTZ=") {
				project.Project = strings.Replace(strings.Replace(info[i], "VTZ=", "", -1), "\r", "", -1)
			} else if strings.Contains(info[i], "Date=") {
				project.ProjectDate = strings.Replace(strings.Replace(info[i], "Date=", "", -1), "\r", "", -1)
			}
		}

		connection.Close()

		return project, nil
	}

	return Project{}, errors.New("File ~.LocalInfo.vtpage does not exist")
}
