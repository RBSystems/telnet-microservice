package helpers

import "errors"

func CheckRequest(request Request) (Request, error) {
	if len(request.Address) < 1 || len(request.Command) < 1 {
		return Request{}, errors.New("Body must contain at least Address and Command tokens")
	}

	if len(request.Port) < 1 {
		request.Port = "23"
	}

	if request.Prompt == "" {
		prompt, err := GetPrompt(request.Address)
		if err != nil {
			return Request{}, errors.New("Error contacting host: " + err.Error())
		}

		request.Prompt = prompt
	}

	return request, nil
}
