/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"

	"golang.org/x/crypto/scrypt"
)

var w http.ResponseWriter
var r http.Request

func setupTestSuite() {
	setupDB("sqlite3", "test.db")
	db.Where("1 LIKE 1").Delete(Tube{})
	db.Where("1 LIKE 1").Delete(Tag{})
	db.Where("1 LIKE 1").Delete(Actor{})
	db.Where("1 LIKE 1").Delete(Video{})
	db.Where("1 LIKE 1").Delete(User{})
}

func doRequest(verb string, route string, body io.Reader) *httptest.ResponseRecorder {
	m := setupRouter()
	user := User{Name: "me"}
	token := getToken(user)
	request, _ := http.NewRequest(verb, route, body)
	request.Header.Set("Authorization", "Bearer "+string(token))
	request.Header.Set("X-Auth-Token", "test-token")
	response := httptest.NewRecorder()
	m.ServeHTTP(response, request)

	return response
}

func createTubes(n int) {
	i := 0
	for i < n {
		x := strconv.Itoa(i)
		db.Create(&Tube{
			Name: "Test" + x,
			URL:  "http://www.name" + x + ".com"})
		i += 1
	}
}

func createTags(n int) {
	i := 0
	for i < n {
		x := strconv.Itoa(i)
		db.Create(&Tag{
			Name: "Test" + x,
		})
		i += 1
	}
}

func createActors(n int) {
	i := 0
	for i < n {
		x := strconv.Itoa(i)
		db.Create(&Actor{
			Name: "Test" + x,
		})
		i += 1
	}
}

func createVideos(n int) {
	i := 0
	for i < n {
		x := strconv.Itoa(i)
		db.Create(&Video{
			Title: "Test" + x,
		})
		i += 1
	}
}

func createUsers(n int) {
	i := 0

	salt := make([]byte, SaltSize)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		// Raise error
	}

	hash, err := scrypt.Key([]byte("testpwd"), salt, 16384, 8, 1, HashSize)
	if err != nil {
		// Raise error
	}

	for i < n {
		x := strconv.Itoa(i)
		db.Create(&User{
			Name:     "Test" + x,
			UserName: "test" + x,
			Salt:     base64.StdEncoding.EncodeToString(salt),
			Password: base64.StdEncoding.EncodeToString(hash),
		})
		i += 1
	}
}
