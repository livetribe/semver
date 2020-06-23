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
	"encoding/json"
	"strconv"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"l7e.io/semver/v1"
)

func TestJSON(t *testing.T) {
	Convey("Test marshall", t, func() {
		s := "3.1.4-alpha.1.5.9+build.2.6.5"
		v := semver.New(s)
		j, err := json.Marshal(v)
		So(err, ShouldBeNil)
		So(string(j), ShouldEqual, strconv.Quote(s))
	})
	Convey("Test unmarshall", t, func() {
		Convey("valid version", func() {
			s := "3.1.4-alpha.1.5.9+build.2.6.5"
			var v semver.Version
			err := json.Unmarshal([]byte(strconv.Quote(s)), &v)
			So(err, ShouldBeNil)
			So(v.String(), ShouldEqual, s)
		})
		Convey("invalid version", func() {
			s := "3.1.4.1.5.9.2.6.5-other-digits-of-pi"
			var v semver.Version
			err := json.Unmarshal([]byte(strconv.Quote(s)), &v)
			So(err, ShouldNotBeNil)
		})
		Convey("unmarshal a number", func() {
			var v semver.Version
			err := json.Unmarshal([]byte("1234"), &v)
			So(err, ShouldNotBeNil)
		})
	})
}
