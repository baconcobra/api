/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"crypto/subtle"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/scrypt"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// Handlers
var GetTokenHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	user := User{}
	vars := mux.Vars(r)
	db.Find(&user, vars["username"])
	if user.UserName != vars["username"] {
		// TODO : Return proper errors
		return
	}

	if !isValidPassword(user, vars["password"]) {
		// TODO : Return proper errors
		return
	}

	token := getToken(user)
	w.Write(token)
})

func isValidPassword(u User, pw string) bool {
	userpass, err := base64.StdEncoding.DecodeString(u.Password)
	if err != nil {
		return false
	}

	usersalt, err := base64.StdEncoding.DecodeString(u.Salt)
	if err != nil {
		return false
	}

	hash, err := scrypt.Key([]byte(pw), usersalt, 16384, 8, 1, HashSize)
	if err != nil {
		return false
	}

	// Compare in constant time to mitigate timing attacks
	if subtle.ConstantTimeCompare(userpass, hash) == 1 {
		return true
	}

	return false
}

func getToken(user User) []byte {
	secret := []byte(os.Getenv("AUTH_CLIENT_SECRET"))

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":       user.ID,
		"name":     user.Name,
		"username": user.UserName,
		"role":     user.Role,
		"nbf":      time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(secret)

	if err != nil {
		log.Fatal(err)
	}

	return []byte(tokenString)
}

var jwtMiddleware = jwtmiddleware.New(jwtmiddleware.Options{
	ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("AUTH_CLIENT_SECRET")), nil
	},
})

func currentUser(r *http.Request) (u User) {
	user := context.Get(r, "user").(*jwt.Token)
	claims, ok := user.Claims.(jwt.MapClaims)

	if ok {
		u.ID = claims["id"].(uint)
		u.Name = claims["name"].(string)
		u.UserName = claims["username"].(string)
		u.Role = claims["role"].(string)
	}

	return u
}
