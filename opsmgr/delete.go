package opsmgr

import "net/http"

// DeleteUnusedTiles asks OpsMgr to delete any uploaded tiles that are not being used/installed
func (opsmgr OpsMgr) DeleteUnusedTiles() (err error) {
	req, err := http.NewRequest("DELETE", opsmgr.apiURL("/api/products"), nil)
	if err != nil {
		return
	}
	req.SetBasicAuth(opsmgr.Username, opsmgr.Password)

	resp, err := opsmgr.httpClient().Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	return
}
