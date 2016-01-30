package opsmgr

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

// Products describes products and all their uploaded version numbers
type Products map[string]productVersions

type productVersions []string

// GetProducts gets the current product/versions that have been uploaded to OpsMgr
func (opsmgr OpsMgr) GetProducts() (products Products, err error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: opsmgr.SkipSSLVerification},
	}
	client := &http.Client{Transport: tr}

	req, err := http.NewRequest("GET", opsmgr.apiURL("/api/products"), nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(opsmgr.Username, opsmgr.Password)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	productsResp := []struct {
		Name    string `json:"name"`
		Version string `json:"product_version"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&productsResp)

	products = Products{}
	for _, productVersion := range productsResp {
		name := productVersion.Name
		versions := products[name]
		if versions == nil {
			versions = productVersions{}
		}
		versions = append(versions, productVersion.Version)
		products[name] = versions
	}

	return
}

func (opsmgr OpsMgr) apiURL(path string) string {
	return fmt.Sprintf("%s%s", opsmgr.URL, path)
}
