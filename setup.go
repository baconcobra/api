/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

func setupRouter() *mux.Router {
	r := mux.NewRouter()

	// Status
	r.Handle("/status", StatusHandler).Methods("GET")

	// Tags
	r.Handle("/tags", jwtMiddleware.Handler(TagsGetHandler)).Methods("GET")
	r.Handle("/tags", jwtMiddleware.Handler(TagsPostHandler)).Methods("POST")
	r.Handle("/tags/{id}", jwtMiddleware.Handler(TagsPatchHandler)).Methods("PATCH")
	r.Handle("/tags/{id}", jwtMiddleware.Handler(TagGetHandler)).Methods("GET")
	r.Handle("/tags/{id}", jwtMiddleware.Handler(TagDeleteHandler)).Methods("DELETE")

	// Actors
	r.Handle("/actors", jwtMiddleware.Handler(ActorsGetHandler)).Methods("GET")
	r.Handle("/actors", jwtMiddleware.Handler(ActorsPostHandler)).Methods("POST")
	r.Handle("/actors/{id}", jwtMiddleware.Handler(ActorsPatchHandler)).Methods("PATCH")
	r.Handle("/actors/{id}", jwtMiddleware.Handler(ActorGetHandler)).Methods("GET")
	r.Handle("/actors/{id}", jwtMiddleware.Handler(ActorDeleteHandler)).Methods("DELETE")

	// Videos
	r.Handle("/videos", jwtMiddleware.Handler(VideosGetHandler)).Methods("GET")
	r.Handle("/videos", jwtMiddleware.Handler(VideosPostHandler)).Methods("POST")
	r.Handle("/videos/{id}", jwtMiddleware.Handler(VideosPatchHandler)).Methods("PATCH")
	r.Handle("/videos/{id}", jwtMiddleware.Handler(VideoGetHandler)).Methods("GET")
	r.Handle("/videos/{id}", jwtMiddleware.Handler(VideoDeleteHandler)).Methods("DELETE")
	r.Handle("/videos/searches", jwtMiddleware.Handler(NotImplemented)).Methods("POST")

	// Tubes
	r.Handle("/tubes", jwtMiddleware.Handler(TubesGetHandler)).Methods("GET")
	r.Handle("/tubes", jwtMiddleware.Handler(TubesPostHandler)).Methods("POST")
	r.Handle("/tubes/{id}", jwtMiddleware.Handler(TubesPatchHandler)).Methods("PATCH")
	r.Handle("/tubes/{id}", jwtMiddleware.Handler(TubeGetHandler)).Methods("GET")
	r.Handle("/tubes/{id}", jwtMiddleware.Handler(TubeDeleteHandler)).Methods("DELETE")

	// Users
	r.Handle("/users", jwtMiddleware.Handler(UsersGetHandler)).Methods("GET")
	r.Handle("/users", jwtMiddleware.Handler(UsersPostHandler)).Methods("POST")
	r.Handle("/users/{id}", jwtMiddleware.Handler(UsersPatchHandler)).Methods("PATCH")
	r.Handle("/users/{id}", jwtMiddleware.Handler(UserGetHandler)).Methods("GET")
	r.Handle("/users/{id}", jwtMiddleware.Handler(UserDeleteHandler)).Methods("DELETE")

	// Auth
	r.Handle("/auth", GetTokenHandler).Methods("POST")

	return r
}

func setupDB(connector string, database string) {
	db, err = gorm.Open(connector, database)
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Tube{})
	db.AutoMigrate(&Tag{})
	db.AutoMigrate(&Actor{})
	db.AutoMigrate(&Video{})
	db.AutoMigrate(&User{})
}
