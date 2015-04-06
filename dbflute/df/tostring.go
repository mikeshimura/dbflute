package df

import (
	"strconv"
)

func ToStringInterface(val interface{}) string {
	var res string = ""
	switch val.(type) {
	case string:
		res = val.(string)
	case int64:
		res = strconv.Itoa(int(val.(int64)))
	case float64:
		res = strconv.FormatFloat(val.(float64), 'f', -1, 64)
		//todo datetime
	}
	return res
}
