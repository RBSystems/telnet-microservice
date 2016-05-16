package helpers

import (
	"errors"
	"regexp"

	"github.com/ziutek/telnet"
)

func GetPrompt(address string) (string, error) {
	connection, err := telnet.Dial("tcp", address+":23")
	if err != nil {
		return "", errors.New("Error contacting host: " + err.Error())
	}

	defer connection.Close()
	connection.SetUnixWriteMode(true) // Convert any '\n' (LF) to '\r\n' (CR LF)

	_, err = connection.Write([]byte("\n\n")) // Send two LF so the return contains the prompt
	if err != nil {
		return "", err
	}

	// Dynamically get the prompt
	connection.SkipUntil(">")
	promptBytes, err := connection.ReadUntil(">")
	if err != nil {
		return "", err
	}

	regex := "\\S.*?>"
	re := regexp.MustCompile(regex)

	prompt := string(re.Find(promptBytes))

	return prompt, nil
}
