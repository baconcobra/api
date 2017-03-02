/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetVideos(t *testing.T) {
	Convey("Given no videos on the database", t, func() {
		setupTestSuite()
		Convey("When I call GET /videos", func() {

			response := doRequest("GET", "/videos", nil)

			Convey("Then I should get an empty list of videos", func() {
				body := response.Body.String()
				gt := GetVideos{}
				json.Unmarshal([]byte(body), &gt)
				So(len(gt.Videos), ShouldEqual, 0)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})

	Convey("Given 250 videos on the database", t, func() {
		setupTestSuite()
		createVideos(250)
		Convey("When I call GET /videos", func() {

			response := doRequest("GET", "/videos", nil)

			Convey("Then I should get a list of 100 videos", func() {
				body := response.Body.String()
				gt := GetVideos{}
				json.Unmarshal([]byte(body), &gt)
				So(len(gt.Videos), ShouldEqual, 100)
				So(len(gt.Nav.Next), ShouldEqual, 1)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestPostVideos(t *testing.T) {
	Convey("Given no videos on the database", t, func() {
		setupTestSuite()
		Convey("When I call POST /videos", func() {
			t := Video{Title: "foo"}
			data, err := json.Marshal(t)
			if err != nil {
				panic(err.Error())
			}
			body := bytes.NewBuffer(data)
			response := doRequest("POST", "/videos", body)

			Convey("Then a new video should be added to db", func() {
				videos := []Video{}
				db.Find(&videos)
				So(len(videos), ShouldEqual, 1)
			})

			Convey("Then I should get an empty list of videos", func() {
				body := response.Body.String()
				video := Video{}
				json.Unmarshal([]byte(body), &video)
				So(t.Title, ShouldEqual, video.Title)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestGetVideo(t *testing.T) {
	Convey("Given a video exists on the db", t, func() {
		setupTestSuite()
		t := Video{Title: "test"}
		db.Create(&t)
		db.First(&t)
		Convey("When I call GET /video/{id}", func() {
			id := fmt.Sprint(t.ID)
			response := doRequest("GET", "/videos/"+id, nil)

			Convey("Then I should get the video details", func() {
				body := response.Body.String()
				video := Video{}
				json.Unmarshal([]byte(body), &video)
				So(t.Title, ShouldEqual, video.Title)
				So(t.ID, ShouldEqual, video.ID)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestPatchVideo(t *testing.T) {
	Convey("Given a video exists on the db", t, func() {
		setupTestSuite()
		t := Video{Title: "test"}
		db.Create(&t)
		db.First(&t)
		Convey("When I call PATCH /video/{id}", func() {
			id := fmt.Sprint(t.ID)
			t.Title = "modified"

			data, err := json.Marshal(t)
			if err != nil {
				panic(err.Error())
			}
			body := bytes.NewBuffer(data)
			response := doRequest("PATCH", "/videos/"+id, body)

			Convey("Then I should get the video details", func() {
				body := response.Body.String()
				video := Video{}
				json.Unmarshal([]byte(body), &video)
				So(video.Title, ShouldEqual, t.Title)
				So(video.ID, ShouldEqual, t.ID)
			})

			Convey("Then it should be stored on db", func() {
				video := Video{}
				db.Find(&video, id)
				So(video.Title, ShouldEqual, t.Title)
				So(video.ID, ShouldEqual, t.ID)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestDeleteVideo(t *testing.T) {
	Convey("Given a video exists on the db", t, func() {
		setupTestSuite()
		t := Video{Title: "test"}
		db.Create(&t)
		db.First(&t)
		Convey("When I call DELETE /video/{id}", func() {
			id := fmt.Sprint(t.ID)
			response := doRequest("DELETE", "/videos/"+id, nil)

			Convey("Then it should not exist anymore", func() {
				video := Video{}
				db.Find(&video, id)
				So(video.ID, ShouldEqual, 0)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestExtendedPostVideo(t *testing.T) {
	Convey("Given no videos on the database", t, func() {
		setupTestSuite()
		Convey("When I call POST /videos with extended data", func() {
			tags := make([]Tag, 1)
			tags[0] = Tag{Name: "foo_tag"}
			actors := make([]Actor, 1)
			actors[0] = Actor{Name: "foo_actor"}
			t := Video{Title: "foo", Tags: tags, Actors: actors}
			data, err := json.Marshal(t)
			if err != nil {
				panic(err.Error())
			}
			body := bytes.NewBuffer(data)
			response := doRequest("POST", "/videos", body)

			Convey("Then a new video should be added to db", func() {
				videos := []Video{}
				db.Find(&videos)
				So(len(videos), ShouldEqual, 1)
			})

			Convey("Then a new tag should be added to db", func() {
				tags := []Tag{}
				db.Find(&tags)
				So(len(tags), ShouldEqual, 1)
				So(tags[0].Name, ShouldEqual, "foo_tag")
			})

			Convey("Then a new actor should be added to db", func() {
				actors := []Actor{}
				db.Find(&actors)
				So(len(actors), ShouldEqual, 1)
				So(actors[0].Name, ShouldEqual, "foo_actor")
			})

			Convey("Then I should get an empty list of videos", func() {
				body := response.Body.String()
				video := Video{}
				json.Unmarshal([]byte(body), &video)
				So(t.Title, ShouldEqual, video.Title)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}
