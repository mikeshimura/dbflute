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
	"errors"
	"fmt"
	"github.com/mikeshimura/dbflute/log"
	"reflect"
)

type Behavior interface {
	GetBaseBehavior() *BaseBehavior
	ReadNextVal(tx *sql.Tx) int64
	GetDBMeta() *DBMeta
}
type BaseBehavior struct {
	BehaviorCommandInvoker *BehaviorCommandInvoker
	TableDbName            string
	Behavior               *Behavior
}

func (b *BaseBehavior) DoSelectCount(cb interface{},
	tx *sql.Tx) (reult int64, errrtn error) {
	var err error
	defer func() {
		errx := recover()
		if errx != nil {
			errrtn = errors.New(errx.(string))
		}
	}()
	cbbase := reflect.ValueOf(cb).Elem().FieldByName("BaseConditionBean").Interface()
	var base *BaseConditionBean = cbbase.(*BaseConditionBean)
	(*base.SqlClause).GetBaseSqlClause().selectClauseType = Create_SelectClauseType(
		SelectClauseType_UNIQUE_COUNT)
	cmd := b.CreateSelectListCBCommand(cb, "D_Int64", tx)
	var behcmd BehaviorCommand = cmd
	invres := b.Invoke(&behcmd)
	if invres == nil {
		return 0, err
	}
	res:=invres.(*ListResultBean)
	return (res.List.Get(0)).(*D_Int64).GetValue(), err
}

func (b *BaseBehavior) DoSelectList(cb interface{}, entityType string,
	tx *sql.Tx) (bean *ListResultBean, errrtn error) {
	var err error
	defer func() {
		errx := recover()
		if errx != nil {
			errrtn = errors.New(errx.(string))
		}
	}()
	cmd := b.CreateSelectListCBCommand(cb, entityType, tx)
	var behcmd BehaviorCommand = cmd
	invres := b.Invoke(&behcmd)
	if invres == nil {
		return nil, err
	}
	return invres.(*ListResultBean), err
}
func (b *BaseBehavior) ReadNextVal(tx *sql.Tx) int64 {

	return -1
}
func (b *BaseBehavior) DoOutsideSql() *OutsideSqlBasicExecutor {
	return b.BehaviorCommandInvoker.createOutsideSqlBasicExecutor(
		b.AsTableDbName(), b.Behavior)
}
func (b *BaseBehavior) AsTableDbName() string {
	return b.TableDbName
}

func (b *BaseBehavior) DoSelectNextVal(tx *sql.Tx) int64 {
	invres := b.Invoke(b.createSelectNextValCommand(tx))
	res := invres.(*ListResultBean)
	var ent *D_Int64 = (res.List.Get(0)).(*D_Int64)
	return ent.value
}
func (b *BaseBehavior) DoDelete(entity *Entity, option *DeleteOption,
	tx *sql.Tx, ctx *Context) (no int64, errrtn error) {
	var err error
	defer func() {
		errx := recover()
		if errx != nil {
			errrtn = errors.New(errx.(string))
		}
	}()
	res := b.processBeforeDelete(entity, option, tx, ctx)

	if !res {
		return 0, err
	}
	var invres interface{}
	invres = b.Invoke(b.createDeleteEntityCommand(entity, option, tx))
	return invres.(int64), err
}
func (b *BaseBehavior) DoInsert(entity *Entity, option *InsertOption,
	tx *sql.Tx, ctx *Context) (no int64, errrtn error) {
	var err error
	defer func() {
		errx := recover()
		if errx != nil {
			fmt.Printf("errx%v %T \n", errx, errx)
			errrtn = errors.New(errx.(string))
		}
	}()
	res := b.processBeforeInsert(entity, option, tx, ctx)

	if !res {
		return 0, err
	}
	var invres interface{}
	invres = b.Invoke(b.createInsertEntityCommand(entity, option, tx))
	return invres.(int64), err
}
func (b *BaseBehavior) DoUpdate(entity *Entity, option *UpdateOption,
	tx *sql.Tx, ctx *Context) (no int64, errrtn error) {
	var err error
	defer func() {
		errx := recover()
		if errx != nil {
			errrtn = errors.New(errx.(string))
		}
	}()
	res := b.processBeforeUpdate(entity, option, ctx)
	if !res {
		return 0, err
	}
	var invres interface{}
	invres = b.Invoke(b.createUpdateEntityCommand(entity, option, tx))
	if invres == nil {
		return 0, err
	}
	return invres.(int64), err
}
func (b *BaseBehavior) createDeleteEntityCommand(entity *Entity, option *DeleteOption,
	tx *sql.Tx) *BehaviorCommand {
	cmd := new(DeleteEntityCommand)
	cmd.entity = entity
	cmd.tx = tx
	cmd.Behavior = b.Behavior
	b.XsetupSelectCommand(&cmd.BaseBehaviorCommand)
	var bc BehaviorCommand = cmd
	cmd.BehaviorCommand = &bc
	return &bc
}
func (b *BaseBehavior) createInsertEntityCommand(entity *Entity,
	option *InsertOption, tx *sql.Tx) *BehaviorCommand {
	cmd := new(InsertEntityCommand)
	cmd.entity = entity
	cmd.tx = tx
	cmd.Behavior = b.Behavior
	b.XsetupSelectCommand(&cmd.BaseBehaviorCommand)
	var bc BehaviorCommand = cmd
	cmd.BehaviorCommand = &bc
	return &bc
}
func (b *BaseBehavior) createSelectNextValCommand(tx *sql.Tx) *BehaviorCommand {
	cmd := new(SelectNextValCommand)
	b.XsetupSelectCommand(&cmd.BaseBehaviorCommand)
	cmd.dbmeta = (*b.Behavior).GetDBMeta()
	cmd.Behavior = b.Behavior
	cmd.tx = tx
	var bc BehaviorCommand = cmd
	return &bc
}
func (b *BaseBehavior) createUpdateEntityCommand(entity *Entity, option *UpdateOption,
	tx *sql.Tx) *BehaviorCommand {

	//	        assertBehaviorCommandInvoker("createUpdateEntityCommand");
	cmd := new(UpdateEntityCommand)
	var behavior BehaviorCommand = cmd
	cmd.BehaviorCommand = &behavior
	cmd.tx = tx
	b.xsetupEntityCommand(&cmd.BaseEntityCommand, entity, (*entity).AsTableDbName())
	//        cmd.setUpdateOption(option);
	var bcmd BehaviorCommand = cmd
	return &bcmd
}
func (b *BaseBehavior) processBeforeDelete(entity *Entity, option *DeleteOption,
	tx *sql.Tx, ctx *Context) bool {

	//        filterEntityOfDelete(entity, option);
	//        assertEntityOfDelete(entity, option);
	//        return true;
	b.frameworkFilterEntityOfDelete(entity, option, tx, ctx)
	if !(*(*entity).GetDBMeta()).HasIdentity() {
		b.assertEntityNotNullAndHasPrimaryKeyValue(entity)
	}
	return true
}
func (b *BaseBehavior) processBeforeInsert(entity *Entity, option *InsertOption,
	tx *sql.Tx, ctx *Context) bool {
	//        assertEntityNotNull(entity); // primary key is checked later
	//        frameworkFilterEntityOfInsert(entity, option);
	//        filterEntityOfInsert(entity, option);
	//        assertEntityOfInsert(entity, option);
	//        // check primary key after filtering at an insert process
	//        // because a primary key value may be set in filtering process
	//        // (for example, sequence)
	//        if (!entity.getDBMeta().hasIdentity()) { // identity does not need primary key value here
	//            assertEntityNotNullAndHasPrimaryKeyValue(entity);
	//        }
	//        return true;
	if entity == nil {
		panic("Entity Null")
	}
	b.frameworkFilterEntityOfInsert(entity, option, tx, ctx)
	if !(*(*entity).GetDBMeta()).HasIdentity() {
		b.assertEntityNotNullAndHasPrimaryKeyValue(entity)
	}

	return true
}
func (b *BaseBehavior) frameworkFilterEntityOfDelete(entity *Entity,
	option *DeleteOption, tx *sql.Tx, ctx *Context) {
	b.setupCommonColumnOfUpdateIfNeeds(entity, ctx)
	return
}
func (b *BaseBehavior) frameworkFilterEntityOfInsert(entity *Entity,
	option *InsertOption, tx *sql.Tx, ctx *Context) {
	b.injectSequenceToPrimaryKeyIfNeeds(entity, tx)
	b.setupCommonColumnOfInsertIfNeeds(entity, ctx)
}
func (b *BaseBehavior) setupCommonColumnOfInsertIfNeeds(
	entity *Entity, ctx *Context) {
	(*CommonColumnAutoSetupper_I).HandleCommonColumnOfInsertIfNeeds(entity, ctx)
}
func (b *BaseBehavior) setupCommonColumnOfUpdateIfNeeds(
	entity *Entity, ctx *Context) {
	(*CommonColumnAutoSetupper_I).HandleCommonColumnOfUpdateIfNeeds(entity, ctx)
}
func (b *BaseBehavior) injectSequenceToPrimaryKeyIfNeeds(entity *Entity,
	tx *sql.Tx) {
	dbmeta := (*entity).GetDBMeta()

	if !(*dbmeta).HasSequence() || (*dbmeta).HasCompoundPrimaryKey() ||
		(*entity).HasPrimaryKeyValue() || (*dbmeta).HasIdentity() {
		return
	}
	// basically property(column) type is same as next value type
	// so there is NOT type conversion cost when writing to the entity
	col := ((*dbmeta).GetPrimaryUniqueInfo().UniqueColumnList.Get(0)).(*ColumnInfo)
	nextVal := (*b.Behavior).ReadNextVal(tx)
	log.InternalDebug("next val " + fmt.Sprintf("%v", nextVal))
	SetEntityValue(entity, col.PropertyName, nextVal)
	return
}
func (b *BaseBehavior) processBeforeUpdate(entity *Entity,
	option *UpdateOption, ctx *Context) bool {
	b.assertEntityNotNullAndHasPrimaryKeyValue(entity)
	b.frameworkFilterEntityOfUpdate(entity, option, ctx)
	//未実装
	//        filterEntityOfUpdate(entity, option);
	return true
}
func (b *BaseBehavior) frameworkFilterEntityOfUpdate(
	entity *Entity, option *UpdateOption, ctx *Context) {
	b.setupCommonColumnOfUpdateIfNeeds(entity, ctx)
}
func (b *BaseBehavior) assertEntityNotNullAndHasPrimaryKeyValue(
	entity *Entity) {
	if entity == nil {
		panic("Entity nil")
	}
	//通常はNULL TYPEで無いのでCK 不要
	//	if !(*entity).HasPrimaryKeyValue() {
	//		return errors.New("EntityPrimaryKeyNotFound")
	//	}

	//        b.assertEntityOfUpdate(entity, option);
	return
}
func (b *BaseBehavior) GetBehaviorCommandInvoker() *BehaviorCommandInvoker {
	return b.BehaviorCommandInvoker
}
func (b *BaseBehavior) Invoke(cmd *BehaviorCommand) interface{} {
	log.InternalDebug("Invoke")
	return b.BehaviorCommandInvoker.Invoke(cmd)
}
func (b *BaseBehavior) CreateSelectListCBCommand(cb interface{},
	entityType string, tx *sql.Tx) *SelectListCBCommand {
	//assert 省略
	cmd := new(SelectListCBCommand)
	var behavior BehaviorCommand = cmd
	cmd.BehaviorCommand = &behavior
	b.XsetupSelectCommand(&cmd.BaseBehaviorCommand)
	cmd.ConditionBean = cb
	cmd.EntityType = entityType
	cmd.Behavior = b.Behavior
	cmd.tx = tx
	return cmd
}

func (b *BaseBehavior) XsetupSelectCommand(cmd *BaseBehaviorCommand) {
	cmd.TableDbName = (*b.Behavior).GetBaseBehavior().AsTableDbName()
	(*b.Behavior).GetBaseBehavior().GetBehaviorCommandInvoker().
		InjectComponentProperty(cmd)
}

type InvokerAssistant interface {
	Create()
	GetStatementFactory() *StatementFactory
	GetDBMetaProvider() *DBMetaProvider
	GetDBCurrent() *DBCurrent
	AssistOutsideSqlExecutorFactory() *OutsideSqlExecutorFactory
}

func (b *BaseBehavior) CreateBehaviorCommandInvoker() {
	b.BehaviorCommandInvoker = new(BehaviorCommandInvoker)
}

func (b *BaseBehavior) xsetupEntityCommand(cmd *BaseEntityCommand, entity *Entity,
	tableDbName string) {
	cmd.entity = entity
	cmd.TableDbName = tableDbName
	(*b.Behavior).GetBaseBehavior().GetBehaviorCommandInvoker().
		InjectComponentProperty(&cmd.BaseBehaviorCommand)
}

type StatementFactory interface {
	PrepareStatement(orgSql string, tx *sql.Tx, dbc *DBCurrent) *sql.Stmt
	ModifyBindVariables(bindVariables *List, bindVariableTypes *StringList) *List
}

type TnStatementFactoryImpl struct {
}

func (t *TnStatementFactoryImpl) ModifyBindVariables(bindVariables *List,
	bindVariableTypes *StringList) *List {
	if bindVariables == nil {
		return bindVariables
	}
	//convert df.NullString to sql.NullString
	//	for i, item := range bindVariables.data {
	//		stype := GetType(item)
	//		if stype == "df.NullString" {
	//			var dns NullString = item.(NullString)
	//			ns := new(sql.NullString)
	//			ns.Valid = dns.Valid
	//			ns.String = dns.String
	//			bindVariables.data[i] = ns
	//		}
	//	}
	return bindVariables
}
func (t *TnStatementFactoryImpl) PrepareStatement(orgSql string, tx *sql.Tx,
	dbc *DBCurrent) *sql.Stmt {
	sql := t.modifySql(orgSql, dbc)
	//fmt.Printf("sql %s tx %v %T\n", sql, tx, tx)
	stmt, errs := tx.Prepare(sql)
	if errs != nil {
		panic(errs.Error() + ":" + sql)
	}
	return stmt
}
func (t *TnStatementFactoryImpl) modifySql(sql string, dbc *DBCurrent) string {
	if (*dbc.DBWay).GetPlaceholderType() == "$1" {
		return ModifyToDollerPlaceholder(sql)
	}
	return sql
}

type TnBeanMetaDataFactory struct {
}

var Gopath string

type BhvUtil struct {
	entityMap map[string]func() *Entity
}

var BhvUtil_I *BhvUtil

func (b *BhvUtil) GetEntityAndInterfaceArray(name string) (
	interface{}, []interface{}) {
	entityp := b.entityMap[name]()
	var entity Entity = *entityp
	return entity, entity.GetAsInterfaceArray()

}
func (b *BhvUtil) SetUp() {
	b.entityMap = make(map[string]func() *Entity)
}
func (b *BhvUtil) AddEntity(ename string, ef func() *Entity) {
	log.InternalDebug("AddEntity :" + ename)
	b.entityMap[ename] = ef
}
func (b *BhvUtil) GetListResultBean(rows *sql.Rows, entity string,
	sqlClause interface{}) *ListResultBean {
		
	list := new(ListResultBean)
	list.New()
	for rows.Next() {
		//table, array := (*behavior).GetEntityAndInterfaceArray(t.ResultType)
		table, array := b.GetEntityAndInterfaceArray(entity)
		err := rows.Scan(array...)
		if err != nil {
			panic(err.Error())
		}
		list.List.Add(table)
	}
	list.AllRecordCount = list.List.Size()
	log.InternalDebug("entity:" + entity)
	tmap := (*DBMetaProvider_I).TableDbNameInstanceMap[entity]
	if tmap == nil {
		list.TableDbName = "__" + entity
	} else {
		list.TableDbName = (*tmap).GetTableDbName()
	}

	return list
}
