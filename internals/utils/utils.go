package utils

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, status int, data any) error {
	js, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return err
	}

	js = append(js, '\n')
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}
