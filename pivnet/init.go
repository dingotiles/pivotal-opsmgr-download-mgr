package pivnet

import "os"

// PivNet is configuration for a target Pivotal Network access account
type PivNet struct {
	APIToken string
}

// NewPivNet creates a new PivNet struct
func NewPivNet() PivNet {
	return PivNet{
		APIToken: os.Getenv("PIVOTAL_NETWORK_TOKEN"),
	}
}
