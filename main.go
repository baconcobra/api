/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/lib/pq"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var err error

func main() {
	if os.Getenv("DATABASE_URL") != "" {
		url := os.Getenv("DATABASE_URL")
		connection, _ := pq.ParseURL(url)
		connection += " sslmode=require"

		setupDB("postgres", connection)
	} else {
		setupDB("sqlite3", "dev.db")
	}
	r := setupRouter()
	http.ListenAndServe(":"+os.Getenv("PORT"), handlers.LoggingHandler(os.Stdout, r))
}

var NotImplemented = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Not Implemented"))
})

var StatusHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API is up and running"))
})
