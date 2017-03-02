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

type Video struct {
	gorm.Model
	Title        string     `json:"title"`
	URL          string     `json:"url"`
	ExtID        string     `json:"extid"`
	Duration     string     `json:"duration"`
	Rating       int        `json:"rating"`
	Embed        string     `json:"embed"`
	SmallImages  string     `json:"small_images"`
	MediumImages string     `json:"medium_images"`
	BigImages    string     `json:"big_images"`
	Uuid         string     `json:"uuid"`
	Views        int        `json:"views"`
	MasterImage  string     `json:"master_image"`
	Sexuality    string     `json:"sexuality"`
	Tags         []Tag      `json:"tags" gorm:"many2many:video_tags;"`
	Actors       []Actor    `json:"actors" gorm:"many2many:video_actors;"`
	Tube         Tube       `json:"tube"`
	Uploaded     *time.Time `json:"uploaded"`
}

type GetVideos struct {
	Nav    Navigation `json:"nav"`
	Videos []Video    `json:"videos"`
}

var VideosGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	limit := 100
	page := 0
	videos := []Video{}
	db.Limit(limit).Find(&videos).Offset(page * limit)
	nav := getNavigation(len(videos), page, limit)

	response, _ := json.Marshal(GetVideos{Nav: nav, Videos: videos})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
})

var VideoGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var video Video
	getVideo(r, &video)

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(video)
	w.Write([]byte(response))
})

var VideosPostHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var t Video
	mapVideo(r, &t)
	db.Create(&t)
	response, _ := json.Marshal(t)
	w.Write([]byte(response))
})

var VideosPatchHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var video Video
	var updatedVideo Video
	getVideo(r, &video)
	mapVideo(r, &updatedVideo)

	video.Title = updatedVideo.Title

	db.Save(&video)
	response, _ := json.Marshal(video)
	w.Write([]byte(response))
})

var VideoDeleteHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var video Video
	getVideo(r, &video)
	db.Delete(&video)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(""))
})

func getVideo(r *http.Request, video *Video) {
	vars := mux.Vars(r)
	db.Find(&video, vars["id"])
}

func mapVideo(r *http.Request, t *Video) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)
	if err != nil {
		log.Println("Invalid input")

	}
}
