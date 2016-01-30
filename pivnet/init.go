package pivnet

import (
	"os"

	"github.com/dingodb/pivotal-opsmgr-download-mgr/marketplaces"
)

// PivNet is configuration for a target Pivotal Network access account
type PivNet struct {
	apiToken     string
	productTiles marketplaces.ProductTiles
}

// NewPivNet creates a new PivNet struct
func NewPivNet() *PivNet {
	return &PivNet{
		apiToken: os.Getenv("PIVOTAL_NETWORK_TOKEN"),
	}
}

// Name returns the common name for PivNet
func (pivnetAPI *PivNet) Name() string {
	return "Pivotal Network"
}

// Slug returns the internal slug name for PivNet
func (pivnetAPI *PivNet) Slug() string {
	return "pivnet"
}

// ProductTiles returns the fetched Product Tiles for PivNet
func (pivnetAPI *PivNet) ProductTiles() marketplaces.ProductTiles {
	return pivnetAPI.productTiles
}
