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

	//Cheap cop-out way to deal with getting the version of the touchpanels. Split out xmodem into own endpoint?
	//TODO: Figure out a better way to handle this.
	if strings.EqualFold(req.Command, "xget ~.LocalInfo.vtpage") {
		resp, err := getProjectInfo(req, conn)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error with contacting host %s", err.Error())
			return
		}

		fmt.Fprintf(w, "%s", strings.TrimSpace(string(resp)))
		return
	}

	conn.SetReadDeadline(time.Now().Add(45 * time.Second))

	if req.Prompt == "" {
		p, err := getPrompt(req, conn)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error with contacting host %s", err.Error())
			return
		}

		req.Prompt = p
	}

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

func getPromptHandler(c web.C, w http.ResponseWriter, r *http.Request) {
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

	p, err := getPrompt(req, conn)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error with contacting host %s", err.Error())
		return
	}

	req.Prompt = p

	b, _ := json.Marshal(req)

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, "%s", string(b))
}

func getProjectInfo(req request, conn *telnet.Conn) (string, error) {
	fmt.Printf("%s Getting project info...\n", req.IPAddress)

	defer conn.Close()

	if req.Prompt == "" {
		prompt, _ := getPrompt(req, conn)
		req.Prompt = prompt
	}
	conn.Write([]byte("udir \\romdisk\\user\\display\\\n\n"))
	conn.SkipUntil(req.Prompt) // Skip to the first prompt delimiter

	resp1, err := conn.ReadUntil(req.Prompt) // Read until the second prompt delimiter (provided by sending two commands in sendCommand)

	if err != nil {
		return "", err
	}

	fmt.Printf("%s %s\n", req.IPAddress, resp1)

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
	fmt.Printf("%s skipped\n", req.IPAddress)
	resp, err = conn.ReadUntil("[END_INFO]", "Panel", "not")

	if err != nil {
		return "", err
	}

	fmt.Printf("%s Response: %s\n", req.IPAddress, string(resp))

	conn.Close() //actively close the xmodem connection.

	return string(resp), nil
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

	if strings.Contains(response, "IP Table:") { //remove the first pieces so all we have
		//is the actual table.
		response = strings.Split(response, "IP Table:")[1]
	}
	response = strings.TrimSpace(response)

	lines := strings.Split(response, "\n") //get each line. The first is the header, each subsequent line is an entry in the table

	fmt.Printf("ResponseString:\n %s\n", response)

	fmt.Printf("Length: %v\n", len(lines))

	var toReturn IPTable

	for i := 1; i < len(lines); i++ { //start at one so we skip the headers.
		entries := strings.Fields(lines[i]) //psplit on whitespace

		fmt.Printf("Fields: %+v\n", entries)

		var toAdd IPEntry

		if len(entries) == 0 { //skip empty lines
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
		default: //We don't recognize this IPtable. If we already have entries, just skip it, otherwise throw a fit.
			if len(toReturn.Entries) == 0 {
				return IPTable{}, errors.New("Unrecognized IP Table returned.\n")
			}
			continue
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
	goji.Post("/sendCommand/getPrompt/", getPromptHandler)
	goji.Post("/sendCommand/getPrompt", getPromptHandler)
	goji.Serve()
}
