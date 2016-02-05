package opsmgr

import (
	"fmt"
	"strings"
)

// Director connection details for BOSH director managed by OpsMgr
type Director struct {
	URL      string
	IP       string
	Username string
	Password string
}

// GetDirector discovers the connection credentials to OpsMgr Director
func (opsmgr *OpsMgr) GetDirector() (director *Director, err error) {
	settings, err := opsmgr.GetInstallationSettings()
	if err != nil {
		return
	}

	for _, product := range settings.Products {
		if strings.HasPrefix(product.InstallationName, "p-bosh") {
			director = &Director{}

			soloJob := product.Jobs[0]
			for _, property := range soloJob.Properties {
				if property.Identifier == "director_credentials" {
					cred := property.Value.(map[string]interface{})
					director.Username = cred["identity"].(string)
					director.Password = cred["password"].(string)
					break
				}
			}

			for _, ips := range product.Ips {
				director.IP = ips[0]
				break
			}

			director.URL = fmt.Sprintf("https://%s:25555", director.IP)
		}
	}

	return
}
