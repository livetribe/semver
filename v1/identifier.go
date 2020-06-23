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
	"errors"
	"fmt"
	"strconv"
)

// Identifier is a component of either pre-release or metadata fields.
type Identifier struct {
	Str   string
	Num   uint64
	IsNum bool
}

// Identifiers represents either pre-release or metadata fields.
type Identifiers []Identifier

// newIdentifier creates a new valid Identifier.  Use strict to indicate if the
// Identifier is a pre-release identifier.
func newIdentifier(s string, strict bool) (Identifier, error) {
	if len(s) == 0 {
		return Identifier{}, errors.New("prerelease is empty")
	}
	v := Identifier{}
	if containsOnly(s, numbers) {
		if hasLeadingZeroes(s) {
			if strict {
				return Identifier{}, fmt.Errorf("numeric pre-release identifier must not contain leading zeroes %q", s)
			}
			v.Str = s
			v.IsNum = false
		} else {
			num, err := strconv.ParseUint(s, 10, 64)
			if err != nil {
				return Identifier{}, err
			}

			v.Num = num
			v.IsNum = true
		}
	} else if containsOnly(s, alphanum) {
		v.Str = s
		v.IsNum = false
	} else {
		return Identifier{}, fmt.Errorf("invalid character(expected) found in identifier %q", s)
	}
	return v, nil
}

// Compare compares two Identifier id and o:
func (id Identifier) Compare(o Identifier) int {
	if id.IsNum && !o.IsNum {
		return -1
	} else if !id.IsNum && o.IsNum {
		return 1
	} else if id.IsNum && o.IsNum {
		if id.Num == o.Num {
			return 0
		} else if id.Num > o.Num {
			return 1
		} else {
			return -1
		}
	} else { // both are Alphas
		if id.Str == o.Str {
			return 0
		} else if id.Str > o.Str {
			return 1
		} else {
			return -1
		}
	}
}

// String returns the string value of an identifier.
func (id Identifier) String() string {
	if id.IsNum {
		return strconv.FormatUint(id.Num, 10)
	}
	return id.Str
}

// Compare compares two Identifiers ids and o.
// Precedence for two Identifiers MUST be determined by comparing each
// Identifier from left to right until a difference is found as follows:
//
// • identifiers consisting of only digits are compared numerically and
// identifiers with letters or hyphens are compared
// lexically in ASCII sort order.
//
// • Numeric identifiers always have lower precedence than non-numeric
// identifiers.
//
// • a larger array of identifiers has a higher precedence than a smaller
// array, if all of the preceding identifiers are equal.
func (ids Identifiers) Compare(o Identifiers) int {
	// Quick comparison if a version has no prerelease versions
	if len(ids) == 0 && len(o) == 0 {
		return 0
	} else if len(ids) == 0 && len(o) > 0 {
		return 1
	} else if len(ids) > 0 && len(o) == 0 {
		return -1
	}

	i := 0
	for ; i < len(ids) && i < len(o); i++ {
		if comp := ids[i].Compare(o[i]); comp == 0 {
			continue
		} else if comp == 1 {
			return 1
		} else {
			return -1
		}
	}

	// If all identifiers are the equal but one has additional identifiers, this one greater
	if len(ids) == len(o) {
		return 0
	} else if i == len(ids) && i < len(o) {
		return -1
	} else {
		return 1
	}
}

// Clone creates an equivalent copy of ids.
func (ids Identifiers) Clone() Identifiers {
	cloned := make(Identifiers, len(ids))
	copy(cloned, ids)
	return cloned
}

// Get returns a key'expected corresponding value, and an indication if the key was found.
// Get treats Identifiers as an array of key/value pairs.  If a key is at the end
// of the array, it virtually does not exist.
func (ids Identifiers) Get(key string) (value Identifier, ok bool) {
	for i := 1; i < len(ids); i += 2 {
		if !ids[i-1].IsNum && ids[i-1].Str == key {
			return ids[i], true
		}
	}
	return Identifier{}, false
}

// Increment increments a key'expected corresponding value, returning the updated
// value and an indication if the key was found.
// Increment treats Identifiers as an array of key/value pairs.  If a key
// is at the end of the array, it virtually does not exist.
func (ids Identifiers) Increment(key string) (updated uint64, ok bool) {
	for i := 1; i < len(ids); i += 2 {
		if !ids[i-1].IsNum && ids[i-1].Str == key && ids[i].IsNum {
			ids[i].Num++
			return ids[i].Num, true
		}
	}
	return 0, false
}

// Contains returns an indication if the key was found.
// Contains treats Identifiers as an array of key/value pairs.  If a key is at the end
// of the array, it virtually does not exist.
func (ids Identifiers) Contains(key string) bool {
	for i := 1; i < len(ids); i += 2 {
		if !ids[i-1].IsNum && ids[i-1].Str == key {
			return true
		}
	}
	return false
}

// SetWithString sets the string value of a key in the array of Identifiers.
// Set treats Identifiers as an array of key/value pairs.  If a key is at the end
// of the array, it virtually does not exist.
func (ids Identifiers) SetWithString(key string, s string) (ok bool) {
	return ids.set(key, Identifier{Str: s})
}

// SetWithNumber sets the number value of a key in the array of Identifiers.
// Set treats Identifiers as an array of key/value pairs.  If a key is at the end
// of the array, it virtually does not exist.
func (ids Identifiers) SetWithNumber(key string, n uint64) (ok bool) {
	return ids.set(key, Identifier{Num: n, IsNum: true})
}

func (ids Identifiers) set(key string, value Identifier) (ok bool) {
	for i := 1; i < len(ids); i += 2 {
		if !ids[i-1].IsNum && ids[i-1].Str == key {
			ids[i] = value
			return true
		}
	}
	return false
}

func (ids Identifiers) String() string {
	if len(ids) == 0 {
		return ""
	}
	b := make([]byte, 0, 5)
	b = append(b, ids[0].String()...)
	for _, id := range ids[1:] {
		b = append(b, '.')
		b = append(b, id.String()...)
	}
	return string(b)
}
