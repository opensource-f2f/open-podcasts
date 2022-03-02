package v1alpha1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

func sendJSON(api, payload string) (err error) {
	var req *http.Request
	var resp *http.Response
	if req, err = http.NewRequest(http.MethodPost, api, bytes.NewBuffer(json.RawMessage([]byte(payload)))); err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	if resp, err = client.Do(req); err != nil {
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	return
}

func escapeString(str string) string {
	// Escape double quotes (")
	return strings.Replace(str, `"`, `\"`, -1)
}
