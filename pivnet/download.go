package pivnet

import (
	"fmt"
	"io"
	"net/http"

	"github.com/dingodb/pivotal-opsmgr-download-mgr/marketplaces"
)

// DownloadProductTileFile accepts the EULA & downloads a product's .pivotal tile file to a io.Writer
func (pivnetAPI *PivNet) DownloadProductTileFile(tile marketplaces.ProductTile, out io.Writer) (err error) {
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

	fmt.Println("Fetching temporary download URL for", tile.TileName, "via", tile.ProductFileURL)
	req, err = http.NewRequest("POST", tile.ProductFileURL, nil)
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
	downloadLocation, err := resp.Location()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Downloading", tile.TileName, "from", downloadLocation)
	req, err = http.NewRequest("GET", downloadLocation.String(), nil)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err.Error())
	}

	return
}
