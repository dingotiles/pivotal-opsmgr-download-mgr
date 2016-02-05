package opsmgr

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"

	"github.com/hashicorp/go-version"
)

// Products describes products and all their uploaded version numbers
type Products map[string]*Product

// GetProducts gets the current product/versions that have been uploaded to OpsMgr
func (opsmgr *OpsMgr) GetProducts() (products *Products, err error) {
	req, err := http.NewRequest("GET", opsmgr.apiURL("/api/products"), nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(opsmgr.Username, opsmgr.Password)

	resp, err := opsmgr.httpClient().Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	productsResp := []struct {
		Name    string `json:"name"`
		Version string `json:"product_version"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&productsResp)

	products = &Products{}
	for _, productVersion := range productsResp {
		name := productVersion.Name
		if name == "p-bosh" {
			continue
		}
		product := (*products)[name]
		if product == nil {
			product = &Product{}
		}
		product.Name = name
		version, err := version.NewVersion(productVersion.Version)
		if err != nil {
			fmt.Printf("Error parsing product version %s for %s product\n", productVersion.Version, name)
			product.RawVersions = append(product.RawVersions, productVersion.Version)
		} else {
			product.Versions = append(product.Versions, version)
		}
		(*products)[name] = product
	}

	// Sort the product versions so we can determine latest
	// TODO: sorting non-semver (e.g. p-bosh '1.6.7.0')
	for _, product := range *products {
		if len(product.Versions) > 0 {
			sort.Sort(version.Collection(product.Versions))
			product.LatestVersion = product.Versions[len(product.Versions)-1].String()
		} else {
			product.LatestVersion = product.RawVersions[0]
		}
		fmt.Println("Latest version", product.Name, product.LatestVersion)
	}

	return
}
