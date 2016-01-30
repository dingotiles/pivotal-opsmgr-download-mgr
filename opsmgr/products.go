package opsmgr

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
)

// ProductsResponse lists the product/version that have been uploaded to OpsMgr
type ProductsResponse []struct {
	Name           string `json:"name"`
	ProductVersion string `json:"product_version"`
}

// GetProducts gets the current product/versions that have been uploaded to OpsMgr
func (opsmgr OpsMgr) GetProducts() (products *ProductsResponse, err error) {
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

	products = &ProductsResponse{}
	err = json.NewDecoder(resp.Body).Decode(products)

	return
}

func (opsmgr OpsMgr) apiURL(path string) string {
	return fmt.Sprintf("%s%s", opsmgr.URL, path)
}
