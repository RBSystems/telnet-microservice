package helpers

import (
	"errors"
	"strings"
)

func GetIPTable(output []byte) (IPTable, error) {
	table := string(output) // Convert the []byte to a string

	if strings.Contains(table, "IP Table:") { // Remove the first pieces so all we have is the actual table
		table = strings.Split(table, "IP Table:")[1]
	}

	table = strings.TrimSpace(table)
	lines := strings.Split(table, "\n")
	response := IPTable{}

	for i := 1; i < len(lines); i++ { // Start at 1 so we skip the headers
		entries := strings.Fields(lines[i]) // Split on whitespace
		toAdd := IPTableEntry{}

		if len(entries) == 0 { // Skip empty lines
			continue
		}

		switch len(entries) {
		case 5: // If there are 5 entries, assume DevID isn't there
			toAdd = IPTableEntry{CipID: entries[0], Type: entries[1], Status: entries[2], Port: entries[3], Address: entries[4]}
		case 6: // If there are 6 entries, DevID is there
			toAdd = IPTableEntry{CipID: entries[0], Type: entries[1], Status: entries[2], DevID: entries[3], Port: entries[4], Address: entries[5]}
		default: // We don't recognize this IP table--if we already have entries, just skip it, otherwise throw a fit
			if len(response.Entries) == 0 {
				return IPTable{}, errors.New("Unrecognized IP Table returned.\n")
			}

			continue
		}

		response.Entries = append(response.Entries, toAdd)
	}

	return response, nil
}
