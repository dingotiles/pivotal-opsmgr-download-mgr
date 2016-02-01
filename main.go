package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"

	"github.com/codegangsta/martini-contrib/render"
	"github.com/dingodb/pivotal-opsmgr-download-mgr/marketplaces"
	"github.com/dingodb/pivotal-opsmgr-download-mgr/opsmgr"
	"github.com/dingodb/pivotal-opsmgr-download-mgr/starkandwayne"
	"github.com/go-martini/martini"
)

func main() {
	catalogs := marketplaces.NewMarketplaces()
	var products *opsmgr.Products
	loadingCatalogs := true

	go func() {
		opsmgrAPI := opsmgr.NewOpsMgr()
		fmt.Printf("Fetching uploaded products from OpsMgr %s...\n", opsmgrAPI.URL)
		var err error
		products, err = opsmgrAPI.GetProducts()

		// Errors:
		// - no VPN or bad URL - "Get https://10.58.111.65/api/products: dial tcp 10.58.111.65:443: i/o timeout"
		// - bad connection - "Get https://10.58.111.65/api/products: net/http: TLS handshake timeout"
		// - need to skip SSL validation - "Get https://10.58.111.65/api/products: x509: cannot validate certificate for 10.58.111.65 because it doesn't contain any IP SANs"
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		starkCatalog := starkandwayne.NewStarkAndWayne()
		catalogs[starkCatalog.Slug()] = starkCatalog
		// catalog := pivnet.NewPivNet()
		// catalogs[catalog.Slug()] = catalog

		for _, catalog := range catalogs {
			fmt.Printf("Fetching available product tiles from %s...\n", catalog.Name())
			err := catalog.UpdateProductTiles()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Println(catalog.Name(), catalog.ProductTiles())
		}

		products.DetermineMarketplaceMappings(catalogs)
		loadingCatalogs = false
	}()

	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		if products == nil {
			r.HTML(200, "loading", nil)
		} else {
			r.HTML(200, "index", struct {
				OpsMgrProducts     *opsmgr.Products
				PivNetTiles        marketplaces.ProductTiles
				StarkAndWayneTiles marketplaces.ProductTiles
				LoadingCatalogs    bool
			}{products, catalogs["pivnet"].ProductTiles(), catalogs["starkandwayne"].ProductTiles(), loadingCatalogs})
		}
	})
	m.Get("/install/:marketplace/:tilename", func(params martini.Params, r render.Render) {
		marketplaceSlug := params["marketplace"]
		r.Redirect("/")

		catalog := catalogs[marketplaceSlug]
		if catalog == nil {
			fmt.Println("Unknown :marketplace slug", marketplaceSlug)
			return
		}
		tileSlug := params["tilename"]
		tile := catalog.LookupProductTile(tileSlug)
		if tile == nil {
			fmt.Printf("Unknown %s product %s\n", marketplaceSlug, tileSlug)
			return
		}

		fmt.Println(tile)
		buffer := &bytes.Buffer{}
		body := bufio.NewWriter(buffer)
		catalog.DownloadProductTileFile(tile, body)
		fmt.Printf("Downloaded %v, size %d\n", tile, buffer.Len())

		r.Redirect("/")
	})
	m.Run()
}
