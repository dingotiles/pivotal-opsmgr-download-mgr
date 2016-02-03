package pivnet

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/dingodb/pivotal-opsmgr-download-mgr/marketplaces"
)

// DownloadProductTileFile accepts the EULA & starts the download for a PivNet product file
func (pivnetAPI *PivNet) DownloadProductTileFile(tile *marketplaces.ProductTile) (resp *http.Response, err error) {
	return pivnetAPI.acceptEULAAndDownload(tile.TileName, tile.EULAAcceptanceURL, tile.ProductFileURL)
}

// DownloadProductStemcellFile accepts the EULA & starts the download for a PivNet product stemcell
func (pivnetAPI *PivNet) DownloadProductStemcellFile(stemcell *marketplaces.ProductStemcell) (resp *http.Response, err error) {
	return pivnetAPI.acceptEULAAndDownload(
		fmt.Sprintf("Stemcell v%s", stemcell.Version),
		stemcell.EULAAcceptanceURL,
		stemcell.ProductFileURL)
}

func (pivnetAPI *PivNet) acceptEULAAndDownload(title string, eulaAcceptanceURL string, productFileURL string) (resp *http.Response, err error) {
	fmt.Println("Accepting EULA for", title, "via", eulaAcceptanceURL)
	req, err := http.NewRequest("POST", eulaAcceptanceURL, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", pivnetAPI.apiToken))

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	downloadURL := fmt.Sprintf("%s/download", productFileURL)
	req, err = http.NewRequest("POST", downloadURL, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", pivnetAPI.apiToken))

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	dump, err := httputil.DumpResponse(resp, false)
	if err == nil {
		fmt.Println(string(dump))
	}

	return
}
