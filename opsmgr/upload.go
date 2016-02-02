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
func (opsmgr OpsMgr) UploadProductFile(tile *marketplaces.ProductTile, downloadResponse *http.Response) (err error) {
	readPipe, writePipe := io.Pipe()
	writer := multipart.NewWriter(writePipe)

	curlBoundary := "------------------------9f5f4cd1c5c81384"
	err = writer.SetBoundary(curlBoundary)
	if err != nil {
		return
	}

	go func() {
		fmt.Printf("create a multipart filter to 'pass through' the data...\n")
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", `form-data; name="product[file]"`)
		h.Set("Content-Type", "application/zip")

		// h.Set("Content-Length", fmt.Sprintf("%d", tile.TileSize))
		part, err := writer.CreatePart(h)
		if err != nil {
			return
		}
		defer writePipe.Close()

		fmt.Printf("copy the download file to the multipart\n")
		io.Copy(part, downloadResponse.Body)
		downloadResponse.Body.Close()
		writer.Close()
	}()

	req, err := http.NewRequest("POST", opsmgr.apiURL("/api/products"), readPipe)
	if err != nil {
		return
	}
	req.SetBasicAuth(opsmgr.Username, opsmgr.Password)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	fmt.Printf("start the 'cross load'...\n")
	dump, err := httputil.DumpRequest(req, true)
	if err == nil {
		fmt.Println(string(dump[:500]))
	}

	uploadResponse, err := opsmgr.httpClient().Do(req)
	if err != nil {
		fmt.Printf("error running upload: %s\n", err)
		return
	}
	fmt.Printf("upload response: %v\n", uploadResponse)
	dump, err = httputil.DumpResponse(uploadResponse, false)
	if err == nil {
		fmt.Println(string(dump))
	}
	return
}