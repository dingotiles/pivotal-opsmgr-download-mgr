package opsmgr

import (
	"regexp"

	"github.com/hashicorp/go-version"
)

// Product includes the uploaded product tile versions, and reference to the marketplace/tile name
type Product struct {
	Name                   string
	Versions               []*version.Version
	RawVersions            []string // if version not semver (e.g. 1.2.3.4)
	LatestVersion          string
	Marketplace            string
	MarketplaceProductName string
	MarketplaceTileName    string
	MarketplaceTileVersion string
}

// NewTileAvailable is true if MarketplaceTileVersion indicates a newer version than LatestVersion
func (product *Product) NewTileAvailable() bool {
	if product.MarketplaceTileVersion == "" {
		return false
	}
	buildSuffix, _ := regexp.Compile("-build.*")
	uploadedVersion := buildSuffix.ReplaceAllString(product.LatestVersion, "")
	marketplaceVersion := buildSuffix.ReplaceAllString(product.MarketplaceTileVersion, "")
	return uploadedVersion != marketplaceVersion
}
