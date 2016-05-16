package helpers

import (
	"fmt"
	"strings"
	"time"

	"github.com/ziutek/telnet"
)

func GetProjectInfo(req Request, conn *telnet.Conn) (string, error) {
	fmt.Printf("%s Getting project info...\n", req.Address)

	defer conn.Close()

	if req.Prompt == "" {
		prompt, _ := GetPrompt(req.Address)
		req.Prompt = prompt
	}

	conn.Write([]byte("udir \\romdisk\\user\\display\\\n\n"))
	conn.SkipUntil(req.Prompt) // Skip to the first prompt delimiter

	resp1, err := conn.ReadUntil(req.Prompt) // Read until the second prompt delimiter (provided by sending two commands in sendCommand)
	if err != nil {
		return "", err
	}

	fmt.Printf("%s %s\n", req.Address, resp1)

	if !strings.Contains(string(resp1), ".vtpage") {
		return "File ~.LocalInfo.vtpage does not exist.\n", nil
	}

	conn.Write([]byte("cd \\romdisk\\user\\display\\\n"))
	conn.SkipUntil(req.Prompt) // Skip to the first prompt delimiter

	resp, err := conn.ReadUntil(req.Prompt) // Read until the second prompt delimiter (provided by sending two commands in sendCommand)
	if err != nil {
		return "", err
	}

	fmt.Printf("%s\n", resp)
	conn.SetReadDeadline(time.Now().Add(2 * time.Minute))
	conn.Write([]byte(req.Command + "\n\n"))

	conn.SkipUntil("[BEGIN_INFO]", "ERROR")
	fmt.Printf("%s skipped\n", req.Address)
	resp, err = conn.ReadUntil("[END_INFO]", "Panel", "not")
	if err != nil {
		return "", err
	}

	fmt.Printf("%s Response: %s\n", req.Address, string(resp))

	conn.Close() // Actively close the xmodem connection

	return string(resp), nil
}
