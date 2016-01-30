package opsmgr

import (
	"encoding/json"
	"log"
	"os"
	"strings"
)

// OpsMgr is configuration for a target OpsMgr deployment
type OpsMgr struct {
	URL                 string
	SkipSSLVerification bool
	Username            string
	Password            string
	PivNetMapping       MarketplaceMapping
}

// MarketplaceMapping maps an opsmgr product name to a non-standard/guessable marketplace tile name
type MarketplaceMapping map[string]string

// NewOpsMgr creates a new OpsMgr struct
func NewOpsMgr() OpsMgr {
	skipSSLVerification := false
	if os.Getenv("OPSMGR_SKIP_SSL_VERIFICATION") != "" || os.Getenv("OPSMGR_SKIP_SSL_VALIDATION") != "" {
		skipSSLVerification = true
	}
	mappingStr := os.Getenv("PIVNET_PRODUCT_TO_TILE_MAPPING")
	pivnetMapping := MarketplaceMapping{}
	if mappingStr != "" {
		err := json.NewDecoder(strings.NewReader(mappingStr)).Decode(&pivnetMapping)
		if err != nil {
			log.Fatalf("Invalid JSON array in $PIVNET_PRODUCT_TO_TILE_MAPPING: %s", mappingStr)
		}
	}
	return OpsMgr{
		URL:                 os.Getenv("OPSMGR_URL"),
		SkipSSLVerification: skipSSLVerification,
		Username:            os.Getenv("OPSMGR_USERNAME"),
		Password:            os.Getenv("OPSMGR_PASSWORD"),
		PivNetMapping:       pivnetMapping,
	}
}
