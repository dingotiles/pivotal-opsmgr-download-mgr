package pivnet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/dingodb/pivotal-opsmgr-download-mgr/marketplaces"
	"github.com/dustin/go-humanize"
)

// LookupProductTile tries to match an Opsmgr Product name with a PivNet product/release/.pivotal tile
func (pivnetAPI *PivNet) LookupProductTile(opsMgrProductName string) *marketplaces.ProductTile {
	for _, product := range pivnetAPI.productTiles {
		if product.TileName == opsMgrProductName {
			return product
		}
	}
	return nil
}

// LookupStemcell tries to match an Opsmgr Product name with a PivNet product/release/.pivotal tile
func (pivnetAPI *PivNet) LookupStemcell(version string) *marketplaces.ProductStemcell {
	for _, stemcell := range pivnetAPI.stemcells {
		if stemcell.Version == version {
			return stemcell
		}
	}
	return nil
}

func (pivnetAPI *PivNet) updateProductTileInfo(tile *marketplaces.ProductTile) (err error) {
	req, err := http.NewRequest("GET", pivnetAPI.apiURL(fmt.Sprintf("/products/%s/releases", tile.Slug)), nil)
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

	releasesResp := struct {
		Releases []struct {
			ID              int    `json:"id"`
			Version         string `json:"version"`
			ReleaseType     string `json:"release_type"`
			ReleaseDate     string `json:"release_date"`
			ReleaseNotesURL string `json:"release_notes_url"`
			Availability    string `json:"availability"`
			Description     string `json:"description"`
			Eula            struct {
				ID    int    `json:"id"`
				Slug  string `json:"slug"`
				Name  string `json:"name"`
				Links struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"_links"`
			} `json:"eula"`
			EndOfSupportDate string `json:"end_of_support_date"`
			Eccn             string `json:"eccn"`
			LicenseException string `json:"license_exception"`
			Controlled       bool   `json:"controlled"`
			Links            struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				EulaAcceptance struct {
					Href string `json:"href"`
				} `json:"eula_acceptance"`
				ProductFiles struct {
					Href string `json:"href"`
				} `json:"product_files"`
				FileGroups struct {
					Href string `json:"href"`
				} `json:"file_groups"`
				UserGroups struct {
					Href string `json:"href"`
				} `json:"user_groups"`
			} `json:"_links"`
		} `json:"releases"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"_links"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&releasesResp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 1. find latest release
	latestReleaseDatedReleaseID := 0
	latestReleaseDate := "0000-00-00"
	latestReleaseVersion := ""
	for _, release := range releasesResp.Releases {
		if strings.Compare(latestReleaseDate, release.ReleaseDate) < 0 {
			latestReleaseDate = release.ReleaseDate
			latestReleaseDatedReleaseID = release.ID
			latestReleaseVersion = release.Version
		}
	}
	fmt.Printf("Latest release for %s is '%s' date %s with ID %d\n", tile.Slug, latestReleaseVersion, latestReleaseDate, latestReleaseDatedReleaseID)
	// Skip if product has no releases and hence no product_files which might be .pivotal tiles
	if latestReleaseDate == "0000-00-00" {
		return
	}

	tile.EULAAcceptanceURL = pivnetAPI.apiURL(fmt.Sprintf("/products/%s/releases/%d/eula_acceptance", tile.Slug, latestReleaseDatedReleaseID))

	// 2. look at product_files for one with aws_object_key ~= /<name>-<version>.pivotal/
	req, err = http.NewRequest("GET", pivnetAPI.apiURL(fmt.Sprintf("/products/%s/releases/%d/product_files", tile.Slug, latestReleaseDatedReleaseID)), nil)
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

	productFilesResp := struct {
		ProductFiles []struct {
			ID           int    `json:"id"`
			AwsObjectKey string `json:"aws_object_key"`
			FileVersion  string `json:"file_version"`
			Name         string `json:"name"`
			Links        struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				Download struct {
					Href string `json:"href"`
				} `json:"download"`
				SignatureFileDownload struct {
					Href interface{} `json:"href"`
				} `json:"signature_file_download"`
			} `json:"_links"`
		} `json:"product_files"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"_links"`
	}{}

	err = json.NewDecoder(resp.Body).Decode(&productFilesResp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	r, _ := regexp.Compile("([a-zA-Z_-]+)-(v?[0-9][a-zA-Z0-9._-]*)\\.pivotal")

	productFileID := 0

	// 3. if so, then it is a Tile; and deduce its product TileName & TileVersion
	for _, productFile := range productFilesResp.ProductFiles {
		fmt.Println("Checking if product file", productFile.AwsObjectKey, "is a .pivotal file...")
		tileTokens := r.FindStringSubmatch(productFile.AwsObjectKey)
		if len(tileTokens) == 3 {
			productFileID = productFile.ID
			tile.Tile = true
			tile.ReleaseDate = latestReleaseDate
			tile.TileName = strings.ToLower(tileTokens[1])
			tile.TileVersion = tileTokens[2]
			tile.ProductFileName = productFile.AwsObjectKey
			fmt.Printf("Found tile ID %d, named %s %s for product %s\n", productFileID, tile.TileName, tile.TileVersion, tile.Slug)
			break
		}
	}
	// Get product size
	if productFileID > 0 {
		fmt.Println("Looking up file size for product file", productFileID)
		tile.ProductFileURL = pivnetAPI.apiURL(fmt.Sprintf("/products/%s/releases/%d/product_files/%d", tile.Slug, latestReleaseDatedReleaseID, productFileID))
		req, err = http.NewRequest("GET", tile.ProductFileURL, nil)
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

		productFileResp := struct {
			ProductFile struct {
				ID                 int           `json:"id"`
				AwsObjectKey       string        `json:"aws_object_key"`
				Description        string        `json:"description"`
				DocsURL            string        `json:"docs_url"`
				FileType           string        `json:"file_type"`
				FileVersion        string        `json:"file_version"`
				IncludedFiles      []interface{} `json:"included_files"`
				Md5                string        `json:"md5"`
				Name               string        `json:"name"`
				Platforms          []interface{} `json:"platforms"`
				ReleasedAt         string        `json:"released_at"`
				Size               uint64        `json:"size"`
				SystemRequirements []interface{} `json:"system_requirements"`
				Links              struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
					Download struct {
						Href string `json:"href"`
					} `json:"download"`
					SignatureFileDownload struct {
						Href interface{} `json:"href"`
					} `json:"signature_file_download"`
				} `json:"_links"`
			} `json:"product_file"`
		}{}

		err = json.NewDecoder(resp.Body).Decode(&productFileResp)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		tile.TileSize = productFileResp.ProductFile.Size
		tile.TileHumanSize = humanize.Bytes(productFileResp.ProductFile.Size)
	}
	return
}

func (pivnetAPI *PivNet) updateStemcellsInfo(cpi string) (err error) {
	fmt.Println("Updating list of available stemcells")
	req, err := http.NewRequest("GET", pivnetAPI.apiURL(fmt.Sprintf("/products/stemcells/releases")), nil)
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

	releasesResp := struct {
		Releases []struct {
			ID              int    `json:"id"`
			Version         string `json:"version"`
			ReleaseType     string `json:"release_type"`
			ReleaseDate     string `json:"release_date"`
			ReleaseNotesURL string `json:"release_notes_url"`
			Availability    string `json:"availability"`
			Description     string `json:"description"`
			Eula            struct {
				ID    int    `json:"id"`
				Slug  string `json:"slug"`
				Name  string `json:"name"`
				Links struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
				} `json:"_links"`
			} `json:"eula"`
			EndOfSupportDate string `json:"end_of_support_date"`
			Eccn             string `json:"eccn"`
			LicenseException string `json:"license_exception"`
			Controlled       bool   `json:"controlled"`
			Links            struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				EulaAcceptance struct {
					Href string `json:"href"`
				} `json:"eula_acceptance"`
				ProductFiles struct {
					Href string `json:"href"`
				} `json:"product_files"`
				FileGroups struct {
					Href string `json:"href"`
				} `json:"file_groups"`
				UserGroups struct {
					Href string `json:"href"`
				} `json:"user_groups"`
			} `json:"_links"`
		} `json:"releases"`
		Links struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
		} `json:"_links"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&releasesResp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for _, stemcellRelease := range releasesResp.Releases {
		stemcell := &marketplaces.ProductStemcell{
			Slug:            "stemcells",
			MarketplaceSlug: pivnetAPI.Slug(),
			CPI:             cpi,
			Version:         stemcellRelease.Version,
		}
		stemcell.EULAAcceptanceURL = pivnetAPI.apiURL(fmt.Sprintf("/products/stemcells/releases/%d/eula_acceptance", stemcellRelease.ID))

		// 2. look at product_files for one with aws_object_key ~= /bosh-stemcell-*-<cpi>*.tgz/
		req, err = http.NewRequest("GET", pivnetAPI.apiURL(fmt.Sprintf("/products/stemcells/releases/%d/product_files", stemcellRelease.ID)), nil)
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

		productFilesResp := struct {
			ProductFiles []struct {
				ID           int    `json:"id"`
				AwsObjectKey string `json:"aws_object_key"`
				FileVersion  string `json:"file_version"`
				Name         string `json:"name"`
				Links        struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
					Download struct {
						Href string `json:"href"`
					} `json:"download"`
					SignatureFileDownload struct {
						Href interface{} `json:"href"`
					} `json:"signature_file_download"`
				} `json:"_links"`
			} `json:"product_files"`
			Links struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
			} `json:"_links"`
		}{}

		err = json.NewDecoder(resp.Body).Decode(&productFilesResp)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		stemcell.ReleaseDate = stemcellRelease.ReleaseDate

		// bosh-stemcell-3160-vsphere-esxi-ubuntu-trusty-go_agent.tgz
		r, _ := regexp.Compile(fmt.Sprintf("bosh-stemcell-.*-%s-.*\\.tgz", stemcell.CPI))

		var wg sync.WaitGroup

		for _, pf := range productFilesResp.ProductFiles {
			productFile := pf
			wg.Add(1)
			go func() {
				defer wg.Done()

				fmt.Println("Checking if product file", productFile.AwsObjectKey, "is a", cpi, "stemcell...")
				tileTokens := r.FindStringSubmatch(productFile.AwsObjectKey)
				if len(tileTokens) > 0 {
					stemcell.ProductFileName = productFile.AwsObjectKey
					stemcell.ProductFileURL = pivnetAPI.apiURL(fmt.Sprintf("/products/%s/releases/%d/product_files/%d", stemcell.Slug, stemcellRelease.ID, productFile.ID))
					fmt.Printf("Found stemcell %s\n", stemcell.ProductFileName)
					return
				}
			}()
		}
		wg.Wait()

		pivnetAPI.stemcells = append(pivnetAPI.stemcells, stemcell)
	}

	return
}
