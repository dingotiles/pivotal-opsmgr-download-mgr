package opsmgr

import "github.com/dingodb/pivotal-opsmgr-download-mgr/marketplaces"

// DetermineMarketplaceMappings determines which marketplace and what its tile name is on that marketplace
func (products *Products) DetermineMarketplaceMappings(catalogs marketplaces.Marketplaces) {
	for _, product := range *products {
		product.DetermineMarketplaceTile(catalogs)
	}
}

// DetermineMarketplaceTile determine which marketplace and what its tile name is on that marketplace
func (product *Product) DetermineMarketplaceTile(catalogs marketplaces.Marketplaces) {
	product.Marketplace = "unknown"
	product.MarketplaceTileName = product.Name
	for _, marketplace := range catalogs {
		tile := marketplace.LookupProductTile(product.Name)
		if tile != nil {
			product.Marketplace = marketplace.Slug()
			product.MarketplaceProductName = tile.Slug
			product.MarketplaceTileName = tile.TileName
			product.MarketplaceTileVersion = tile.TileVersion
		}
	}
}
