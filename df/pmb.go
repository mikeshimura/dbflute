/*
 * Copyright 2014-2015 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language
 * governing permissions and limitations under the License.
 */
package df

import (
"database/sql"
)

type BasePmb struct{
	
}
func (b *BasePmb) CheckAndComvertEmptyToNull(value interface{} ) interface{}{
		if DBCurrent_I.EmptyStringParameterAllowed {
		return value
	}
	switch value.(type) {
	case sql.NullString:
		var nstr sql.NullString = value.(sql.NullString)
		if nstr.Valid && nstr.String == "" {
			nstr.Valid = false
		}
		return nstr
	case *sql.NullString:
		var nstr *sql.NullString = value.(*sql.NullString)
		if nstr.Valid && nstr.String == "" {
			nstr.Valid = false
		}
		return nstr
	case string:
		var str string = value.(string)
		if str == "" {
			var null sql.NullString
			null.Valid = false
			return null
		}
	case *string:
		var strx string = *value.(*string)
		if strx == "" {
			var null sql.NullString
			null.Valid = false
			return null
		}
	default:
		panic("This type not supported :" + GetType(value))
	}
	return value
}
