package helpers

import "errors"

func CheckRequest(request Request) (Request, error) {
	if len(request.Address) < 1 || len(request.Command) < 1 || len(request.Prompt) < 1 {
		return Request{}, errors.New("Body must contain Address, Command, and Prompt tokens")
	}

	if len(request.Port) < 1 {
		request.Port = "23"
	}

	return request, nil
}
