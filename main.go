package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/martini-contrib/render"
	"github.com/dingodb/pivotal-opsmgr-download-mgr/marketplaces"
	"github.com/dingodb/pivotal-opsmgr-download-mgr/opsmgr"
	"github.com/dingodb/pivotal-opsmgr-download-mgr/pivnet"
	"github.com/go-martini/martini"
)

func main() {
	opsmgrAPI := opsmgr.NewOpsMgr()
	fmt.Printf("Fetching uploaded products from OpsMgr %s...\n", opsmgrAPI.URL)
	products, err := opsmgrAPI.GetProducts()

	// Errors:
	// - no VPN or bad URL - "Get https://10.58.111.65/api/products: dial tcp 10.58.111.65:443: i/o timeout"
	// - bad connection - "Get https://10.58.111.65/api/products: net/http: TLS handshake timeout"
	// - need to skip SSL validation - "Get https://10.58.111.65/api/products: x509: cannot validate certificate for 10.58.111.65 because it doesn't contain any IP SANs"
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	catalogs := marketplaces.NewMarketplaces()

	catalog := pivnet.NewPivNet()
	catalogs[catalog.Slug()] = catalog

	for _, marketplace := range catalogs {
		fmt.Printf("Fetching available product tiles from %s...\n", marketplace.Name())
		err := marketplace.UpdateProductTiles()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(marketplace.Name(), marketplace.ProductTiles())
	}

	products.DetermineMarketplaceMappings(catalogs)

	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", struct {
			OpsMgrProducts *opsmgr.Products
			PivNetTiles    marketplaces.ProductTiles
		}{products, catalogs["pivnet"].ProductTiles()})
	})
	m.Run()
}
