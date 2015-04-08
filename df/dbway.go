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
	"fmt"
	"reflect"
)

type DBWay interface {
	BuildSequenceNextValSql(s string) string
	GetIdentitySelectSql() string
	IsBlockCommentSupported() bool
	IsLineCommentSupported() bool
	IsScrollableCursorSupported() bool
	GetOriginalWildCardList() []string
	IsUniqueConstraintException(sqlState string, errorCode int64) bool
	Connect(v interface{}) string
	GetPlaceholderType() string
	GetStringConnector() *StringConnector
}

//DBWay interfaceを実装
type WayOfPostgreSQL struct {
	stringConnector *StringConnector
}

func (w *WayOfPostgreSQL) GetStringConnector() *StringConnector {
	if w.stringConnector == nil {
		ssc := new(StandardStringConnector)
		var sc StringConnector = ssc
		w.stringConnector = &sc
	}
	return w.stringConnector
}
func (w *WayOfPostgreSQL) BuildSequenceNextValSql(sequenceName string) string {
	return "select nextval ('" + sequenceName + "')"
}
func (w *WayOfPostgreSQL) GetIdentitySelectSql() string {
	return ""
}
func (w *WayOfPostgreSQL) IsBlockCommentSupported() bool {
	return true
}
func (w *WayOfPostgreSQL) IsLineCommentSupported() bool {
	return true
}
func (w *WayOfPostgreSQL) IsScrollableCursorSupported() bool {
	return true
}
func (w *WayOfPostgreSQL) GetOriginalWildCardList() []string {
	return []string{}
}
func (w *WayOfPostgreSQL) IsUniqueConstraintException(sqlState string, errorCode int64) bool {
	return sqlState == "23505"
}
func (w *WayOfPostgreSQL) Connect(v interface{}) string {
	buf := ""
	s := reflect.ValueOf(v)
	for i := 0; i < s.Len(); i++ {
		if len(buf) > 0 {
			buf += " || "
		}
		buf += fmt.Sprint(s.Index(i).Interface())
	}
	return buf
}
func (w *WayOfPostgreSQL) GetPlaceholderType() string {
	return "$1"
}

type WayOfMySQL struct {
	stringConnector *StringConnector
}

func (w *WayOfMySQL) GetStringConnector() *StringConnector {
	if w.stringConnector == nil {
		ssc := new(MysqlStringConnector)
		var sc StringConnector = ssc
		w.stringConnector = &sc
	}
	return w.stringConnector
}
func (w *WayOfMySQL) BuildSequenceNextValSql(sequenceName string) string {
	return ""
}
func (w *WayOfMySQL) GetIdentitySelectSql() string {
	return "SELECT LAST_INSERT_ID()"
}
func (w *WayOfMySQL) IsBlockCommentSupported() bool {
	return true
}
func (w *WayOfMySQL) IsLineCommentSupported() bool {
	return true
}
func (w *WayOfMySQL) IsScrollableCursorSupported() bool {
	return true
}
func (w *WayOfMySQL) GetOriginalWildCardList() []string {
	return []string{}
}
func (w *WayOfMySQL) IsUniqueConstraintException(sqlState string, errorCode int64) bool {
	return errorCode == 1062
}
func (w *WayOfMySQL) Connect(v interface{}) string {
	buf := ""
	s := reflect.ValueOf(v)
	buf +="concat("
	for i := 0; i < s.Len(); i++ {
		if len(buf) > 0 {
			buf += ", "
		}
		buf += fmt.Sprint(s.Index(i).Interface())
	}
	return buf+")"
}
func (w *WayOfMySQL) GetPlaceholderType() string {
	return "?"
}
type WayOfSQLServer struct {
	stringConnector *StringConnector
}

func (w *WayOfSQLServer) GetStringConnector() *StringConnector {
	if w.stringConnector == nil {
		ssc := new(MysqlStringConnector)
		var sc StringConnector = ssc
		w.stringConnector = &sc
	}
	return w.stringConnector
}
func (w *WayOfSQLServer) BuildSequenceNextValSql(sequenceName string) string {
	return ""
}
func (w *WayOfSQLServer) GetIdentitySelectSql() string {
	return "select @@identity"
}
func (w *WayOfSQLServer) IsBlockCommentSupported() bool {
	return true
}
func (w *WayOfSQLServer) IsLineCommentSupported() bool {
	return true
}
func (w *WayOfSQLServer) IsScrollableCursorSupported() bool {
	return true
}
func (w *WayOfSQLServer) GetOriginalWildCardList() []string {
	return []string{"[","]"}
}
func (w *WayOfSQLServer) IsUniqueConstraintException(sqlState string, errorCode int64) bool {
	return errorCode == 2627
}
func (w *WayOfSQLServer) Connect(v interface{}) string {
	buf := ""
	s := reflect.ValueOf(v)
	for i := 0; i < s.Len(); i++ {
		if len(buf) > 0 {
			buf += " + "
		}
		buf += fmt.Sprint(s.Index(i).Interface())
	}
	return buf
}
func (w *WayOfSQLServer) GetPlaceholderType() string {
	return "?"
}
