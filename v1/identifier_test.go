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

func TestIdentifier(t *testing.T) {
	Convey("Test String()", t, func() {
		Convey("non-numeric", func() {
			id, err := semver.NewIdentifierForTest("a-b", true)
			So(err, ShouldBeNil)
			So(id.String(), ShouldEqual, "a-b")
		})
		Convey("numeric", func() {
			id, err := semver.NewIdentifierForTest("12", true)
			So(err, ShouldBeNil)
			So(id.String(), ShouldEqual, "12")
		})
	})

	Convey("Test Compare()", t, func() {
		Convey("non-numeric", func() {
			ab, _ := semver.NewIdentifierForTest("a-b", true)
			ac, _ := semver.NewIdentifierForTest("a-c", true)
			So(ab.Compare(ab), ShouldEqual, 0)
			So(ab.Compare(ac), ShouldBeLessThan, 0)
			So(ac.Compare(ab), ShouldBeGreaterThan, 0)
		})
		Convey("numeric", func() {
			id12, _ := semver.NewIdentifierForTest("12", true)
			id13, _ := semver.NewIdentifierForTest("13", true)
			So(id12.Compare(id12), ShouldEqual, 0)
			So(id12.Compare(id13), ShouldBeLessThan, 0)
			So(id13.Compare(id12), ShouldBeGreaterThan, 0)
		})
		Convey("numeric and non-numeric", func() {
			id12, _ := semver.NewIdentifierForTest("12", true)
			ab, _ := semver.NewIdentifierForTest("a-b", true)
			So(id12.Compare(ab), ShouldBeLessThan, 0)
			So(ab.Compare(id12), ShouldBeGreaterThan, 0)
		})
	})

	Convey("Test Clone()", t, func() {
		var ids semver.Identifiers

		for _, s := range []string{"a-b", "12"} {
			id, _ := semver.NewIdentifierForTest(s, true)
			ids = append(ids, id)
		}
		cloned := ids.Clone()

		So(cloned, ShouldNotEqual, ids)
		So(cloned, ShouldResemble, ids)
	})

	Convey("Test map-like features", t, func() {
		Convey("Get", func() {
			var ids semver.Identifiers

			for _, s := range []string{"a", "12", "b", "c", "3", "d", "e"} {
				id, _ := semver.NewIdentifierForTest(s, true)
				ids = append(ids, id)
			}

			v, ok := ids.Get("a")
			So(ok, ShouldBeTrue)
			So(v, ShouldResemble, semver.Identifier{Num: 12, IsNum: true})

			v, ok = ids.Get("b")
			So(ok, ShouldBeTrue)
			So(v, ShouldResemble, semver.Identifier{Str: "c"})

			_, ok = ids.Get("c")
			So(ok, ShouldBeFalse)

			_, ok = ids.Get("foo")
			So(ok, ShouldBeFalse)

			_, ok = ids.Get("3")
			So(ok, ShouldBeFalse)

			_, ok = ids.Get("e")
			So(ok, ShouldBeFalse)
		})

		Convey("Increment", func() {
			var ids semver.Identifiers

			for _, s := range []string{"a", "12", "b", "c", "3", "d", "e"} {
				id, _ := semver.NewIdentifierForTest(s, true)
				ids = append(ids, id)
			}

			v, ok := ids.Increment("a")
			So(ok, ShouldBeTrue)
			So(v, ShouldEqual, 13)
			So(ids.String(), ShouldEqual, "a.13.b.c.3.d.e")

			_, ok = ids.Increment("b")
			So(ok, ShouldBeFalse)

			_, ok = ids.Increment("c")
			So(ok, ShouldBeFalse)

			_, ok = ids.Increment("foo")
			So(ok, ShouldBeFalse)

			_, ok = ids.Increment("3")
			So(ok, ShouldBeFalse)

			_, ok = ids.Increment("e")
			So(ok, ShouldBeFalse)
		})

		Convey("Contains", func() {
			var ids semver.Identifiers

			for _, s := range []string{"a", "12", "b", "c", "3", "d", "e"} {
				id, _ := semver.NewIdentifierForTest(s, true)
				ids = append(ids, id)
			}

			ok := ids.Contains("a")
			So(ok, ShouldBeTrue)

			ok = ids.Contains("b")
			So(ok, ShouldBeTrue)

			ok = ids.Contains("c")
			So(ok, ShouldBeFalse)

			ok = ids.Contains("foo")
			So(ok, ShouldBeFalse)

			ok = ids.Contains("3")
			So(ok, ShouldBeFalse)

			ok = ids.Contains("e")
			So(ok, ShouldBeFalse)
		})
	})

	Convey("Set", t, func() {
		var ids semver.Identifiers

		for _, s := range []string{"a", "12", "b", "c", "3", "d", "e"} {
			id, _ := semver.NewIdentifierForTest(s, true)
			ids = append(ids, id)
		}
		So(ids.String(), ShouldEqual, "a.12.b.c.3.d.e")

		ok := ids.SetWithNumber("a", 42)
		So(ok, ShouldBeTrue)
		So(ids.String(), ShouldEqual, "a.42.b.c.3.d.e")

		ok = ids.SetWithString("a", "bar")
		So(ok, ShouldBeTrue)
		So(ids.String(), ShouldEqual, "a.bar.b.c.3.d.e")

		ok = ids.SetWithString("b", "car")
		So(ok, ShouldBeTrue)
		So(ids.String(), ShouldEqual, "a.bar.b.car.3.d.e")

		ok = ids.SetWithString("car", "cdr")
		So(ok, ShouldBeFalse)
		So(ids.String(), ShouldEqual, "a.bar.b.car.3.d.e")

		ok = ids.SetWithString("zoo", "cat")
		So(ok, ShouldBeFalse)
		So(ids.String(), ShouldEqual, "a.bar.b.car.3.d.e")

		ok = ids.SetWithString("3", "number")
		So(ok, ShouldBeFalse)
		So(ids.String(), ShouldEqual, "a.bar.b.car.3.d.e")

		ok = ids.SetWithString("e", "end")
		So(ok, ShouldBeFalse)
		So(ids.String(), ShouldEqual, "a.bar.b.car.3.d.e")
	})
}

func TestIdentifiers(t *testing.T) {
	Convey("Test identifiers", t, func() {
		Convey("Test empty identifiers", func() {
			ids := make(semver.Identifiers, 0)

			So(ids.String(), ShouldEqual, "")
		})

		Convey("Test compare identifiers", func() {
			Convey("Test equal identifiers", func() {
				var this semver.Identifiers
				var that semver.Identifiers

				for _, s := range []string{"a-b", "12"} {
					id, _ := semver.NewIdentifierForTest(s, true)
					this = append(this, id)
					that = append(that, id)
				}

				So(this.Compare(that), ShouldEqual, 0)
				So(that.Compare(this), ShouldEqual, 0)
			})

			Convey("numeric vs non-numeric", func() {
				var numeric semver.Identifiers
				var nonNumeric semver.Identifiers

				for _, s := range []string{"a-b", "12"} {
					id, _ := semver.NewIdentifierForTest(s, true)
					numeric = append(numeric, id)
					nonNumeric = append(nonNumeric, id)
				}
				id, _ := semver.NewIdentifierForTest("123456789", true)
				numeric = append(numeric, id)
				id, _ = semver.NewIdentifierForTest("how-now-brow-cow", true)
				nonNumeric = append(nonNumeric, id)
				id, _ = semver.NewIdentifierForTest("end", true)
				numeric = append(numeric, id)
				nonNumeric = append(nonNumeric, id)

				So(numeric.Compare(nonNumeric), ShouldBeLessThan, 0)
				So(nonNumeric.Compare(numeric), ShouldBeGreaterThan, 0)
			})

			Convey("non-pre-release takes precedence over pre-release", func() {
				var notPreRelease = make(semver.Identifiers, 0)
				var preRelease semver.Identifiers

				for _, s := range []string{"a-b", "12"} {
					id, _ := semver.NewIdentifierForTest(s, true)
					preRelease = append(preRelease, id)
				}

				So(notPreRelease.Compare(preRelease), ShouldBeGreaterThan, 0)
				So(preRelease.Compare(notPreRelease), ShouldBeLessThan, 0)
			})

			Convey("A larger array of pre-release identifiers has a higher precedence than a smaller array, if all of the preceding identifiers are equal", func() {
				var smaller semver.Identifiers
				var larger semver.Identifiers

				for _, s := range []string{"a-b", "12"} {
					id, _ := semver.NewIdentifierForTest(s, true)
					smaller = append(smaller, id)
					larger = append(larger, id)
				}
				id, _ := semver.NewIdentifierForTest("extra", true)
				larger = append(larger, id)

				So(smaller.Compare(larger), ShouldBeLessThan, 0)
				So(larger.Compare(smaller), ShouldBeGreaterThan, 0)
			})
		})
	})
}
