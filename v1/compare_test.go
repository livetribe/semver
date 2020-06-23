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

func TestCompare(t *testing.T) {
	Convey("Test Compare()", t, func() {
		Convey("test equality", func() {
			this := semver.New("1.2.3-alpha.0")
			that := semver.New("1.2.3-alpha.0")
			So(this.Compare(that), ShouldEqual, 0)
		})
		Convey("test equality different metadata", func() {
			this := semver.New("1.2.3-alpha.0+build.001")
			that := semver.New("1.2.3-alpha.0+build.002")
			So(this.Compare(that), ShouldEqual, 0)
		})
		Convey("pre-release is less than release", func() {
			this := semver.New("1.2.3")
			that := semver.New("1.2.3-alpha.0")
			So(this.Compare(that), ShouldBeGreaterThan, 0)
			So(that.Compare(this), ShouldBeLessThan, 0)
		})
		Convey("patch increment", func() {
			this := semver.New("1.2.4")
			that := semver.New("1.2.3")
			So(this.Compare(that), ShouldBeGreaterThan, 0)
			So(that.Compare(this), ShouldBeLessThan, 0)
		})
		Convey("minor increment", func() {
			this := semver.New("1.3.0")
			that := semver.New("1.2.3")
			So(this.Compare(that), ShouldBeGreaterThan, 0)
			So(that.Compare(this), ShouldBeLessThan, 0)
		})
		Convey("major increment", func() {
			this := semver.New("2.0.0")
			that := semver.New("1.2.3")
			So(this.Compare(that), ShouldBeGreaterThan, 0)
			So(that.Compare(this), ShouldBeLessThan, 0)
		})
	})

	Convey("Test 'operators'", t, func() {
		Convey("test Equals", func() {
			Convey("equal", func() {
				this := semver.New("1.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.Equals(that), ShouldBeTrue)
			})
			Convey("different build metadata", func() {
				this := semver.New("1.2.3-alpha.0+build.001")
				that := semver.New("1.2.3-alpha.0+build.001")
				So(this.Equals(that), ShouldBeTrue)
			})
			Convey("mismatched major", func() {
				this := semver.New("9.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.Equals(that), ShouldBeFalse)
			})
			Convey("mismatched minor", func() {
				this := semver.New("1.9.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.Equals(that), ShouldBeFalse)
			})
			Convey("mismatched patch", func() {
				this := semver.New("1.2.9-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.Equals(that), ShouldBeFalse)
			})
			Convey("mismatched pre-release", func() {
				Convey("number identifier", func() {
					this := semver.New("1.2.3-alpha.9")
					that := semver.New("1.2.3-alpha.0")
					So(this.Equals(that), ShouldBeFalse)
				})
				Convey("string identifier", func() {
					this := semver.New("1.2.3-beta.0")
					that := semver.New("1.2.3-alpha.0")
					So(this.Equals(that), ShouldBeFalse)
				})
			})
		})
		Convey("test EQ", func() {
			Convey("equal", func() {
				this := semver.New("1.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.EQ(that), ShouldBeTrue)
			})
			Convey("different build metadata", func() {
				this := semver.New("1.2.3-alpha.0+build.001")
				that := semver.New("1.2.3-alpha.0+build.001")
				So(this.EQ(that), ShouldBeTrue)
			})
			Convey("mismatched major", func() {
				this := semver.New("9.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.EQ(that), ShouldBeFalse)
			})
			Convey("mismatched minor", func() {
				this := semver.New("1.9.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.EQ(that), ShouldBeFalse)
			})
			Convey("mismatched patch", func() {
				this := semver.New("1.2.9-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.EQ(that), ShouldBeFalse)
			})
			Convey("mismatched pre-release", func() {
				Convey("number identifier", func() {
					this := semver.New("1.2.3-alpha.9")
					that := semver.New("1.2.3-alpha.0")
					So(this.EQ(that), ShouldBeFalse)
				})
				Convey("string identifier", func() {
					this := semver.New("1.2.3-beta.0")
					that := semver.New("1.2.3-alpha.0")
					So(this.EQ(that), ShouldBeFalse)
				})
			})
		})
		Convey("test NE", func() {
			Convey("equal", func() {
				this := semver.New("1.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.NE(that), ShouldBeFalse)
			})
			Convey("different build metadata", func() {
				this := semver.New("1.2.3-alpha.0+build.001")
				that := semver.New("1.2.3-alpha.0+build.001")
				So(this.NE(that), ShouldBeFalse)
			})
			Convey("mismatched major", func() {
				this := semver.New("9.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.NE(that), ShouldBeTrue)
			})
			Convey("mismatched minor", func() {
				this := semver.New("1.9.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.NE(that), ShouldBeTrue)
			})
			Convey("mismatched patch", func() {
				this := semver.New("1.2.9-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.NE(that), ShouldBeTrue)
			})
			Convey("mismatched pre-release", func() {
				Convey("number identifier", func() {
					this := semver.New("1.2.3-alpha.9")
					that := semver.New("1.2.3-alpha.0")
					So(this.NE(that), ShouldBeTrue)
				})
				Convey("string identifier", func() {
					this := semver.New("1.2.3-beta.0")
					that := semver.New("1.2.3-alpha.0")
					So(this.NE(that), ShouldBeTrue)
				})
			})
		})
		Convey("test GT", func() {
			Convey("equal", func() {
				this := semver.New("1.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.GT(that), ShouldBeFalse)
			})
			Convey("mismatched major", func() {
				this := semver.New("9.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.GT(that), ShouldBeTrue)
				So(that.GT(this), ShouldBeFalse)
			})
			Convey("mismatched minor", func() {
				this := semver.New("1.9.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.GT(that), ShouldBeTrue)
				So(that.GT(this), ShouldBeFalse)
			})
			Convey("mismatched patch", func() {
				this := semver.New("1.2.9-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.GT(that), ShouldBeTrue)
				So(that.GT(this), ShouldBeFalse)
			})
			Convey("mismatched pre-release", func() {
				Convey("number identifier", func() {
					this := semver.New("1.2.3-alpha.9")
					that := semver.New("1.2.3-alpha.0")
					So(this.GT(that), ShouldBeTrue)
					So(that.GT(this), ShouldBeFalse)
				})
				Convey("string identifier", func() {
					this := semver.New("1.2.3-beta.0")
					that := semver.New("1.2.3-alpha.0")
					So(this.GT(that), ShouldBeTrue)
					So(that.GT(this), ShouldBeFalse)
				})
			})
		})
		Convey("test GTE", func() {
			Convey("equal", func() {
				this := semver.New("1.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.GTE(that), ShouldBeTrue)
			})
			Convey("mismatched major", func() {
				this := semver.New("9.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.GTE(that), ShouldBeTrue)
				So(that.GTE(this), ShouldBeFalse)
			})
			Convey("mismatched minor", func() {
				this := semver.New("1.9.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.GTE(that), ShouldBeTrue)
				So(that.GTE(this), ShouldBeFalse)
			})
			Convey("mismatched patch", func() {
				this := semver.New("1.2.9-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.GTE(that), ShouldBeTrue)
				So(that.GTE(this), ShouldBeFalse)
			})
			Convey("mismatched pre-release", func() {
				Convey("number identifier", func() {
					this := semver.New("1.2.3-alpha.9")
					that := semver.New("1.2.3-alpha.0")
					So(this.GTE(that), ShouldBeTrue)
					So(that.GTE(this), ShouldBeFalse)
				})
				Convey("string identifier", func() {
					this := semver.New("1.2.3-beta.0")
					that := semver.New("1.2.3-alpha.0")
					So(this.GTE(that), ShouldBeTrue)
					So(that.GTE(this), ShouldBeFalse)
				})
			})
		})
		Convey("test GE", func() {
			Convey("equal", func() {
				this := semver.New("1.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.GE(that), ShouldBeTrue)
			})
			Convey("mismatched major", func() {
				this := semver.New("9.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.GE(that), ShouldBeTrue)
				So(that.GE(this), ShouldBeFalse)
			})
			Convey("mismatched minor", func() {
				this := semver.New("1.9.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.GE(that), ShouldBeTrue)
				So(that.GE(this), ShouldBeFalse)
			})
			Convey("mismatched patch", func() {
				this := semver.New("1.2.9-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.GE(that), ShouldBeTrue)
				So(that.GE(this), ShouldBeFalse)
			})
			Convey("mismatched pre-release", func() {
				Convey("number identifier", func() {
					this := semver.New("1.2.3-alpha.9")
					that := semver.New("1.2.3-alpha.0")
					So(this.GE(that), ShouldBeTrue)
					So(that.GE(this), ShouldBeFalse)
				})
				Convey("string identifier", func() {
					this := semver.New("1.2.3-beta.0")
					that := semver.New("1.2.3-alpha.0")
					So(this.GE(that), ShouldBeTrue)
					So(that.GE(this), ShouldBeFalse)
				})
			})
		})
		Convey("test LT", func() {
			Convey("equal", func() {
				this := semver.New("1.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.LT(that), ShouldBeFalse)
			})
			Convey("mismatched major", func() {
				this := semver.New("9.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.LT(that), ShouldBeFalse)
				So(that.LT(this), ShouldBeTrue)
			})
			Convey("mismatched minor", func() {
				this := semver.New("1.9.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.LT(that), ShouldBeFalse)
				So(that.LT(this), ShouldBeTrue)
			})
			Convey("mismatched patch", func() {
				this := semver.New("1.2.9-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.LT(that), ShouldBeFalse)
				So(that.LT(this), ShouldBeTrue)
			})
			Convey("mismatched pre-release", func() {
				Convey("number identifier", func() {
					this := semver.New("1.2.3-alpha.9")
					that := semver.New("1.2.3-alpha.0")
					So(this.LT(that), ShouldBeFalse)
					So(that.LT(this), ShouldBeTrue)
				})
				Convey("string identifier", func() {
					this := semver.New("1.2.3-beta.0")
					that := semver.New("1.2.3-alpha.0")
					So(this.LT(that), ShouldBeFalse)
					So(that.LT(this), ShouldBeTrue)
				})
			})
		})
		Convey("test LTE", func() {
			Convey("equal", func() {
				this := semver.New("1.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.LTE(that), ShouldBeTrue)
			})
			Convey("mismatched major", func() {
				this := semver.New("9.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.LTE(that), ShouldBeFalse)
				So(that.LTE(this), ShouldBeTrue)
			})
			Convey("mismatched minor", func() {
				this := semver.New("1.9.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.LTE(that), ShouldBeFalse)
				So(that.LTE(this), ShouldBeTrue)
			})
			Convey("mismatched patch", func() {
				this := semver.New("1.2.9-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.LTE(that), ShouldBeFalse)
				So(that.LTE(this), ShouldBeTrue)
			})
			Convey("mismatched pre-release", func() {
				Convey("number identifier", func() {
					this := semver.New("1.2.3-alpha.9")
					that := semver.New("1.2.3-alpha.0")
					So(this.LTE(that), ShouldBeFalse)
					So(that.LTE(this), ShouldBeTrue)
				})
				Convey("string identifier", func() {
					this := semver.New("1.2.3-beta.0")
					that := semver.New("1.2.3-alpha.0")
					So(this.LTE(that), ShouldBeFalse)
					So(that.LTE(this), ShouldBeTrue)
				})
			})
		})
		Convey("test LE", func() {
			Convey("equal", func() {
				this := semver.New("1.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.LE(that), ShouldBeTrue)
			})
			Convey("mismatched major", func() {
				this := semver.New("9.2.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.LE(that), ShouldBeFalse)
				So(that.LE(this), ShouldBeTrue)
			})
			Convey("mismatched minor", func() {
				this := semver.New("1.9.3-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.LE(that), ShouldBeFalse)
				So(that.LE(this), ShouldBeTrue)
			})
			Convey("mismatched patch", func() {
				this := semver.New("1.2.9-alpha.0")
				that := semver.New("1.2.3-alpha.0")
				So(this.LE(that), ShouldBeFalse)
				So(that.LE(this), ShouldBeTrue)
			})
			Convey("mismatched pre-release", func() {
				Convey("number identifier", func() {
					this := semver.New("1.2.3-alpha.9")
					that := semver.New("1.2.3-alpha.0")
					So(this.LE(that), ShouldBeFalse)
					So(that.LE(this), ShouldBeTrue)
				})
				Convey("string identifier", func() {
					this := semver.New("1.2.3-beta.0")
					that := semver.New("1.2.3-alpha.0")
					So(this.LE(that), ShouldBeFalse)
					So(that.LE(this), ShouldBeTrue)
				})
			})
		})
	})
}
