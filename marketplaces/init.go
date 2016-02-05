package marketplaces

import (
	"net/http"

	"github.com/cloudfoundry-community/gogobosh/models"
)

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
	ProductFileName   string
}

// ProductStemcells is a list of stemcells for target CPI
type ProductStemcells []*ProductStemcell

// ProductStemcell presents a stemcell file that can be downloaded from a marketplace
type ProductStemcell struct {
	Slug              string
	MarketplaceSlug   string
	CPI               string
	Version           string
	Uploaded          bool
	ReleaseDate       string
	EULAAcceptanceURL string
	ProductFileURL    string
	ProductFileName   string
}

// Marketplace is an interface to PivNet or StarkAndWayneMarketplace
type Marketplace interface {
	Name() string
	Slug() string
	UpdateProductCatalog() error
	ProductTiles() ProductTiles
	ProductStemcells() ProductStemcells
	LookupProductTile(productName string) *ProductTile
	LookupStemcell(version string) *ProductStemcell
	DownloadProductTileFile(tile *ProductTile) (*http.Response, error)
	DownloadProductStemcellFile(stemcell *ProductStemcell) (*http.Response, error)
	DetermineStemcellsUploaded(directorStemcells models.Stemcells) (err error)
}
