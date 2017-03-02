/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package main

import (
	"strconv"
)

type Navigation struct {
	Prev string `json:"prev,omitempty"`
	Next string `json:"next,omitemtpy"`
}

func getNavigation(n, page, limit int) Navigation {
	nav := Navigation{}
	if n == limit {
		nav.Next = strconv.Itoa(page + 1)
	}
	if page > 0 {
		nav.Prev = strconv.Itoa(page - 1)
	}
	return nav
}
