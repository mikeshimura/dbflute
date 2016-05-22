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
	//	"container/list"
	"github.com/mikeshimura/dbflute/log"
	"reflect"
	"strings"
//	"fmt"
)

var DBMetaProvider_I *DBMetaProvider

type DBMeta interface {
	GetProjectName() string
	GetDbCurrent() *DBCurrent
	GetTableDbName() string
	GetTableDispName() string
	GetTablePropertyName() string
	GetTableSqlName() *TableSqlName
	GetTableAlias() string
	GetPrimaryInfo() *PrimaryInfo
	GetColumnInfoList() *List
	GetColumnInfoMap() map[string]int
	GetColumnInfoByPropertyName(propetyName string) *ColumnInfo
	HasVersionNo() bool
	GetVersionNoColumnInfo() *ColumnInfo
	HasSequence() bool
	GetSequenceName() string
	GetSequenceIncrementSize() int
	GetSequenceCacheSize() int
	HasCommonColumn() bool
	GetCommonColumnInfoList() *List
	GetCommonColumnInfoBeforeInsertList() *List
	GetCommonColumnInfoBeforeUpdateList() *List
	GetPropertyType(propertyName string) *TnPropertyType
	HasIdentity() bool
	GetPrimaryUniqueInfo() *UniqueInfo
	GetSequenceNextValSql() string
	HasPrimaryKey() bool
	HasCompoundPrimaryKey() bool
	FindForeignInfo(foreignPropertyName string) *ForeignInfo
	CreateForeignInfoMap()
}

type BaseDBMeta struct {
	TableDbName       string
	TableDispName     string
	TablePropertyName string
	TableSqlName      *TableSqlName
	TableAlias        string
	ColumnInfoList    *List
	//ColumnInfoFlexibleMap        *StringKeyMap
	ColumnInfoMap                    map[string]int
	PrimaryInfo                      *PrimaryInfo
	VersionNoFlag                    bool
	VersionNoColumnInfo              *ColumnInfo
	SequenceFlag                     bool
	SequenceName                     string
	SequenceIncrementSize            int
	SequenceCacheSize                int
	CommonColumnFlag                 bool
	CommonColumnInfoList             *List
	CommonColumnInfoBeforeInsertList *List
	CommonColumnInfoBeforeUpdateList *List
	PrimaryKey                       bool
	CompoundPrimaryKey               bool
	Identity                         bool
	DBMeta                           *DBMeta
	ForeignInfoMap                   map[string]*ForeignInfo
}

func (b *BaseDBMeta) FindForeignInfo(foreignPropertyName string) *ForeignInfo {
	if b.ForeignInfoMap == nil {
		(*b.DBMeta).CreateForeignInfoMap()
	}
	fi:=b.ForeignInfoMap[foreignPropertyName]
	return fi
}
func (b *BaseDBMeta) Cfi(consName string, propName string, columns []*ColumnInfo,
	relno int, oneToOne bool, bizOne bool, asOne bool, addFK bool,
	fixedCond string, dynParamList *StringList, fixedInl bool,
	revsName string) *ForeignInfo {

	fi := new(ForeignInfo)
	fi.ConstraintName = consName
	fi.ForeignPropertyName = propName
	fi.ColumnInfos = columns
	fi.RelationNo = relno
	fi.OneToOne = oneToOne
	fi.BizOneToOne = bizOne
	fi.ReferrerAsOne = asOne
	fi.AdditionalFK = addFK
	fi.FixedCondition = fixedCond
	fi.DynamicParameterList = dynParamList
	fi.FixedInline = fixedInl
	fi.ReversePropertyName = revsName
	return fi
}

func (b *BaseDBMeta) CreateForeignInfoMap() {
	//Dummy Implementation
	return
}

func (b *BaseDBMeta) HasPrimaryKey() bool {
	return b.PrimaryKey
}
func (b *BaseDBMeta) HasCompoundPrimaryKey() bool {
	return b.CompoundPrimaryKey
}
func (b *BaseDBMeta) GetSequenceNextValSql() string {
	if !b.HasSequence() {
		return ""
	}
	sql := (*(*b.DBMeta).GetDbCurrent().DBWay).
		BuildSequenceNextValSql((*b.DBMeta).GetSequenceName())
	return sql

}

func (b *BaseDBMeta) GetPrimaryUniqueInfo() *UniqueInfo {
	if b.PrimaryInfo != nil {
		return b.PrimaryInfo.UniqueInfo
	}
	return nil
}
func (b *BaseDBMeta) HasIdentity() bool {
	return b.Identity
}

func (b *BaseDBMeta) GetColumnInfoMap() map[string]int {
	return b.ColumnInfoMap
}
func (b *BaseDBMeta) GetColumnInfoByPropertyName(propetyName string) *ColumnInfo {
	ci, ok := b.ColumnInfoMap[propetyName]
	if !ok {
		return nil
	}
	return (b.ColumnInfoList.Get(ci)).(*ColumnInfo)
}
func (b *BaseDBMeta) GetPropertyType(propertyName string) *TnPropertyType {
	pt := new(TnPropertyType)

	columnno, ok := b.ColumnInfoMap[propertyName]
	if ok {
		columnInfox := b.ColumnInfoList.Get(columnno)
		var columnInfo *ColumnInfo = columnInfox.(*ColumnInfo)
		pt.ColumnDbName = columnInfo.ColumnDbName
		pt.ColumnSqlName = columnInfo.ColumnSqlName
		pt.EntityColumnInfo = columnInfo
		pt.persistent = true
		pt.propetyName = columnInfo.PropertyName
		pt.GoType = columnInfo.GoType
	}
	if b.PrimaryInfo.UniqueInfo.Primary {
		for _, colInfo := range b.PrimaryInfo.UniqueInfo.UniqueColumnList.data {
			var ci *ColumnInfo = colInfo.(*ColumnInfo)
			if ci.PropertyName == propertyName {
				pt.primaryKey = true
			}
		}
	}
	return pt
}

func (b *BaseDBMeta) GetColumnInfoList() *List {
	return b.ColumnInfoList
}

//func (b *BaseDBMeta) GetColumnInfoFlexibleMap() *StringKeyMap {
//	return b.ColumnInfoFlexibleMap
//}
//
//func (b *BaseDBMeta) CreateColumnInfoFlexibleMap()  {
//	b.ColumnInfoFlexibleMap=CreateAsFlexible()
//	for _,col:= range b.ColumnInfoList.data{
//		var ci *ColumnInfo
//		ci = col.(*ColumnInfo)
//		b.ColumnInfoFlexibleMap.Put(ci.ColumnDbName,ci)
//	}
//}

func (b *BaseDBMeta) GetPrimaryInfo() *PrimaryInfo {
	return b.PrimaryInfo
}

func (b *BaseDBMeta) GetTableDbName() string {
	return b.TableDbName
}

func (b *BaseDBMeta) GetTableDispName() string {
	return b.TableDispName
}

func (b *BaseDBMeta) GetTablePropertyName() string {
	return b.TablePropertyName
}

func (b *BaseDBMeta) GetTableSqlName() *TableSqlName {
	return b.TableSqlName
}

func (b *BaseDBMeta) GetTableAlias() string {
	return b.TableAlias
}
func (b *BaseDBMeta) HasVersionNo() bool {
	return b.VersionNoFlag
}
func (b *BaseDBMeta) GetVersionNoColumnInfo() *ColumnInfo {
	return b.VersionNoColumnInfo
}
func (b *BaseDBMeta) HasSequence() bool {
	return b.SequenceFlag
}
func (b *BaseDBMeta) GetSequenceName() string {
	return b.SequenceName
}
func (b *BaseDBMeta) GetSequenceIncrementSize() int {
	return b.SequenceIncrementSize
}
func (b *BaseDBMeta) GetSequenceCacheSize() int {
	return b.SequenceCacheSize
}
func (b *BaseDBMeta) HasCommonColumn() bool {
	return b.CommonColumnFlag
}
func (b *BaseDBMeta) GetCommonColumnInfoList() *List {
	return b.CommonColumnInfoList
}
func (b *BaseDBMeta) GetCommonColumnInfoBeforeInsertList() *List {
	return b.CommonColumnInfoBeforeInsertList
}
func (b *BaseDBMeta) GetCommonColumnInfoBeforeUpdateList() *List {
	return b.CommonColumnInfoBeforeUpdateList
}

type StringKeyMap struct {
	SearchMap map[string]interface{}
	Flexible  bool
}

func CreateAsCaseInsensitive() *StringKeyMap {
	sm := new(StringKeyMap)
	sm.SearchMap = make(map[string]interface{})
	sm.Flexible = false
	return sm
}

func CreateAsFlexible() *StringKeyMap {
	sm := new(StringKeyMap)
	sm.SearchMap = make(map[string]interface{})
	sm.Flexible = true
	return sm
}

func (s *StringKeyMap) Get(key interface{}) interface{} {
	stringKey := s.ConvertStringKey(key)
	return s.SearchMap[stringKey]
}

func (s *StringKeyMap) Put(key string, value interface{}) {
	stringKey := s.ConvertStringKey(key)
	s.SearchMap[stringKey] = value
}

func (s *StringKeyMap) Remove(key interface{}) {
	stringKey := s.ConvertStringKey(key)
	delete(s.SearchMap, stringKey)
}

func (s *StringKeyMap) PutAll(m map[string]interface{}) {
	for key, value := range m {
		stringKey := s.ConvertStringKey(key)
		s.Put(stringKey, value)
	}
}

func (s *StringKeyMap) ContainKey(key interface{}) bool {
	stringKey := s.ConvertStringKey(key)
	_, ok := s.SearchMap[stringKey]
	return ok
}

func (s *StringKeyMap) Clear() {
	s.SearchMap = make(map[string]interface{})
}

func (s *StringKeyMap) ConvertStringKey(value interface{}) string {
	switch value.(type) {
	case string:
		return strings.ToLower(s.RemoveConnector(value.(string)))
	default:
		return ""
	}
}

func (s *StringKeyMap) Size() int {
	return len(s.SearchMap)
}

func (s *StringKeyMap) IsEmpty() bool {
	return len(s.SearchMap) == 0
}

func (s *StringKeyMap) KeyMap() []string {
	res := make([]string, len(s.SearchMap))
	i := 0
	for k, _ := range s.SearchMap {
		res[i] = k
		i++
	}
	return res
}

func (s *StringKeyMap) Values() []interface{} {
	res := make([]interface{}, len(s.SearchMap))
	i := 0
	for _, v := range s.SearchMap {
		res[i] = v
		i++
	}
	return res
}

func (s *StringKeyMap) ContainsValue(value interface{}) bool {
	for _, v := range s.SearchMap {
		if reflect.DeepEqual(v, value) {
			return true
		}
	}
	return false
}

func (s *StringKeyMap) RemoveConnector(str string) string {
	if str == "" {
		return ""
	}
	if s.Flexible == false {
		return str
	}
	if IsSingleQuoted(str) {
		str = UnquoteSingle(str)
	} else if IsDoubleQuoted(str) {
		str = UnquoteDouble(str)
	}
	str = strings.Replace(str, "_", "", -1)
	str = strings.Replace(str, "-", "", -1)
	str = strings.Replace(str, " ", "", -1)
	return str
}

func IsSingleQuoted(str string) bool {
	return len(str) > 1 && strings.Index(str, "'") == 0 &&
		strings.LastIndex(str, "'") == len(str)-1
}

func IsDoubleQuoted(str string) bool {
	return len(str) > 1 && strings.Index(str, "\"") == 0 &&
		strings.LastIndex(str, "\"") == len(str)-1
}

func UnquoteSingle(str string) string {
	if !IsSingleQuoted(str) {
		return str
	}
	return string(str[1 : len(str)-1])
}

func UnquoteDouble(str string) string {
	if !IsDoubleQuoted(str) {
		return str
	}
	return string(str[1 : len(str)-1])
}

type ColumnInfo struct {
	DbMeta                    *DBMeta
	ColumnDbName              string
	ColumnSqlName             *ColumnSqlName
	ColumnSynonym             string
	ColumnAlias               string
	PropertyName              string
	ObjectNativeType          string
	PropertyAccessType        string
	Primary                   bool
	AutoIncrement             bool
	NotNull                   bool
	ColumnDbType              string
	ColumnSize                int64
	DecimalDigits             int64
	DefaultValue              string
	GoType                    string //int64 float64 bool string Numeric および 各 Null型
	Seq                       int
	IsCommonColumn            bool
	OptimistickLock           string
	CommentForDBMetaSetting   string
	ForeignPropertyName       string
	ReferrerPropertyName      string
	ClassificationMetaSetting string
	CanBeColumnNullObject     bool
}

func CCI(DBMeta *DBMeta, ColumnDbName string, ColumnSqlName *ColumnSqlName,
	ColumnSynonym string, ColumnAlias string, ObjectNativeType string,
	PropertyName string, PropertyAccessType string, Primary bool,
	AutoIncrement bool, NotNull bool, ColumnDbType string,
	ColumnSize int64, DecimalDigits int64, DefaultValue string,
	IsCommonColumn bool, OptimistickLock string,
	CommentForDBMetaSetting string, ForeignPropertyName string,
	ReferrerPropertyName string, ClassificationMetaSetting string,
	CanBeColumnNullObject bool, GoType string) *ColumnInfo {

	log.InternalDebug("ColumnSqlName :" + ColumnSqlName.ColumnSqlName)
	ci := new(ColumnInfo)
	ci.DbMeta = DBMeta
	ci.ColumnDbName = ColumnDbName
	ci.ColumnSqlName = ColumnSqlName
	ci.ColumnSynonym = ColumnSynonym
	ci.ColumnAlias = ColumnAlias
	ci.PropertyName = PropertyName
	ci.ObjectNativeType = ObjectNativeType
	ci.PropertyAccessType = PropertyAccessType
	ci.Primary = Primary
	ci.AutoIncrement = AutoIncrement
	ci.NotNull = NotNull
	ci.ColumnDbType = ColumnDbType
	ci.ColumnSize = ColumnSize
	ci.DecimalDigits = DecimalDigits
	ci.DefaultValue = DefaultValue
	ci.GoType = GoType
	ci.IsCommonColumn = IsCommonColumn
	ci.OptimistickLock = OptimistickLock
	ci.CommentForDBMetaSetting = CommentForDBMetaSetting
	ci.ForeignPropertyName = ForeignPropertyName
	ci.ReferrerPropertyName = ReferrerPropertyName
	ci.ClassificationMetaSetting = ClassificationMetaSetting
	ci.CanBeColumnNullObject = CanBeColumnNullObject
	return ci
}

type UniqueInfo struct {
	DbMeta           *DBMeta
	UniqueColumnList *List
	Primary          bool
}

type PrimaryInfo struct {
	UniqueInfo *UniqueInfo
}
type TableSqlName struct {
	TableSqlName        string
	CorrespondingDbName string
}
type DBMetaInstanceHandler struct {
	TableDbNameInstanceMap map[string]*DBMeta
}

var DBMetaInstanceHandler_I *DBMetaInstanceHandler

func CreateDBMetaInstanceHandle() {
	DBMetaInstanceHandler_I = new(DBMetaInstanceHandler)
	DBMetaInstanceHandler_I.TableDbNameInstanceMap = make(map[string]*DBMeta)
}

type TnPropertyType struct {
	ColumnDbName     string
	ColumnSqlName    *ColumnSqlName
	EntityColumnInfo *ColumnInfo
	persistent       bool
	primaryKey       bool
	propetyName      string
	GoType           string
}

type ForeignInfo struct {
	ConstraintName       string
	ForeignPropertyName  string
	ColumnInfos          []*ColumnInfo
	RelationNo           int
	OneToOne             bool
	BizOneToOne          bool
	ReferrerAsOne        bool
	AdditionalFK         bool
	FixedCondition       string
	DynamicParameterList *StringList
	FixedInline          bool
	ReversePropertyName  string
}
func (p *ForeignInfo)IsPureFK() bool{
	return !p.AdditionalFK && !p.ReferrerAsOne
}
func (p *ForeignInfo)IsNotNullFKColumn() bool{
	for i:=0;i<len(p.ColumnInfos);i+=2{
		localColumnInfo:=p.ColumnInfos[i]
		if localColumnInfo.NotNull==false{
			return false
		}
	}
	return true
}
