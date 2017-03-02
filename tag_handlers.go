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

type Tag struct {
	gorm.Model
	Name string `json:"name"`
}

type GetTags struct {
	Nav  Navigation `json:"nav"`
	Tags []Tag      `json:"tags"`
}

var TagsGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	limit := 100
	page := 0
	tags := []Tag{}
	db.Limit(limit).Find(&tags).Offset(page * limit)
	nav := getNavigation(len(tags), page, limit)

	response, _ := json.Marshal(GetTags{Nav: nav, Tags: tags})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
})

var TagGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var tag Tag
	getTag(r, &tag)

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(tag)
	w.Write([]byte(response))
})

var TagsPostHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var t Tag
	mapTag(r, &t)
	db.Create(&t)
	response, _ := json.Marshal(t)
	w.Write([]byte(response))
})

var TagsPatchHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var tag Tag
	var updatedTag Tag
	getTag(r, &tag)
	mapTag(r, &updatedTag)

	tag.Name = updatedTag.Name

	db.Save(&tag)
	response, _ := json.Marshal(tag)
	w.Write([]byte(response))
})

var TagDeleteHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var tag Tag
	getTag(r, &tag)
	db.Delete(&tag)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(""))
})

func getTag(r *http.Request, tag *Tag) {
	vars := mux.Vars(r)
	db.Find(&tag, vars["id"])
}

func mapTag(r *http.Request, t *Tag) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)
	if err != nil {
		log.Println("Invalid input")

	}
}
