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

func TestGetUsers(t *testing.T) {
	Convey("Given no users on the database", t, func() {
		setupTestSuite()
		Convey("When I call GET /users", func() {

			response := doRequest("GET", "/users", nil)

			Convey("Then I should get an empty list of users", func() {
				body := response.Body.String()
				gt := GetUsers{}
				json.Unmarshal([]byte(body), &gt)
				So(len(gt.Users), ShouldEqual, 0)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})

	Convey("Given 250 users on the database", t, func() {
		setupTestSuite()
		createUsers(250)
		Convey("When I call GET /users", func() {

			response := doRequest("GET", "/users", nil)

			Convey("Then I should get a list of 100 users", func() {
				body := response.Body.String()
				gt := GetUsers{}
				json.Unmarshal([]byte(body), &gt)
				So(len(gt.Users), ShouldEqual, 100)
				So(len(gt.Nav.Next), ShouldEqual, 1)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestPostUsers(t *testing.T) {
	Convey("Given no users on the database", t, func() {
		setupTestSuite()
		Convey("When I call POST /users", func() {
			t := User{Name: "foo"}
			data, err := json.Marshal(t)
			if err != nil {
				panic(err.Error())
			}
			body := bytes.NewBuffer(data)
			response := doRequest("POST", "/users", body)

			Convey("Then a new user should be added to db", func() {
				users := []User{}
				db.Find(&users)
				So(len(users), ShouldEqual, 1)
			})

			Convey("Then I should get an empty list of users", func() {
				body := response.Body.String()
				user := User{}
				json.Unmarshal([]byte(body), &user)
				So(t.Name, ShouldEqual, user.Name)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestGetUser(t *testing.T) {
	Convey("Given a user exists on the db", t, func() {
		setupTestSuite()
		t := User{Name: "test"}
		db.Create(&t)
		db.First(&t)
		Convey("When I call GET /user/{id}", func() {
			id := fmt.Sprint(t.ID)
			response := doRequest("GET", "/users/"+id, nil)

			Convey("Then I should get the user details", func() {
				body := response.Body.String()
				user := User{}
				json.Unmarshal([]byte(body), &user)
				So(t.Name, ShouldEqual, user.Name)
				So(t.ID, ShouldEqual, user.ID)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestPatchUser(t *testing.T) {
	Convey("Given a user exists on the db", t, func() {
		setupTestSuite()
		t := User{Name: "test"}
		db.Create(&t)
		db.First(&t)
		Convey("When I call PATCH /user/{id}", func() {
			id := fmt.Sprint(t.ID)
			t.Name = "modified"

			data, err := json.Marshal(t)
			if err != nil {
				panic(err.Error())
			}
			body := bytes.NewBuffer(data)
			response := doRequest("PATCH", "/users/"+id, body)

			Convey("Then I should get the user details", func() {
				body := response.Body.String()
				user := User{}
				json.Unmarshal([]byte(body), &user)
				So(user.Name, ShouldEqual, t.Name)
				So(user.ID, ShouldEqual, t.ID)
			})

			Convey("Then it should be stored on db", func() {
				user := User{}
				db.Find(&user, id)
				So(user.Name, ShouldEqual, t.Name)
				So(user.ID, ShouldEqual, t.ID)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}

func TestDeleteUser(t *testing.T) {
	Convey("Given a user exists on the db", t, func() {
		setupTestSuite()
		t := User{Name: "test"}
		db.Create(&t)
		db.First(&t)
		Convey("When I call DELETE /user/{id}", func() {
			id := fmt.Sprint(t.ID)
			response := doRequest("DELETE", "/users/"+id, nil)

			Convey("Then it should not exist anymore", func() {
				user := User{}
				db.Find(&user, id)
				So(user.ID, ShouldEqual, 0)
			})

			Convey("And I should get a 200 response", func() {
				status := response.Code
				So(status, ShouldEqual, 200)
			})
		})
	})
}
