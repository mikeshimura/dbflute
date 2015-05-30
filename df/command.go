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
	"database/sql"
	//	"errors"
	"fmt"
	"github.com/mikeshimura/dbflute/log"
	"reflect"
	"strings"
	"time"
)

type BehaviorCommandInvoker struct {
	InvokerAssistant *InvokerAssistant
}

func (b *BehaviorCommandInvoker) createOutsideSqlBasicExecutor(
	tableDbName string, bhv *Behavior) *OutsideSqlBasicExecutor {
	//	        final OutsideSqlExecutorFactory factory = _invokerAssistant.assistOutsideSqlExecutorFactory();
	factory := (*b.InvokerAssistant).AssistOutsideSqlExecutorFactory()
	dbdef := (*b.InvokerAssistant).GetDBCurrent().DBDef

	//未実装
	//config =(*b.InvokerAssistant). _invokerAssistant.assistDefaultStatementConfig();
	return (*factory).CreateBasic(b, bhv, tableDbName, dbdef, nil, nil) // for an entry instance

}

func (b *BehaviorCommandInvoker) InjectComponentProperty(cmd *BaseBehaviorCommand) {
	cmd.StatementFactory = (*b.InvokerAssistant).GetStatementFactory()
	cmd.DBMetaProvider = (*b.InvokerAssistant).GetDBMetaProvider()

}
func (b *BehaviorCommandInvoker) FindSqlExecution(cmd *BehaviorCommand) *SqlExecution {
	//        final boolean logEnabled = isLogEnabled();
	//        SqlExecution execution = null;
	//        try {
	//            final String key = behaviorCommand.buildSqlExecutionKey();
	key := (*cmd).BuildSqlExecutionKey()
	log.InternalDebug("KEY =" + key)
	//            execution = getSqlExecution(key);
	//Cache 取り敢えず未実装
	/////execution:=b.GetSqlExecution(key)
	//            if (execution == null) {
	//                long beforeCmd = 0;
	//                if (logEnabled) {
	//                    beforeCmd = systemTime();
	//                }
	//                SqlExecutionCreator creator = behaviorCommand.createSqlExecutionCreator();
	//	creator,err1 := (*cmd).CreateSqlExecutionCreator()
	//		if err1!=nil{
	//		return nil,err1
	//	}

	execution := b.GetOrCreateSqlExecution(key, cmd, (*cmd).GetConditionBean(),
		(*cmd).GetEntityBean())
	return execution
	//                execution = getOrCreateSqlExecution(key, creator);
	//                if (logEnabled) {
	//                    final long afterCmd = systemTime();
	//                    if (beforeCmd != afterCmd) {
	//                        logSqlExecution(behaviorCommand, execution, beforeCmd, afterCmd);
	//                    }
	//                }
	//            }
	//            return execution;
	//        } finally {
	//            if (logEnabled) {
	//                logInvocation(behaviorCommand, false);
	//            }
	//            readyInvokePath(behaviorCommand);
	//        }

	//return nil

}
func (b *BehaviorCommandInvoker) GetOrCreateSqlExecution(key string,
	creator *BehaviorCommand, conditionBean interface{},
	entityBean interface{}) *SqlExecution {
	//	        SqlExecution execution = null;
	//        synchronized (_executionCacheLock) {
	//            execution = getSqlExecution(key);
	//            if (execution != null) {
	//                // previous thread might have initialized
	//                // or reading might failed by same-time writing
	//                return execution;
	//            }
	//            if (isLogEnabled()) {
	//                log("...Initializing sqlExecution for the key '" + key + "'");
	//            }
	//            execution = executionCreator.createSqlExecution();
	//            assertCreatorReturnExecution(key, executionCreator, execution);
	//            _executionMap.put(key, execution);
	//        }
	//        toBeDisposable(); // for HotDeploy
	//        return execution;
	execution := (*creator).CreateSqlExecution(conditionBean, entityBean)
	return execution
}
func (b *BehaviorCommandInvoker) DispatchInvoking(cmd *BehaviorCommand) interface{} {
	//        if (behaviorCommand.isInitializeOnly()) {
	//            initializeSqlExecution(behaviorCommand);
	//            return null; // The end! (Initialize Only)
	//        }
	//        behaviorCommand.beforeGettingSqlExecution();
	(*cmd).BeforeGettingSqlExecution()
	//        SqlExecution execution = findSqlExecution(behaviorCommand);
	execution := b.FindSqlExecution(cmd)
	// - - - - - - - - - - -
	// Execute SQL Execution
	// - - - - - - - - - - -
	//        final SqlResultHandler sqlResultHander = getSqlResultHander();
	//        final boolean hasSqlResultHandler = sqlResultHander != null;
	///// ResultHandler必要性確認
	//        final long before = deriveCommandBeforeAfterTimeIfNeeds(logEnabled, hasSqlResultHandler);
	//        Long after = null;
	//        Object ret = null;
	//        RuntimeException cause = null;
	//        try {
	//// LiTblCBが帰る
	args := (*cmd).GetSqlExecutionArgument()
	log.InternalDebug("BehaviorCommandInvoker DispatchInvoking Execute")
	var tx *sql.Tx = (*cmd).GetTx()
	log.InternalDebug(fmt.Sprintf("command 102 %v%T\n", tx, tx))
	log.InternalDebug(fmt.Sprintf("command 102 %v%T\n", execution, execution))
	before := time.Now()
	ret := (*execution).Execute(args, tx, (*cmd).GetBehavior())
	after := time.Now()
	//            final Object[] args = behaviorCommand.getSqlExecutionArgument();
	//args:=cmd.GetSqlExecutionArgument()
	//            ret = executeSql(execution, args);
	//
	//            final Class<?> retType = behaviorCommand.getCommandReturnType();
	//            assertRetType(retType, ret);
	//
	//            after = deriveCommandBeforeAfterTimeIfNeeds(logEnabled, hasSqlResultHandler);
	//            if (logEnabled) {
	//                logResult(behaviorCommand, retType, ret, before, after);
	//            }
	if !LogStop {
		b.logResult(cmd, after.Sub(before), (*cmd).GetEntityType(), ret)
	}
	//
	//            ret = convertReturnValueIfNeeds(ret, retType);
	//        } catch (RuntimeException e) {
	//            try {
	//                handleExecutionException(e); // always throw
	//            } catch (RuntimeException handled) {
	//                cause = handled;
	//                throw handled;
	//            }
	//        } finally {
	//            behaviorCommand.afterExecuting();

	// - - - - - - - - - - - -
	// Call the handler back!
	// - - - - - - - - - - - -
	//            if (hasSqlResultHandler) {
	//                callbackSqlResultHanler(behaviorCommand, sqlResultHander, ret, before, after, cause);
	//            }
	//        }

	// - - - - - - - - -
	// Cast and Return!
	// - - - - - - - - -
	//        @SuppressWarnings("unchecked")
	//        final RESULT result = (RESULT) ret;
	//        return result;

	return ret
}
func (b *BehaviorCommandInvoker) logResult(cmd *BehaviorCommand,
	elapse time.Duration, entityType string, ret interface{}) {
	//	        final BehaviorResultBuilder behaviorResultBuilder = createBehaviorResultBuilder();
	//        final String resultExp = behaviorResultBuilder.buildResultExp(retType, ret, before, after);
	//        log(resultExp);
	//        log(" ");
	brb := new(BehaviorResultBuilder)
	resultExp := brb.buildResultExp(elapse, entityType, ret)
	DFLog(resultExp)
}
func (b *BehaviorCommandInvoker) Invoke(cmd *BehaviorCommand) interface{} {
	//        RuntimeException cause = null;
	//        RESULT result = null;
	//        try {
	//            final ResourceContext parentContext = getParentContext();
	//            initializeContext();
	//            setupResourceContext(behaviorCommand, parentContext);
	//            processBeforeHook(behaviorCommand);
	//            result = dispatchInvoking(behaviorCommand);
	rc := new(ResourceContext)
	(*cmd).SetResourceContext(rc)
	return b.DispatchInvoking(cmd)
}

type DeleteEntityCommand struct {
	BaseEntityCommand
	deleteOption *DeleteOption
}

func (s *DeleteEntityCommand) CreateSqlExecution(a interface{}, b interface{}) *SqlExecution {
	dcmd := new(TnDeleteEntityStaticCommand)
	dcmd.StatementFactory = s.StatementFactory
	dcmd.rc = s.rc
	dcmd.propertyNames = s.CreatePropertyNames()
	dcmd.targetDBMeta = (*s.entity).GetDBMeta()
	dcmd.optimisticLockHandling = true
	dcmd.setupDeleteSql()

	var se SqlExecution = dcmd
	dcmd.sqlExecution = &se
	return &se
}
func (s *DeleteEntityCommand) GetCommandName() string {
	return "delete"
}
func (s *DeleteEntityCommand) GetSqlExecutionArgument() []interface{} {
	return []interface{}{s.entity, s.deleteOption}
}

type InsertEntityCommand struct {
	BaseEntityCommand
	insertOption *InsertOption
}

func (s *InsertEntityCommand) GetSqlExecutionArgument() []interface{} {
	return []interface{}{s.entity, s.insertOption}
}

func (s *InsertEntityCommand) CreateSqlExecution(a interface{}, b interface{}) *SqlExecution {
	dcmd := new(TnInsertEntityDynamicCommand)
	dcmd.StatementFactory = s.StatementFactory
	dcmd.rc = s.rc
	dcmd.propertyNames = s.CreatePropertyNames()
	dcmd.targetDBMeta = (*s.entity).GetDBMeta()
	var se SqlExecution = dcmd
	dcmd.sqlExecution = &se
	return &se
}

func (s *InsertEntityCommand) GetCommandName() string {
	return "insert"
}

type SelectNextValCommand struct {
	BaseBehaviorCommand
	dbmeta *DBMeta
}

func (s *SelectNextValCommand) GetCommandName() string {
	return "selectNextVal"
}
func (s *SelectNextValCommand) GetEntityType() string {
	return "D_Int64"
}
func (s *SelectNextValCommand) BeforeGettingSqlExecution() {

}
func (s *SelectNextValCommand) BuildSqlExecutionKey() string {
	//        assertStatus("buildSqlExecutionKey");
	return s.TableDbName + ":" + s.GetCommandName() + "()"
}

func (s *SelectNextValCommand) CreateSqlExecution(a interface{}, b interface{}) *SqlExecution {
	//        assertStatus("createSelectNextValExecution");
	//        final DBMeta dbmeta = _dbmeta;
	s.assertTableHasSequence()
	sql := s.getSequenceNextValSql() // filtered later
	//fmt.Println("SQL:" + sql)
	//        assertSequenceReturnsNotNull(sql, dbmeta);
	//
	//        // handling for sequence cache
	//        final SequenceCache sequenceCache = findSequenceCache(dbmeta);
	//        sql = prepareSequenceCache(sql, sequenceCache);
	//
	//        return createSequenceExecution(handler, sql, sequenceCache);
	exe := new(SelectNextValExecution)
	exe.rc = s.rc
	exe.StatementFactory = s.StatementFactory
	exe.rootNode = exe.AnalyzeTwoWaySql(sql)

	exe.ResultType = "D_Int64"

	//fmt.Printf("RootNode %v\n", exe.rootNode)
	var se SqlExecution = exe
	exe.sqlExecution = &se
	return &se
}
func (s *SelectNextValCommand) getSequenceNextValSql() string {
	return (*s.dbmeta).GetSequenceNextValSql()
}
func (s *SelectNextValCommand) GetSqlExecutionArgument() []interface{} {
	return []interface{}{}
}
func (s *SelectNextValCommand) assertTableHasSequence() {
	if !(*s.dbmeta).HasSequence() {
		panic("If it uses sequence, the table should be related to a sequence: table=" +
			(*s.dbmeta).GetTableDbName())
	}
	return
}

type SelectListCBCommand struct {
	AbstractSelectCBCommand
	EntityType string
}

func (b *SelectListCBCommand) GetEntityType() string {
	return b.EntityType
}
func (b *SelectListCBCommand) GetEntityBean() interface{} {
	return b.EntityType
}
func (s *SelectListCBCommand) CreateSqlExecution(cb interface{}, entity interface{}) *SqlExecution {
	//                TnBeanMetaData bmd = createBeanMetaData();
	//                TnResultSetHandler handler = createBeanListResultSetHandler(bmd);
	//                return createSelectCBExecution(_conditionBean.getClass(), handler);

	//    }
	return s.CreateSelectCBExecution(cb, entity.(string))
}
func (s *SelectListCBCommand) CreateSelectCBExecution(
	cb interface{}, entity string) *SqlExecution {
	//    protected SelectCBExecution createSelectCBExecution(Class<? extends ConditionBean> cbType, TnResultSetHandler handler) {
	//        return newSelectCBExecution(createBeanArgNameTypeMap(cbType), handler);
	amap := s.createBeanArgNameTypeMap(cb)
	return s.NewSelectCBExecution(amap, entity, cb)

}
func (s *SelectListCBCommand) NewSelectCBExecution(
	amap map[string]string, entity interface{}, cb interface{}) *SqlExecution {
	se := new(SelectCBExecution)
	var sqlExecution SqlExecution = se
	se.sqlExecution = &sqlExecution
	se.rc = s.rc
	se.StatementFactory = s.StatementFactory
	se.ResultType = entity.(string)
	cbname, ok := amap["pmb"]
	if ok {
		se.ArgNames = []string{"pmb"}
		se.ArgTypes = []string{cbname}
	}
	var sqe SqlExecution = se
	return &sqe
}
func (s *SelectListCBCommand) GetCommandName() string {
	return "selectList"
}

func (s *SelectListCBCommand) BuildSqlExecutionKey() string {
	entityName := s.EntityType
	var cmd BehaviorCommand = s
	return s.BuildSqlExecutionKeySuper(&cmd) + ":" + entityName
}

func (s *SelectListCBCommand) BeforeGettingSqlExecution() {
	cbc := s.rc.GetConditionBeanContext()
	cbc.ConditionBean = s.ConditionBean
	//Fetch Assist 未実装
	//	        assertStatus("beforeGettingSqlExecution");
	//        final ConditionBean cb = _conditionBean;
	//        FetchAssistContext.setFetchBeanOnThread(cb);
	//        ConditionBeanContext.setConditionBeanOnThread(cb);
}

type AbstractSelectCBCommand struct {
	BaseBehaviorCommand
	ConditionBean interface{}
}

func (b *AbstractSelectCBCommand) GetConditionBean() interface{} {
	return b.ConditionBean
}
func (b *AbstractSelectCBCommand) GetSqlExecutionArgument() []interface{} {
	i := make([]interface{}, 1)
	i[0] = b.ConditionBean
	return i
}
func (b *AbstractSelectCBCommand) BuildSqlExecutionKeySuper(
	cmd *BehaviorCommand) string {
	//fmt.Printf("ConditionBean 1 %v%T", b.ConditionBean, b.ConditionBean)
	v := reflect.ValueOf(b.ConditionBean).Elem().
		FieldByName("BaseConditionBean").Interface()
	cbname := reflect.ValueOf(v).MethodByName("GetName").
		Call([]reflect.Value{})
	arg := cbname[0]
	return b.TableDbName + ":" + (*cmd).GetCommandName() + "(" +
		arg.String() + ")"
}

type BehaviorCommand interface {
	GetCommandName() string
	BeforeGettingSqlExecution()
	BuildSqlExecutionKey() string
	GetSqlExecutionArgument() []interface{}
	GetConditionBean() interface{}
	GetEntityBean() interface{}
	GetTx() *sql.Tx
	GetBehavior() *Behavior
	GetEntityType() string
	GetResourceContext() *ResourceContext
	SetResourceContext(rc *ResourceContext)
	CreateSqlExecution(cb interface{}, entity interface{}) *SqlExecution
}
type BaseBehaviorCommand struct {
	TableDbName            string
	StatementFactory       *StatementFactory
	BeanMetaDataFactory    *TnBeanMetaDataFactory
	DBMetaProvider         *DBMetaProvider
	Behavior               *Behavior
	tx                     *sql.Tx
	behaviorCommandInvoker *BehaviorCommandInvoker
	rc                     *ResourceContext
	BehaviorCommand        *BehaviorCommand
}

func (b *BaseBehaviorCommand) SetResourceContext(rc *ResourceContext) {
	b.rc = rc
}
func (b *BaseBehaviorCommand) GetResourceContext() *ResourceContext {
	return b.rc
}
func (b *BaseBehaviorCommand) GetEntityType() string {
	return ""
}
func (b *BaseBehaviorCommand) GetConditionBean() interface{} {
	return nil
}
func (b *BaseBehaviorCommand) GetEntityBean() interface{} {
	return ""
}
func (b *BaseBehaviorCommand) GetBehavior() *Behavior {
	return b.Behavior
}
func (b *BaseBehaviorCommand) GetTx() *sql.Tx {
	return b.tx
}
func (b *BaseBehaviorCommand) createBeanArgNameTypeMap(
	pmbTypeObj interface{}) map[string]string {
	amap := make(map[string]string)
	if pmbTypeObj == nil {
		return amap
	}
	stype := GetType(pmbTypeObj)

	// stype *df.ConditionBean -> xxxxCB
	if stype == "*df.ConditionBean" {
		var cb *ConditionBean = pmbTypeObj.(*ConditionBean)
		var cbx ConditionBean = *cb
		stype = GetType(cbx)
	}
	amap["pmb"] = stype
	return amap
}

type QueryUpdateCBCommand struct {
	AbstractQueryEntityCBCommand
	option *UpdateOption
}

func (u *QueryUpdateCBCommand) GetCommandName() string {
	return "queryUpdate"
}
func (s *QueryUpdateCBCommand) CreateSqlExecution(cb interface{}, entity interface{}) *SqlExecution {
	entityx := entity.(*Entity)
	dbmeta := (*entityx).GetDBMeta()
	return s.createQueryEntityCBExecution(dbmeta)
}
func (u *QueryUpdateCBCommand) createQueryEntityCBExecution(dbmeta *DBMeta) *SqlExecution {
	propertyNames := u.getPersistentPropertyNames(dbmeta)
	propertyNames = propertyNames
	tnCommand := u.createQueryUpdateDynamicCommand(propertyNames, dbmeta)
	var sqle SqlExecution = tnCommand
	return &sqle
}
func (u *QueryUpdateCBCommand) createQueryUpdateDynamicCommand(
	propertyNames *StringList, dbmeta *DBMeta) *TnQueryUpdateDynamicCommand {
	cmd := new(TnQueryUpdateDynamicCommand)
	var sqlExecution SqlExecution = cmd
	cmd.sqlExecution = &sqlExecution
	cmd.rc = u.rc
	cmd.StatementFactory = u.StatementFactory
	cmd.targetDBMeta = dbmeta
	cmd.propertyNames = propertyNames
	return cmd
}

func (u *QueryUpdateCBCommand) GetSqlExecutionArgument() []interface{} {
	return []interface{}{u.entity, u.ConditionBean, u.option}
}
type QueryDeleteCBCommand struct {
	AbstractQueryEntityCBCommand
	option *DeleteOption
	entityType string
}
func (s * QueryDeleteCBCommand) CreateSqlExecution(cb interface{}, entity interface{}) *SqlExecution {
	dbmeta := DBMetaProvider_I.TableDbNameInstanceMap[s.entityType]
	return s.createQueryEntityCBExecution(dbmeta)
}
func (s * QueryDeleteCBCommand) createQueryEntityCBExecution(dbmeta *DBMeta) *SqlExecution {
	tnCommand := new(TnQueryDeleteDynamicCommand)
	tnCommand.ResultType=s.entityType
	tnCommand.StatementFactory=s.StatementFactory
	var sqle SqlExecution = tnCommand
	return &sqle
}
func (u *QueryDeleteCBCommand) GetCommandName() string {
	return "queryDelete"
}
func (u *QueryDeleteCBCommand) GetSqlExecutionArgument() []interface{} {
	return []interface{}{u.ConditionBean, u.option}
}
func (u *QueryDeleteCBCommand) BuildSqlExecutionKey() string{
	return ""
}
type AbstractQueryEntityCBCommand struct {
	BaseEntityCommand
	ConditionBean interface{}
}

func (b *AbstractQueryEntityCBCommand) GetConditionBean() interface{} {
	return b.ConditionBean
}

type UpdateEntityCommand struct {
	BaseEntityCommand
	option *UpdateOption
}

func (u *UpdateEntityCommand) GetSqlExecutionArgument() []interface{} {
	return []interface{}{u.entity, u.option}
}
func (u *UpdateEntityCommand) GetCommandName() string {
	return "update"
}

func (u *UpdateEntityCommand) CreateSqlExecution(
	cb interface{}, entity interface{}) *SqlExecution {
	//        return new SqlExecutionCreator() {
	//            public SqlExecution createSqlExecution() {
	//                final TnBeanMetaData bmd = createBeanMetaData();
	//                return createUpdateEntitySqlExecution(bmd);
	entityx := entity.(*Entity)
	dbmeta := (*entityx).GetDBMeta()
	return u.CreateUpdateEntitySqlExecution(dbmeta)
}

func (u *UpdateEntityCommand) CreateUpdateEntitySqlExecution(
	dbmeta *DBMeta) *SqlExecution {
	propertyNames := u.getPersistentPropertyNames(dbmeta)
	propertyNames = propertyNames
	tnCommand := u.createUpdateEntityDynamicCommand(propertyNames, dbmeta)
	var sqle SqlExecution = tnCommand
	return &sqle
}
func (u *UpdateEntityCommand) createUpdateEntityDynamicCommand(
	propertyNames *StringList, dbmeta *DBMeta) *TnUpdateEntityDynamicCommand {
	//   final TnUpdateEntityDynamicCommand cmd = newUpdateEntityDynamicCommand();
	cmd := new(TnUpdateEntityDynamicCommand)
	var sqlExecution SqlExecution = cmd
	cmd.sqlExecution = &sqlExecution
	cmd.rc = u.rc
	cmd.StatementFactory = u.StatementFactory
	cmd.targetDBMeta = dbmeta
	cmd.propertyNames = propertyNames
	cmd.optimisticLockHandling = u.isOptimisticLockHandling()
	cmd.versionNoAutoIncrementOnMemory = u.isVersionNoAutoIncrementOnMemory()
	return cmd
}
func (u *UpdateEntityCommand) isOptimisticLockHandling() bool {
	return true
}
func (u *UpdateEntityCommand) isVersionNoAutoIncrementOnMemory() bool {
	return u.isOptimisticLockHandling()
}

type BaseEntityCommand struct {
	BaseBehaviorCommand
	entity *Entity
}

func (s *BaseEntityCommand) CreatePropertyNames() *StringList {
	var propertyList = new(StringList)
	for _, ci := range (*(*s.entity).GetDBMeta()).GetColumnInfoList().data {
		var columnInfo *ColumnInfo = ci.(*ColumnInfo)
		propertyList.Add(columnInfo.PropertyName)
	}
	return propertyList
}
func (b *BaseEntityCommand) GetEntityBean() interface{} {
	return b.entity
}
func (b *BaseEntityCommand) xsetupEntityCommand(entity *Entity) {
	b.entity = entity
}

func (b *BaseEntityCommand) getPersistentPropertyNames(
	dbmeta *DBMeta) *StringList {
	columnInfoList := (*dbmeta).GetColumnInfoList()
	propertyNameList := new(StringList)
	for _, columnInfo := range columnInfoList.data {
		var ci *ColumnInfo = columnInfo.(*ColumnInfo)
		propertyNameList.Add(ci.PropertyName)
	}
	log.InternalDebug(fmt.Sprintf("propertyNameList %v\n",
		propertyNameList.data))
	return propertyNameList
}
func (b *BaseEntityCommand) BeforeGettingSqlExecution() {

}
func (b *BaseEntityCommand) BuildSqlExecutionKey() string {
	//	    assertStatus("buildSqlExecutionKey");
	entityName := (*(*b.entity).GetDBMeta()).GetTablePropertyName()

	return b.TableDbName + ":" + (*b.BehaviorCommand).GetCommandName() +
		"(" + entityName + ")"

}

type OutsideSqlSelectListCommand struct {
	AbstractOutsideSqlSelectCommand
	entityType string
}

func (o *OutsideSqlSelectListCommand) GetCommandName() string {
	return "selectList"
}
func (o *OutsideSqlSelectListCommand) GetEntityType() string {
	return o.entityType
}

type AbstractOutsideSqlSelectCommand struct {
	AbstractOutsideSqlCommand
}

func (a *AbstractOutsideSqlSelectCommand) GetSqlExecutionArgument() []interface{} {
	//	var pmbi interface{} = reflect.ValueOf(*a.pmb)
	//	fmt.Printf("pmbi %v %T\n",pmbi,pmbi)
	return []interface{}{a.pmb}
}
func (a *AbstractOutsideSqlSelectCommand) BeforeGettingSqlExecution() {

	//	        assertStatus("beforeGettingSqlExecution");
	//        OutsideSqlContext.setOutsideSqlContextOnThread(createOutsideSqlContext());
	//
	//        // set up fetchNarrowingBean
	//        final Object pmb = _parameterBean;
	//        final OutsideSqlOption option = _outsideSqlOption;
	//        setupFetchBean(pmb, option);
	cbc := a.rc.GetOutsideSqlContext()
	//var pmbi interface{} = reflect.ValueOf(*a.pmb).Interface()
	cbc.Pmb = a.pmb
	a.outsideSqlContext = cbc
}
func (a *AbstractOutsideSqlSelectCommand) BuildSqlExecutionKey() string {

	//        assertStatus("buildSqlExecutionKey");
	//        return generateSpecifiedOutsideSqlUniqueKey();
	return a.generateSpecifiedOutsideSqlUniqueKey()
}
func (a *AbstractOutsideSqlSelectCommand) generateSpecifiedOutsideSqlUniqueKey() string {
	methodName := (*a.BehaviorCommand).GetCommandName()
	path := a.OutsideSqlPath
	pmb := a.pmb
	option := a.OutsideSqlOption
	resultType := (*a.BehaviorCommand).GetEntityType()
	return (*a.rc).GetOutsideSqlContext().generateSpecifiedOutsideSqlUniqueKey(
		methodName, path, pmb, option, resultType)
	//        return OutsideSqlContext.generateSpecifiedOutsideSqlUniqueKey(methodName, path, pmb, option, resultType);
}
func (a *AbstractOutsideSqlSelectCommand) CreateSqlExecution(
	outsideSqlContext interface{}, entity interface{}) *SqlExecution {
	//	                final OutsideSqlContext outsideSqlContext = OutsideSqlContext.getOutsideSqlContextOnThread();
	//                return createOutsideSqlSelectExecution(outsideSqlContext);
	return a.createOutsideSqlSelectExecution(outsideSqlContext.(*OutsideSqlContext))
}
func (a *AbstractOutsideSqlSelectCommand) createOutsideSqlSelectExecution(
	outsideSqlContext *OutsideSqlContext) *SqlExecution {
	pmb := outsideSqlContext.Pmb
	suffix := a.buildDbmsSuffix()
	a.outsideSqlContext.OutsideSqlPath = a.OutsideSqlPath
	sql := outsideSqlContext.readFilteredOutsideSql(suffix)

	//
	//        // - - - - - - - - - - - - -
	//        // Create ResultSetHandler.
	//        // - - - - - - - - - - - - -
	//        final TnResultSetHandler handler = createOutsideSqlSelectResultSetHandler();
	//
	//        // - - - - - - - - - - -
	//        // Create SqlExecution.
	//        // - - - - - - - - - - -
	//        final OutsideSqlSelectExecution execution = createOutsideSqlSelectExecution(pmb, sql, handler);
	execution := a.createOutsideSqlSelectExecutionSub(pmb, sql)
	//        execution.setRemoveBlockComment(isRemoveBlockComment(outsideSqlContext));
	//        execution.setRemoveLineComment(isRemoveLineComment(outsideSqlContext));
	//        execution.setFormatSql(outsideSqlContext.isFormatSql());
	//        execution.setOutsideSqlFilter(_outsideSqlFilter);
	var exei SqlExecution = execution
	return &exei
}
func (a *AbstractOutsideSqlSelectCommand) createOutsideSqlSelectExecutionSub(
	pmb interface{}, sql string) *OutsideSqlSelectExecution {
	//        final Map<String, Class<?>> argNameTypeMap = createBeanArgNameTypeMap(pmbTypeObj);
	//        return newOutsideSqlSelectExecution(argNameTypeMap, sql, handler);
	argNameTypeMap := a.createBeanArgNameTypeMap(pmb)
	ex := new(OutsideSqlSelectExecution)
	var sqlExecution SqlExecution = ex
	ex.sqlExecution = &sqlExecution
	ex.rc = a.rc
	ex.StatementFactory = a.StatementFactory
	ex.ResultType = (*a.BehaviorCommand).GetEntityType()
	ex.IsBlockeNullParameter = true
	pmbname, ok := argNameTypeMap["pmb"]
	if ok {
		ex.ArgNames = []string{"pmb"}
		ex.ArgTypes = []string{pmbname}
	}
	analyzer := ex.CreateSqlAnalyzer(sql)
	rn := analyzer.Analyze()
	ex.rootNode = rn
	if ex.rootNode == nil {
		panic("rootNode NIL")
	}
	return ex
}

//func (a *AbstractOutsideSqlSelectCommand)createBeanArgNameTypeMap(pmbTypeObj interface{})map[string]string{
//	amap:=make(map[string]string)
//	amap["pmb"]=GetType(pmbTypeObj)
//	return amap
//}

type AbstractOutsideSqlCommand struct {
	BaseBehaviorCommand
	OutsideSqlPath           string
	OutsideSqlOption         *OutsideSqlOption
	CurrentDBDef             *DBDef
	outsideSqlContextFactory *OutsideSqlContextFactory
	pmb                      interface{}
	outsideSqlContext        *OutsideSqlContext
}

func (a *AbstractOutsideSqlCommand) buildDbmsSuffix() string {
	//	        assertOutsideSqlBasic("buildDbmsSuffix");
	//        final String productName = _currentDBDef.code();
	//        return (productName != null ? "_" + productName.toLowerCase() : "");
	productName := (*a.CurrentDBDef).Code()
	return "_" + strings.ToLower(productName)
}
func (a *AbstractOutsideSqlCommand) GetConditionBean() interface{} {
	return a.outsideSqlContext
}

type OutsideSqlExecuteCommand struct {
	AbstractOutsideSqlCommand
}

func (s *OutsideSqlExecuteCommand) GetEntityType() string {
	return "D_Int64"
}
func (a *OutsideSqlExecuteCommand) GetSqlExecutionArgument() []interface{} {
	//	var pmbi interface{} = reflect.ValueOf(*a.pmb)
	//	fmt.Printf("pmbi %v %T\n",pmbi,pmbi)
	return []interface{}{a.pmb}
}
func (o *OutsideSqlExecuteCommand) GetCommandName() string {
	return "execute"
}
func (a *OutsideSqlExecuteCommand) CreateSqlExecution(
	outsideSqlContext interface{}, entity interface{}) *SqlExecution {
	//	                final OutsideSqlContext outsideSqlContext = OutsideSqlContext.getOutsideSqlContextOnThread();
	//                return createOutsideSqlSelectExecution(outsideSqlContext);
	return a.createOutsideSqlExecuteExecution(outsideSqlContext.(*OutsideSqlContext))
}
func (a *OutsideSqlExecuteCommand) createOutsideSqlExecuteExecution(
	outsideSqlContext *OutsideSqlContext) *SqlExecution {
	pmb := outsideSqlContext.Pmb
	suffix := a.buildDbmsSuffix()
	a.outsideSqlContext.OutsideSqlPath = a.OutsideSqlPath
	sql := outsideSqlContext.readFilteredOutsideSql(suffix)
	execution := a.createOutsideSqlExecuteExecutionSub(pmb, sql)
	//	        final Object pmb = outsideSqlContext.getParameterBean();
	//        final String suffix = buildDbmsSuffix();
	//        final String sql = outsideSqlContext.readFilteredOutsideSql(_sqlFileEncoding, suffix);
	//
	//        final OutsideSqlExecuteExecution execution = createOutsideSqlExecuteExecution(pmb, sql);
	//        execution.setOutsideSqlFilter(_outsideSqlFilter);
	//        execution.setRemoveBlockComment(isRemoveBlockComment(outsideSqlContext));
	//        execution.setRemoveLineComment(isRemoveLineComment(outsideSqlContext));
	//        execution.setFormatSql(outsideSqlContext.isFormatSql());
	//        return execution;
	var exei SqlExecution = execution
	return &exei
}
func (a *OutsideSqlExecuteCommand) createOutsideSqlExecuteExecutionSub(
	pmb interface{}, sql string) *OutsideSqlExecuteExecution {
	argNameTypeMap := a.createBeanArgNameTypeMap(pmb)
	ex := new(OutsideSqlExecuteExecution)
	var sqlExecution SqlExecution = ex
	ex.sqlExecution = &sqlExecution
	ex.rc = a.rc
	ex.StatementFactory = a.StatementFactory
	pmbname, ok := argNameTypeMap["pmb"]
	if ok {
		ex.ArgNames = []string{"pmb"}
		ex.ArgTypes = []string{pmbname}
	}
	analyzer := ex.CreateSqlAnalyzer(sql)
	rn := analyzer.Analyze()
	ex.rootNode = rn
	if ex.rootNode == nil {
		panic("rootNode NIL")
	}
	return ex
}
func (a *OutsideSqlExecuteCommand) BuildSqlExecutionKey() string {

	//        assertStatus("buildSqlExecutionKey");
	//        return generateSpecifiedOutsideSqlUniqueKey();
	return a.generateSpecifiedOutsideSqlUniqueKey()
}
func (a *OutsideSqlExecuteCommand) generateSpecifiedOutsideSqlUniqueKey() string {
	methodName := (*a.BehaviorCommand).GetCommandName()
	path := a.OutsideSqlPath
	pmb := a.pmb
	option := a.OutsideSqlOption
	return (*a.rc).GetOutsideSqlContext().generateSpecifiedOutsideSqlUniqueKey(
		methodName, path, pmb, option, "")
	//        return OutsideSqlContext.generateSpecifiedOutsideSqlUniqueKey(methodName, path, pmb, option, resultType);
}
func (a *OutsideSqlExecuteCommand) BeforeGettingSqlExecution() {

	//	        assertStatus("beforeGettingSqlExecution");
	//        OutsideSqlContext.setOutsideSqlContextOnThread(createOutsideSqlContext());
	//
	//        // set up fetchNarrowingBean
	//        final Object pmb = _parameterBean;
	//        final OutsideSqlOption option = _outsideSqlOption;
	//        setupFetchBean(pmb, option);
	cbc := a.rc.GetOutsideSqlContext()
	//var pmbi interface{} = reflect.ValueOf(*a.pmb).Interface()
	cbc.Pmb = a.pmb
	a.outsideSqlContext = cbc
}

// BehaviorCommandInvoker #InvokerAssistant @DispatchInvoking @FindSqlExecution @GetOrCreateSqlExecution @InjectComponentProperty @Invoke @createOutsideSqlBasicExecutor @logResult
//
// OutsideSqlSelectListCommand #entityType @GetCommandName @GetEntityType
// AbstractOutsideSqlSelectCommand @BeforeGettingSqlExecution @BuildSqlExecutionKey @CreateSqlExecution @CreateSqlExecutionCreator @GetSqlExecutionArgument @buildDbmsSuffix @createOutsideSqlSelectExecution @createOutsideSqlSelectExecutionSub @generateSpecifiedOutsideSqlUniqueKey
// AbstractOutsideSqlCommand #CurrentDBDef #OutsideSqlOption #OutsideSqlPath #outsideSqlContext #outsideSqlContextFactory #pmb @GetConditionBean
// BaseBehaviorCommand #BeanMetaDataFactory #Behavior #DBMetaProvider #StatementFactory #TableDbName #behaviorCommandInvoker #rc #topcommand #tx @GetBehavior @GetConditionBean @GetEntityBean @GetEntityType @GetResourceContext @GetTx @SetResourceContext @createBeanArgNameTypeMap
//
// SelectListCBCommand #EntityType @BeforeGettingSqlExecution @BuildSqlExecutionKey @CreateSelectCBExecution @CreateSqlExecution @CreateSqlExecutionCreator @GetCommandName @GetEntityBean @GetEntityType @NewSelectCBExecution
// AbstractSelectCBCommand #ConditionBean @BuildSqlExecutionKeySuper @GetConditionBean @GetSqlExecutionArgument
// BaseBehaviorCommand #BeanMetaDataFactory #Behavior #DBMetaProvider #StatementFactory #TableDbName #behaviorCommandInvoker #rc #topcommand #tx @GetBehavior @GetConditionBean @GetEntityBean @GetEntityType @GetResourceContext @GetTx @SetResourceContext @createBeanArgNameTypeMap
//
// UpdateEntityCommand #option @CreateSqlExecution @CreateSqlExecutionCreator @CreateUpdateEntitySqlExecution @GetCommandName @GetSqlExecutionArgument @createUpdateEntityDynamicCommand @isOptimisticLockHandling @isVersionNoAutoIncrementOnMemory
// BaseEntityCommand #entity @BeforeGettingSqlExecution @BuildSqlExecutionKey @GetEntityBean @getPersistentPropertyNames @xsetupEntityCommand
// BaseBehaviorCommand #BeanMetaDataFactory #Behavior #DBMetaProvider #StatementFactory #TableDbName #behaviorCommandInvoker #rc #topcommand #tx @GetBehavior @GetConditionBean @GetEntityBean @GetEntityType @GetResourceContext @GetTx @SetResourceContext @createBeanArgNameTypeMap
//
