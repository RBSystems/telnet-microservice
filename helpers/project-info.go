package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/ziutek/telnet"
)

func GetProjectInfo(request Request, connection *telnet.Conn) (string, error) {
	defer connection.Close()

	connection.Write([]byte("udir \\romdisk\\user\\display\\\n\n"))
	connection.SkipUntil(request.Prompt) // Skip to the first prompt delimiter

	resp1, err := connection.ReadUntil(request.Prompt) // Read until the second prompt delimiter (provided by sending two commands in sendCommand)
	if err != nil {
		return "", err
	}

	fmt.Printf("STUFF %s %s\n", request.Address, resp1)

	if !strings.Contains(string(resp1), ".vtpage") {
		return "File ~.LocalInfo.vtpage does not exist", nil
	}

	connection.Write([]byte("cd \\romdisk\\user\\display\\\n"))
	connection.SkipUntil(request.Prompt) // Skip to the first prompt delimiter

	resp, err := connection.ReadUntil(request.Prompt) // Read until the second prompt delimiter (provided by sending two commands in sendCommand)
	if err != nil {
		return "", err
	}

	fmt.Printf("%s\n", resp)
	connection.SetReadDeadline(time.Now().Add(2 * time.Minute))
	connection.Write([]byte(request.Command + "\n\n"))

	connection.SkipUntil("[BEGIN_INFO]", "ERROR")
	fmt.Printf("%s skipped\n", request.Address)
	resp, err = connection.ReadUntil("[END_INFO]", "Panel", "not")
	if err != nil {
		return "", err
	}

	fmt.Printf("%s Response: %s\n", request.Address, string(resp))

	// connection.Close() // Actively close the xmodem connection

	return strings.TrimSpace(string(resp)), nil
}
