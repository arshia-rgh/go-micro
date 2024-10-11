package main

import (
	"logger/data"
	"net/http"
)

type JSONPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) WriteLog(w http.ResponseWriter, r *http.Request) {
	var jsonPayload JSONPayload

	_ = app.readJSON(w, r, &jsonPayload)

	event := data.LogEntry{
		Name: jsonPayload.Name,
		Data: jsonPayload.Data,
	}

	ID, err := event.Insert()
	if err != nil {
		_ = app.errorJSON(w, err)
		return
	}

	resp := jsonResponse{
		Error:   false,
		Message: "logged",
		Data:    ID,
	}

	_ = app.writeJSON(w, http.StatusCreated, resp)
}
