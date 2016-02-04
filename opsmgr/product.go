package opsmgr

import "github.com/hashicorp/go-version"

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
