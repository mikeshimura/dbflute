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
	"database/sql/driver"
	"errors"
//	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	MYSQL_DEFAULT_DATE_FORMAT      = "2006-01-02"
	MYSQL_DEFAULT_TIME_FORMAT      = "15:04:05.000000"
	MYSQL_DEFAULT_TIMESTAMP_FORMAT = "2006-01-02 15:04:05.000000"
)

//Sql Numeric Custom Implementation
type Numeric struct {
	IntValue int64 //total value 10.50 -> 1050
	DecPoint int   //decimal point numeric(10,2) -> 2
}
func CreateNumeric(value string) (*Numeric, error) {
	var valuebu string = value
	num := new(Numeric)
	pos := strings.Index(value, ".")
	if pos == -1 {
		num.IntValue = 0
	} else {
		num.DecPoint = len(value) - pos -1
		value = value[0:pos] + value[pos+1:]
	}

	cv, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, errors.New("Numericに変換出来ません:" + valuebu)
	}
	num.IntValue = cv
	return num, nil
}
// Scan implements the Scanner interface.
func (nn *Numeric) Scan(value interface{}) error {
	if value == nil {
		err := errors.New("df001:Numeric Null Error")
		return err
	}
	var svalue = string(value.([]byte))
	var err error
	pos := strings.Index(svalue, ".")
	if nn.DecPoint > 0 {
		if pos != len(svalue)-nn.DecPoint-1 {
			err = errors.New("df002:Numeric DecPoint position unmatch:" + svalue)
			return err
		}
	} else {
		if pos > -1 {
			err = errors.New("df003:Numeric DecPoint=0 but . found:" + svalue)
			return err
		}
	}
	nn.IntValue, err = getIntFromSvalue(svalue)
	if err != nil {
		return err
	}
	return nil
}

// return Int part Dec part values and valid value
func (nn Numeric) GetValues() (int64, int64) {
	ivalue := nn.IntValue / nn.GetDivValue()
	dvalue := nn.IntValue - ivalue*nn.GetDivValue()
	return ivalue, dvalue
}
func (nn Numeric) String() string {
	ivalue, dvalue := nn.GetValues()
	return strconv.Itoa(int(ivalue)) + "." + strconv.Itoa(int(dvalue))
}

// numeric DecPoint=2 -> return 100
func (nn Numeric) GetDivValue() int64 {
	var v int64 = 1
	for i := 0; i < nn.DecPoint; i++ {
		v *= 10
	}
	return v
}
func (nn Numeric) Value() (driver.Value, error) {
	ivalue, dvalue := nn.GetValues()
	return strconv.Itoa(int(ivalue)) + "." + strconv.Itoa(int(dvalue)), nil
}

//type NullString struct {
//	String string //total value 10.50 -> 1050
//	Valid  bool   // Valid is true if Value not null
//}
//
//func (nn *NullString) Scan(value interface{}) error {
//	if value == nil {
//		nn.Valid = false
//		return nil
//	}
//	nn.Valid = true
//	switch value.(type) {
//	case string:
//		var svalue = value.(string)
//		nn.String = svalue
//		return nil
//	case time.Time:
//		var tvalue time.Time = value.(time.Time)
//		nn.String = tvalue.Format(DISP_SQL_DEFAULT_DATE_FORMAT)
//		return nil
//	case []uint8:
//		nn.String = (string)(value.([]uint8))
//		return nil
//	}
//	panic("this type not supported :" + fmt.Sprintf("%T", value))
//	return nil
//}
//func (nn *NullString) Value() (driver.Value, error) {
//	if !nn.Valid {
//		return nil, nil
//	}
//	var sv []uint8
//	sv = ([]uint8)(nn.String)
//	return sv, nil
//}

type NullNumeric struct {
	IntValue int64 //total value 10.50 -> 1050
	DecPoint int   //decimal point numeric(10,2) -> 2
	Valid    bool  // Valid is true if Value not null
}
func CreateNullNumeric(value string) (*NullNumeric, error) {
	num,err:=CreateNumeric(value)
	if err != nil{
		return nil,err
	}
	nnum:=new(NullNumeric)
	nnum.Valid=true
	nnum.IntValue=num.IntValue
	nnum.DecPoint=num.DecPoint
	return nnum,nil
}
func (nn *NullNumeric) String() string {
	if !nn.Valid {
		return "null"
	}
	ivalue, dvalue, _ := nn.GetValues()
	return strconv.Itoa(int(ivalue)) + "." + strconv.Itoa(int(dvalue))
}
func (nn *NullNumeric) Value() (driver.Value, error) {
	if !nn.Valid {
		return nil, nil
	}
	ivalue, dvalue, _ := nn.GetValues()
	return strconv.Itoa(int(ivalue)) + "." + strconv.Itoa(int(dvalue)), nil
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
			err = errors.New("df002:Numeric DecPoint position unmatch:" + svalue)
			return err
		}
	} else {
		if pos > -1 {
			err = errors.New("df003Numeric DecPoint=0 but . found:" + svalue)
			return err
		}
	}
	nn.IntValue, err = getIntFromSvalue(svalue)
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
	ivalue := nn.IntValue / nn.GetDivValue()
	dvalue := nn.IntValue - ivalue*nn.GetDivValue()
	return ivalue, dvalue, nn.Valid
}

// numeric DecPoint=2 -> return 100
func (nn NullNumeric) GetDivValue() int64 {
	var v int64 = 1
	for i := 0; i < nn.DecPoint; i++ {
		v *= 10
	}
	return v
}

func getIntFromSvalue(svalue string) (int64, error) {
	split := strings.Split(svalue, ".")
	var sint = ""
	if len(split) == 2 {
		sint = split[0] + split[1]
	} else {
		sint = split[0]
	}

	intValue, err := strconv.ParseInt(sint, 10, 64)
	if err != nil {
		return 0, err
	}
	return intValue, nil
}

type Date struct {
	Date time.Time
}

func (nt *Date) Scan(value interface{}) error {
	time, valid := value.(time.Time)
	nt.Date = time
	if !valid {
		err := errors.New("Date Error")
		return err
	}
	return nil
}
func (nt *Date) Value() (driver.Value, error) {
	return nt.Date, nil
}

type NullDate struct {
	Date  time.Time
	Valid bool
}

func (nt *NullDate) Scan(value interface{}) error {
	nt.Date, nt.Valid = value.(time.Time)
	return nil
}

func (nt NullDate) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Date, nil
}

type Timestamp struct {
	Timestamp time.Time
}

func (nt *Timestamp) Scan(value interface{}) error {
	time, valid := value.(time.Time)
	nt.Timestamp = time
	if !valid {
		err := errors.New("Timestamp Error")
		return err
	}
	return nil
}
func (nt *Timestamp) Value() (driver.Value, error) {
	return nt.Timestamp, nil
}

type NullTimestamp struct {
	Timestamp time.Time
	Valid     bool
}

func (nt *NullTimestamp) Scan(value interface{}) error {
	nt.Timestamp, nt.Valid = value.(time.Time)
	return nil
}

func (nt NullTimestamp) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Timestamp, nil
}

type MysqlDate struct {
	Date time.Time
}

// Scan implements the Scanner interface.
func (nn *MysqlDate) Scan(value interface{}) error {
	if value == nil {
		err := errors.New("MysqlDate Null Error")
		return err
	}
	var svalue = string(value.([]byte))
	var err error
	nn.Date, err = time.Parse(MYSQL_DEFAULT_DATE_FORMAT[0:len(svalue)], svalue)
	return err
}

func (nn MysqlDate) String() string {
	return nn.Date.Format(MYSQL_DEFAULT_DATE_FORMAT)
}

func (nn MysqlDate) Value() (driver.Value, error) {
	return nn.Date.Format(MYSQL_DEFAULT_DATE_FORMAT), nil
}

type MysqlTime struct {
	Time time.Time
}

// Scan implements the Scanner interface.
func (nn *MysqlTime) Scan(value interface{}) error {
	if value == nil {
		err := errors.New("MysqlTime Null Error")
		return err
	}
	var svalue = string(value.([]byte))
	var err error
	nn.Time, err = time.Parse(MYSQL_DEFAULT_TIME_FORMAT[0:len(svalue)], svalue)
	return err
}

func (nn MysqlTime) String() string {
	return nn.Time.Format(MYSQL_DEFAULT_TIME_FORMAT)
}

func (nn MysqlTime) Value() (driver.Value, error) {
	return nn.Time.Format(MYSQL_DEFAULT_TIME_FORMAT), nil
}

type MysqlTimestamp struct {
	Timestamp time.Time
}

// Scan implements the Scanner interface.
func (nn *MysqlTimestamp) Scan(value interface{}) error {
	if value == nil {
		err := errors.New("MysqlTimestamp Null Error")
		return err
	}
	var svalue = string(value.([]byte))
	var err error
	nn.Timestamp, err = time.Parse(MYSQL_DEFAULT_TIMESTAMP_FORMAT[0:len(svalue)], svalue)
	return err
}

func (nn MysqlTimestamp) String() string {
	return nn.Timestamp.Format(MYSQL_DEFAULT_TIMESTAMP_FORMAT)
}

func (nn MysqlTimestamp) Value() (driver.Value, error) {
	return nn.Timestamp.Format(MYSQL_DEFAULT_TIMESTAMP_FORMAT), nil
}

type MysqlNullDate struct {
	Date  time.Time
	Valid bool
}

// Scan implements the Scanner interface.
func (nn *MysqlNullDate) Scan(value interface{}) error {
	if value == nil {
		nn.Valid = false
		return nil
	}
	nn.Valid = true
	var svalue = string(value.([]byte))
	var err error
	nn.Date, err = time.Parse(MYSQL_DEFAULT_DATE_FORMAT[0:len(svalue)], svalue)
	return err
}

func (nn MysqlNullDate) String() string {
	if nn.Valid == false {
		return "null"
	}
	return nn.Date.Format(MYSQL_DEFAULT_DATE_FORMAT)
}

func (nn MysqlNullDate) Value() (driver.Value, error) {
	if nn.Valid == false {
		return nil, nil
	}
	return nn.Date.Format(MYSQL_DEFAULT_DATE_FORMAT), nil
}

type MysqlNullTime struct {
	Time  time.Time
	Valid bool
}

// Scan implements the Scanner interface.
func (nn *MysqlNullTime) Scan(value interface{}) error {
	if value == nil {
		nn.Valid = false
		return nil
	}
	nn.Valid = true
	var svalue = string(value.([]byte))
	var err error
	nn.Time, err = time.Parse(MYSQL_DEFAULT_TIME_FORMAT[0:len(svalue)], svalue)
	return err
}

func (nn MysqlNullTime) String() string {
	if nn.Valid == false {
		return "null"
	}
	return nn.Time.Format(MYSQL_DEFAULT_TIME_FORMAT)
}

func (nn MysqlNullTime) Value() (driver.Value, error) {
	if nn.Valid == false {
		return nil, nil
	}
	return nn.Time.Format(MYSQL_DEFAULT_TIME_FORMAT), nil
}

type MysqlNullTimestamp struct {
	Timestamp time.Time
	Valid     bool
}

// Scan implements the Scanner interface.
func (nn *MysqlNullTimestamp) Scan(value interface{}) error {
	if value == nil {
		nn.Valid = false
		return nil
	}
	nn.Valid = true
	var svalue = string(value.([]byte))
	var err error
	nn.Timestamp, err = time.Parse(MYSQL_DEFAULT_TIMESTAMP_FORMAT[0:len(svalue)], svalue)
	return err
}

func (nn MysqlNullTimestamp) String() string {
	if nn.Valid == false {
		return "null"
	}
	return nn.Timestamp.Format(MYSQL_DEFAULT_TIMESTAMP_FORMAT)
}

func (nn MysqlNullTimestamp) Value() (driver.Value, error) {
	if nn.Valid == false {
		return nil, nil
	}
	return nn.Timestamp.Format(MYSQL_DEFAULT_TIMESTAMP_FORMAT), nil
}
