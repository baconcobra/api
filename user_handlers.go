/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/scrypt"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

const (
	// SaltSize is the size of the salt in bits
	SaltSize = 32
	// HashSize is the size of the hash in bits
	HashSize = 64
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	UserName string `json:"username"`
	Password string `json:"password"`
	Salt     string `json:"salt"`
	Role     string `json:"role"`
}

type GetUsers struct {
	Nav   Navigation `json:"nav"`
	Users []User     `json:"users"`
}

var UsersGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	limit := 100
	page := 0
	users := []User{}
	db.Limit(limit).Find(&users).Offset(page * limit)
	nav := getNavigation(len(users), page, limit)

	response, _ := json.Marshal(GetUsers{Nav: nav, Users: users})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(response))
})

var UserGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var user User
	getUser(r, &user)

	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(user)
	w.Write([]byte(response))
})

var UsersPostHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var t User
	mapUser(r, &t)

	salt := make([]byte, SaltSize)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		// Raise error
	}

	hash, err := scrypt.Key([]byte(t.Password), salt, 16384, 8, 1, HashSize)
	if err != nil {
		// Raise error
	}

	// Create a base64 string of the binary salt and hash for storage
	t.Salt = base64.StdEncoding.EncodeToString(salt)
	t.Password = base64.StdEncoding.EncodeToString(hash)

	db.Create(&t)
	response, _ := json.Marshal(t)
	w.Write([]byte(response))
})

var UsersPatchHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var user User
	var updatedUser User
	getUser(r, &user)
	mapUser(r, &updatedUser)

	user.Name = updatedUser.Name

	db.Save(&user)
	response, _ := json.Marshal(user)
	w.Write([]byte(response))
})

var UserDeleteHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	var user User
	getUser(r, &user)
	db.Delete(&user)

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(""))
})

func getUser(r *http.Request, user *User) {
	vars := mux.Vars(r)
	db.Find(&user, vars["id"])
}

func mapUser(r *http.Request, t *User) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)
	if err != nil {
		log.Println("Invalid input")
	}
	t.Salt = ""
	t.Role = ""
}
