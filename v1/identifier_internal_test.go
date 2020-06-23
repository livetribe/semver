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

package semver

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func NewIdentifierForTest(s string, strict bool) (Identifier, error) {
	return newIdentifier(s, strict)
}

func TestNewIdentifier(t *testing.T) {
	Convey("Test pre-release identifiers", t, func() {
		Convey("Simple text identifier", func() {
			id, err := newIdentifier("a-b", true)
			So(err, ShouldBeNil)
			So(id, ShouldResemble, Identifier{Str: "a-b"})
		})
		Convey("Simple numeric identifier", func() {
			id, err := newIdentifier("12", true)
			So(err, ShouldBeNil)
			So(id, ShouldResemble, Identifier{Num: 12, IsNum: true})
		})
		Convey("Identifiers MUST comprise only ASCII alphanumerics and hyphen [0-9A-Za-z-]", func() {
			_, err := newIdentifier("ab#c", true)
			So(err, ShouldNotBeNil)
		})
		Convey("Identifiers MUST NOT be empty", func() {
			_, err := newIdentifier("", true)
			So(err, ShouldNotBeNil)
		})
		Convey("Numeric identifiers MUST NOT include leading zeroes", func() {
			_, err := newIdentifier("012", true)
			So(err, ShouldNotBeNil)
		})
		Convey("Handle strconv.ParseUint(expected, 10, 64) errors", func() {
			// 18446744073709551615 is max uint64
			_, err := newIdentifier("18446744073709551616", true)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Test build metadata identifiers", t, func() {
		Convey("Simple text identifier", func() {
			id, err := newIdentifier("a-b", false)
			So(err, ShouldBeNil)
			So(id, ShouldResemble, Identifier{Str: "a-b"})
		})
		Convey("Simple numeric identifier", func() {
			id, err := newIdentifier("12", false)
			So(err, ShouldBeNil)
			So(id, ShouldResemble, Identifier{Num: 12, IsNum: true})
		})
		Convey("Identifiers MUST comprise only ASCII alphanumerics and hyphen [0-9A-Za-z-]", func() {
			_, err := newIdentifier("ab#c", false)
			So(err, ShouldNotBeNil)
		})
		Convey("Identifiers MUST NOT be empty", func() {
			_, err := newIdentifier("", false)
			So(err, ShouldNotBeNil)
		})
		Convey("Numeric identifiers in build metadata can include leading zeroes and are treated as strings", func() {
			id, err := newIdentifier("012", false)
			So(err, ShouldBeNil)
			So(id, ShouldResemble, Identifier{Str: "012"})
			So(id.String(), ShouldEqual, "012")
		})
		Convey("Handle strconv.ParseUint(expected, 10, 64) errors", func() {
			// 18446744073709551615 is max uint64, add 1 to cause "overflow"
			_, err := newIdentifier("18446744073709551616", false)
			So(err, ShouldNotBeNil)
		})
	})
}
