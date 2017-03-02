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

func TestGetTags(t *testing.T) {
	Convey("Given no tags on the database", t, func() {
		setupTestSuite()
		Convey("When I call GET /tags", func() {

			response := doRequest("GET", "/tags", nil)

			Convey("Then I should get an empty list of tags", func() {
				body := response.Body.String()
				gt := GetTags{}
				json.Unmarshal([]byte(body), &gt)
				So(len(gt.Tags), ShouldEqual, 0)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})

	Convey("Given 250 tags on the database", t, func() {
		setupTestSuite()
		createTags(250)
		Convey("When I call GET /tags", func() {

			response := doRequest("GET", "/tags", nil)

			Convey("Then I should get a list of 100 tags", func() {
				body := response.Body.String()
				gt := GetTags{}
				json.Unmarshal([]byte(body), &gt)
				So(len(gt.Tags), ShouldEqual, 100)
				So(len(gt.Nav.Next), ShouldEqual, 1)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestPostTags(t *testing.T) {
	Convey("Given no tags on the database", t, func() {
		setupTestSuite()
		Convey("When I call POST /tags", func() {
			t := Tag{Name: "foo"}
			data, err := json.Marshal(t)
			if err != nil {
				panic(err.Error())
			}
			body := bytes.NewBuffer(data)
			response := doRequest("POST", "/tags", body)

			Convey("Then a new tag should be added to db", func() {
				tags := []Tag{}
				db.Find(&tags)
				So(len(tags), ShouldEqual, 1)
			})

			Convey("Then I should get an empty list of tags", func() {
				body := response.Body.String()
				tag := Tag{}
				json.Unmarshal([]byte(body), &tag)
				So(t.Name, ShouldEqual, tag.Name)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestGetTag(t *testing.T) {
	Convey("Given a tag exists on the db", t, func() {
		setupTestSuite()
		t := Tag{Name: "test"}
		db.Create(&t)
		db.First(&t)
		Convey("When I call GET /tag/{id}", func() {
			id := fmt.Sprint(t.ID)
			response := doRequest("GET", "/tags/"+id, nil)

			Convey("Then I should get the tag details", func() {
				body := response.Body.String()
				tag := Tag{}
				json.Unmarshal([]byte(body), &tag)
				So(t.Name, ShouldEqual, tag.Name)
				So(t.ID, ShouldEqual, tag.ID)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestPatchTag(t *testing.T) {
	Convey("Given a tag exists on the db", t, func() {
		setupTestSuite()
		t := Tag{Name: "test"}
		db.Create(&t)
		db.First(&t)
		Convey("When I call PATCH /tag/{id}", func() {
			id := fmt.Sprint(t.ID)
			t.Name = "modified"

			data, err := json.Marshal(t)
			if err != nil {
				panic(err.Error())
			}
			body := bytes.NewBuffer(data)
			response := doRequest("PATCH", "/tags/"+id, body)

			Convey("Then I should get the tag details", func() {
				body := response.Body.String()
				tag := Tag{}
				json.Unmarshal([]byte(body), &tag)
				So(tag.Name, ShouldEqual, t.Name)
				So(tag.ID, ShouldEqual, t.ID)
			})

			Convey("Then it should be stored on db", func() {
				tag := Tag{}
				db.Find(&tag, id)
				So(tag.Name, ShouldEqual, t.Name)
				So(tag.ID, ShouldEqual, t.ID)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestDeleteTag(t *testing.T) {
	Convey("Given a tag exists on the db", t, func() {
		setupTestSuite()
		t := Tag{Name: "test"}
		db.Create(&t)
		db.First(&t)
		Convey("When I call DELETE /tag/{id}", func() {
			id := fmt.Sprint(t.ID)
			response := doRequest("DELETE", "/tags/"+id, nil)

			Convey("Then it should not exist anymore", func() {
				tag := Tag{}
				db.Find(&tag, id)
				So(tag.ID, ShouldEqual, 0)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}
