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
	"database/sql/driver"
	"fmt"
)

// Scan implements the database/sql.Scanner interface.
func (v *Version) Scan(src interface{}) error {
	var str string
	switch src := src.(type) {
	case string:
		str = src
	case []byte:
		str = string(src)
	default:
		return fmt.Errorf("cannot convert %T to string", src)
	}

	return v.Set(str)
}

// Value implements the database/sql/driver.Valuer interface.
func (v Version) Value() (driver.Value, error) {
	return v.String(), nil
}
