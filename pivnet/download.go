package pivnet

import (
	"fmt"
	"io"
	"net/http"

	"github.com/dingodb/pivotal-opsmgr-download-mgr/marketplaces"
)

// DownloadProductTileFile accepts the EULA & downloads a product's .pivotal tile file to a io.Writer
func (pivnetAPI *PivNet) DownloadProductTileFile(tile *marketplaces.ProductTile, out io.Writer) (err error) {
	fmt.Println("Accepting EULA for", tile.TileName, "via", tile.EULAAcceptanceURL)
	req, err := http.NewRequest("POST", tile.EULAAcceptanceURL, nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", pivnetAPI.apiToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	downloadURL := fmt.Sprintf("%s/download", tile.ProductFileURL)
	fmt.Println("Fetching temporary download URL for", tile.TileName, "via", downloadURL)
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
	defer resp.Body.Close()
	fmt.Println("Storing", tile.TileName, "into", out)
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	return
}
