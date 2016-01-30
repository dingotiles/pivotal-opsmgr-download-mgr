package marketplaces

// Marketplaces references all configured Marketplace APIs, such as PivNet, indexed by Slug
type Marketplaces map[string]Marketplace

// NewMarketplaces creates new Marketplaces struct
func NewMarketplaces() Marketplaces {
	return Marketplaces{}
}

// ProductTiles is a list of product tiles
type ProductTiles []string

// ProductTile represents a catalog listing for a tile
type ProductTile struct {
	Slug string
}

// Marketplace is an interface to PivNet or StarkAndWayneMarketplace
type Marketplace interface {
	Name() string
	Slug() string
	UpdateProductTiles() error
	ProductTiles() ProductTiles
	LookupProductTile(productName string) *ProductTile
}
