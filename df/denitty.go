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

)

type D_Int64 struct {
	value int64
	BaseEntity
}

func (l *D_Int64) GetValue() int64 {
	return l.value
}

func (t *D_Int64) GetAsInterfaceArray() []interface{} {
	i := make([]interface{}, 1)
	i[0] = &(t.value)
	return i
}

func (t *D_Int64) AsTableDbName() string {
	return "D_Int64"
}

func (t *D_Int64) HasPrimaryKeyValue() bool {
	return false
}
func (t *D_Int64) SetValue(value int64) {
	t.AddPropertyName("value")
	t.value = value
}

func (t *D_Int64) SetUp() {

}
func (t *D_Int64) GetDBMeta() *DBMeta {
	return DBMetaInstanceHandler_I.TableDbNameInstanceMap[t.AsTableDbName()]
}
type D_Int64Dbm_T struct {
	BaseDBMeta
	ColumnValue        *ColumnInfo
}

func (b *D_Int64Dbm_T) GetProjectName() string {
	return DBCurrent_I.ProjectName
}

func (b *D_Int64Dbm_T) GetDbCurrent() *DBCurrent {
	return DBCurrent_I
}

var D_Int64Dbm *D_Int64Dbm_T

func Create_D_Int64Dbm() {
	D_Int64Dbm = new(D_Int64Dbm_T)
	D_Int64Dbm.TableDbName = "d_Int64"
	D_Int64Dbm.TableDispName = "d_Int64"
	D_Int64Dbm.TablePropertyName = "d_Int64"
	D_Int64Dbm.TableSqlName = new(TableSqlName)
	D_Int64Dbm.TableSqlName.TableSqlName = "d_Int64"
	D_Int64Dbm.TableSqlName.CorrespondingDbName = D_Int64Dbm.TableDbName

	var dm DBMeta
	dm = D_Int64Dbm
	valueSqlName := new(ColumnSqlName)
	valueSqlName.ColumnSqlName = "value"
	valueSqlName.IrregularChar = false
	D_Int64Dbm.ColumnValue = CCI(&dm, "value", valueSqlName, "", "", "Integer.class","value", "", false, false, true, "int4", 10, 0, "",false,"","","","","",false, "int64")

	D_Int64Dbm.ColumnInfoList = new(List)
	D_Int64Dbm.ColumnInfoList.Add(D_Int64Dbm.ColumnValue)


	D_Int64Dbm.ColumnInfoMap=make(map[string]int)
	D_Int64Dbm.ColumnInfoMap["value"]=0

	ui := new(UniqueInfo)
	ui.DbMeta = &dm
	ui.Primary = true
	ui.UniqueColumnList = new(List)
	ui.UniqueColumnList.Add(D_Int64Dbm.ColumnValue)

	D_Int64Dbm.PrimaryInfo = nil
	D_Int64Dbm.VersionNoFlag = false
	D_Int64Dbm.SequenceFlag = false
	D_Int64Dbm.SequenceName = ""
	D_Int64Dbm.SequenceIncrementSize = 0
	D_Int64Dbm.SequenceCacheSize = 0
	D_Int64Dbm.CommonColumnFlag = false
	D_Int64Dbm.CommonColumnInfoList = new(List)
	D_Int64Dbm.CommonColumnInfoBeforeInsertList = new(List)
	D_Int64Dbm.CommonColumnInfoBeforeUpdateList = new(List)
	var dmap DBMeta = D_Int64Dbm
	DBMetaInstanceHandler_I.TableDbNameInstanceMap["d_Int64"] = &dmap
}
