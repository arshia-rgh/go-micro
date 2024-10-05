package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	const maxBytes = 1048576 // 1 Megabyte

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)

	err := dec.Decode(data)

	if err != nil {
		return err
	}

	// check if there is two JSON values
	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("body must have only a single JSON value")
	}

	return nil
}
