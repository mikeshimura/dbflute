package sqltype

import (
	"errors"
	"strconv"
	"strings"
)
//Sql Numeric Custom Implementation
type NullNumeric struct {
	IntValue int64 //total value 10.50 -> 1050
	DecPoint int   //decimal point numeric(10,2) -> 2
	Valid    bool  // Valid is true if Value not null
}

// Scan implements the Scanner interface.
func (nn *NullNumeric) Scan(value interface{}) error {
	if value == nil {
		nn.IntValue = 0
		nn.Valid = false
		return nil
	}
	nn.Valid = true
	var svalue = string(value.([]byte))
	var err error
	pos := strings.Index(svalue, ".")
	if nn.DecPoint > 0 {
		if pos != len(svalue)-nn.DecPoint-1 {
			err = errors.New("Numeric DecPoint position unmatch:" + svalue)
			return err
		}
	} else {
		if pos > -1 {
			err = errors.New("Numeric DecPoint=0 but . found:" + svalue)
			return err
		}
	}
	split := strings.Split(svalue, ".")
	var sint = ""
	if len(split) == 2 {
		sint = split[0] + split[1]
	} else {
		sint = split[0]
	}

	nn.IntValue, err = strconv.ParseInt(sint, 10, 64)
	if err != nil {
		return err
	}

	return nil
}

// return Int part Dec part values and valid value
func (nn NullNumeric) GetValues() (int64, int64, bool) {
	if !nn.Valid {
		return int64(0), int64(0), nn.Valid
	}
	ivalue:=nn.IntValue/nn.GetDivValue()
	dvalue:=nn.IntValue - ivalue * nn.GetDivValue()
	return ivalue,dvalue, nn.Valid
}
// numeric DecPoint=2 -> return 100
func (nn NullNumeric) GetDivValue()(int64){
		var v int64 = 1
		for i:=0;i<nn.DecPoint;i++{
			v *=10
		}
		return v
}