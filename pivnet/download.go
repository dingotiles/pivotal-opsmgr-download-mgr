package pivnet

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/dingodb/pivotal-opsmgr-download-mgr/marketplaces"
)

// DownloadProductTileFile accepts the EULA & starts the download for a PivNet product file
func (pivnetAPI *PivNet) DownloadProductTileFile(tile *marketplaces.ProductTile) (resp *http.Response, err error) {
	fmt.Println("Accepting EULA for", tile.TileName, "via", tile.EULAAcceptanceURL)
	req, err := http.NewRequest("POST", tile.EULAAcceptanceURL, nil)
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

	downloadURL := fmt.Sprintf("%s/download", tile.ProductFileURL)
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
