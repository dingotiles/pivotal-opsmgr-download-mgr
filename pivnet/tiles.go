package pivnet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/dingodb/pivotal-opsmgr-download-mgr/marketplaces"
)

// UpdateProductCatalog fetches available Product Tiles from Pivotal Network
func (pivnetAPI *PivNet) UpdateProductCatalog() (err error) {
	req, err := http.NewRequest("GET", pivnetAPI.apiURL("/products"), nil)
	if err != nil {
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", pivnetAPI.apiToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	productsResp := struct {
		Products []struct {
			ID      int    `json:"id"`
			Slug    string `json:"slug"`
			Name    string `json:"name"`
			LogoURL string `json:"logo_url"`
			Links   struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				Releases struct {
					Href string `json:"href"`
				} `json:"releases"`
				ProductFiles struct {
					Href string `json:"href"`
				} `json:"product_files"`
				FileGroups struct {
					Href string `json:"href"`
				} `json:"file_groups"`
			} `json:"_links"`
		} `json:"products"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"_links"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&productsResp)
	if err != nil {
		return
	}

	var wg sync.WaitGroup
	for _, p := range productsResp.Products {
		product := p
		wg.Add(1)
		go func() {
			defer wg.Done()
			if product.Slug == "stemcells" {
				pivnetAPI.updateStemcellsInfo("vsphere")
			} else {
				tile := &marketplaces.ProductTile{Slug: product.Slug, MarketplaceSlug: pivnetAPI.Slug()}
				pivnetAPI.updateProductTileInfo(tile)
				pivnetAPI.productTiles = append(pivnetAPI.productTiles, tile)
			}
		}()
	}
	wg.Wait()

	return
}

func (pivnetAPI PivNet) apiURL(path string) string {
	return pivnetAPI.apiEndpoint + path
}
