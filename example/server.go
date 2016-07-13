package main

import (
	"github.com/PeerRails/auth-kite"
	"net/http"
)

func main() {
	//postgres://dev:dev@localhost/omckonrails-dev?sslmode=disable
	http.HandleFunc("/", authkite.PingHandler)
	http.HandleFunc("/auth", authkite.AuthKeyHandler)
	http.ListenAndServe(":3000", nil)
}
