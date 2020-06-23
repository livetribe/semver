/*
 * Copyright (c) 2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package semver_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"l7e.io/semver/v1"
)

type scanTest struct {
	name        string
	val         interface{}
	shouldError bool
	expected    string
}

var scanTests = []scanTest{
	{"scan string", "1.2.3-alpha.0+build.012", false, "1.2.3-alpha.0+build.012"},
	{"scan bytes", []byte("1.2.3-alpha.0+build.012"), false, "1.2.3-alpha.0+build.012"},
	{"scan invalid integer", 7, true, ""},
	{"scan invalid float", 7e4, true, ""},
	{"scan boolean", true, true, ""},
}

func TestScanString(t *testing.T) {
	Convey("Test Scan()", t, func() {
		for _, tc := range scanTests {
			Convey(tc.name, func() {
				v := &semver.Version{}
				err := v.Scan(tc.val)
				if tc.shouldError {
					So(err, ShouldNotBeNil)
				} else {
					So(err, ShouldBeNil)
					val, err := v.Value()
					So(err, ShouldBeNil)
					So(val, ShouldEqual, tc.expected)
				}
			})
		}
	})
}
