package pivnet

import "github.com/dingodb/pivotal-opsmgr-download-mgr/marketplaces"

// LookupProductTile tries to match an Opsmgr Product name with a PivNet product/release/.pivotal tile
func (pivnetAPI *PivNet) LookupProductTile(opsMgrProductName string) (tile *marketplaces.ProductTile) {
	for _, product := range pivnetAPI.productTiles {
		if product.TileName == opsMgrProductName {
			return product
		}
	}
	return nil
}
