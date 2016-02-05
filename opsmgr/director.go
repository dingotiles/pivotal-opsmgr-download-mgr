package opsmgr

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry-community/gogobosh"
	"github.com/cloudfoundry-community/gogobosh/api"
	"github.com/cloudfoundry-community/gogobosh/models"
	"github.com/cloudfoundry-community/gogobosh/net"
)

// Director connection details for BOSH director managed by OpsMgr
type Director struct {
	Target   string
	IP       string
	Username string
	Password string
}

// GetDirectorConfig discovers the connection credentials to OpsMgr Director
func (opsmgr *OpsMgr) GetDirectorConfig() (director *Director, err error) {
	if opsmgr.InstallationSettings == nil {
		err = opsmgr.GetInstallationSettings()
		if err != nil {
			return
		}
	}

	for _, product := range opsmgr.InstallationSettings.Products {
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

			director.Target = fmt.Sprintf("https://%s:25555", director.IP)
		}
	}

	return
}

// GetStemcells fetches uploaded stemcells from OpsMgr Director
func (opsmgr *OpsMgr) GetStemcells() (stemcells models.Stemcells, err error) {
	fmt.Println("Fetching OpsMgr Director stemcells...")
	directorConfig, err := opsmgr.GetDirectorConfig()
	if err != nil {
		return
	}
	director := gogobosh.NewDirector(directorConfig.Target, directorConfig.Username, directorConfig.Password)
	repo := api.NewBoshDirectorRepository(&director, net.NewDirectorGateway())
	stemcells, apiResponse := repo.GetStemcells()
	if apiResponse.IsNotSuccessful() {
		return stemcells, fmt.Errorf("Director API error: %s", apiResponse.Message)
	}
	return
}
