package main

import (
	"encoding/json"
	"github.com/op/go-logging"
	"net/http"
)

var log = logging.MustGetLogger("main.log")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

type Pong struct {
	Text   string `json:"text"`
	Status string `json:"status"`
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Request Ping")
	w.Header().Set("Content-Type", "application/json")
	pong, err := json.Marshal(&Pong{Text: "pong", Status: "OK"})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(pong)
}

func main() {
	http.HandleFunc("/", pingHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
