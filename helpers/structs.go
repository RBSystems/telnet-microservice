package helpers

// IPTable represents an IPTable returned from a crestron device
type IPTable struct {
	Entries []IPTableEntry `json:"entries"`
}

// IPTableEntry represents a single entry in the IPTable
type IPTableEntry struct {
	CipID   string `json:"cipID"`
	Type    string `json:"type"`
	Status  string `json:"status"`
	DevID   string `json:"devID"`
	Port    string `json:"port"`
	Address string `json:"address"`
}

// Request represents the body that's POSTed to the service
type Request struct {
	Address string `json:"address"`
	Port    string `json:"port"`
	Command string `json:"command"`
	Prompt  string `json:"prompt"`
}

// Project represents loaded project information for a touchpanel
type Project struct {
	Project     string `json:"project"`
	ProjectDate string `json:"projectDate"`
}
