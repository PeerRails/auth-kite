package main

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	pongJSON          = `{"text":"pong","status":"OK"}`
	paramErrorJSON    = `{"error":true,"message":"Invalid Parameters","code":403}`
	keyQuery          = "key=keykeykey"
	testNotFoundError = `{"error":true,"message":"Not Found","code":404}`
)

func TestPing(t *testing.T) {
	Convey("Start request", t, func() {
		resp := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}
		http.DefaultServeMux.ServeHTTP(resp, req)
		w := httptest.NewRecorder()

		pingHandler(w, req)

		Convey("home should return 200", func() {
			So(w.Code, ShouldEqual, http.StatusOK)
		})

		Convey("home should return OK", func() {
			So(w.Body.String(), ShouldEqual, pongJSON)
		})

		fmt.Printf("%d - %s", w.Code, w.Body.String())
	})
}

func TestAuthKey(t *testing.T) {
	Convey("Request", t, func() {
		resp := httptest.NewRecorder()
		Convey("with no parameters", func() {
			req, err := http.NewRequest("GET", "/auth", nil)
			if err != nil {
				t.Fatal(err)
			}
			http.DefaultServeMux.ServeHTTP(resp, req)
			w := httptest.NewRecorder()
			authKeyHandler(w, req)

			Convey("should return error", func() {
				So(w.Code, ShouldEqual, http.StatusForbidden)
				So(w.Body.String(), ShouldEqual, paramErrorJSON)
				fmt.Printf("%d - %s", w.Code, w.Body.String())
			})
		})
		Convey("with right params", func() {
			req, err := http.NewRequest("GET", "/auth?key=keykeykey", nil)
			if err != nil {
				t.Fatal(err)
			}
			http.DefaultServeMux.ServeHTTP(resp, req)
			w := httptest.NewRecorder()
			authKeyHandler(w, req)

			Convey("should 200", func() {
				So(w.Code, ShouldEqual, http.StatusOK)
				fmt.Printf("%d - %s", w.Code, w.Body.String())
			})
		})
	})
}

func TestErrorHandler(t *testing.T) {
	Convey("Error", t, func() {
		resp := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/error", nil)
		http.DefaultServeMux.ServeHTTP(resp, req)
		w := httptest.NewRecorder()
		Convey("should raise 404 error", func() {
			errorHandler(w, req, notFoundError)
			Convey("httpCode should be 404", func() {
				So(w.Code, ShouldEqual, http.StatusNotFound)
			})
			Convey("Message should be Not Found", func() {
				So(w.Body.String(), ShouldEqual, testNotFoundError)
			})
		})
	})
}
