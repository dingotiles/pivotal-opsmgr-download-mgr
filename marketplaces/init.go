package marketplaces

import "net/http"

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
	Slug              string
	MarketplaceSlug   string
	Tile              bool
	ReleaseDate       string
	TileName          string
	TileVersion       string
	TileSize          uint64
	TileHumanSize     string
	EULAAcceptanceURL string
	ProductFileURL    string
}

// Marketplace is an interface to PivNet or StarkAndWayneMarketplace
type Marketplace interface {
	Name() string
	Slug() string
	UpdateProductTiles() error
	ProductTiles() ProductTiles
	LookupProductTile(productName string) *ProductTile
	DownloadProductTileFile(tile *ProductTile) (*http.Response, error)
}
