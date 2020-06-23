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
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	v122 = New("1.2.2")
	v123 = New("1.2.3")
	v124 = New("1.2.4")
)

type wildcardTypeTest struct {
	input        string
	wildcardType wildcardType
}

type comparatorTest struct {
	input      string
	comparator func(comparator) bool
}

func TestParseComparator(t *testing.T) {
	tests := []comparatorTest{
		{">", testGT},
		{">=", testGE},
		{"<", testLT},
		{"<=", testLE},
		{"", testEQ},
		{"=", testEQ},
		{"==", testEQ},
		{"!=", testNE},
		{"!", testNE},
		{"-", nil},
		{"<==", nil},
		{"<<", nil},
		{">>", nil},
	}

	Convey("Test comparator parsing", t, func() {
		for _, tc := range tests {
			Convey(tc.input, func() {
				c := parseComparator(tc.input)
				if tc.comparator == nil {
					So(c, ShouldBeNil)
				} else {
					So(tc.comparator(c), ShouldBeTrue)
				}
			})
		}
	})
}

func TestGetWildcardType(t *testing.T) {
	tests := []wildcardTypeTest{
		{"x", majorWildcard},
		{"1.x", minorWildcard},
		{"1.2.x", patchWildcard},
		{"fo.o.b.ar", noneWildcard},
	}

	Convey("Test wildcard types", t, func() {
		for _, tc := range tests {
			Convey(tc.input, func() {
				o := getWildcardType(tc.input)
				So(o, ShouldEqual, tc.wildcardType)
			})
		}
	})
}

func testEQ(f comparator) bool {
	return f(v122, v122) && !f(v122, v123)
}

func testNE(f comparator) bool {
	return !f(v122, v122) && f(v122, v123)
}

func testGT(f comparator) bool {
	return f(v123, v122) && f(v124, v123) && !f(v122, v123) && !f(v122, v122)
}

func testGE(f comparator) bool {
	return f(v123, v122) && f(v124, v123) && !f(v122, v123)
}

func testLT(f comparator) bool {
	return f(v122, v123) && f(v123, v124) && !f(v123, v122) && !f(v122, v122)
}

func testLE(f comparator) bool {
	return f(v122, v123) && f(v123, v124) && !f(v123, v122)
}

func TestSplitAndTrim(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"1.2.3 1.2.3", []string{"1.2.3", "1.2.3"}},
		{"     1.2.3     1.2.3     ", []string{"1.2.3", "1.2.3"}},       // Spaces
		{"  >=   1.2.3   <=  1.2.3   ", []string{">=1.2.3", "<=1.2.3"}}, // Spaces between operator and version
		{"1.2.3 || >=1.2.3 <1.2.3", []string{"1.2.3", "||", ">=1.2.3", "<1.2.3"}},
		{"      1.2.3      ||     >=1.2.3     <1.2.3    ", []string{"1.2.3", "||", ">=1.2.3", "<1.2.3"}},
	}

	Convey("Test split and trim", t, func() {
		for _, tc := range tests {
			Convey(tc.input, func() {
				p := splitAndTrim(tc.input)
				So(p, ShouldResemble, tc.expected)
			})
		}
	})
}

func TestSplitComparatorVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{">1.2.3", []string{">", "1.2.3"}},
		{">=1.2.3", []string{">=", "1.2.3"}},
		{"<1.2.3", []string{"<", "1.2.3"}},
		{"<=1.2.3", []string{"<=", "1.2.3"}},
		{"1.2.3", []string{"", "1.2.3"}},
		{"=1.2.3", []string{"=", "1.2.3"}},
		{"==1.2.3", []string{"==", "1.2.3"}},
		{"!=1.2.3", []string{"!=", "1.2.3"}},
		{"!1.2.3", []string{"!", "1.2.3"}},
		{"error", nil},
	}

	Convey("Test split comparator version", t, func() {
		for _, tc := range tests {
			Convey(tc.input, func() {
				op, v, err := splitComparatorVersion(tc.input)
				if tc.expected == nil {
					So(err, ShouldNotBeNil)
				} else {
					So(err, ShouldBeNil)
					So(op, ShouldEqual, tc.expected[0])
					So(v, ShouldEqual, tc.expected[1])
				}
			})
		}
	})
}

func TestBuildVersionRange(t *testing.T) {
	tests := []struct {
		op       string
		v        string
		c        func(comparator) bool
		expected string
	}{
		{">", "1.2.3", testGT, "1.2.3"},
		{">=", "1.2.3", testGE, "1.2.3"},
		{"<", "1.2.3", testLT, "1.2.3"},
		{"<=", "1.2.3", testLE, "1.2.3"},
		{"", "1.2.3", testEQ, "1.2.3"},
		{"=", "1.2.3", testEQ, "1.2.3"},
		{"==", "1.2.3", testEQ, "1.2.3"},
		{"!=", "1.2.3", testNE, "1.2.3"},
		{"!", "1.2.3", testNE, "1.2.3"},
		{">>", "1.2.3", nil, ""},  // Invalid comparator
		{"=", "invalid", nil, ""}, // Invalid version
	}

	Convey("Test build version range", t, func() {
		for _, tc := range tests {
			r, err := buildVersionRange(tc.op, tc.v)
			if tc.c == nil {
				So(err, ShouldNotBeNil)
			} else {
				So(err, ShouldBeNil)
				So(r, ShouldNotBeNil)
				So(r.v, ShouldResemble, New(tc.expected))
				So(r.c, ShouldNotBeNil)
				So(tc.c(r.c), ShouldBeTrue)
			}
		}
	})
}

func TestSplitORParts(t *testing.T) {
	tests := []struct {
		input    []string
		expected [][]string
	}{
		{[]string{">1.2.3", "||", "<1.2.3", "||", "=1.2.3"}, [][]string{
			{">1.2.3"},
			{"<1.2.3"},
			{"=1.2.3"},
		}},
		{[]string{">1.2.3", "<1.2.3", "||", "=1.2.3"}, [][]string{
			{">1.2.3", "<1.2.3"},
			{"=1.2.3"},
		}},
		{[]string{">1.2.3", "||"}, nil},
		{[]string{"||", ">1.2.3"}, nil},
	}

	Convey("Test split or conditionals", t, func() {
		for _, tc := range tests {
			Convey(strings.Join(tc.input, " "), func() {
				o, err := splitORParts(tc.input)
				if tc.expected == nil {
					So(err, ShouldNotBeNil)
				} else {
					So(err, ShouldBeNil)
					So(o, ShouldResemble, tc.expected)
				}
			})
		}
	})
}

func TestCreateVersionFromWildcard(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.2.x", "1.2.0"},
		{"1.x", "1.0.0"},
	}

	Convey("Creating version from wildcard", t, func() {
		for _, tc := range tests {
			Convey(tc.input, func() {
				v := createVersionFromWildcard(tc.input)
				So(v, ShouldEqual, tc.expected)
			})
		}
	})
}

func TestIncrementMajorVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.2.3", "2.2.3"},
		{"1.2", "2.2"},
		{"foo.bar", ""},
	}

	Convey("Test increment major version", t, func() {
		for _, tc := range tests {
			Convey(tc.input, func() {
				v, _ := incrementMajorVersion(tc.input)
				So(v, ShouldEqual, tc.expected)
			})
		}
	})
}

func TestIncrementMinorVersion(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1.2.3", "1.3.3"},
		{"1.2", "1.3"},
		{"foo.bar", ""},
	}

	Convey("Test increment minor version", t, func() {
		for _, tc := range tests {
			Convey(tc.input, func() {
				p, _ := incrementMinorVersion(tc.input)
				So(p, ShouldEqual, tc.expected)
			})
		}
	})
}

func TestExpandWildcardVersion(t *testing.T) {
	tests := []struct {
		i [][]string
		o [][]string
	}{
		{[][]string{{"foox"}}, nil},
		{[][]string{{">=1.2.x"}}, [][]string{{">=1.2.0"}}},
		{[][]string{{"<=1.2.x"}}, [][]string{{"<1.3.0"}}},
		{[][]string{{">1.2.x"}}, [][]string{{">=1.3.0"}}},
		{[][]string{{"<1.2.x"}}, [][]string{{"<1.2.0"}}},
		{[][]string{{"!=1.2.x"}}, [][]string{{"<1.2.0", ">=1.3.0"}}},
		{[][]string{{">=1.x"}}, [][]string{{">=1.0.0"}}},
		{[][]string{{"<=1.x"}}, [][]string{{"<2.0.0"}}},
		{[][]string{{">1.x"}}, [][]string{{">=2.0.0"}}},
		{[][]string{{"<1.x"}}, [][]string{{"<1.0.0"}}},
		{[][]string{{"!=1.x"}}, [][]string{{"<1.0.0", ">=2.0.0"}}},
		{[][]string{{"1.2.x"}}, [][]string{{">=1.2.0", "<1.3.0"}}},
		{[][]string{{"1.x"}}, [][]string{{">=1.0.0", "<2.0.0"}}},
	}

	Convey("Test expand wildcard version", t, func() {
		for _, tc := range tests {
			Convey(strings.Join(tc.i[0], " "), func() {
				o, err := expandWildcardVersion(tc.i)
				if tc.o == nil {
					So(err, ShouldNotBeNil)
				} else {
					So(err, ShouldBeNil)
					So(o, ShouldResemble, tc.o)
				}
			})
		}
	})
}

func TestVersionRangeToRange(t *testing.T) {
	vr := versionRange{
		v: New("1.2.3"),
		c: compLT,
	}

	Convey("Test version range to range", t, func() {
		rf := vr.rangeFunc()
		So(rf(New("1.2.2")), ShouldBeTrue)
		So(rf(New("1.2.3")), ShouldBeFalse)
	})
}
