package main

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Error bool   `json:"error"`
	Msg   string `json:"msg"`
	Data  any    `json:"data,omitempty"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error: false,
		Msg:   "Broker service is running",
	}

	out, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)
}
