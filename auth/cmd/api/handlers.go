package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	type payLoad struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var payload payLoad
	err := app.readJSON(w, r, &payload)
	if err != nil {
		_ = app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	user, err := app.Repo.GetByEmail(payload.Email)

	if err != nil {
		_ = app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := app.Repo.PasswordsMatches(payload.Password, *user)
	if err != nil || !valid {
		_ = app.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	err = app.logRequest("authentication", fmt.Sprintf("%v logged in ", user.Email))

	if err != nil {
		_ = app.errorJSON(w, err)
	}

	response := jsonResponse{
		Error:   false,
		Message: "logged in",
		Data:    user,
	}

	_ = app.writeJSON(w, http.StatusAccepted, response)

}

func (app *Config) logRequest(name, data string) error {
	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}

	entry.Name = name
	entry.Data = data

	jsonData, _ := json.Marshal(entry)

	req, err := http.NewRequest("POST", "http://logger-service:8080/log", bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	_, err = app.Client.Do(req)

	return err

}
