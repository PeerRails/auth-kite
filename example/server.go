package main

import (
	"github.com/PeerRails/auth-kite"
	"net/http"
	"os"
)

func main() {
	//postgres://dev:dev@localhost/omckonrails-dev?sslmode=verify-full
	Db, _ := authkite.PrepareDatabase(os.Getenv("DATABASE_URL"))
	_ = Db.Ping
	http.HandleFunc("/", authkite.PingHandler)
	http.HandleFunc("/auth", authkite.AuthKeyHandler)
	http.ListenAndServe(":3000", nil)
}
