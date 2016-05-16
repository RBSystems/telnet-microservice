package helpers

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/ziutek/telnet"
)

func SendCommand(c echo.Context) (string, error) {
	var conn *telnet.Conn

	req := Request{}
	c.Bind(&req)

	fmt.Printf("%+v\n", req)

	if len(req.Port) < 1 {
		req.Port = "23"
	}

	conn, err := telnet.Dial("tcp", req.Address+":"+req.Port)
	if err != nil {
		return "", errors.New("Error contacting host: " + err.Error())
	}

	defer conn.Close()
	conn.SetUnixWriteMode(true) // Convert any '\n' (LF) to '\r\n' (CR LF)

	// Cop-out way to deal with getting the version of the touchpanels--split out xmodem into its own endpoint?
	// TODO: Figure out a better way to handle this
	if strings.EqualFold(req.Command, "xget ~.LocalInfo.vtpage") {
		resp, err := GetProjectInfo(req, conn)
		if err != nil {
			return "", errors.New("Error contacting host: " + err.Error())
		}

		return strings.TrimSpace(string(resp)), nil
	}

	conn.SetReadDeadline(time.Now().Add(45 * time.Second))

	if req.Prompt == "" {
		p, err := GetPrompt(req.Address)
		if err != nil {
			return "", errors.New("Error contacting host: " + err.Error())
		}

		req.Prompt = p
	}

	_, err = conn.Write([]byte(req.Command + "\n\n")) // Send a second newline so we get the prompt
	if err != nil {
		return "", errors.New("Error contacting host: " + err.Error())

	}

	err = conn.SkipUntil(req.Prompt) // Skip to the first prompt delimiter
	if err != nil {
		return "", errors.New("Error contacting host: " + err.Error())
	}

	response, err := conn.ReadUntil(req.Prompt) // Read until the second prompt delimiter (provided by sending two commands in sendCommand)
	if err != nil {
		return "", errors.New("Error contacting host: " + err.Error())
	}

	response = response[:len(response)-len(req.Prompt)] // Ghetto trim the prompt off the response
	response = response[len(req.Command):]

	switch req.Command {
	case "iptable":
		ipTable, err := GetIPTable(string(response))
		if err != nil {
			return "", errors.New("Error: " + err.Error())
		}

		fmt.Printf("STUFF: %+v\n", ipTable)

		return "", nil

		// return string(ipTable), nil
	default:
		return strings.TrimSpace(string(response)), nil
	}

	return "", nil
}
