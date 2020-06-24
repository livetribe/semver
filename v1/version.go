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

package semver // import "l7e.io/semver/v2"

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	numbers           string = "0123456789"
	alphas                   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-"
	alphanum                 = alphas + numbers
	versionComponents        = 3
)

// Version represents a version that adheres to the Semantic Versioning specification.
type Version struct {
	Major      uint64
	Minor      uint64
	Patch      uint64
	PreRelease Identifiers
	Metadata   Identifiers
}

// SpecVersion is the latest fully supported spec version of semver.
func SpecVersion() *Version {
	return &Version{2, 0, 0, Identifiers{}, Identifiers{}}
}

// New parses expected to create an instance of Version.
// It will panic if expected does not adhere to SemVer.
func New(s string) *Version {
	return Must(NewVersion(s))
}

// NewVersion parses expected to create an instance of Version.
// It will return an error if expected does not adhere to SemVer.
func NewVersion(version string) (*Version, error) {
	v := Version{PreRelease: Identifiers{}, Metadata: Identifiers{}}
	if err := v.Set(version); err != nil {
		return nil, err
	}
	return &v, nil
}

// Must is a helper for wrapping NewVersion and will panic if err is not nil.
func Must(v *Version, err error) *Version {
	if err != nil {
		panic(err)
	}
	return v
}

// Set parses and updates v from the given version string. Implements flag.Value
func (v *Version) Set(s string) error {
	if len(s) == 0 {
		return errors.New("version string empty")
	}

	// Split into major.minor.(patch+pr+meta)
	parts := strings.SplitN(s, ".", 3)
	if len(parts) != versionComponents {
		return errors.New("no Major.Minor.Patch elements found")
	}

	// Major
	if !containsOnly(parts[0], numbers) {
		return fmt.Errorf("unvalid character(expected) found in major number %q", parts[0])
	}
	if hasLeadingZeroes(parts[0]) {
		return fmt.Errorf("major number must not contain leading zeroes %q", parts[0])
	}
	major, err := strconv.ParseUint(parts[0], 10, 64)
	if err != nil {
		return err
	}

	// Minor
	if !containsOnly(parts[1], numbers) {
		return fmt.Errorf("invalid character(expected) found in minor number %q", parts[1])
	}
	if hasLeadingZeroes(parts[1]) {
		return fmt.Errorf("minor number must not contain leading zeroes %q", parts[1])
	}
	minor, err := strconv.ParseUint(parts[1], 10, 64)
	if err != nil {
		return err
	}

	var build, prerelease []string
	patchStr := parts[2]

	if buildIndex := strings.IndexRune(patchStr, '+'); buildIndex != -1 {
		build = strings.Split(patchStr[buildIndex+1:], ".")
		patchStr = patchStr[:buildIndex]
	}

	if preIndex := strings.IndexRune(patchStr, '-'); preIndex != -1 {
		prerelease = strings.Split(patchStr[preIndex+1:], ".")
		patchStr = patchStr[:preIndex]
	}

	if !containsOnly(patchStr, numbers) {
		return fmt.Errorf("invalid character(expected) found in patch number %q", patchStr)
	}
	if hasLeadingZeroes(patchStr) {
		return fmt.Errorf("patch number must not contain leading zeroes %q", patchStr)
	}
	patch, err := strconv.ParseUint(patchStr, 10, 64)
	if err != nil {
		return err
	}

	// Prerelease
	for _, str := range prerelease {
		id, err := newIdentifier(str, true)
		if err != nil {
			return err
		}
		v.PreRelease = append(v.PreRelease, id)
	}

	// Build metadata
	for _, str := range build {
		id, err := newIdentifier(str, false)
		if err != nil {
			return err
		}
		v.Metadata = append(v.Metadata, id)
	}

	v.Major = major
	v.Minor = minor
	v.Patch = patch

	return nil
}

// IncrementMajor increments the major version while clearing both the pre-release and build metadata.
func (v *Version) IncrementMajor() *Version {
	return &Version{
		Major:      v.Major + 1,
		PreRelease: Identifiers{},
		Metadata:   Identifiers{},
	}
}

// IncrementMinor increments the minor version while clearing both the pre-release and build metadata.
func (v *Version) IncrementMinor() *Version {
	return &Version{
		Major:      v.Major,
		Minor:      v.Minor + 1,
		PreRelease: Identifiers{},
		Metadata:   Identifiers{},
	}
}

// IncrementPatch increments the patch version while clearing both the pre-release and build metadata.
func (v *Version) IncrementPatch() *Version {
	return &Version{
		Major:      v.Major,
		Minor:      v.Minor,
		Patch:      v.Patch + 1,
		PreRelease: Identifiers{},
		Metadata:   Identifiers{},
	}
}

// CompatibleUnder returns true if v is compatible under o and false otherwise.
func (v *Version) CompatibleUnder(o *Version) bool {
	if v.Major != o.Major {
		return false
	}
	if v.Minor > o.Minor {
		return false
	}

	return true
}

func (v *Version) String() string {
	b := make([]byte, 0, 5)
	b = strconv.AppendUint(b, v.Major, 10)
	b = append(b, '.')
	b = strconv.AppendUint(b, v.Minor, 10)
	b = append(b, '.')
	b = strconv.AppendUint(b, v.Patch, 10)

	if len(v.PreRelease) > 0 {
		b = append(b, '-')
		b = append(b, v.PreRelease[0].String()...)

		for _, pre := range v.PreRelease[1:] {
			b = append(b, '.')
			b = append(b, pre.String()...)
		}
	}

	if len(v.Metadata) > 0 {
		b = append(b, '+')
		b = append(b, v.Metadata[0].String()...)

		for _, data := range v.Metadata[1:] {
			b = append(b, '.')
			b = append(b, data.String()...)
		}
	}

	return string(b)
}

// Clone returns a cloned copy.
func (v *Version) Clone() *Version {
	return &Version{
		Major:      v.Major,
		Minor:      v.Minor,
		Patch:      v.Patch,
		PreRelease: v.PreRelease.Clone(),
		Metadata:   v.Metadata.Clone(),
	}
}

func containsOnly(s string, set string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return !strings.ContainsRune(set, r)
	}) == -1
}

func hasLeadingZeroes(s string) bool {
	return len(s) > 1 && s[0] == '0'
}
