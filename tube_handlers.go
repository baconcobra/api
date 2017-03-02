/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Tube struct {
	gorm.Model
	Name string `json:"name"`
	URL  string `json:"url"`
}

type GetTubes struct {
	Nav   Navigation `json:"nav"`
	Tubes []Tube     `json:"tubes"`
}

var TubesGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	limit := 100
	page := 0
	tubes := []Tube{}
	db.Limit(limit).Find(&tubes).Offset(page * limit)
	nav := getNavigation(len(tubes), page, limit)

	response, _ := json.Marshal(GetTubes{Nav: nav, Tubes: tubes})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
})

var TubeGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var tube Tube
	getTube(r, &tube)

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(tube)
	w.Write([]byte(response))
})

var TubesPostHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var t Tube
	mapBody(r, &t)
	db.Create(&t)
	response, _ := json.Marshal(t)
	w.Write([]byte(response))
})

var TubesPatchHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var tube Tube
	var updatedTube Tube
	getTube(r, &tube)
	mapBody(r, &updatedTube)

	tube.Name = updatedTube.Name
	tube.URL = updatedTube.URL

	db.Save(&tube)
	response, _ := json.Marshal(tube)
	w.Write([]byte(response))
})

var TubeDeleteHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var tube Tube
	getTube(r, &tube)
	db.Delete(&tube)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(""))
})

func getTube(r *http.Request, tube *Tube) {
	vars := mux.Vars(r)
	db.Find(&tube, vars["id"])
}

func mapBody(r *http.Request, t *Tube) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)
	if err != nil {
		log.Println("Invalid input")

	}
}
