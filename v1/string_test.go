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

func TestVersionString(t *testing.T) {
	Convey("Test VersionString semantics", t, func() {
		Convey("Validate valid version string", func() {
			valid := semver.VersionString("1.2.3-alpha.0+build.012")
			So(valid.Validate(), ShouldBeNil)
		})

		Convey("Validate invalid version string", func() {
			invalid := semver.VersionString("How now brown cow")
			So(invalid.Validate(), ShouldNotBeNil)
		})

		Convey("Invalid version strings should panic when converted to a Version", func() {
			invalid := semver.VersionString("How now brown cow")
			So(func() {
				invalid.Version()
			}, ShouldPanic)
		})

		Convey("Calls with no errors should not panic", func() {
			s := "1.2.3-alpha.1+build.001"
			vs := semver.VersionString(s)
			v := vs.Version()
			So(v, ShouldResemble, semver.New(s))
			So(v.String(), ShouldEqual, s)
		})

		Convey("Round trip test", func() {
			s := "1.2.3-alpha.1+build.001"
			vs := semver.VersionString(s)
			v := vs.Version()
			So(v, ShouldResemble, semver.New(s))
			So(v.String(), ShouldEqual, s)
			So(v.VersionString(), ShouldEqual, vs)
		})
	})
}
