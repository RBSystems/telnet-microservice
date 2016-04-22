package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/ziutek/telnet"
)

func getPrompt(req request, conn *telnet.Conn) (string, error) {

	_, err := conn.Write([]byte("\n\n"))

	if err != nil {
		return "", err
	}

	//Dynamically get the prompt
	conn.SkipUntil(">")
	promptBytes, err := conn.ReadUntil(">")

	if err != nil {
		return "", err
	}
	regex := "\\S.*?>"

	re := regexp.MustCompile(regex)

	prompt := string(re.Find(promptBytes))

	return prompt, nil
}

func sendCommand(c web.C, w http.ResponseWriter, r *http.Request) {
	bits, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not read request body: %s\n", err.Error())
		return
	}

	var req request
	err = json.Unmarshal(bits, &req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error with the request body: %s", err.Error())
		return
	}

	var conn *telnet.Conn

	if req.Port == "" {
		conn, err = telnet.Dial("tcp", req.IPAddress+":41795")
	} else {
		conn, err = telnet.Dial("tcp", req.IPAddress+":"+req.Port)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error with contacting host %s", err.Error())
		return
	}

	defer conn.Close()
	conn.SetUnixWriteMode(true) // Convert any '\n' (LF) to '\r\n' (CR LF) This is apparently very important

	if req.Prompt == "" {
		p, err := getPrompt(req, conn)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error with contacting host %s", err.Error())
			return
		}

		req.Prompt = p
	}

	conn.SetReadDeadline(time.Now().Add(45 * time.Second))

	_, err = conn.Write([]byte(req.Command + "\n\n")) // Send a second newline so we get the prompt

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error with contacting host %s", err.Error())
		return
	}
	err = conn.SkipUntil(req.Prompt) // Skip to the first prompt delimiter

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error with contacting host %s", err.Error())
		return
	}

	response, err := conn.ReadUntil(req.Prompt) // Read until the second prompt delimiter (provided by sending two commands in sendCommand)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	response = response[:len(response)-len(req.Prompt)] // Ghetto trim the prompt off the response
	response = response[len(req.Command):]

	switch req.Command {
	case "iptable":
		ipTable, err := getIPTable(string(response))

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error: %s", err.Error())
			return
		}

		bits, err := json.Marshal(ipTable)
		fmt.Fprintf(w, "%s", bits)
		return
	default:
		fmt.Fprintf(w, "%s", strings.TrimSpace(string(response)))
		return

	}

}

func sendCommandConfirm(c web.C, w http.ResponseWriter, r *http.Request) {
	bits, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Could not read request body: %s\n", err.Error())
		return
	}

	var req request
	err = json.Unmarshal(bits, &req)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Error with the request body: %s", err.Error())
		return
	}

	if len(req.Port) < 1 {
		req.Port = "41795"
	}
	err = sendCommandWithConfirm(req.Command, req.IPAddress, req.Port)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error: %s", err.Error())
		return
	}

	fmt.Fprintf(w, "Success!")
}

func sendCommandWithConfirm(command string, ipAddress string, port string) error {
	var conn *telnet.Conn

	conn, err := telnet.Dial("tcp", ipAddress+":"+port)
	if err != nil {
		return err
	}

	conn.SetUnixWriteMode(true) // Convert any '\n' (LF) to '\r\n' (CR LF) This is apparently very important

	_, err = conn.Write([]byte(command + "\n"))

	if err != nil {
		return err
	}

	time.Sleep(1000 * time.Millisecond) //Wait for the prompt to appear

	conn.Write([]byte("y")) //send the yes confirmation

	return nil
}

func getIPTable(response string) (IPTable, error) {
	// Parse the response to build the IP Table.

	lines := strings.Split(response, "\n") //get each line. The first is the header, each subsequent line is an entry in the table

	fmt.Printf("ResponseString: %s\n", response)

	fmt.Printf("\nLines: %v\n Length: %v\n", lines, len(lines))

	var toReturn IPTable

	for i := 1; i < len(lines); i++ {
		entries := strings.Fields(lines[i]) //psplit on whitespace

		var toAdd IPEntry

		if len(entries) == 0 {
			continue
		}

		fmt.Printf("Entries: %+v\n", entries)
		fmt.Printf("Length Entries: %v \n", len(entries))

		switch len(entries) {
		case 5: //There are 5 entries, assume DevID isn't there.
			fmt.Printf("Adding Entry: %v\n", entries)
			toAdd = IPEntry{CipID: entries[0], Type: entries[1], Status: entries[2], Port: entries[3], IPAddressSitename: entries[4]}
		case 6: //There are 6 entries, DevID is there.
			toAdd = IPEntry{CipID: entries[0], Type: entries[1], Status: entries[2], DevID: entries[3], Port: entries[4], IPAddressSitename: entries[5]}
		default: //We don't recognize this IPtable
			return IPTable{}, errors.New("Unrecognized IP Table returned.\n")
		}

		toReturn.Entries = append(toReturn.Entries, toAdd)
		fmt.Printf("ToReturn: %+v\n", toReturn)
	}

	return toReturn, nil
}

func main() {
	goji.Post("/sendCommand", sendCommand)
	goji.Post("/sendCommand/", sendCommand)
	goji.Post("/sendCommandConfirm", sendCommandConfirm)
	goji.Post("/sendCommandConfirm/", sendCommandConfirm)
	goji.Serve()
}
