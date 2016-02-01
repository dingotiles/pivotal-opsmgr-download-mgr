package starkandwayne

import (
	"github.com/dingodb/pivotal-opsmgr-download-mgr/marketplaces"
)

// StarkAndWayne is configuration for a target Pivotal Network access account
type StarkAndWayne struct {
	// apiToken     string
	productTiles marketplaces.ProductTiles
}

// NewStarkAndWayne creates a new StarkAndWayne struct
func NewStarkAndWayne() *StarkAndWayne {
	return &StarkAndWayne{
	// apiToken: os.Getenv("PIVOTAL_NETWORK_TOKEN"),
	}
}

// Name returns the common name for StarkAndWayne
func (starkAndWayneAPI *StarkAndWayne) Name() string {
	return "Stark & Wayne Marketplace"
}

// Slug returns the internal slug name for StarkAndWayne
func (starkAndWayneAPI *StarkAndWayne) Slug() string {
	return "starkandwayne"
}

// ProductTiles returns the fetched Product Tiles for StarkAndWayne
func (starkAndWayneAPI *StarkAndWayne) ProductTiles() marketplaces.ProductTiles {
	return starkAndWayneAPI.productTiles
}
