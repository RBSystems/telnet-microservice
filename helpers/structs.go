package helpers

//IPTable represents an IPTable returend from a crestron device
type IPTable struct {
	Entries []IPTableEntry
}

//IPTableEntry represents a single entry in the IPTable
type IPTableEntry struct {
	CipID           string `json:"CIP_ID"`
	Type            string
	Status          string
	DevID           string
	Port            string
	AddressSitename string
}

type Request struct {
	Address string `json:"address"`
	Port    string `json:"port"`
	Command string `json:"command"`
	Prompt  string `json:"prompt"`
}
