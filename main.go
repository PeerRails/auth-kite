package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/op/go-logging"
	"net/http"
	"os"
)

var log = logging.MustGetLogger("main.log")
var Db *sql.DB

//var format = logging.MustStringFormatter(
//	`%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
//)

type Pong struct {
	Text   string `json:"text"`
	Status string `json:"status"`
}

type Key struct {
	Key     string `json:"key"`
	Expired bool   `json:"expired"`
}

type ErrorMessage struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var (
	internalServerError = &ErrorMessage{Error: true, Message: "Internal Server Error", Code: http.StatusInternalServerError}
	notFoundError       = &ErrorMessage{Error: true, Message: "Not Found", Code: http.StatusNotFound}
	forbiddenError      = &ErrorMessage{Error: true, Message: "Forbidden", Code: http.StatusForbidden}
	invalidParamError   = &ErrorMessage{Error: true, Message: "Invalid Parameters", Code: http.StatusForbidden}
)

func PrepareDatabase(database_url string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", database_url)
	if err != nil {
		log.Fatal(err)
		return db, err
	}
	defer db.Close()
	return db, err
}

func init() {
	database_url := os.Getenv("DATABASE_URL")
	if database_url == "" {
		database_url = "postgres://test:test@pghost/keys_test?sslmode=verify-full"
	}
	var err error
	Db, err = PrepareDatabase(database_url)
	if err != nil {
		panic(err)
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Request Ping")
	w.Header().Set("Content-Type", "application/json")
	pong, _ := json.Marshal(&Pong{Text: "pong", Status: "OK"})

	w.Header().Set("Content-Type", "application/json")
	w.Write(pong)
}

func authKeyHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Request Auth Key")
	w.Header().Set("Content-Type", "application/json")
	keyParam := r.URL.Query().Get("key")

	if keyParam == "" {
		errorHandler(w, r, invalidParamError)
		return
	}

	key, _ := json.Marshal(&Key{Key: keyParam, Expired: false})
	/*
		if err != nil {
			errorHandler(w, r, internalServerError)
			return
		}
	*/

	w.Header().Set("Content-Type", "application/json")
	w.Write(key)
}

func errorHandler(w http.ResponseWriter, r *http.Request, err *ErrorMessage) {
	w.WriteHeader(err.Code)
	error_json, _ := json.Marshal(err)
	w.Write(error_json)
}

func main() {
	http.HandleFunc("/", pingHandler)
	http.HandleFunc("/auth", authKeyHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
	query := "CREATE TABLE IF EXISTS `keys` (`id` INTEGER PRIMARY KEY AUTOINCREMENT,`key` VARCHAR(64) NULL,`user_id` INTEGER,`expires_at` DATE NULL,`created_at` DATE NULL);"
	_, _ = Db.Exec(query)
}
