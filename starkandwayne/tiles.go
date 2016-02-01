package starkandwayne

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/dingodb/pivotal-opsmgr-download-mgr/marketplaces"
)

// UpdateProductTiles fetches available Product Tiles from Pivotal Network
func (starkandwayneAPI *StarkAndWayne) UpdateProductTiles() (err error) {
	req, err := http.NewRequest("GET", starkandwayneAPI.apiURL("/products"), nil)
	if err != nil {
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	// req.Header.Set("Authorization", fmt.Sprintf("Token %s", starkandwayneAPI.apiToken))

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
	fmt.Printf("HELLO\n%+v", productsResp)
	if err != nil {
		return
	}

	for _, product := range productsResp.Products {
		tile := &marketplaces.ProductTile{Slug: product.Slug, MarketplaceSlug: starkandwayneAPI.Slug()}
		starkandwayneAPI.updateProductTileInfo(tile)

		starkandwayneAPI.productTiles = append(starkandwayneAPI.productTiles, tile)
	}

	return
}

func (starkandwayneAPI StarkAndWayne) apiURL(path string) string {
	return fmt.Sprintf("http://localhost:3001/api/v2%s", path)
}