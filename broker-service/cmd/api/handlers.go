package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error: false,
		Msg:   "Broker service is running",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var data RequestPayload

	err := app.readJSON(w, r, &data)
	if err != nil {
		app.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	switch data.Action {
	case "auth":
		app.Authenticate(w, data.Auth)
	default:
		app.errorJSON(w, errors.New("invalid action"), http.StatusBadRequest)
	}
}

func (app *Config) Authenticate(w http.ResponseWriter, a AuthPayload) {
	//create json to sen to auth microsrvc
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	//call service
	request, err := http.NewRequest("POST", "http://localhost:8081/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errorJSON(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJSON(w, errors.New("error calling authentication service"), http.StatusInternalServerError)
		return
	}

	var jsonFromService jsonResponse

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, errors.New(jsonFromService.Msg), http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Msg = "Auth successful"
	payload.Data = jsonFromService.Data

	err = app.writeJSON(w, http.StatusAccepted, payload)
	if err != nil {
		app.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

}
