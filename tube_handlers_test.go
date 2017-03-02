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

func TestGetTubes(t *testing.T) {
	Convey("Given no tubes on the database", t, func() {
		setupTestSuite()
		Convey("When I call GET /tubes", func() {

			response := doRequest("GET", "/tubes", nil)

			Convey("Then I should get an empty list of tubes", func() {
				body := response.Body.String()
				gt := GetTubes{}
				json.Unmarshal([]byte(body), &gt)
				So(len(gt.Tubes), ShouldEqual, 0)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})

	Convey("Given 250 tubes on the database", t, func() {
		setupTestSuite()
		createTubes(250)
		Convey("When I call GET /tubes", func() {

			response := doRequest("GET", "/tubes", nil)

			Convey("Then I should get a list of 100 tubes", func() {
				body := response.Body.String()
				gt := GetTubes{}
				json.Unmarshal([]byte(body), &gt)
				So(len(gt.Tubes), ShouldEqual, 100)
				So(len(gt.Nav.Next), ShouldEqual, 1)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestPostTubes(t *testing.T) {
	Convey("Given no tubes on the database", t, func() {
		setupTestSuite()
		Convey("When I call POST /tubes", func() {
			t := Tube{Name: "foo", URL: "http://bar.go"}
			data, err := json.Marshal(t)
			if err != nil {
				panic(err.Error())
			}
			body := bytes.NewBuffer(data)
			response := doRequest("POST", "/tubes", body)

			Convey("Then a new tube should be added to db", func() {
				tubes := []Tube{}
				db.Find(&tubes)
				So(len(tubes), ShouldEqual, 1)
			})

			Convey("Then I should get an empty list of tubes", func() {
				body := response.Body.String()
				tube := Tube{}
				json.Unmarshal([]byte(body), &tube)
				So(t.Name, ShouldEqual, tube.Name)
				So(t.URL, ShouldEqual, tube.URL)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestGetTube(t *testing.T) {
	Convey("Given a tube exists on the db", t, func() {
		setupTestSuite()
		t := Tube{Name: "test", URL: "http://test.com"}
		db.Create(&t)
		db.First(&t)
		Convey("When I call GET /tube/{id}", func() {
			id := fmt.Sprint(t.ID)
			response := doRequest("GET", "/tubes/"+id, nil)

			Convey("Then I should get the tube details", func() {
				body := response.Body.String()
				tube := Tube{}
				json.Unmarshal([]byte(body), &tube)
				So(t.Name, ShouldEqual, tube.Name)
				So(t.URL, ShouldEqual, tube.URL)
				So(t.ID, ShouldEqual, tube.ID)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestPatchTube(t *testing.T) {
	Convey("Given a tube exists on the db", t, func() {
		setupTestSuite()
		t := Tube{Name: "test", URL: "http://test.com"}
		db.Create(&t)
		db.First(&t)
		Convey("When I call PATCH /tube/{id}", func() {
			id := fmt.Sprint(t.ID)
			t.Name = "modified"
			t.URL = "http://modified.com"

			data, err := json.Marshal(t)
			if err != nil {
				panic(err.Error())
			}
			body := bytes.NewBuffer(data)
			response := doRequest("PATCH", "/tubes/"+id, body)

			Convey("Then I should get the tube details", func() {
				body := response.Body.String()
				tube := Tube{}
				json.Unmarshal([]byte(body), &tube)
				So(tube.Name, ShouldEqual, t.Name)
				So(tube.URL, ShouldEqual, t.URL)
				So(tube.ID, ShouldEqual, t.ID)
			})

			Convey("Then it should be stored on db", func() {
				tube := Tube{}
				db.Find(&tube, id)
				So(tube.Name, ShouldEqual, t.Name)
				So(tube.URL, ShouldEqual, t.URL)
				So(tube.ID, ShouldEqual, t.ID)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestDeleteTube(t *testing.T) {
	Convey("Given a tube exists on the db", t, func() {
		setupTestSuite()
		t := Tube{Name: "test", URL: "http://test.com"}
		db.Create(&t)
		db.First(&t)
		Convey("When I call DELETE /tube/{id}", func() {
			id := fmt.Sprint(t.ID)
			response := doRequest("DELETE", "/tubes/"+id, nil)

			Convey("Then it should not exist anymore", func() {
				tube := Tube{}
				db.Find(&tube, id)
				So(tube.ID, ShouldEqual, 0)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}
