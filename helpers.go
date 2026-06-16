package main

import (
	"encoding/json"
	"net/http"
)

type envelope map[string]any

func (app *application) readJSON(data []byte, dst any) error {
	err := json.Unmarshal(data, dst)

	if err != nil {
		return err
	}

	return nil
}

func (app *application) writeJSON(w http.ResponseWriter, r *http.Request, status int, data envelope) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
