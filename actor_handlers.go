/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Actor struct {
	gorm.Model
	Name        string    `json:"name"`
	Measures    string    `json:"measures"`
	Height      int       `json:"height"`
	DoB         time.Time `json:"data_of_birth"`
	Description string    `json:"description"`
	Twitter     string    `json:"twitter"`
}

type GetActors struct {
	Nav    Navigation `json:"nav"`
	Actors []Actor    `json:"actors"`
}

var ActorsGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	limit := 100
	page := 0
	actors := []Actor{}
	db.Limit(limit).Find(&actors).Offset(page * limit)
	nav := getNavigation(len(actors), page, limit)

	response, _ := json.Marshal(GetActors{Nav: nav, Actors: actors})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
})

var ActorGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var actor Actor
	getActor(r, &actor)

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(actor)
	w.Write([]byte(response))
})

var ActorsPostHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var t Actor
	mapActor(r, &t)
	db.Create(&t)
	response, _ := json.Marshal(t)
	w.Write([]byte(response))
})

var ActorsPatchHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var actor Actor
	var updatedActor Actor
	getActor(r, &actor)
	mapActor(r, &updatedActor)

	actor.Name = updatedActor.Name

	db.Save(&actor)
	response, _ := json.Marshal(actor)
	w.Write([]byte(response))
})

var ActorDeleteHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var actor Actor
	getActor(r, &actor)
	db.Delete(&actor)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(""))
})

func getActor(r *http.Request, actor *Actor) {
	vars := mux.Vars(r)
	db.Find(&actor, vars["id"])
}

func mapActor(r *http.Request, t *Actor) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)
	if err != nil {
		log.Println("Invalid input")
	}
}
