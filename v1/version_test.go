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
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"l7e.io/semver/v1"
)

func TestSpecVersion(t *testing.T) {
	Convey("Spec version is 2.0.0.", t, func() {
		v := semver.SpecVersion()
		So(v, ShouldResemble, semver.New("2.0.0"))
	})
}

func TestMust(t *testing.T) {
	Convey("Test construction with Must()", t, func() {
		Convey("Calls with errors should panic", func() {
			So(func() {
				semver.Must(nil, errors.New("oi vay"))
			}, ShouldPanic)
		})

		Convey("Calls with no errors should not panic", func() {
			expected := semver.SpecVersion()
			v := semver.Must(expected, nil)
			So(v == expected, ShouldBeTrue)
		})
	})
}

func TestSet(t *testing.T) {
	Convey("Test initialization with NewVersion/Set", t, func() {
		Convey("Test vanilla construction", func() {
			s := "1.2.3-alpha.0+build.012.code.42"
			v, err := semver.NewVersion(s)
			So(err, ShouldBeNil)
			So(v.String(), ShouldEqual, s)
			So(v.Major, ShouldEqual, 1)
			So(v.Minor, ShouldEqual, 2)
			So(v.Patch, ShouldEqual, 3)
			So(v.PreRelease.String(), ShouldEqual, "alpha.0")
			So(v.Metadata.String(), ShouldEqual, "build.012.code.42")
		})

		Convey("Empty strings are not allowed", func() {
			v, err := semver.NewVersion("")
			So(err, ShouldNotBeNil)
			So(v, ShouldBeNil)
		})

		Convey("Test major component", func() {
			Convey("must be a number", func() {
				v, err := semver.NewVersion("a.2.3-alpha.0+build.012")
				So(err, ShouldNotBeNil)
				So(v, ShouldBeNil)
			})
			Convey("cannot contain leading zero", func() {
				v, err := semver.NewVersion("01.2.3-alpha.0+build.012")
				So(err, ShouldNotBeNil)
				So(v, ShouldBeNil)
			})
			Convey("cannot be larger than max uint64", func() {
				// 18446744073709551615 is max uint64, add 1 to cause "overflow"
				v, err := semver.NewVersion("18446744073709551616.2.3-alpha.0+build.012")
				So(err, ShouldNotBeNil)
				So(v, ShouldBeNil)
			})
		})

		Convey("Test minor component", func() {
			Convey("must be a number", func() {
				v, err := semver.NewVersion("1.b.3-alpha.0+build.012")
				So(err, ShouldNotBeNil)
				So(v, ShouldBeNil)
			})
			Convey("cannot contain leading zero", func() {
				v, err := semver.NewVersion("1.02.3-alpha.0+build.012")
				So(err, ShouldNotBeNil)
				So(v, ShouldBeNil)
			})
			Convey("cannot be larger than max uint64", func() {
				// 18446744073709551615 is max uint64, add 1 to cause "overflow"
				v, err := semver.NewVersion("1.18446744073709551616.3-alpha.0+build.012")
				So(err, ShouldNotBeNil)
				So(v, ShouldBeNil)
			})
		})

		Convey("Test patch component", func() {
			Convey("must be a number", func() {
				v, err := semver.NewVersion("1.2.c-alpha.0+build.012")
				So(err, ShouldNotBeNil)
				So(v, ShouldBeNil)
			})
			Convey("cannot contain leading zero", func() {
				v, err := semver.NewVersion("1.2.03-alpha.0+build.012")
				So(err, ShouldNotBeNil)
				So(v, ShouldBeNil)
			})
			Convey("cannot be larger than max uint64", func() {
				// 18446744073709551615 is max uint64, add 1 to cause "overflow"
				v, err := semver.NewVersion("1.2.18446744073709551616-alpha.0+build.012")
				So(err, ShouldNotBeNil)
				So(v, ShouldBeNil)
			})
		})

		Convey("Test pre-release component", func() {
			Convey("must contain valid identifiers", func() {
				// 01 is an invalid pre-release identifier
				v, err := semver.NewVersion("1.2.3-alpha.01+build.012")
				So(err, ShouldNotBeNil)
				So(v, ShouldBeNil)
			})
		})

		Convey("Test build metadata component", func() {
			Convey("must contain valid identifiers", func() {
				// .. indicates an invalid empty identifier
				v, err := semver.NewVersion("1.2.3-alpha.1+build..012")
				So(err, ShouldNotBeNil)
				So(v, ShouldBeNil)
			})
		})
	})
}

func TestIncrement(t *testing.T) {
	Convey("Test increment", t, func() {
		Convey("increment major", func() {
			s := "1.2.3-alpha.0+build.012"
			v := semver.New(s)
			inc := v.IncrementMajor()

			So(v.String(), ShouldEqual, s)
			So(v.Major, ShouldEqual, 1)
			So(v.Minor, ShouldEqual, 2)
			So(v.Patch, ShouldEqual, 3)
			So(v.PreRelease.String(), ShouldEqual, "alpha.0")
			So(v.Metadata.String(), ShouldEqual, "build.012")

			So(inc.String(), ShouldEqual, "2.0.0")
			So(inc.Major, ShouldEqual, 2)
			So(inc.Minor, ShouldEqual, 0)
			So(inc.Patch, ShouldEqual, 0)
			So(inc.PreRelease.String(), ShouldEqual, "")
			So(inc.Metadata.String(), ShouldEqual, "")
		})
		Convey("increment minor", func() {
			s := "1.2.3-alpha.0+build.012"
			v := semver.New(s)
			inc := v.IncrementMinor()

			So(v.String(), ShouldEqual, s)
			So(v.Major, ShouldEqual, 1)
			So(v.Minor, ShouldEqual, 2)
			So(v.Patch, ShouldEqual, 3)
			So(v.PreRelease.String(), ShouldEqual, "alpha.0")
			So(v.Metadata.String(), ShouldEqual, "build.012")

			So(inc.String(), ShouldEqual, "1.3.0")
			So(inc.Major, ShouldEqual, 1)
			So(inc.Minor, ShouldEqual, 3)
			So(inc.Patch, ShouldEqual, 0)
			So(inc.PreRelease.String(), ShouldEqual, "")
			So(inc.Metadata.String(), ShouldEqual, "")
		})
		Convey("increment patch", func() {
			s := "1.2.3-alpha.0+build.012"
			v := semver.New(s)
			inc := v.IncrementPatch()

			So(v.String(), ShouldEqual, s)
			So(v.Major, ShouldEqual, 1)
			So(v.Minor, ShouldEqual, 2)
			So(v.Patch, ShouldEqual, 3)
			So(v.PreRelease.String(), ShouldEqual, "alpha.0")
			So(v.Metadata.String(), ShouldEqual, "build.012")

			So(inc.String(), ShouldEqual, "1.2.4")
			So(inc.Major, ShouldEqual, 1)
			So(inc.Minor, ShouldEqual, 2)
			So(inc.Patch, ShouldEqual, 4)
			So(inc.PreRelease.String(), ShouldEqual, "")
			So(inc.Metadata.String(), ShouldEqual, "")
		})
	})
}

func TestCompatibleUnder(t *testing.T) {
	Convey("Test TestCompatibleUnder()", t, func() {
		Convey("Same major version", func() {
			v := semver.New("1.2.3-alpha.1+build.012")
			Convey("major different", func() {
				So(semver.New("2.0.0").CompatibleUnder(v), ShouldBeFalse)
			})
			Convey("minor lower", func() {
				So(semver.New("1.1.0").CompatibleUnder(v), ShouldBeTrue)
			})
			Convey("minor same, patch lower", func() {
				So(semver.New("1.2.2").CompatibleUnder(v), ShouldBeTrue)
			})
			Convey("minor same, patch same", func() {
				So(semver.New("1.2.3").CompatibleUnder(v), ShouldBeTrue)
			})
			Convey("minor same, patch higher", func() {
				So(semver.New("1.2.4").CompatibleUnder(v), ShouldBeTrue)
			})
			Convey("minor higher", func() {
				So(semver.New("1.3.0").CompatibleUnder(v), ShouldBeFalse)
			})
		})
	})
}

func TestClone(t *testing.T) {
	Convey("Test Clone()", t, func() {
		v := semver.New("1.2.3-alpha.1+build.012")
		cloned := v.Clone()
		So(cloned != v, ShouldBeTrue)
		So(cloned, ShouldResemble, v)
	})
}
