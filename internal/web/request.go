package web

import (
	"encoding/json"
	"net/http"
)

func DecodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
