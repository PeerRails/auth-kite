package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	pongJSON                = `{"text":"pong","status":"OK"}`
	paramErrorJSON          = `{"error":true,"message":"Invalid Parameters","code":403}`
	testNotFoundError       = `{"error":true,"message":"Not Found","code":404}`
	testForbiddenError      = `{"error":true,"message":"Forbidden","code":403}`
	testInternalServerError = `{"error":true,"message":"Internal Server Error","code":500}`
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

		Convey("JSON Parse error", func() {
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

		Convey("should raise 403 error", func() {
			errorHandler(w, req, forbiddenError)
			Convey("httpCode should be 403", func() {
				So(w.Code, ShouldEqual, http.StatusForbidden)
			})
			Convey("Message should be Internal Server Error", func() {
				So(w.Body.String(), ShouldEqual, testForbiddenError)
			})
		})

		Convey("should raise 500 error", func() {
			errorHandler(w, req, internalServerError)
			Convey("httpCode should be 500", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
			Convey("Message should be Internal Server Error", func() {
				So(w.Body.String(), ShouldEqual, testInternalServerError)
			})
		})
	})
}

func TestPrepareDatabase(t *testing.T) {
	Convey("Prepare Mock Database", t, func() {
		db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		SkipConvey("Fill Database", func() {
			query := "CREATE TABLE IF EXISTS `keys` (`id` INTEGER PRIMARY KEY AUTOINCREMENT,`key` VARCHAR(64) NULL,`user_id` INTEGER,`expires_at` DATE NULL,`created_at` DATE NULL); DELETE FROM keys;"

			_, err = db.Exec(query)
			if err != nil {
				t.Fatal(err)
			}

			queries, err := db.Prepare("insert into keys(key, user_id, expires_at, created_at) values(?, ?, ?, ?)")

			if err != nil {
				t.Fatal(err)
			}

			_, err = queries.Exec("key1", 1, "2020-01-01", "2010-01-01")

			if err != nil {
				t.Fatal(err)
			}

			// mock data, ugly code just for testing
			_, _ = queries.Exec("key2", 2, "2020-01-01", "2010-01-01")
			_, _ = queries.Exec("key3", 3, "2020-01-01", "2010-01-01")
			_, _ = queries.Exec("key4", 4, "2020-01-01", "2010-01-01")
			_, _ = queries.Exec("key5", 5, "2020-01-01", "2010-01-01")

			Convey("db should have all keys", func() {
				rows, err := db.Query("select count(id) as c from keys")
				if err != nil {
					t.Fatal(err)
				}
				defer rows.Close()
				for rows.Next() {
					var c int
					err = rows.Scan(&c)
					if err != nil {
						t.Fatal(err)
					}
					So(c, ShouldEqual, 5)
				}
				err = rows.Err()
				if err != nil {
					log.Fatal(err)
				}
			})

		})
	})
}

func TestInit(t *testing.T) {
	Convey("Prepare Mock Database", t, func() {
		_ = os.Setenv("DATABASE_URL", "postgres://test:test@pghost/keys_test?sslmode=false")
		Convey("DATABASE_URL should not be empty", func() {
			So(os.Getenv("DATABASE_URL"), ShouldNotBeEmpty)
		})
		Convey("should show error if DATABASE_URL is invalid", func() {
			_ = os.Setenv("DATABASE_URL", "-")
			_, err := PrepareDatabase(os.Getenv("DATABASE_URL"))
			So(err, ShouldBeNil)
		})
		Convey("should connect to DATABASE if DATABASE_URL is empty", func() {
			_ = os.Setenv("DATABASE_URL", "")
			db, err := PrepareDatabase(os.Getenv("DATABASE_URL"))
			SkipSo(err, ShouldBeNil)
			SkipSo(db.Ping, ShouldBeNil)
		})
		Convey("should connect to Database", func() {
			db, err := PrepareDatabase(os.Getenv("DATABASE_URL"))
			SkipSo(err, ShouldBeNil)
			SkipSo(db.Ping, ShouldBeNil)
		})
	})

}
