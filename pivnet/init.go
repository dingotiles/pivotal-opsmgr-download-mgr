package pivnet

import "github.com/dingodb/pivotal-opsmgr-download-mgr/marketplaces"

// PivNet is configuration for a target Pivotal Network access account
type PivNet struct {
	apiToken     string
	apiEndpoint  string
	productTiles marketplaces.ProductTiles
	stemcells    marketplaces.ProductStemcells
}

// NewPivNet creates a new PivNet struct
func NewPivNet(apiToken string, apiEndpoint string) *PivNet {
	return &PivNet{
		apiToken:    apiToken,
		apiEndpoint: apiEndpoint,
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

// ProductStemcells returns the available stemcells from PivNet
func (pivnetAPI *PivNet) ProductStemcells() marketplaces.ProductStemcells {
	return pivnetAPI.stemcells
}
