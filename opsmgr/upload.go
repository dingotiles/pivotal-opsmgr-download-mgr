package opsmgr

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httputil"
	"net/textproto"

	"github.com/dingodb/pivotal-opsmgr-download-mgr/marketplaces"
)

// UploadProductFile uploads a .pivotal file to your OpsMgr
func (opsmgr *OpsMgr) UploadProductFile(tile *marketplaces.ProductTile, downloadResponse *http.Response) (err error) {
	return opsmgr.uploadFile("product", tile.ProductFileName, downloadResponse)
}

// UploadProductStemcell uploads a stemcell to your OpsMgr
func (opsmgr *OpsMgr) UploadProductStemcell(stemcell *marketplaces.ProductStemcell, downloadResponse *http.Response) (err error) {
	return opsmgr.uploadFile("stemcell", stemcell.ProductFileName, downloadResponse)
}

func (opsmgr *OpsMgr) uploadFile(uploadEndpoint string, fileName string, downloadResponse *http.Response) (err error) {
	readPipe, writePipe := io.Pipe()
	writer := multipart.NewWriter(writePipe)

	uploadAPIEndpoint := fmt.Sprintf("/api/%ss", uploadEndpoint)

	go func() {
		fmt.Printf("create a multipart filter to 'pass through' the data...\n")
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s[file]"; filename="%s"`, uploadEndpoint, fileName))
		h.Set("Content-Type", "application/octet-stream")

		part, err := writer.CreatePart(h)
		if err != nil {
			fmt.Printf("error creating part: %s\n", err)
			return
		}
		defer writePipe.Close()

		fmt.Printf("copy the download file to the multipart\n")
		io.Copy(part, downloadResponse.Body)
		downloadResponse.Body.Close()
		writer.Close()
	}()

	req, err := http.NewRequest("POST", opsmgr.apiURL(uploadAPIEndpoint), readPipe)
	if err != nil {
		fmt.Printf("error creating request: %s\n", err)
		return
	}
	req.SetBasicAuth(opsmgr.Username, opsmgr.Password)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	fmt.Printf("dump upload %s request...\n", uploadEndpoint)
	dump, err := httputil.DumpRequest(req, false)
	if err == nil {
		fmt.Println(string(dump))
	}

	fmt.Printf("start the 'cross load'...\n")
	uploadResponse, err := opsmgr.httpClient().Do(req)
	if err != nil {
		fmt.Printf("error running upload: %s\n", err)
		return
	}
	fmt.Printf("upload %s response: %v\n", uploadEndpoint, uploadResponse)
	dump, err = httputil.DumpResponse(uploadResponse, false)
	if err == nil {
		fmt.Println(string(dump))
	}
	return
}
