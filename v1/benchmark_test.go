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

	"l7e.io/semver/v1"
)

func BenchmarkConstruction(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_ = semver.New("1.2.3-alpha.1+build.001")
	}
}

func BenchmarkComparison(b *testing.B) {
	a1 := semver.New("1.2.3-alpha.1+build.001")
	a2 := semver.New("1.2.3-alpha.1+build.001")
	for n := 0; n < b.N; n++ {
		_ = a1.Compare(a2)
	}
}

func BenchmarkString(b *testing.B) {
	a1 := *semver.New("1.2.3-alpha.1+build.001")
	for n := 0; n < b.N; n++ {
		_ = a1.String()
	}
}

func BenchmarkIncrement(b *testing.B) {
	a1 := *semver.New("1.2.3-alpha.1+build.001")
	for n := 0; n < b.N; n++ {
		_ = a1.IncrementMajor()
	}
}

func BenchmarkRangeParseSimple(b *testing.B) {
	const VERSION = ">1.0.0"
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = semver.ParseRange(VERSION)
	}
}

func BenchmarkRangeParseAverage(b *testing.B) {
	const VERSION = ">=1.0.0 <2.0.0"
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = semver.ParseRange(VERSION)
	}
}

func BenchmarkRangeParseComplex(b *testing.B) {
	const VERSION = ">=1.0.0 <2.0.0 || >=3.0.1 <4.0.0 !=3.0.3 || >=5.0.0"
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_, _ = semver.ParseRange(VERSION)
	}
}

func BenchmarkRangeMatchSimple(b *testing.B) {
	const VERSION = ">1.0.0"
	r, _ := semver.ParseRange(VERSION)
	v := semver.New("2.0.0")
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r(v)
	}
}

func BenchmarkRangeMatchAverage(b *testing.B) {
	const VERSION = ">=1.0.0 <2.0.0"
	r, _ := semver.ParseRange(VERSION)
	v := semver.New("1.2.3")
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r(v)
	}
}

func BenchmarkRangeMatchComplex(b *testing.B) {
	const VERSION = ">=1.0.0 <2.0.0 || >=3.0.1 <4.0.0 !=3.0.3 || >=5.0.0"
	r, _ := semver.ParseRange(VERSION)
	v := semver.New("5.0.1")
	b.ReportAllocs()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		r(v)
	}
}
