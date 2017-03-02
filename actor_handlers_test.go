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

func TestGetActors(t *testing.T) {
	Convey("Given no actors on the database", t, func() {
		setupTestSuite()
		Convey("When I call GET /actors", func() {

			response := doRequest("GET", "/actors", nil)

			Convey("Then I should get an empty list of actors", func() {
				body := response.Body.String()
				gt := GetActors{}
				json.Unmarshal([]byte(body), &gt)
				So(len(gt.Actors), ShouldEqual, 0)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})

	Convey("Given 250 actors on the database", t, func() {
		setupTestSuite()
		createActors(250)
		Convey("When I call GET /actors", func() {

			response := doRequest("GET", "/actors", nil)

			Convey("Then I should get a list of 100 actors", func() {
				body := response.Body.String()
				gt := GetActors{}
				json.Unmarshal([]byte(body), &gt)
				So(len(gt.Actors), ShouldEqual, 100)
				So(len(gt.Nav.Next), ShouldEqual, 1)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestPostActors(t *testing.T) {
	Convey("Given no actors on the database", t, func() {
		setupTestSuite()
		Convey("When I call POST /actors", func() {
			t := Actor{Name: "foo"}
			data, err := json.Marshal(t)
			if err != nil {
				panic(err.Error())
			}
			body := bytes.NewBuffer(data)
			response := doRequest("POST", "/actors", body)

			Convey("Then a new actor should be added to db", func() {
				actors := []Actor{}
				db.Find(&actors)
				So(len(actors), ShouldEqual, 1)
			})

			Convey("Then I should get an empty list of actors", func() {
				body := response.Body.String()
				actor := Actor{}
				json.Unmarshal([]byte(body), &actor)
				So(t.Name, ShouldEqual, actor.Name)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestGetActor(t *testing.T) {
	Convey("Given a actor exists on the db", t, func() {
		setupTestSuite()
		t := Actor{Name: "test"}
		db.Create(&t)
		db.First(&t)
		Convey("When I call GET /actor/{id}", func() {
			id := fmt.Sprint(t.ID)
			response := doRequest("GET", "/actors/"+id, nil)

			Convey("Then I should get the actor details", func() {
				body := response.Body.String()
				actor := Actor{}
				json.Unmarshal([]byte(body), &actor)
				So(t.Name, ShouldEqual, actor.Name)
				So(t.ID, ShouldEqual, actor.ID)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestPatchActor(t *testing.T) {
	Convey("Given a actor exists on the db", t, func() {
		setupTestSuite()
		t := Actor{Name: "test"}
		db.Create(&t)
		db.First(&t)
		Convey("When I call PATCH /actor/{id}", func() {
			id := fmt.Sprint(t.ID)
			t.Name = "modified"

			data, err := json.Marshal(t)
			if err != nil {
				panic(err.Error())
			}
			body := bytes.NewBuffer(data)
			response := doRequest("PATCH", "/actors/"+id, body)

			Convey("Then I should get the actor details", func() {
				body := response.Body.String()
				actor := Actor{}
				json.Unmarshal([]byte(body), &actor)
				So(actor.Name, ShouldEqual, t.Name)
				So(actor.ID, ShouldEqual, t.ID)
			})

			Convey("Then it should be stored on db", func() {
				actor := Actor{}
				db.Find(&actor, id)
				So(actor.Name, ShouldEqual, t.Name)
				So(actor.ID, ShouldEqual, t.ID)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestDeleteActor(t *testing.T) {
	Convey("Given a actor exists on the db", t, func() {
		setupTestSuite()
		t := Actor{Name: "test"}
		db.Create(&t)
		db.First(&t)
		Convey("When I call DELETE /actor/{id}", func() {
			id := fmt.Sprint(t.ID)
			response := doRequest("DELETE", "/actors/"+id, nil)

			Convey("Then it should not exist anymore", func() {
				actor := Actor{}
				db.Find(&actor, id)
				So(actor.ID, ShouldEqual, 0)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}
