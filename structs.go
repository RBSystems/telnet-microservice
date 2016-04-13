package main

//IPTable represents an IPTable returend from a crestron device
type IPTable struct {
	Entries []IPEntry
}

//IPEntry represents a single entry in the IPTable
type IPEntry struct {
	CipID             string `json:"CIP_ID"`
	Type              string
	Status            string
	DevID             string
	Port              string
	IPAddressSitename string
}
