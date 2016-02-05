package pivnet

import (
	"fmt"

	"github.com/cloudfoundry-community/gogobosh/models"
)

// DetermineStemcellsUploaded updates each ProductStemcell if it is already uploaded to Director
func (pivnetAPI *PivNet) DetermineStemcellsUploaded(directorStemcells models.Stemcells) (err error) {
	for _, directorStemcell := range directorStemcells {
		for _, productStemcell := range pivnetAPI.ProductStemcells() {
			if directorStemcell.Version == productStemcell.Version {
				productStemcell.Uploaded = true
				fmt.Printf("Stemcell %s %s already uploaded to OpsMgr\n", productStemcell.Slug, productStemcell.Version)
				break
			}
		}
	}
	return
}
