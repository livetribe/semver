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

// Compare tests if v is less than, equal to, or greater than other,
// returning -1, 0, or +1 respectively.
func (v *Version) Compare(other *Version) int {
	if v.Major != other.Major {
		if v.Major > other.Major {
			return 1
		}
		return -1
	}
	if v.Minor != other.Minor {
		if v.Minor > other.Minor {
			return 1
		}
		return -1
	}
	if v.Patch != other.Patch {
		if v.Patch > other.Patch {
			return 1
		}
		return -1
	}

	return v.PreRelease.Compare(other.PreRelease)
}

// Equals checks if v is equal to o.
func (v *Version) Equals(o *Version) bool {
	return v.Compare(o) == 0
}

// EQ checks if v is equal to o.
func (v *Version) EQ(o *Version) bool {
	return v.Compare(o) == 0
}

// NE checks if v is not equal to o.
func (v *Version) NE(o *Version) bool {
	return v.Compare(o) != 0
}

// GT checks if v is greater than o.
func (v *Version) GT(o *Version) bool {
	return v.Compare(o) > 0
}

// GTE checks if v is greater than or equal to o.
func (v *Version) GTE(o *Version) bool {
	return v.Compare(o) >= 0
}

// GE checks if v is greater than or equal to o.
func (v *Version) GE(o *Version) bool {
	return v.Compare(o) >= 0
}

// LT checks if v is less than o.
func (v *Version) LT(o *Version) bool {
	return v.Compare(o) < 0
}

// LTE checks if v is less than or equal to o.
func (v *Version) LTE(o *Version) bool {
	return v.Compare(o) <= 0
}

// LE checks if v is less than or equal to o.
func (v *Version) LE(o *Version) bool {
	return v.Compare(o) <= 0
}
