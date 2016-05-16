package helpers

//IPTable represents an IPTable returend from a crestron device
type IPTable struct {
	Entries []IPTableEntry
}

//IPTableEntry represents a single entry in the IPTable
type IPTableEntry struct {
	CipID             string `json:"CIP_ID"`
	Type              string
	Status            string
	DevID             string
	Port              string
	IPAddressSitename string
}

type Request struct {
	IPAddress string
	Port      string
	Command   string
	Prompt    string
}
