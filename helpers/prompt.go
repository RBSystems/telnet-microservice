package helpers

import (
	"regexp"

	"github.com/ziutek/telnet"
)

func GetPrompt(request Request, connection *telnet.Conn) (string, error) {
	_, err := connection.Write([]byte("\n\n"))

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
