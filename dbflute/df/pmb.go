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