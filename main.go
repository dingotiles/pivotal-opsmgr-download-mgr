package main

import (
	"fmt"
	"os"
	"time"

	"github.com/cloudfoundry-community/gogobosh/models"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/dingodb/pivotal-opsmgr-download-mgr/marketplaces"
	"github.com/dingodb/pivotal-opsmgr-download-mgr/opsmgr"
	"github.com/dingodb/pivotal-opsmgr-download-mgr/pivnet"
	"github.com/go-martini/martini"
	"github.com/hashicorp/go-version"
)

var products *opsmgr.Products
var directorStemcells models.Stemcells
var catalogs marketplaces.Marketplaces
var opsmgrAPI *opsmgr.OpsMgr
var loadingCatalogs bool

func downloadAndUploadTile(opsmgrAPI *opsmgr.OpsMgr, catalog marketplaces.Marketplace, tile *marketplaces.ProductTile) {
	fmt.Printf("starting download...\n")
	downloadResponse, err := catalog.DownloadProductTileFile(tile)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = opsmgrAPI.UploadProductFile(tile, downloadResponse)
	if err != nil {
		fmt.Println(err)
		return
	}

	products, err = opsmgrAPI.GetProducts()
	products.DetermineMarketplaceMappings(catalogs)

	// Errors:
	// - no VPN or bad URL - "Get https://10.58.111.65/api/products: dial tcp 10.58.111.65:443: i/o timeout"
	// - bad connection - "Get https://10.58.111.65/api/products: net/http: TLS handshake timeout"
	// - need to skip SSL validation - "Get https://10.58.111.65/api/products: x509: cannot validate certificate for 10.58.111.65 because it doesn't contain any IP SANs"
	if err != nil {
		fmt.Println(err)
		return
	}
}

func downloadAndUploadStemcell(opsmgrAPI *opsmgr.OpsMgr, catalog marketplaces.Marketplace, stemcell *marketplaces.ProductStemcell) {
	fmt.Printf("starting stemcell download...\n")
	downloadResponse, err := catalog.DownloadProductStemcellFile(stemcell)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = opsmgrAPI.UploadProductStemcell(stemcell, downloadResponse)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func uploadNewStemcells(opsmgrAPI *opsmgr.OpsMgr, catalog marketplaces.Marketplace, stemcells marketplaces.ProductStemcells) {
	latestStemcellVersionUploaded, _ := version.NewVersion("0.0.0")
	for _, productStemcell := range stemcells {
		stemcellVersion, err := version.NewVersion(productStemcell.Version)
		if err != nil {
			fmt.Println("Couldn't parse stemcell version into semver", productStemcell)
			continue
		}
		if productStemcell.Uploaded && latestStemcellVersionUploaded.LessThan(stemcellVersion) {
			latestStemcellVersionUploaded = stemcellVersion
		}
	}

	// Upload any product stemcell that is newer
	for _, s := range stemcells {
		productStemcell := s
		stemcellVersion, _ := version.NewVersion(productStemcell.Version)
		if latestStemcellVersionUploaded.LessThan(stemcellVersion) {
			go func() {
				fmt.Println("Uploading stemcell", productStemcell)
				downloadAndUploadStemcell(opsmgrAPI, catalog, productStemcell)
			}()
		}
	}
}

func refreshOpsMgr() {
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

	directorStemcells, err = opsmgrAPI.GetStemcells()
	if err != nil {
		fmt.Println(err)
	}

	products.DetermineMarketplaceMappings(catalogs)
}

func refreshMarketplaceCatalogs() {
	catalog := pivnet.NewPivNet(os.Getenv("PIVOTAL_NETWORK_TOKEN"), "https://network.pivotal.io/api/v2")
	catalogs[catalog.Slug()] = catalog

	for _, catalog := range catalogs {
		fmt.Printf("Fetching available product tiles & stemcells from %s...\n", catalog.Name())
		err := catalog.UpdateProductCatalog()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println(catalog.Name(), catalog.ProductTiles())
		products.DetermineMarketplaceMappings(catalogs)
	}

	catalog.DetermineStemcellsUploaded(directorStemcells)
	uploadNewStemcells(opsmgrAPI, catalog, catalog.ProductStemcells())

	loadingCatalogs = false
}

func main() {
	opsmgrAPI = opsmgr.NewOpsMgr()
	catalogs = marketplaces.NewMarketplaces()
	loadingCatalogs = true

	refreshRate := 5 * time.Second

	go func() {
		refreshOpsMgr()
		time.Sleep(refreshRate)
	}()
	go func() {
		refreshMarketplaceCatalogs()
		time.Sleep(refreshRate)
	}()

	m := martini.Classic()
	m.Use(render.Renderer(render.Options{Layout: "layout"}))

	m.Get("/", func(r render.Render) {
		if products == nil {
			r.HTML(200, "loading", nil)
		} else {
			r.HTML(200, "index", struct {
				OpsMgrProducts  *opsmgr.Products
				PivNetTiles     marketplaces.ProductTiles
				PivNetStemcells marketplaces.ProductStemcells
				LoadingCatalogs bool
			}{
				products,
				catalogs["pivnet"].ProductTiles(),
				catalogs["pivnet"].ProductStemcells(),
				loadingCatalogs,
			})
		}
	})
	m.Get("/director", func(r render.Render) {
		director, err := opsmgrAPI.GetDirectorConfig()

		r.HTML(200, "director", struct {
			Director *opsmgr.Director
			Error    error
		}{director, err})
	})

	m.Get("/install/:marketplace/tile/:tilename", func(params martini.Params, r render.Render) {
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
		downloadAndUploadTile(opsmgrAPI, catalog, tile)
	})

	m.Get("/install/:marketplace/stemcell/:stemcell", func(params martini.Params, r render.Render) {
		marketplaceSlug := params["marketplace"]
		r.Redirect("/")

		catalog := catalogs[marketplaceSlug]
		if catalog == nil {
			fmt.Println("Unknown :marketplace slug", marketplaceSlug)
			return
		}
		stemcellVersion := params["stemcell"]
		stemcell := catalog.LookupStemcell(stemcellVersion)
		if stemcell == nil {
			fmt.Printf("Unknown %s stemcell version %s\n", marketplaceSlug, stemcellVersion)
			return
		}
		fmt.Println("stemcell", stemcell)
		downloadAndUploadStemcell(opsmgrAPI, catalog, stemcell)
	})

	m.Get("/deleteunused", func(r render.Render) {
		err := opsmgrAPI.DeleteUnusedTiles()
		if err != nil {
			fmt.Println("/deleteunused", err)
		}
		r.Redirect("/")
	})

	m.Run()
}
