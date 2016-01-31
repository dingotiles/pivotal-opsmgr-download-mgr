package marketplaces

// Marketplaces references all configured Marketplace APIs, such as PivNet, indexed by Slug
type Marketplaces map[string]Marketplace

// NewMarketplaces creates new Marketplaces struct
func NewMarketplaces() Marketplaces {
	return Marketplaces{}
}

// ProductTiles is a list of product tiles
type ProductTiles []*ProductTile

// ProductTile represents a catalog listing for a tile (or items that aren't .pivotal tiles)
type ProductTile struct {
	Slug               string
	Tile               bool
	TileName           string
	TileVersion        string
	TileSize           int64
	TileProductFileURL string
}

// Marketplace is an interface to PivNet or StarkAndWayneMarketplace
type Marketplace interface {
	Name() string
	Slug() string
	UpdateProductTiles() error
	ProductTiles() ProductTiles
	LookupProductTile(productName string) *ProductTile
}
