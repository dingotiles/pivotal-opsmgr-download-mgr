package opsmgr

import "os"

// OpsMgr is configuration for a target OpsMgr deployment
type OpsMgr struct {
	URL                 string
	SkipSSLVerification bool
	Username            string
	Password            string
}

// NewOpsMgr creates a new OpsMgr struct
func NewOpsMgr() OpsMgr {
	skipSSLVerification := false
	if os.Getenv("OPSMGR_SKIP_SSL_VERIFICATION") != "" || os.Getenv("OPSMGR_SKIP_SSL_VALIDATION") != "" {
		skipSSLVerification = true
	}
	return OpsMgr{
		URL:                 os.Getenv("OPSMGR_URL"),
		SkipSSLVerification: skipSSLVerification,
		Username:            os.Getenv("OPSMGR_USERNAME"),
		Password:            os.Getenv("OPSMGR_PASSWORD"),
	}
}
