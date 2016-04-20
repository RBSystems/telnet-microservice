package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
	"github.com/ziutek/telnet"
)

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

	conn.SetUnixWriteMode(true) // Convert any '\n' (LF) to '\r\n' (CR LF) This is apparently very important

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error with contacting host %s", err.Error())
		return
	}

	value, err := conn.Write([]byte(req.Command + "\n")) // Send two commands so we get a second prompt to use as a delimiter

	time.Sleep(1000 * time.Millisecond)

	conn.Write([]byte("y"))
	fmt.Printf("%v\n", value)

	fmt.Printf("\n")
	//value, err = conn.Write([]byte("hostname"))
	//conn.
	fmt.Printf("%v\n", value)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error with contacting host %s", err.Error())
		return
	}
	//err = conn.SkipUntil(req.Prompt) // Skip to the first prompt delimiter

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Error with contacting host %s", err.Error())
		return
	}

	response, err := conn.ReadUntil(req.Prompt) // Read until the second prompt delimiter (provided by sending two commands in sendCommand)
	//esponse2, err := conn.ReadUntil(req.Prompt)

	conn.Close()

	response = response[:len(response)-10] // Ghetto trim the prompt off the response

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
		fmt.Fprintf(w, "%s", string(response))
		return

	}

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
	goji.Serve()
}
