package helpers

import (
	"errors"
	"fmt"
	"strings"
)

func GetIPTable(response string) (IPTable, error) {
	// Parse the response to build the IP Table

	if strings.Contains(response, "IP Table:") { // Remove the first pieces so all we have is the actual table
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

		toAdd := IPTableEntry{}

		if len(entries) == 0 { //skip empty lines
			continue
		}

		fmt.Printf("Entries: %+v\n", entries)
		fmt.Printf("Length Entries: %v \n", len(entries))

		switch len(entries) {
		case 5: //There are 5 entries, assume DevID isn't there.
			fmt.Printf("Adding Entry: %v\n", entries)
			toAdd = IPTableEntry{CipID: entries[0], Type: entries[1], Status: entries[2], Port: entries[3], AddressSitename: entries[4]}
		case 6: //There are 6 entries, DevID is there.
			toAdd = IPTableEntry{CipID: entries[0], Type: entries[1], Status: entries[2], DevID: entries[3], Port: entries[4], AddressSitename: entries[5]}
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
