package main

import (
	"net/http"
)

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error: false,
		Msg:   "Broker service is running",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}
