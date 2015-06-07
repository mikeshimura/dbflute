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
	"fmt"
	"github.com/lib/pq"
	"time"
)

type Stack struct {
	data     []interface{}
	position int
}

func (s *Stack) Push(value interface{}) {
	if s.data == nil {
		s.data = make([]interface{}, 0, 100)
	}
	if cap(s.data) < len(s.data)+1 {
		newSlice := make([]interface{}, len(s.data), cap(s.data)*2)
		copy(newSlice, s.data)
		s.data = newSlice
	}
	s.data = append(s.data, value)
	s.position++
}
func (s *Stack) Pop() interface{} {
	if s.position == 0 {
		return nil
	}
	s.position--
	value := s.data[s.position]
	s.data = s.data[:s.position]
	return value

}
func CreateList() *List {
	l := new(List)
	l.data = make([]interface{}, 0, 100)
	return l
}

type List struct {
	data []interface{}
}

func (s *List) Add(value interface{}) {
	if s.data == nil {
		s.data = make([]interface{}, 0, 100)
	}
	if cap(s.data) < len(s.data)+1 {
		newSlice := make([]interface{}, len(s.data), cap(s.data)*2)
		copy(newSlice, s.data)
		s.data = newSlice
	}
	s.data = append(s.data, value)
}
func (s *List) Get(i int) interface{} {
	if i < len(s.data) {
		return s.data[i]
	} else {
		return nil
	}
}
func (s *List) Size() int {
	return len(s.data)
}
func (s *List) GetAsArray() []interface{} {
	if s.data == nil {
		s.data = make([]interface{}, 0, 100)
	}
	return s.data
}
func GetType(o interface{}) string {
	return fmt.Sprintf("%T", o)
}
func CreateStringList() *StringList {
	l := new(StringList)
	l.data = make([]string, 0, 100)
	return l
}

type StringList struct {
	data []string
}

func (s *StringList) Add(value string) {
	if s.data == nil {
		s.data = make([]string, 0, 50)
	}
	if cap(s.data) < len(s.data)+1 {
		newSlice := make([]string, len(s.data), cap(s.data)*2)
		copy(newSlice, s.data)
		s.data = newSlice
	}
	s.data = append(s.data, value)
}
func (s *StringList) Get(i int) string {
	if i < len(s.data) {
		return s.data[i]
	} else {
		return ""
	}
}
func (s *StringList) GetAsArray() []string {
	if s.data == nil {
		s.data = make([]string, 0, 50)
	}
	return s.data
}
func (s *StringList) Size() int {
	return len(s.data)
}
func CreateNullString(s string) sql.NullString {
	var ns sql.NullString
	ns.Valid = true
	ns.String = s
	return ns
}
func CreateNullInt64(v int64) sql.NullInt64 {
	var ni sql.NullInt64
	ni.Valid = true
	ni.Int64 = v
	return ni
}
func CreateNullFloat64(v float64) sql.NullFloat64 {
	var ni sql.NullFloat64
	ni.Valid = true
	ni.Float64 = v
	return ni
}
func CreateNullBool(v bool) sql.NullBool {
	var nb sql.NullBool
	nb.Valid = true
	nb.Bool = v
	return nb
}
func CreateNullTime(v time.Time) pq.NullTime {
	var nt pq.NullTime
	nt.Valid = true
	nt.Time = v
	return nt
}
func CreateNullDate(v time.Time) NullDate {
	var nt NullDate
	nt.Valid = true
	nt.Date = v
	return nt
}
func CreateNullTimestamp(v time.Time) NullTimestamp {
	var nt NullTimestamp
	nt.Valid = true
	nt.Timestamp = v
	return nt
}
func CreateDate(v time.Time) Date {
	var nt Date
	nt.Date = v
	return nt
}
func CreateTimestamp(v time.Time) Timestamp {
	var nt Timestamp
	nt.Timestamp = v
	return nt
}
func CreateMysqlTimestamp(v time.Time) MysqlTimestamp {
	var nt MysqlTimestamp
	nt.Timestamp = v
	return nt
}
func CreateMysqlTime(v time.Time) MysqlTime {
	var nt MysqlTime
	nt.Time = v
	return nt
}
func CreateMysqlDate(v time.Time) MysqlDate {
	var nt MysqlDate
	nt.Date = v
	return nt
}
