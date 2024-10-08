package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type RequestPayload struct {
	Action string `json:"action"`
	Auth   AuthPayload
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)

}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)

	default:
		app.errorJSON(w, errors.New("unknown action"))
	}

}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.Marshal(a)

	req, err := http.NewRequest("POST", "https://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}

	res, err := client.Do(req)

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("invalid credentials"))
		return
	}

	var jsonFromAuth jsonResponse

	err = json.NewDecoder(res.Body).Decode(&jsonFromAuth)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if jsonFromAuth.Error {
		app.errorJSON(w, errors.New(fmt.Sprint(jsonFromAuth.Message)), http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Authenticated !",
		Data:    jsonFromAuth.Data,
	}

	app.writeJSON(w, http.StatusAccepted, payload)
}
