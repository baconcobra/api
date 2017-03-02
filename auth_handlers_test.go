/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAuthInvalidUser(t *testing.T) {
	SkipConvey("Given a user 'andrew' exists", t, func() {
		setupTestSuite()
		Convey("When I call POST /auth with non 'joy' username", func() {

			response := doRequest("POST", "/auth", nil)

			Convey("Then I should get a 403 response", func() {
				status := response.Code
				So(status, ShouldEqual, 403)
			})
		})
	})
}
