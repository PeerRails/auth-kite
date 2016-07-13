package authkite

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"github.com/op/go-logging"
	"net/http"
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
	return db, err
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Request Ping")
	w.Header().Set("Content-Type", "application/json")
	pong, _ := json.Marshal(&Pong{Text: "pong", Status: "OK"})

	w.Header().Set("Content-Type", "application/json")
	w.Write(pong)
}

func AuthKeyHandler(w http.ResponseWriter, r *http.Request) {
	log.Info("Request Auth Key")
	w.Header().Set("Content-Type", "application/json")
	keyParam := r.URL.Query().Get("key")
	log.Info("Key: ?", keyParam)

	if keyParam == "" {
		ErrorHandler(w, r, invalidParamError)
		return
	}

	var keyjson []byte
	var key string
	err := Db.QueryRow("SELECT key FROM keys WHERE key = '?'", keyParam).Scan(&key)
	if err == nil {
		ErrorHandler(w, r, invalidParamError)
		return
	} else {
		if key != "" {
			keyjson, _ = json.Marshal(&Key{Key: key, Expired: false})
			w.Write(keyjson)
		} else {
			ErrorHandler(w, r, invalidParamError)
			return
		}
	}
}

func ErrorHandler(w http.ResponseWriter, r *http.Request, err *ErrorMessage) {
	w.WriteHeader(err.Code)
	error_json, _ := json.Marshal(err)
	w.Write(error_json)
}
