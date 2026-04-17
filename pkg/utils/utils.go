package utils

import (
	"encoding/json"
	"io"
	"net/http"
)

// ParseBody accepts the request body and unmarshals it into the provided interface
// the reason why interface is used is to allow any type to be passed in
// You must pass a pointer (&user)
func ParseBody(r *http.Request, x interface{}) {
	if body, err := io.ReadAll(r.Body); err == nil {
		//Take the raw JSON and turn it into a Go struct
		if err := json.Unmarshal([]byte(body), x); err != nil {
			return
		}
	}
}
