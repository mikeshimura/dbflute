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
	//"container/list"
	//	"bytes"
	"database/sql"
	//"errors"
	"fmt"
	"github.com/mikeshimura/dbflute/log"
	//	"path/filepath"
	//	"reflect"
	"strconv"
	"strings"
)

type ParameterHandler interface {
	execute(bindVariables *List, bindVariableTypes *StringList, tx *sql.Tx, behavior *Behavior) interface{}
}
type TnBasicUpdateHandler struct {
	TnBasicParameterHandler
}

func (t *TnBasicUpdateHandler) execute(bindVariables *List, bindVariableTypes *StringList, tx *sql.Tx, behavior *Behavior) interface{} {
	dbc := (*(*behavior).GetBaseBehavior().GetBehaviorCommandInvoker().InvokerAssistant).GetDBCurrent()

	ps := (*t.statementFactory).PrepareStatement(t.sql, tx, dbc)
	defer ps.Close()
	bindVar := (*t.statementFactory).ModifyBindVariables(bindVariables, bindVariableTypes)
	//ps := t.prepareStatement(tx, dbc)
	//	ns:=new(sql.NullString)
	//	ns.Valid=true
	//	ns.String="2"
	//	var itest interface{}=ns
	t.logSql(bindVariables, bindVariableTypes)
	res, err := tx.Stmt(ps).Exec(bindVar.data...)
	if err != nil {
		panic(err.Error())
	}

	updateno, _ := res.RowsAffected()
	log.InternalDebug(fmt.Sprintln("result no:", updateno))

	//        logSql(args, argTypes);
	//        PreparedStatement ps = null;
	//        try {
	//            ps = prepareStatement(conn);
	//            bindArgs(conn, ps, args, argTypes);
	//            return queryResult(ps);
	//        } catch (SQLException e) {
	//            final SQLExceptionResource resource = createSQLExceptionResource();
	//            resource.setNotice("Failed to execute the SQL for select.");
	//            handleSQLException(e, resource);
	//            return null; // unreachable
	//        } finally {
	//            close(ps);
	//        }
	return updateno
}

type TnBasicSelectHandler struct {
	TnBasicParameterHandler
	ResultType string
	SqlClause  interface{}
}

func (t *TnBasicSelectHandler) execute(bindVariables *List, bindVariableTypes *StringList, tx *sql.Tx, behavior *Behavior) interface{} {
	dbc := (*(*behavior).GetBaseBehavior().GetBehaviorCommandInvoker().InvokerAssistant).GetDBCurrent()

	ps := (*t.statementFactory).PrepareStatement(t.sql, tx, dbc)
	defer ps.Close()
	bindVar := (*t.statementFactory).ModifyBindVariables(bindVariables, bindVariableTypes)
	//ps := t.prepareStatement(tx, dbc)
	//	ns:=new(sql.NullString)
	//	ns.Valid=true
	//	ns.String="2"
	//	var itest interface{}=ns
	t.logSql(bindVariables, bindVariableTypes)
	var rows *sql.Rows
	var err error
	if bindVariables == nil {
		rows, err = tx.Stmt(ps).Query()
		if err != nil {
			panic(err.Error())
		}
	} else {
		rows, err = tx.Stmt(ps).Query(bindVar.data...)
		if err != nil {
			panic(err.Error())
		}
	}

	log.InternalDebug(fmt.Sprintln("ResultType:", t.ResultType))
	rh := new(ResultSetHandler)
	//l := BhvUtil_I.GetListResultBean(rows, t.ResultType,t.SqlClause)
	l := rh.GetListResultBean(rows, t.ResultType, t.SqlClause)
	log.InternalDebug(fmt.Sprintln("result no:", l.List.Size()))
	log.InternalDebug(fmt.Sprintf("data %v\n", l.List.Get(0)))
	//        logSql(args, argTypes);
	//        PreparedStatement ps = null;
	//        try {
	//            ps = prepareStatement(conn);
	//            bindArgs(conn, ps, args, argTypes);
	//            return queryResult(ps);
	//        } catch (SQLException e) {
	//            final SQLExceptionResource resource = createSQLExceptionResource();
	//            resource.setNotice("Failed to execute the SQL for select.");
	//            handleSQLException(e, resource);
	//            return null; // unreachable
	//        } finally {
	//            close(ps);
	//        }
	return l
}
func (t *TnBasicSelectHandler) logSql(args *List, argTypes *StringList) {
	//	        final boolean logEnabled = isLogEnabled();
	//        final boolean hasSqlFireHook = hasSqlFireHook();
	//        final boolean hasSqlLog = hasSqlLogHandler();
	//        final boolean hasSqlResult = hasSqlResultHandler();
	//        final Object sqlLogRegistry = getSqlLogRegistry();
	//        final boolean hasRegistry = sqlLogRegistry != null;
	//
	//        if (logEnabled || hasSqlFireHook || hasSqlLog || hasSqlResult || hasRegistry) {
	//            if (isInternalDebugEnabled()) {
	//                final String determination = logEnabled + ", " + hasSqlFireHook + ", " + hasSqlLog + ", "
	//                        + hasSqlResult + ", " + hasRegistry;
	//                _log.debug("...Logging SQL by " + determination);
	//            }
	//            if (processBeforeLogging(args, argTypes, logEnabled, hasSqlFireHook, hasSqlLog, hasSqlResult,
	//                    sqlLogRegistry)) {
	//                return; // processed by anyone
	//            }
	//            doLogSql(args, argTypes, logEnabled, hasSqlFireHook, hasSqlLog, hasSqlResult, sqlLogRegistry);
	//        }
	if !LogStop {
		t.doLogSql(args, argTypes)
	}
}

func ModifyToDollerPlaceholder(sql string) string {
	parai := 1
	for pos := strings.Index(sql, "?"); pos > -1; pos = strings.Index(sql, "?") {
		sql = sql[0:pos] + "$" + strconv.Itoa(parai) + sql[pos+1:]
		parai++
	}
	return sql
}

type TnBasicParameterHandler struct {
	TnAbstractBasicSqlHandler
}
type BasicSqlHander interface {
	processSuccess(entity *Entity, updateno int64, lastInsertId int64)
	setupBindVariables(entity *Entity)
}
type TnAbstractBasicSqlHandler struct {
	statementFactory *StatementFactory
	sql              string
	BasicSqlHander   *BasicSqlHander
}

func (t *TnAbstractBasicSqlHandler) logSql(args *List, argTypes *StringList) {
	//        final boolean logEnabled = isLogEnabled();
	//        final boolean hasSqlFireHook = hasSqlFireHook();
	//        final boolean hasSqlLog = hasSqlLogHandler();
	//        final boolean hasSqlResult = hasSqlResultHandler();
	//        final Object sqlLogRegistry = getSqlLogRegistry();
	//        final boolean hasRegistry = sqlLogRegistry != null;
	//
	//        if (logEnabled || hasSqlFireHook || hasSqlLog || hasSqlResult || hasRegistry) {
	//            if (isInternalDebugEnabled()) {
	//                final String determination = logEnabled + ", " + hasSqlFireHook + ", " + hasSqlLog + ", "
	//                        + hasSqlResult + ", " + hasRegistry;
	//                _log.debug("...Logging SQL by " + determination);
	//            }
	//            if (processBeforeLogging(args, argTypes, logEnabled, hasSqlFireHook, hasSqlLog, hasSqlResult,
	//                    sqlLogRegistry)) {
	//                return; // processed by anyone
	//            }
	//            doLogSql(args, argTypes, logEnabled, hasSqlFireHook, hasSqlLog, hasSqlResult, sqlLogRegistry);
	//        }
	if !LogStop {
		t.doLogSql(args, argTypes)
	}
}
func (t *TnAbstractBasicSqlHandler) doLogSql(args *List, argTypes *StringList) {
	//        final boolean hasRegistry = sqlLogRegistry != null;
	//        final String firstDisplaySql;
	//        if (logEnabled || hasRegistry) { // build at once
	//            if (isInternalDebugEnabled()) {
	//                _log.debug("...Building DisplaySql by " + logEnabled + ", " + hasRegistry);
	//            }
	//            firstDisplaySql = buildDisplaySql(_sql, args);
	firstDisplaySql := t.buildDisplaySql(t.sql, args)
	DFLog(firstDisplaySql)
	//            if (logEnabled) {
	//                logDisplaySql(firstDisplaySql);
	//            }
	//            if (hasRegistry) { // S2Container provides
	//                pushToSqlLogRegistry(args, argTypes, firstDisplaySql, sqlLogRegistry);
	//            }
	//        } else {
	//            firstDisplaySql = null;
	//        }
	//        if (hasSqlFireHook || hasSqlLog || hasSqlResult) { // build lazily
	//            if (isInternalDebugEnabled()) {
	//                _log.debug("...Handling SqlFireHook or SqlLog or SqlResult by " + hasSqlFireHook + ", " + hasSqlLog
	//                        + ", " + hasSqlResult);
	//            }
	//            final SqlLogInfo sqlLogInfo = prepareSqlLogInfo(args, argTypes, firstDisplaySql);
	//            if (sqlLogInfo != null) { // basically true (except override)
	//                if (hasSqlLog) {
	//                    getSqlLogHander().handle(sqlLogInfo);
	//                }
	//                if (hasSqlFireHook) {
	//                    saveHookSqlLogInfo(sqlLogInfo);
	//                }
	//                if (hasSqlResult) {
	//                    saveResultSqlLogInfo(sqlLogInfo);
	//                }
	//            }
	//        }

}
func (t *TnAbstractBasicSqlHandler) buildDisplaySql(sql string, args *List) string {
	log.InternalDebug("sql :" + sql)
	db := new(DisplaySqlBuilder)
	return db.BuildDisplaySql(sql, args)
}

type TnInsertEntityHandler struct {
	TnAbstractEntityHandler
}

func (t *TnInsertEntityHandler) setupBindVariables(entity *Entity) {
	t.setupInsertBindVariables(entity)
}
func (t *TnInsertEntityHandler) processSuccess(entity *Entity, updateno int64, lastInsertId int64) {
	//	        super.processSuccess(bean, ret);
	t.doProcessIdentity(entity, lastInsertId) //
	//     	new IdentityProcessCallback() {
	//            public void callback(TnIdentifierGenerator generator) {
	//                if (generator.isPrimaryKey() && isPrimaryKeyIdentityDisabled()) {
	//                    return;
	//                }
	//                generator.setIdentifier(bean, _dataSource);
	//            }
	//        });
	//        updateVersionNoIfNeed(bean);
	//        updateTimestampIfNeed(bean);
}
func (t *TnInsertEntityHandler) doProcessIdentity(entity *Entity, lastInsertId int64) {
	dbmeta := (*entity).GetDBMeta()
	dbway := (*dbmeta).GetDbCurrent().DBWay
	if (*dbmeta).HasIdentity() == false {
		return
	}
	sql := (*dbway).GetIdentitySelectSql()
	if sql == "" {
		return
	}
	primary := (*dbmeta).GetPrimaryInfo()
	primaryCol := (primary.UniqueInfo.UniqueColumnList.Get(0)).(*ColumnInfo)
	if primaryCol.AutoIncrement {
		propertyName := primaryCol.PropertyName
		SetEntityValue(entity, propertyName, lastInsertId)
	}
}

type TnUpdateEntityHandler struct {
	TnAbstractEntityHandler
}

func (t *TnUpdateEntityHandler) setupBindVariables(entity *Entity) {
	t.setupUpdateBindVariables(entity)
}
func (t *TnUpdateEntityHandler) processSuccess(entity *Entity, updateno int64, lastInsertId int64) {
	if t.newVersionNoList != nil && t.newVersionNoList.Size() > 0 {
		ok := SetEntityValue(entity, (*(*entity).GetDBMeta()).GetVersionNoColumnInfo().
			PropertyName, t.newVersionNoList.Get(0))
		if !ok {
		}
	}
}

type TnDeleteEntityHandler struct {
	TnAbstractEntityHandler
}

func (t *TnDeleteEntityHandler) setupBindVariables(entity *Entity) {
	t.setupDeleteBindVariables(entity)
}

type TnAbstractEntityHandler struct {
	TnAbstractBasicSqlHandler
	optimisticLockHandling         bool
	versionNoAutoIncrementOnMemory bool
	updateoption                   *UpdateOption
	insertoption                   *InsertOption
	boundPropTypes                 *List
	newVersionNoList               *List
	bindVariables                  *List
	bindVariableValueTypes         *StringList
}

func (t *TnAbstractEntityHandler) setupBindVariables(entity *Entity) {

}
func (t *TnAbstractEntityHandler) processSuccess(entity *Entity, updateno int64, lastInsertId int64) {

}
func (t *TnAbstractEntityHandler) setupDeleteBindVariables(entity *Entity) {
	//	       final List<Object> varList = new ArrayList<Object>();
	//        final List<ValueType> varValueTypeList = new ArrayList<ValueType>();
	//        addAutoUpdateWhereBindVariables(varList, varValueTypeList, bean);
	//        _bindVariables = varList.toArray();
	//        _bindVariableValueTypes = (ValueType[]) varValueTypeList.toArray(new ValueType[varValueTypeList.size()]);
	varList := new(List)
	varValueTypeList := new(StringList)
	t.addAutoUpdateWhereBindVariables(varList, varValueTypeList, entity)
	t.bindVariables = varList
	t.bindVariableValueTypes = varValueTypeList
}

func (t *TnAbstractEntityHandler) execute(entity *Entity, tx *sql.Tx) int64 {
	//        processBefore(bean);
	(*t.BasicSqlHander).setupBindVariables(entity)
	t.logSql(t.bindVariables, t.bindVariableValueTypes)
	ps := (*t.statementFactory).PrepareStatement(t.sql, tx, (*(*entity).GetDBMeta()).GetDbCurrent())
	defer ps.Close()
	bindVar := (*t.statementFactory).ModifyBindVariables(t.bindVariables, t.bindVariableValueTypes)
	//        RuntimeException sqlEx = null;
	//        final int ret;
	//        try {
	//            bindArgs(conn, ps, _bindVariables, _bindVariableValueTypes);
	//            ret = executeUpdate(ps);
	//fmt.Printf("tx %v ps %v bindVar %v\n", tx, ps, bindVar)
	res, err := tx.Stmt(ps).Exec(bindVar.data...)
	if err != nil {
		panic(err.Error())
	}
	updateno, _ := res.RowsAffected()
	idno, _ := res.LastInsertId()
	log.InternalDebug(fmt.Sprintf("idno %d \n", idno))
	t.handleUpdateResultWithOptimisticLock(entity, updateno)
	//        } catch (RuntimeException e) {
	//            // not SQLFailureException because the JDBC wrapper may throw an other exception
	//            sqlEx = e;
	//            throw e;
	//        } finally {
	//            close(ps);
	//            processFinally(bean, sqlEx);
	//        }
	//        // a value of exclusive control column should be synchronized
	//        // after handling optimistic lock
	//        processSuccess(bean, ret);
	//        return ret;
	//fmt.Printf("BasicSqlHander %v entity %v updateno %v\n", t.BasicSqlHander, entity, updateno)
	(*t.BasicSqlHander).processSuccess(entity, updateno, idno)
	return updateno
}

func (t *TnAbstractEntityHandler) handleUpdateResultWithOptimisticLock(entity *Entity, updateno int64) {
	if t.optimisticLockHandling && updateno < 1 { // means no update (contains minus just in case)
		panic("EntityAlreadyUpdatedException")
	}
	return
}
func (t *TnAbstractEntityHandler) setupUpdateBindVariables(entity *Entity) {
	//	        setupUpdateBindVariables(bean);
	//      これはTnAbstractEntityHandlerで実装
	varList := new(List)
	varList = varList
	varValueTypeList := new(StringList)
	varValueTypeList = varValueTypeList
	timestampPropertyName := ""
	timestampPropertyName = timestampPropertyName
	versionNoPropertyName := ""
	if (*(*entity).GetDBMeta()).HasVersionNo() {
		versionNoPropertyName = (*(*entity).GetDBMeta()).GetVersionNoColumnInfo().PropertyName
	}
	versionNoPropertyName = versionNoPropertyName
	for _, ptx := range t.boundPropTypes.data {
		var pt *TnPropertyType = ptx.(*TnPropertyType)
		pt = pt
		log.InternalDebug(fmt.Sprintf("PT %v\n", pt))

		//未実装
		//            if (pt.getPropertyName().equalsIgnoreCase(timestampPropertyName)) {
		//                final Timestamp timestamp = ResourceContext.getAccessTimestamp();
		//                addNewTimestamp(timestamp);
		//                varList.add(timestamp);
		//            } else if (pt.getPropertyName().equalsIgnoreCase(versionNoPropertyName)) {
		if pt.propetyName == versionNoPropertyName {
			if !t.versionNoAutoIncrementOnMemory { // means OnQuery
				continue // because of 'VERSION_NO = VERSION_NO + 1'
			}

			var versionNo int64 = *(GetEntityValue(entity, versionNoPropertyName).(*int64)) + 1
			//                final Object value = pt.getPropertyDesc().getValue(bean); // already null-checked
			//                final long longValue = DfTypeUtil.toPrimitiveLong(value) + 1L;
			//                final Long versionNo = Long.valueOf(longValue);
			if t.newVersionNoList == nil {
				t.newVersionNoList = new(List)
			}
			t.newVersionNoList.Add(versionNo)
			varList.Add(versionNo)

			//            } else if (_updateOption != null && _updateOption.hasStatement(pt.getColumnDbName())) {
			//                continue; // because of 'FOO_COUNT = FOO_COUNT + 1'
		} else {
			varList.Add(GetEntityValue(entity, pt.propetyName))
		}
		varValueTypeList.Add(pt.GoType)
		//        }
	}
	t.addAutoUpdateWhereBindVariables(varList, varValueTypeList, entity)
	t.bindVariables = varList
	t.bindVariableValueTypes = varValueTypeList
}
func (t *TnAbstractEntityHandler) setupInsertBindVariables(entity *Entity) {
	//	        setupInsertBindVariables(bean);
	//      これはTnAbstractEntityHandlerで実装
	varList := new(List)
	varList = varList
	varValueTypeList := new(StringList)
	varValueTypeList = varValueTypeList
	timestampPropertyName := ""
	timestampPropertyName = timestampPropertyName
	versionNoPropertyName := ""
	if (*(*entity).GetDBMeta()).HasVersionNo() {
		versionNoPropertyName = (*(*entity).GetDBMeta()).GetVersionNoColumnInfo().PropertyName
	}
	versionNoPropertyName = versionNoPropertyName
	//fmt.Printf("boundPropTypes %d\n",t.boundPropTypes.Size())
	for _, ptx := range t.boundPropTypes.data {
		var pt *TnPropertyType = ptx.(*TnPropertyType)
		pt = pt
		log.InternalDebug(fmt.Sprintf("PT %v\n", pt))

		//未実装
		//            if (pt.getPropertyName().equalsIgnoreCase(timestampPropertyName)) {
		//                final Timestamp timestamp = ResourceContext.getAccessTimestamp();
		//                addNewTimestamp(timestamp);
		//                varList.add(timestamp);
		//            } else if (pt.getPropertyName().equalsIgnoreCase(versionNoPropertyName)) {
		//fmt.Printf("propetyName %s versionNoPropertyName %s\n", pt.propetyName, versionNoPropertyName)
		if pt.propetyName == versionNoPropertyName {

			var firstNo int64 = VERSION_NO_FIRST_VALUE
			if t.newVersionNoList == nil {
				t.newVersionNoList = new(List)
			}
			t.newVersionNoList.Add(firstNo)
			varList.Add(firstNo)
			//fmt.Printf("version no %d propetyName %s\n", firstNo, pt.propetyName)
		} else {
			varList.Add(GetEntityValue(entity, pt.propetyName))
		}
		varValueTypeList.Add(pt.GoType)
		//        }
	}
	t.bindVariables = varList
	t.bindVariableValueTypes = varValueTypeList
	//fmt.Printf("varList size %d %v\n",varList.Size(),varList)
	//fmt.Printf("varValueTypeList size %d\n",varValueTypeList.Size())
}
func (t *TnAbstractEntityHandler) addAutoUpdateWhereBindVariables(varList *List, varValueTypeList *StringList, entity *Entity) {

	//	        final TnBeanMetaData bmd = getBeanMetaData();
	meta := (*entity).GetDBMeta()
	if (*meta).HasPrimaryKey() {
		for _, cix := range (*meta).GetPrimaryInfo().UniqueInfo.UniqueColumnList.data {
			var ci *ColumnInfo = cix.(*ColumnInfo)
			varList.Add(GetEntityValue(entity, ci.PropertyName))
			varValueTypeList.Add(ci.GoType)
		}
	}
	if t.optimisticLockHandling && (*meta).HasVersionNo() {
		varList.Add(GetEntityValue(entity, (*meta).GetVersionNoColumnInfo().PropertyName))
		varValueTypeList.Add((*meta).GetVersionNoColumnInfo().GoType)
	}

	//        if (_optimisticLockHandling && bmd.hasTimestampPropertyType()) {
	//            final TnPropertyType pt = bmd.getTimestampPropertyType();
	//            final DfPropertyDesc pd = pt.getPropertyDesc();
	//            varList.add(pd.getValue(bean));
	//            varValueTypeList.add(pt.getValueType());
	//        }
	//}
	//        setExceptionMessageSqlArgs(_bindVariables);
}

type TnCommandContextHandler struct {
	TnAbstractBasicSqlHandler
	CommandContext *CommandContext
}

func (p *TnCommandContextHandler) Execute(args []interface{}, tx *sql.Tx,behavior *Behavior)interface{} {

	p.sql=(*p.CommandContext).getSql()
//	fmt.Printf("getBindVariables %v\n",(*p.CommandContext).getBindVariables())
//	p.logSql((*p.CommandContext).getBindVariables(),
//		(*p.CommandContext).getBindVariableTypes())
//	log.Flush()
//	fmt.Printf("statementFactory %v \n",p.statementFactory)
	bindVariables:=(*p.CommandContext).getBindVariables()
	bindVariableTypes:=(*p.CommandContext).getBindVariableTypes()
		dbc := (*(*behavior).GetBaseBehavior().GetBehaviorCommandInvoker().InvokerAssistant).GetDBCurrent()

	ps := (*p.statementFactory).PrepareStatement(p.sql, tx, dbc)
	defer ps.Close()
	bindVar := (*p.statementFactory).ModifyBindVariables(bindVariables, bindVariableTypes)
	p.logSql(bindVariables, bindVariableTypes)
	log.Flush()
	res, err := tx.Stmt(ps).Exec(bindVar.data...)
	if err != nil {
		panic(err.Error())
	}

	updateno, _ := res.RowsAffected()
	log.InternalDebug(fmt.Sprintln("result no:", updateno))
	return updateno
}

// TnBasicSelectHandler #ResultType @execute @logSql
// TnBasicParameterHandler
// TnAbstractBasicSqlHandler #BasicSqlHander #sql #statementFactory @buildDisplaySql @doLogSql @logSql
//
// TnBasicUpdateHandler @execute
// TnBasicParameterHandler
// TnAbstractBasicSqlHandler #BasicSqlHander #sql #statementFactory @buildDisplaySql @doLogSql @logSql
//
// TnDeleteEntityHandler @setupBindVariables
// TnAbstractEntityHandler #bindVariableValueTypes #bindVariables #boundPropTypes #insertoption #newVersionNoList #optimisticLockHandling #updateoption #versionNoAutoIncrementOnMemory @addAutoUpdateWhereBindVariables @execute @handleUpdateResultWithOptimisticLock @processSuccess @setupBindVariables @setupDeleteBindVariables @setupInsertBindVariables @setupUpdateBindVariables
// TnAbstractBasicSqlHandler #BasicSqlHander #sql #statementFactory @buildDisplaySql @doLogSql @logSql
//
// TnInsertEntityHandler @doProcessIdentity @processSuccess @setupBindVariables
// TnAbstractEntityHandler #bindVariableValueTypes #bindVariables #boundPropTypes #insertoption #newVersionNoList #optimisticLockHandling #updateoption #versionNoAutoIncrementOnMemory @addAutoUpdateWhereBindVariables @execute @handleUpdateResultWithOptimisticLock @processSuccess @setupBindVariables @setupDeleteBindVariables @setupInsertBindVariables @setupUpdateBindVariables
// TnAbstractBasicSqlHandler #BasicSqlHander #sql #statementFactory @buildDisplaySql @doLogSql @logSql
//
// TnUpdateEntityHandler @processSuccess @setupBindVariables
// TnAbstractEntityHandler #bindVariableValueTypes #bindVariables #boundPropTypes #insertoption #newVersionNoList #optimisticLockHandling #updateoption #versionNoAutoIncrementOnMemory @addAutoUpdateWhereBindVariables @execute @handleUpdateResultWithOptimisticLock @processSuccess @setupBindVariables @setupDeleteBindVariables @setupInsertBindVariables @setupUpdateBindVariables
// TnAbstractBasicSqlHandler #BasicSqlHander #sql #statementFactory @buildDisplaySql @doLogSql @logSql
//
