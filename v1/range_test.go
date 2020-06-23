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

var (
	v121 = semver.New("1.2.1")
	v122 = semver.New("1.2.2")
	v123 = semver.New("1.2.3")
)

func TestRangeAND(t *testing.T) {
	Convey("Test and range", t, func() {
		gt121 := semver.Range(func(v *semver.Version) bool {
			return v.GT(v121)
		})
		lt123 := semver.Range(func(v *semver.Version) bool {
			return v.LT(v123)
		})
		rf := gt121.AND(lt123)

		So(rf(v121), ShouldBeFalse)
		So(rf(v122), ShouldBeTrue)
		So(rf(v123), ShouldBeFalse)
	})
}

func TestRangeOR(t *testing.T) {
	tests := []struct {
		v        *semver.Version
		expected bool
	}{
		{semver.New("1.2.0"), true},
		{semver.New("1.2.2"), false},
		{semver.New("1.2.4"), true},
	}

	Convey("Test and range", t, func() {
		lt121 := semver.Range(func(v *semver.Version) bool {
			return v.LT(v121)
		})
		gt123 := semver.Range(func(v *semver.Version) bool {
			return v.GT(v123)
		})
		rf := lt121.OR(gt123)

		for _, tc := range tests {
			Convey(tc.v.String(), func() {
				So(rf(tc.v), ShouldEqual, tc.expected)
			})
		}
	})
}

func TestParseRange(t *testing.T) {
	type test struct {
		v        string
		expected bool
	}
	tests := []struct {
		conditional string
		data        []test
	}{
		// Simple expressions
		{">1.2.3", []test{
			{"1.2.2", false},
			{"1.2.3", false},
			{"1.2.4", true},
		}},
		{">=1.2.3", []test{
			{"1.2.3", true},
			{"1.2.4", true},
			{"1.2.2", false},
		}},
		{"<1.2.3", []test{
			{"1.2.2", true},
			{"1.2.3", false},
			{"1.2.4", false},
		}},
		{"<=1.2.3", []test{
			{"1.2.2", true},
			{"1.2.3", true},
			{"1.2.4", false},
		}},
		{"1.2.3", []test{
			{"1.2.2", false},
			{"1.2.3", true},
			{"1.2.4", false},
		}},
		{"=1.2.3", []test{
			{"1.2.2", false},
			{"1.2.3", true},
			{"1.2.4", false},
		}},
		{"==1.2.3", []test{
			{"1.2.2", false},
			{"1.2.3", true},
			{"1.2.4", false},
		}},
		{"!=1.2.3", []test{
			{"1.2.2", true},
			{"1.2.3", false},
			{"1.2.4", true},
		}},
		{"!1.2.3", []test{
			{"1.2.2", true},
			{"1.2.3", false},
			{"1.2.4", true},
		}},
		// Simple Expression errors
		{">>1.2.3", nil},
		{"!!1.2.3", nil},
		{"1.0", nil},
		{"string", nil},
		{"", nil},
		{"fo.ob.ar.x", nil},
		// AND Expressions
		{">1.2.2 <1.2.4", []test{
			{"1.2.2", false},
			{"1.2.3", true},
			{"1.2.4", false},
		}},
		{"<1.2.2 <1.2.4", []test{
			{"1.2.1", true},
			{"1.2.2", false},
			{"1.2.3", false},
			{"1.2.4", false},
		}},
		{">1.2.2 <1.2.5 !=1.2.4", []test{
			{"1.2.2", false},
			{"1.2.3", true},
			{"1.2.4", false},
			{"1.2.5", false},
		}},
		{">1.2.2 <1.2.5 !1.2.4", []test{
			{"1.2.2", false},
			{"1.2.3", true},
			{"1.2.4", false},
			{"1.2.5", false},
		}},
		// OR Expressions
		{">1.2.2 || <1.2.4", []test{
			{"1.2.2", true},
			{"1.2.3", true},
			{"1.2.4", true},
		}},
		{"<1.2.2 || >1.2.4", []test{
			{"1.2.2", false},
			{"1.2.3", false},
			{"1.2.4", false},
		}},
		// Wildcard expressions
		{">1.x", []test{
			{"0.1.9", false},
			{"1.2.6", false},
			{"1.9.0", false},
			{"2.0.0", true},
		}},
		{">1.2.x", []test{
			{"1.1.9", false},
			{"1.2.6", false},
			{"1.3.0", true},
		}},
		// Combined Expressions
		{">1.2.2 <1.2.4 || >=2.0.0", []test{
			{"1.2.2", false},
			{"1.2.3", true},
			{"1.2.4", false},
			{"2.0.0", true},
			{"2.0.1", true},
		}},
		{"1.x || >=2.0.x <2.2.x", []test{
			{"0.9.2", false},
			{"1.2.2", true},
			{"2.0.0", true},
			{"2.1.8", true},
			{"2.2.0", false},
		}},
		{">1.2.2 <1.2.4 || >=2.0.0 <3.0.0", []test{
			{"1.2.2", false},
			{"1.2.3", true},
			{"1.2.4", false},
			{"2.0.0", true},
			{"2.0.1", true},
			{"2.9.9", true},
			{"3.0.0", false},
		}},
	}

	Convey("Test range parsing", t, func() {
		for _, tc := range tests {
			Convey(tc.conditional, func() {
				r, err := semver.ParseRange(tc.conditional)
				if tc.data == nil {
					So(err, ShouldNotBeNil)
				} else {
					for _, td := range tc.data {
						Convey(td.v, func() {
							v := semver.New(td.v)
							So(r(v), ShouldEqual, td.expected)
						})
					}
				}
			})
		}
	})
}

func TestMustParseRange(t *testing.T) {
	Convey("Test MustParseRange", t, func() {
		rf := semver.MustParseRange(">1.2.2 <1.2.4 || >=2.0.0 <3.0.0")
		So(rf(semver.New("1.2.3")), ShouldBeTrue)
	})
}

func TestMustParseRange_panic(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Errorf("Should have panicked")
		}
	}()
	_ = semver.MustParseRange("invalid version")
}
