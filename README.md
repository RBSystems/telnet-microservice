# Telnet Microservice

The service is meant to facilitate telnet connections and commands by proxy. Specifically designed for communication with Crestron systems.

Active endpoints are

* `POST:/sendCommand`
* `POST:/sendCommandConfirm`

both with a JSON payload in the form of
```
{
	"Address":	"Address of target",
	"Port": 		"Optional parameter of port to connect. Defaults to 41795"
	"Command":  	"The command to send the target"
	"Prompt": 	  "Optional string to use as a delimiter when parsing the response."
}
```

If the `Prompt` field is not included the service will attempt to determine the prompt. If unsuccessful will return an error. 

If the `Command` field contains a recognized command, the service will parse and return a JSON response with the parsed results. Otherwise it will return the raw
response.

## sendCommand

Send command is used to send a generic command and return the response. Any telnet command can be sent and the response will be sent back. If the command is one of the recognized commands the output will be parsed and sent back in a json payload, otherwise a raw response will be sent.

## sendCommandConfirm

Send Command with confirm is meant for commands that require a confirmation after command execution. (`initialize` for instance) This endpoint will execute the command, wait for 1 second, and then send the 'y' character. Response is not recorded, just a `Success!` if no error was received.  

## Recognized commands

This is a list of the currently recognized commands with the JSON response that can be expected.

#### iptable
```
{"IPTable": [
		{
			"CIP_ID": "string",
			"Type": "string",
			"Status": "string",
			"DevID": "string",
			"Port": "string",
			"AddressSitename": "string"
		}
	]}
```

It should be noted that DevID will usually be blank.
