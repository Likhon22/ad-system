package utils

import (
	"encoding/json"
	"net/http"
)

func ReadJson[T interface{}](w http.ResponseWriter, r *http.Request, data *T) error {

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}
