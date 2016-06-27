package auth_kite

import (
	"encoding/json"
	"github.com/op/go-logging"
	"net/http"
)

var log = logging.MustGetLogger("main.log")

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

/*
func main() {
	http.HandleFunc("/", pingHandler)
	http.HandleFunc("/auth", authKeyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
*/
