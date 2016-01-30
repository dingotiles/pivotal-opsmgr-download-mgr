package marketplaces

// ProductTiles is a list of product tiles
type ProductTiles []string

// Marketplace is an interface to PivNet or StarkAndWayneMarketplace
type Marketplace interface {
	GetProductTiles() ProductTiles
}
