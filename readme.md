# Telnet Microservice

The service is meant to facilitate telnet connections and commands by proxy. Specifically designed for communication with Crestron systems.

The only active endpoint is `POST:/sendCommand` with a JSON payload in the form of
```
{
	"IPAddress":	"IPAddress of target",
	"Port": 		"Optional parameter of port to connect. Defaults to 41795"
	"Command":  	"The command to send the target"
	"Prompt": 	  "The string to use as a delimiter when parsing the response."
}
```

If the `Command` field contains a recognized command, the service will parse and return a JSON response with the parsed results. Otherwise it will return the raw
response.

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
			"IPAddressSitename": "string"
		}
	]}
```

It should be noted that DevID will usually be blank.
