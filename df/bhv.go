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
)

type Behavior interface {
	GetBaseBehavior() *BaseBehavior
	ReadNextVal(tx *sql.Tx) (int64, error)
	GetDBMeta() *DBMeta
}
type BaseBehavior struct {
	BehaviorCommandInvoker *BehaviorCommandInvoker
	TableDbName            string
	Behavior               *Behavior
}

func (b *BaseBehavior) DoSelectList(cb interface{}, entityType string, tx *sql.Tx) (*ListResultBean, error) {
	cmd := b.CreateSelectListCBCommand(cb, entityType, tx)
	var behcmd BehaviorCommand = cmd
	invres, err := b.Invoke(&behcmd)
	if invres==nil{
		return nil,err
	}
	return invres.(*ListResultBean), err
}
func (b *BaseBehavior) ReadNextVal(tx *sql.Tx) (int64, error) {
	return -1, nil
}
func (b *BaseBehavior) DoOutsideSql() *OutsideSqlBasicExecutor {
	return b.BehaviorCommandInvoker.createOutsideSqlBasicExecutor(b.AsTableDbName(), b.Behavior)
}
func (b *BaseBehavior) AsTableDbName() string {
	return b.TableDbName
}

func (b *BaseBehavior) DoSelectNextVal(tx *sql.Tx) (int64, error) {
	invres, err1 := b.Invoke(b.createSelectNextValCommand(tx))
	if err1 != nil {
		return 0, err1
	}
	res := invres.(*ListResultBean)
	var ent *D_Int64 = (res.List.Get(0)).(*D_Int64)
	return ent.value, nil
}
func (b *BaseBehavior) DoDelete(entity *Entity, option *DeleteOption, tx *sql.Tx) (int64, error) {

	res, err := b.processBeforeDelete(entity, option, tx)

	if err != nil {
		return 0, err
	}
	if !res {
		return 0, nil
	}
	var invres interface{}
	invres, err1 := b.Invoke(b.createDeleteEntityCommand(entity, option, tx))
	if err1 != nil {
		return 0, err1
	}
	return invres.(int64), err1
}
func (b *BaseBehavior) DoInsert(entity *Entity, option *InsertOption, tx *sql.Tx) (int64, error) {

	res, err := b.processBeforeInsert(entity, option, tx)

	if err != nil {
		return 0, err
	}
	if !res {
		return 0, nil
	}
	var invres interface{}
	invres, err1 := b.Invoke(b.createInsertEntityCommand(entity, option, tx))
	if err1 != nil {
		return 0, err1
	}
	return invres.(int64), err1
}
func (b *BaseBehavior) DoUpdate(entity *Entity, option *UpdateOption, tx *sql.Tx) (int64, error) {
	res, err := b.processBeforeUpdate(entity, option)
	if err != nil {
		return 0, err
	}
	if !res {
		return 0, nil
	}
	invres, err1 := b.Invoke(b.createUpdateEntityCommand(entity, option, tx))
	if err1 != nil {
		return 0, err1
	}
	return invres.(int64), err1
}
func (b *BaseBehavior) createDeleteEntityCommand(entity *Entity, option *DeleteOption, tx *sql.Tx) *BehaviorCommand {
	cmd := new(DeleteEntityCommand)
	cmd.entity = entity
	cmd.tx = tx
	cmd.Behavior = b.Behavior
	b.XsetupSelectCommand(&cmd.BaseBehaviorCommand)
	var bc BehaviorCommand = cmd
	cmd.BehaviorCommand = &bc
	return &bc
}
func (b *BaseBehavior) createInsertEntityCommand(entity *Entity, option *InsertOption, tx *sql.Tx) *BehaviorCommand {
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
func (b *BaseBehavior) createUpdateEntityCommand(entity *Entity, option *UpdateOption, tx *sql.Tx) *BehaviorCommand {

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
func (b *BaseBehavior) processBeforeDelete(entity *Entity, option *DeleteOption, tx *sql.Tx) (bool, error) {

	//        filterEntityOfDelete(entity, option);
	//        assertEntityOfDelete(entity, option);
	//        return true;
	err := b.frameworkFilterEntityOfDelete(entity, option, tx)
	if err != nil {
		return true, err
	}
	if !(*(*entity).GetDBMeta()).HasIdentity() {
		b.assertEntityNotNullAndHasPrimaryKeyValue(entity)
	}

	return true, nil
}
func (b *BaseBehavior) processBeforeInsert(entity *Entity, option *InsertOption, tx *sql.Tx) (bool, error) {
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
		return true, errors.New("Entity Null")
	}
	err := b.frameworkFilterEntityOfInsert(entity, option, tx)
	if err != nil {
		return true, err
	}
	if !(*(*entity).GetDBMeta()).HasIdentity() {
		b.assertEntityNotNullAndHasPrimaryKeyValue(entity)
	}

	return true, nil
}
func (b *BaseBehavior) frameworkFilterEntityOfDelete(entity *Entity, option *DeleteOption, tx *sql.Tx) error {
	return nil
}
func (b *BaseBehavior) frameworkFilterEntityOfInsert(entity *Entity, option *InsertOption, tx *sql.Tx) error {
	return b.injectSequenceToPrimaryKeyIfNeeds(entity, tx)
	//        setupCommonColumnOfInsertIfNeeds(entity, option);
}
func (b *BaseBehavior) injectSequenceToPrimaryKeyIfNeeds(entity *Entity, tx *sql.Tx) error {
	dbmeta := (*entity).GetDBMeta()

	if !(*dbmeta).HasSequence() || (*dbmeta).HasCompoundPrimaryKey() || (*entity).HasPrimaryKeyValue()||  (*dbmeta).HasIdentity(){
		return nil
	}
	// basically property(column) type is same as next value type
	// so there is NOT type conversion cost when writing to the entity
	col := ((*dbmeta).GetPrimaryUniqueInfo().UniqueColumnList.Get(0)).(*ColumnInfo)
	nextVal, err := (*b.Behavior).ReadNextVal(tx)
	if err != nil {
		return err
	}
	log.InternalDebug("next val " + fmt.Sprintf("%v", nextVal))
	SetEntityValue(entity, col.PropertyName, nextVal)
	return nil
}
func (b *BaseBehavior) processBeforeUpdate(entity *Entity, option *UpdateOption) (bool, error) {
	err := b.assertEntityNotNullAndHasPrimaryKeyValue(entity)
	if err != nil {
		return false, err
	}
	//未実装 setupCommonColumnOfUpdateIfNeeds
	//	      frameworkFilterEntityOfUpdate(entity, option);
	//未実装
	//        filterEntityOfUpdate(entity, option);
	//未実装
	return true, nil
}
func (b *BaseBehavior) assertEntityNotNullAndHasPrimaryKeyValue(entity *Entity) error {
	if entity == nil {
		return errors.New("Entity nil")
	}
	//通常はNULL TYPEで無いのでCK 不要
	//	if !(*entity).HasPrimaryKeyValue() {
	//		return errors.New("EntityPrimaryKeyNotFound")
	//	}

	//        b.assertEntityOfUpdate(entity, option);
	return nil
}
func (b *BaseBehavior) GetBehaviorCommandInvoker() *BehaviorCommandInvoker {
	return b.BehaviorCommandInvoker
}
func (b *BaseBehavior) Invoke(cmd *BehaviorCommand) (interface{}, error) {
	log.InternalDebug("Invoke")
	return b.BehaviorCommandInvoker.Invoke(cmd)
}
func (b *BaseBehavior) CreateSelectListCBCommand(cb interface{}, entityType string, tx *sql.Tx) *SelectListCBCommand {
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
	(*b.Behavior).GetBaseBehavior().GetBehaviorCommandInvoker().InjectComponentProperty(cmd)
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

func (b *BaseBehavior) xsetupEntityCommand(cmd *BaseEntityCommand, entity *Entity, tableDbName string) {
	cmd.entity = entity
	cmd.TableDbName = tableDbName
	(*b.Behavior).GetBaseBehavior().GetBehaviorCommandInvoker().InjectComponentProperty(&cmd.BaseBehaviorCommand)
}

type StatementFactory interface {
	PrepareStatement(orgSql string, tx *sql.Tx, dbc *DBCurrent) (*sql.Stmt, error)
	ModifyBindVariables(bindVariables *List, bindVariableTypes *StringList) *List
}

type TnStatementFactoryImpl struct {
}

func (t *TnStatementFactoryImpl) ModifyBindVariables(bindVariables *List, bindVariableTypes *StringList) *List {
	if bindVariables==nil{
		return bindVariables
	}
	//convert df.NullString to sql.NullString
	for i,item:=range bindVariables.data{
		stype:=GetType(item)
		if stype=="df.NullString" {
			var dns NullString=item.(NullString)
			ns:=new(sql.NullString)
			ns.Valid=dns.Valid
			ns.String=dns.String
			bindVariables.data[i]=ns
		}
	}
	return bindVariables
}
func (t *TnStatementFactoryImpl) PrepareStatement(orgSql string, tx *sql.Tx, dbc *DBCurrent) (*sql.Stmt, error) {
	sql := t.modifySql(orgSql, dbc)
	//fmt.Printf("sql %s tx %v %T\n", sql, tx, tx)
	stmt, errs := tx.Prepare(sql)
	if errs != nil {
		return nil, errors.New(errs.Error() + ":" + sql)
	}
	return stmt, nil
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

func (b *BhvUtil) GetEntityAndInterfaceArray(name string) (interface{}, []interface{}) {
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
func (b *BhvUtil) GetListResultBean(rows *sql.Rows, entity string) (*ListResultBean, error) {
	list := new(ListResultBean)
	list.New()
	for rows.Next() {
		//table, array := (*behavior).GetEntityAndInterfaceArray(t.ResultType)
		table, array := b.GetEntityAndInterfaceArray(entity)
		err := rows.Scan(array...)
		if err != nil {
			return nil, err
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

	return list, nil
}
