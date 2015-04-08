package df

import (
	//"container/list"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"github.com/mikeshimura/dbflute/log"
	"path/filepath"
	"reflect"
	//	"strconv"
	//	"strings"
)

const (
	VERSION_NO_FIRST_VALUE = 0
)

type SqlExecutionCreator interface {
	CreateSqlExecution(cb interface{}, entity interface{}) (*SqlExecution, error)
}

type SqlExecution interface {
	Execute(args []interface{}, tx *sql.Tx, behavior *Behavior) (interface{}, error)
	GetRootNode(args []interface{}) (*Node, error)
	GetArgNames() []string
	GetArgTypes() []string
	filterExecutedSql(executedSql string) string
	newBasicParameterHandler(executedSql string) *ParameterHandler
}

type TnDeleteEntityStaticCommand struct {
	TnAbstractEntityStaticCommand
}

func (t *TnDeleteEntityStaticCommand) setupSql() error {
	return t.setupDeleteSql()
}

type TnAbstractEntityStaticCommand struct {
	propertyNames *StringList
	TnAbstractBasicSqlCommand
	targetDBMeta                   *DBMeta
	versionNoAutoIncrementOnMemory bool
	optimisticLockHandling         bool
	propertyTypes                  *List
	sql                            string
}

func (t *TnAbstractEntityStaticCommand) Execute(args []interface{}, tx *sql.Tx, behavior *Behavior) (interface{}, error) {
	handler := new(TnDeleteEntityHandler)
	handler.sql = t.sql
	handler.statementFactory = t.StatementFactory
	handler.optimisticLockHandling = t.optimisticLockHandling
	var hd BasicSqlHander = handler
	handler.BasicSqlHander = &hd
	var entity Entity = *args[0].(*Entity)
	return t.doExecute(&entity, handler, tx)
}

func (t *TnAbstractEntityStaticCommand) setupDeleteSql() error {
	//	        checkPrimaryKey();
	//        final StringBuilder sb = new StringBuilder(64);
	//        sb.append("delete from ").append(_targetDBMeta.getTableSqlName());
	//        setupDeleteWhere(sb);
	//        _sql = sb.toString();
	sb := new(bytes.Buffer)
	sb.WriteString("delete from " + (*t.targetDBMeta).GetTableDbName())
	t.setupDeleteWhere(sb)
	t.sql = sb.String()
	log.InternalDebug(t.sql)
	return nil
}
func (t *TnAbstractEntityStaticCommand) setupDeleteWhere(sb *bytes.Buffer) {
	//        final TnBeanMetaData bmd = _beanMetaData;
	sb.WriteString(" where ")
	if (*t.targetDBMeta).HasPrimaryKey() {
		for i, colInfo := range (*t.targetDBMeta).GetPrimaryInfo().UniqueInfo.UniqueColumnList.data {
			var ci *ColumnInfo = colInfo.(*ColumnInfo)
			if i > 0 {
				sb.WriteString(" and ")
			}
			sb.WriteString(ci.ColumnSqlName.ColumnSqlName + " = ?")

		}
	}
	if t.optimisticLockHandling && (*t.targetDBMeta).HasVersionNo() {
		sb.WriteString(" and " +
			(*t.targetDBMeta).GetVersionNoColumnInfo().ColumnSqlName.ColumnSqlName + " = ?")
	}

	//        for (int i = 0; i < bmd.getPrimaryKeySize(); ++i) {
	//            sb.append(bmd.getPrimaryKeySqlName(i)).append(" = ? and ");
	//        }
	//        sb.setLength(sb.length() - 5);
	//        if (_optimisticLockHandling && bmd.hasVersionNoPropertyType()) {
	//            TnPropertyType pt = bmd.getVersionNoPropertyType();
	//            sb.append(" and ").append(pt.getColumnSqlName()).append(" = ?");
	//        }
	//        if (_optimisticLockHandling && bmd.hasTimestampPropertyType()) {
	//            TnPropertyType pt = bmd.getTimestampPropertyType();
	//            sb.append(" and ").append(pt.getColumnSqlName()).append(" = ?");
	//        }
}

type TnInsertEntityDynamicCommand struct {
	TnAbstractEntityDynamicCommand
}

func (t *TnAbstractEntityStaticCommand) doExecute(args *Entity, handler *TnDeleteEntityHandler, tx *sql.Tx) (int64, error) {
	//	        handler.setExceptionMessageSqlArgs(args);
	//        final int rows = handler.execute(args);
	//        return Integer.valueOf(rows);
	return handler.execute(args, tx)
}

func (t *TnInsertEntityDynamicCommand) Execute(args []interface{}, tx *sql.Tx, behavior *Behavior) (interface{}, error) {
	//fmt.Printf("args0 %v %T \n", args[0], args[0])
	//	        if (args == null || args.length == 0) {
	//            String msg = "The argument 'args' should not be null or empty.";
	//            throw new IllegalArgumentException(msg);
	//        }
	//        final Object bean = args[0];
	//        final InsertOption<ConditionBean> option = extractInsertOptionChecked(args);
	//        prepareStatementConfigOnThreadIfExists(option);
	//
	//        final TnBeanMetaData bmd = _beanMetaData;
	//        final TnPropertyType[] propertyTypes = createInsertPropertyTypes(bmd, bean, _propertyNames, option);
	//        final String sql = filterExecutedSql(createInsertSql(bmd, propertyTypes, option));
	//        return doExecute(bean, propertyTypes, sql, option);

	//Dummy
	var option *InsertOption = nil
	var entity *Entity = (args[0]).(*Entity)
	propertyTypes := t.createInsertPropertyTypes(entity, option)
	//fmt.Println("propertyTypes len", propertyTypes.Size())
	sql := t.createInsertSql(entity)
	fsql := t.filterExecutedSql(sql)
	return t.doExecute(entity, propertyTypes, fsql, option, tx)
}

func (t *TnInsertEntityDynamicCommand) doExecute(entity *Entity, propertyTypes *List, sql string, option *InsertOption, tx *sql.Tx) (int64, error) {
	//	        final TnInsertEntityHandler handler = createInsertEntityHandler(propertyTypes, sql, option);
	//        final Object[] realArgs = new Object[] { bean };
	//        handler.setExceptionMessageSqlArgs(realArgs);
	//        final int rows = handler.execute(realArgs);
	//        return Integer.valueOf(rows);
	handler := new(TnInsertEntityHandler)
	handler.statementFactory = t.StatementFactory
	handler.sql = sql
	handler.insertoption = option
	handler.boundPropTypes = propertyTypes
	var base BasicSqlHander = handler
	handler.BasicSqlHander = &base
	return handler.execute(entity, tx)
}
func (t *TnInsertEntityDynamicCommand) createInsertPropertyTypes(entity *Entity, option *InsertOption) *List {
	propetyNames := new(StringList)
	//dbMeta := (*entity).GetDBMeta()
	modifiedSet := t.getModifiedPropertyNames(entity)
	//fmt.Printf("ModifiedSet %v\n", modifiedSet)
	typeList := new(List)
	timestampProp := ""
	timestampProp = timestampProp
	versionNoProp := ""

	if (*(*entity).GetDBMeta()).HasVersionNo() {
		versionNoProp = (*(*entity).GetDBMeta()).GetVersionNoColumnInfo().PropertyName
		versionNoProp = versionNoProp
	}
	primaryKeyset := make(map[string]string)
	primaryKeyset = primaryKeyset
	if (*(*entity).GetDBMeta()).HasPrimaryKey() {
		primaryKeys := (*(*entity).GetDBMeta()).GetPrimaryInfo().UniqueInfo.UniqueColumnList.data
		for _, primaryKey := range primaryKeys {
			var colInfo *ColumnInfo = primaryKey.(*ColumnInfo)
			primaryKeyset[colInfo.PropertyName] = colInfo.PropertyName
		}
	}
	for _, propertyName := range t.propertyNames.data {
		_, ok := primaryKeyset[propertyName]
		if ok {
			if option == nil || !option.isPrimaryKeyIdentityDisabled() {
				if (*(*entity).GetDBMeta()).HasIdentity() {
					continue
				}
			}
			typeList.Add((*(*entity).GetDBMeta()).GetPropertyType(propertyName))
			propetyNames.Add(propertyName)
			continue
		}
		log.InternalDebug(fmt.Sprintf("propertyName %v\n", propertyName))
		//pt:=(*dbMeta).GetPropertyType(propertyName)
		if t.isOptimisticLockProperty(timestampProp, versionNoProp, propertyName) || t.isSpecifiedProperty(option, modifiedSet, propertyName) {
			// Statement
			typeList.Add((*(*entity).GetDBMeta()).GetPropertyType(propertyName))
			propetyNames.Add(propertyName)
			//fmt.Println("typeList Added :" + propertyName)
		}
	}
	t.propertyNames = propetyNames
	return typeList
}
func (t *TnInsertEntityDynamicCommand) isSpecifiedProperty(option *InsertOption, modifiedSet map[string]string, propertyName string) bool {
	//未実装
	//	        if (option != null && option.hasSpecifiedUpdateColumn()) { // BatchUpdate
	//            return option.isSpecifiedUpdateColumn(pt.getColumnDbName());
	//        } else { // EntityUpdate
	//            return isModifiedProperty(modifiedSet, pt); // process for ModifiedColumnUpdate
	//        }

	_, ok := modifiedSet[propertyName]
	//fmt.Printf("propety name %s ok %v\n",propertyName,ok)
	return ok
}
func (t *TnInsertEntityDynamicCommand) checkPropetyList(propertyName string) bool {
	for _, name := range t.propertyNames.data {
		if name == propertyName {
			return true
		}
	}
	return false
}
func (t *TnInsertEntityDynamicCommand) createInsertSql(entity *Entity) string {
	//tableDbName := (*t.targetDBMeta).GetTableDbName()
	columnSb := new(bytes.Buffer)
	valuesSb := new(bytes.Buffer)
	//        for (int i = 0; i < propertyTypes.length; ++i) {
	//            final TnPropertyType pt = propertyTypes[i];
	for _, ci := range (*t.targetDBMeta).GetColumnInfoList().data {

		var colInfo *ColumnInfo = ci.(*ColumnInfo)
		columnSqlName := colInfo.ColumnSqlName.ColumnSqlName
		propertyName := colInfo.PropertyName
		if t.checkPropetyList(propertyName) == false {
			continue
		}
		if columnSb.Len() > 0 {
			columnSb.WriteString(", ")
			valuesSb.WriteString(", ")
		}
		columnSb.WriteString(columnSqlName)
		//columnDbName :=colInfo.ColumnDbName
		//valuesSb.append(encryptIfNeeds(tableDbName, columnDbName, "?"));
		valuesSb.WriteString("?")
	}

	sb := new(bytes.Buffer)
	sb.WriteString("insert into " + (*t.targetDBMeta).GetTableSqlName().TableSqlName)
	sb.WriteString(" (" + columnSb.String() + ")")
	sb.WriteString(Ln + " values (" + valuesSb.String() + ")")
	return sb.String()

}

type AbstractFixedArgExecution struct {
	TnAbstractTwoWaySqlCommand

	ArgNames []string
	ArgTypes []string
}

func (a *AbstractFixedArgExecution) GetArgNames() []string {
	return a.ArgNames
}
func (a *AbstractFixedArgExecution) GetArgTypes() []string {
	return a.ArgTypes
}

type SelectNextValExecution struct {
	SelectSimpleExecution
}

type SelectSimpleExecution struct {
	AbstractFixedSqlExecution
}

func (s *SelectSimpleExecution) newBasicParameterHandler(sql string) *ParameterHandler {
	bh := new(TnBasicSelectHandler)
	bh.sql = sql
	//fmt.Printf("StatementFactoty %v\n", s.StatementFactory)
	bh.statementFactory = s.StatementFactory
	bh.ResultType = s.ResultType
	var ph ParameterHandler = bh
	return &ph
}
func (s *SelectSimpleExecution) IsBlockNullParameter() bool {
	return true
}

type SelectCBExecution struct {
	AbstractFixedArgExecution
	//ResultType string
}

func (s *SelectCBExecution) newBasicParameterHandler(executedSql string) *ParameterHandler {
	handler := new(TnBasicSelectHandler)
	handler.sql = executedSql
	handler.statementFactory = s.StatementFactory
	handler.ResultType = s.ResultType
	var hand ParameterHandler = handler
	return &hand
}
func (s *SelectCBExecution) filterExecutedSql(executedSql string) string {
	//dbflute TnAbstractTwoWaySqlCommand はそのまま Returnしている
	//未実装
	//CallbackContext.isExistSqlStringFilterOnThread()
	return executedSql
}
func (s *SelectCBExecution) GetRootNode(args []interface{}) (*Node, error) {
	sql := s.ExtractTwoWaySql(args)
	node, err := s.AnalyzeTwoWaySql(sql)
	return node, err
}
func (s *SelectCBExecution) ExtractTwoWaySql(args []interface{}) string {
	//        assertArgsValid(args);
	//        final Object firstElement = args[0];
	//        assertObjectNotNull("args[0]", firstElement);
	//        assertFirstElementConditionBean(firstElement);
	//        final ConditionBean cb = (ConditionBean) firstElement;
	//        return cb.getSqlClause().getClause();
	//	fmt.Printf("args %v \n",args[0])
	cb := args[0]
	cbbase := reflect.ValueOf(cb).Elem().FieldByName("BaseConditionBean").Interface()
	cbir := reflect.ValueOf(cbbase).MethodByName("GetSqlClause").Call([]reflect.Value{})
	sqlcr := cbir[0].Elem().MethodByName("GetClause").Call([]reflect.Value{})
	//return (*(*cb).GetSqlClause()).GetClause()
	return sqlcr[0].String()
}

type TnAbstractBasicSqlCommand struct {
	StatementFactory *StatementFactory
	sqlExecution     *SqlExecution
	rc               *ResourceContext
	ResultType       string
}

func (t *TnAbstractBasicSqlCommand) GetRootNode(args []interface{}) (*Node, error) {
	return nil, nil
}
func (t *TnAbstractBasicSqlCommand) GetArgNames() []string {
	return nil
}
func (t *TnAbstractBasicSqlCommand) GetArgTypes() []string {
	return nil
}
func (t *TnAbstractBasicSqlCommand) filterExecutedSql(executedSql string) string {
	return executedSql
}
func (t *TnAbstractBasicSqlCommand) newBasicParameterHandler(executedSql string) *ParameterHandler {
	return nil
}

type TnAbstractTwoWaySqlCommand struct {
	TnAbstractBasicSqlCommand
	IsBlockeNullParameter bool //defult true
}

func (t *TnAbstractTwoWaySqlCommand) Execute(args []interface{}, tx *sql.Tx, behavior *Behavior) (interface{}, error) {
	log.InternalDebug("TnAbstractTwoWaySqlCommand Execute")
	//Current
	//	        final Node rootNode = getRootNode(args);
	rootNode, err := (*t.sqlExecution).GetRootNode(args)
	if err != nil {
		return nil, err
	}
	ctx, err1 := t.apply(rootNode, args, (*t.sqlExecution).GetArgNames(), (*t.sqlExecution).GetArgTypes())
	if err1 != nil {
		return nil, err1
	}
	log.InternalDebug("ctx Sql :" + (*ctx).getSql())
	executedSql := (*t.sqlExecution).filterExecutedSql((*ctx).getSql())
	log.InternalDebug("executedSql :" + executedSql)
	handler := t.createBasicParameterHandler(ctx, executedSql)
	//fmt.Printf("handler %v %T\n", handler, handler)
	//fmt.Printf("behavior %v %T\n", behavior, behavior)
	bindVariables := (*ctx).getBindVariables()
	log.InternalDebug(fmt.Sprintf("bind Variables %v \n", bindVariables))
	bindVariableTypes := (*ctx).getBindVariableTypes()
	bindVariables = bindVariables
	bindVariableTypes = bindVariableTypes

	rtn, err := (*handler).execute(bindVariables, bindVariableTypes, tx, behavior)
	rtn = rtn
	//        return filterReturnValue(handler.execute(bindVariables, bindVariableTypes));

	return rtn, err
}
func (t *TnAbstractTwoWaySqlCommand) createBasicParameterHandler(ctx *CommandContext, executedSql string) *ParameterHandler {
	handler := (*t.sqlExecution).newBasicParameterHandler(executedSql)
	//        final Object[] bindVariables = context.getBindVariables();
	//        handler.setExceptionMessageSqlArgs(bindVariables);
	return handler
}

func (t *TnAbstractTwoWaySqlCommand) apply(rootNode *Node, args []interface{}, argNames []string, argTypes []string) (*CommandContext, error) {
	log.InternalDebug("TnAbstractTwoWaySqlCommand apply")
	log.InternalDebug(fmt.Sprintf("argNames %v argtypes %v \n", argNames, argTypes))
	//fmt.Printf("rootNode %v  \n", rootNode)
	ctx := t.createCommandContext(args, argNames, argTypes)
	err := (*rootNode).accept(ctx, rootNode)
	if err != nil {
		return nil, err
	}
	//fmt.Println("sql apply %s"+(*ctx).getSql())
	return ctx, nil
}
func (t *TnAbstractTwoWaySqlCommand) createCommandContext(args []interface{}, argNames []string, argTypes []string) *CommandContext {
	cr := new(CommandContextCreator)
	cr.argNames = argNames
	cr.argTypes = argTypes
	cc := cr.createCommandContext(args)
	return cc
}
func (t *TnAbstractTwoWaySqlCommand) AnalyzeTwoWaySql(twoWaySql string) (*Node, error) {
	sqlAnalyzer := t.CreateSqlAnalyzer(twoWaySql)
	return sqlAnalyzer.Analyze()
}
func (t *TnAbstractTwoWaySqlCommand) CreateSqlAnalyzer(twoWaySql string) *SqlAnalyzer {
	return (*t.rc).CreateSqlAnalyzer(twoWaySql, t.IsBlockeNullParameter)
}

type PrepareStatement struct {
}
type TnAbstractEntityDynamicCommand struct {
	TnAbstractBasicSqlCommand
	targetDBMeta  *DBMeta
	propertyNames *StringList
}

func (t *TnAbstractEntityDynamicCommand) getModifiedPropertyNames(entity *Entity) map[string]string {
	set := make(map[string]string)
	items := (*entity).GetModifiedPropertyNamesArray()
	for _, item := range items {
		set[item] = item
	}
	return set
}
func (t *TnAbstractEntityDynamicCommand) isOptimisticLockProperty(timestampProp string, versionNoProp string, propertyName string) bool {
	return propertyName == timestampProp || propertyName == versionNoProp
}

type TnUpdateEntityDynamicCommand struct {
	TnAbstractEntityDynamicCommand
	optimisticLockHandling         bool
	versionNoAutoIncrementOnMemory bool
}

func (t *TnUpdateEntityDynamicCommand) Execute(args []interface{}, tx *sql.Tx, behavior *Behavior) (interface{}, error) {
	if args == nil || len(args) == 0 {
		return nil, errors.New("The argument 'args' should not be null or empty.")
	}
	bean := args[0]
	var option *UpdateOption = (args[1]).(*UpdateOption)
	//未実装
	//        final UpdateOption<ConditionBean> option = extractUpdateOptionChecked(args);
	//        prepareStatementConfigOnThreadIfExists(option);
	//
	//        final TnPropertyType[] propertyTypes = createUpdatePropertyTypes(bean, option);
	var entity *Entity = bean.(*Entity)
	propertyTypes := t.createUpdatePropertyTypes(entity, nil)
	if propertyTypes.Size() == 0 {
		//        if (propertyTypes.length == 0) {
		//            if (isLogEnabled()) {
		//                log(createNonUpdateLogMessage(bean));
		//            }
		//            return getNonUpdateReturn();
		//        }
		return 1, nil
	}
	sql, err2 := t.createUpdateSql(propertyTypes, option)
	if err2 != nil {
		return 0, nil
	}
	log.InternalDebug(fmt.Sprintln("sql :" + sql))
	sql2 := t.filterExecutedSql(sql)
	sql2 = sql2
	return t.doExecute(entity, propertyTypes, sql2, option, tx)
}
func (t *TnUpdateEntityDynamicCommand) doExecute(entity *Entity, propertyTypes *List, sql string, option *UpdateOption, tx *sql.Tx) (int64, error) {
	handler := t.createUpdateEntityHandler(propertyTypes, sql, option)
	//        final Object[] realArgs = new Object[] { bean };
	//        handler.setExceptionMessageSqlArgs(realArgs);
	//        final int result = handler.execute(realArgs);
	res, err := handler.execute(entity, tx)
	//        return Integer.valueOf(result);
	return res, err
}

func (t *TnUpdateEntityDynamicCommand) createUpdateEntityHandler(propertyTypes *List, sql string, option *UpdateOption) *TnUpdateEntityHandler {
	handler := new(TnUpdateEntityHandler)
	var bsh BasicSqlHander = handler
	handler.BasicSqlHander = &bsh
	handler.boundPropTypes = propertyTypes
	handler.optimisticLockHandling = t.optimisticLockHandling
	handler.updateoption = option
	handler.statementFactory = t.StatementFactory
	handler.versionNoAutoIncrementOnMemory = t.versionNoAutoIncrementOnMemory
	handler.sql = sql
	return handler
	return nil
}
func (t *TnUpdateEntityDynamicCommand) filterExecutedSql(sql string) string {
	return sql
}
func (t *TnUpdateEntityDynamicCommand) createUpdateSql(propertyTypes *List, option *UpdateOption) (string, error) {
	tableDbName := (*t.targetDBMeta).GetTableDbName()
	//_beanMetaDataとtargetDBMetaの区別不明 要確認
	if (*t.targetDBMeta).HasPrimaryKey() == false {
		return "", errors.New("The table '" + tableDbName + "' should have primary key.")
	}
	sb := new(bytes.Buffer)
	sb.WriteString("update " + (*t.targetDBMeta).GetTableSqlName().TableSqlName + " set ")
	versionNoPropertyName := ""
	if (*t.targetDBMeta).HasVersionNo() {
		versionNoPropertyName = (*t.targetDBMeta).GetVersionNoColumnInfo().PropertyName
	}
	for i, ptx := range propertyTypes.data {
		var pt *TnPropertyType = ptx.(*TnPropertyType)
		columnDbName := pt.ColumnDbName
		columnDbName = columnDbName
		columnSqlName := pt.ColumnSqlName
		propertyName := pt.propetyName
		if i > 0 {
			sb.WriteString(", ")
		}
		if propertyName == versionNoPropertyName {
			if !t.versionNoAutoIncrementOnMemory {
				t.setupVersionNoAutoIncrementOnQuery(sb, columnSqlName)
				continue
			}
		}
		sb.WriteString(columnSqlName.ColumnSqlName + " = ")
		valueExp := ""
		//option 未実装
		//   if (option != nil && option.hasStatement(columnDbName)) {
		//                final String statement = option.buildStatement(columnDbName);
		//                valueExp = encryptIfNeeds(tableDbName, columnDbName, statement);
		//            } else {

		//                valueExp = encryptIfNeeds(tableDbName, columnDbName, "?");
		//            }
		//            sb.append(valueExp);
		valueExp = "?"
		sb.WriteString(valueExp)
		//        }

	}
	sb.WriteString(Ln + " where ")
	if (*t.targetDBMeta).HasPrimaryKey() {
		for i, colInfo := range (*t.targetDBMeta).GetPrimaryInfo().UniqueInfo.UniqueColumnList.data {
			var ci *ColumnInfo = colInfo.(*ColumnInfo)
			if i > 0 {
				sb.WriteString(" and ")
			}
			sb.WriteString(ci.ColumnSqlName.ColumnSqlName + " = ?")

		}
	}
	if t.optimisticLockHandling && (*t.targetDBMeta).HasVersionNo() {
		sb.WriteString(" and " +
			(*t.targetDBMeta).GetVersionNoColumnInfo().ColumnSqlName.ColumnSqlName + " = ?")
	}

	//        if (_optimisticLockHandling && _beanMetaData.hasTimestampPropertyType()) {
	//            final TnPropertyType pt = _beanMetaData.getTimestampPropertyType();
	//            sb.append(" and ").append(pt.getColumnSqlName()).append(" = ?");
	//        }

	return sb.String(), nil
}
func (t *TnUpdateEntityDynamicCommand) setupVersionNoAutoIncrementOnQuery(sb *bytes.Buffer, columnSqlName *ColumnSqlName) {
	sb.WriteString(columnSqlName.ColumnSqlName + " = " + columnSqlName.ColumnSqlName + " + 1")
}
func (t *TnUpdateEntityDynamicCommand) createUpdatePropertyTypes(entity *Entity, option *UpdateOption) *List {
	dbMeta := (*entity).GetDBMeta()
	dbMeta = dbMeta
	modifiedSet := t.getModifiedPropertyNames(entity)
	//fmt.Printf("ModifiedSet %v\n", modifiedSet)
	typeList := new(List)
	timestampProp := ""
	timestampProp = timestampProp
	versionNoProp := ""

	if (*(*entity).GetDBMeta()).HasVersionNo() {
		versionNoProp = (*(*entity).GetDBMeta()).GetVersionNoColumnInfo().PropertyName
		versionNoProp = versionNoProp
	}
	primaryKeyset := make(map[string]string)
	primaryKeyset = primaryKeyset
	if (*(*entity).GetDBMeta()).HasPrimaryKey() {
		primaryKeys := (*(*entity).GetDBMeta()).GetPrimaryInfo().UniqueInfo.UniqueColumnList.data
		for _, primaryKey := range primaryKeys {
			var colInfo *ColumnInfo = primaryKey.(*ColumnInfo)
			primaryKeyset[colInfo.PropertyName] = colInfo.PropertyName
		}
	}
	for _, propertyName := range t.propertyNames.data {
		_, ok := primaryKeyset[propertyName]
		if ok {
			continue
		}
		log.InternalDebug(fmt.Sprintf("propertyName %v\n", propertyName))
		//pt:=(*dbMeta).GetPropertyType(propertyName)
		if t.isOptimisticLockProperty(timestampProp, versionNoProp, propertyName) ||
			t.isSpecifiedProperty(option, modifiedSet, propertyName) ||
			t.isStatementProperty(option, propertyName) { // Statement
			typeList.Add((*(*entity).GetDBMeta()).GetPropertyType(propertyName))
			//fmt.Println("typeList Added :" + propertyName)
		}
	}
	return typeList
}
func (t *TnUpdateEntityDynamicCommand) isStatementProperty(option *UpdateOption, propertyName string) bool {
	//未実装
	//return option != null && option.hasStatement(pt.getColumnDbName());
	return false
}
func (t *TnUpdateEntityDynamicCommand) isSpecifiedProperty(option *UpdateOption, modifiedSet map[string]string, propertyName string) bool {
	//未実装
	//	        if (option != null && option.hasSpecifiedUpdateColumn()) { // BatchUpdate
	//            return option.isSpecifiedUpdateColumn(pt.getColumnDbName());
	//        } else { // EntityUpdate
	//            return isModifiedProperty(modifiedSet, pt); // process for ModifiedColumnUpdate
	//        }

	_, ok := modifiedSet[propertyName]
	//fmt.Printf("propety name %s ok %v\n",propertyName,ok)
	return ok
}

type OutsideSqlBasicExecutor struct {
	behaviorCommandInvoker   *BehaviorCommandInvoker
	tableDbName              string
	currentDBDef             *DBDef
	defaultStatementConfig   *StatementConfig
	outsideSqlOption         *OutsideSqlOption
	outsideSqlContextFactory *OutsideSqlContextFactory
	behavior                 *Behavior
}

func (o *OutsideSqlBasicExecutor) Execute(pmb interface{}, tx *sql.Tx) (int64, error) {

	path := reflect.ValueOf(pmb).MethodByName("GetOutsideSqlPath").Call([]reflect.Value{})[0].Interface().(string)
	return o.doExecute(path, pmb, tx)
}
func (o *OutsideSqlBasicExecutor) doExecute(path string, pmb interface{}, tx *sql.Tx) (int64, error) {
	fullpath, err1 := o.getFullPath(path)
	if err1 != nil {
		return 0, err1
	}
	cmd := o.createExecuteCommand(fullpath, pmb, tx)
	res, err := o.behaviorCommandInvoker.Invoke(cmd)
	if err !=nil{
		return 0,err
	}
	var lres int64 = res.(int64)
	return lres, err
}
func (o *OutsideSqlBasicExecutor) createExecuteCommand(path string, pmb interface{}, tx *sql.Tx) *BehaviorCommand {
	cmd := new(OutsideSqlExecuteCommand)

	var behaviorCommand BehaviorCommand = cmd
	cmd.BehaviorCommand = &behaviorCommand
	o.xsetupCommand(&cmd.AbstractOutsideSqlCommand, path, pmb)
	cmd.tx = tx
	cmd.BaseBehaviorCommand.Behavior = o.behavior
	//        {
	//            final OutsideSqlSelectListCommand<ENTITY> newed = newOutsideSqlSelectListCommand();
	//            cmd = xsetupCommand(newed, path, pmb); // has a little generic headache...
	//        }
	//        cmd.setEntityType(entityType);
	//        return cmd;
	return &behaviorCommand
}
func (o *OutsideSqlBasicExecutor) SelectList(pmb interface{}, tx *sql.Tx) (*ListResultBean, error) {
	if pmb == nil {
		return nil, errors.New("The argument 'pmb' (typed parameter-bean) should not be null.")
	}
	path := reflect.ValueOf(pmb).MethodByName("GetOutsideSqlPath").Call([]reflect.Value{})[0].Interface().(string)
	entityType := reflect.ValueOf(pmb).MethodByName("GetEntityType").Call([]reflect.Value{})[0].Interface().(string)
	return o.doSelectList(path, pmb, entityType, tx)
}
func (o *OutsideSqlBasicExecutor) getFullPath(path string) (string, error) {
	fullpath := filepath.Join(Gopath, "src", path)
	files, _ := filepath.Glob(fullpath)
	//	fmt.Printf("files %v %T %d\n",files,files,len(files))
	if len(files) == 0 {
		return "", errors.New("SQL File Not found. GOPATH NOT SET? " + fullpath)
	}
	return fullpath, nil
}
func (o *OutsideSqlBasicExecutor) doSelectList(path string, pmb interface{}, entityType string, tx *sql.Tx) (*ListResultBean, error) {
	//////////
	fullpath, err1 := o.getFullPath(path)
	if err1 != nil {
		return nil, err1
	}
	//	fmt.Println("PATH :"+Gopath+"/src/"+path)

	//        if (entityType == null) {
	//            String msg = "The argument 'entityType' for result should not be null: path=" + path;
	//            throw new IllegalArgumentException(msg);
	//        }
	//        try {
	//            List<ENTITY> resultList = invoke(createSelectListCommand(path, pmb, entityType));
	//            return createListResultBean(resultList);
	//        } catch (FetchingOverSafetySizeException e) { // occurs only when fetch-bean
	//            throwDangerousResultSizeException(pmb, e);
	//            return null; // unreachable
	//        }
	cmd := o.createSelectListCommand(fullpath, pmb, entityType, tx)
	res, err := o.behaviorCommandInvoker.Invoke(cmd)
	if err != nil {
		return nil, err
	}
	var lres *ListResultBean = res.(*ListResultBean)
	return lres, nil
}
func (o *OutsideSqlBasicExecutor) createSelectListCommand(path string, pmb interface{}, entityType string, tx *sql.Tx) *BehaviorCommand {
	cmd := new(OutsideSqlSelectListCommand)
	var behaviorCommand BehaviorCommand = cmd
	cmd.BehaviorCommand = &behaviorCommand
	o.xsetupCommand(&cmd.AbstractOutsideSqlCommand, path, pmb)
	cmd.entityType = entityType
	cmd.tx = tx
	cmd.BaseBehaviorCommand.Behavior = o.behavior
	//        {
	//            final OutsideSqlSelectListCommand<ENTITY> newed = newOutsideSqlSelectListCommand();
	//            cmd = xsetupCommand(newed, path, pmb); // has a little generic headache...
	//        }
	//        cmd.setEntityType(entityType);
	//        return cmd;
	return &behaviorCommand
}
func (o *OutsideSqlBasicExecutor) xsetupCommand(cmd *AbstractOutsideSqlCommand, path string, pmb interface{}) {
	cmd.TableDbName = o.tableDbName
	o.behaviorCommandInvoker.InjectComponentProperty(&cmd.BaseBehaviorCommand)
	cmd.OutsideSqlPath = path
	cmd.OutsideSqlOption = o.outsideSqlOption
	cmd.CurrentDBDef = o.currentDBDef
	cmd.outsideSqlContextFactory = o.outsideSqlContextFactory
	cmd.pmb = pmb
	//未実装
	//        cmd.setOutsideSqlFilter(_outsideSqlFilter);
	//        return cmd;
}

type OutsideSqlExecutorFactory interface {
	CreateBasic(behaviorCommandInvoker *BehaviorCommandInvoker, bhv *Behavior, tableDbName string, currentDBDef *DBDef, defaultStatementConfig *StatementConfig, outsideSqlOption *OutsideSqlOption) *OutsideSqlBasicExecutor
}

type OutsideSqlOption struct {
	TableDbName string
}

func (o *OutsideSqlOption) GenerateUniqueKey() string {
	//	       return "{" + _pagingRequestType + "/" + _removeBlockComment + "/" + _removeLineComment + "/" + _formatSql + "}";
	return "{}"
}

type StatementConfig interface {
}
type DefaultOutsideSqlExecutorFactory struct {
}

func (d *DefaultOutsideSqlExecutorFactory) CreateBasic(behaviorCommandInvoker *BehaviorCommandInvoker, bhv *Behavior, tableDbName string, currentDBDef *DBDef, defaultStatementConfig *StatementConfig, outsideSqlOption *OutsideSqlOption) *OutsideSqlBasicExecutor {
	outsideSqlContextFactory := d.createOutsideSqlContextFactory()
	//未実装
	//        final OutsideSqlFilter outsideSqlFilter = createOutsideSqlExecutionFilter();
	ex := new(OutsideSqlBasicExecutor)
	ex.behaviorCommandInvoker = behaviorCommandInvoker
	ex.currentDBDef = currentDBDef
	ex.defaultStatementConfig = defaultStatementConfig
	ex.outsideSqlContextFactory = outsideSqlContextFactory
	ex.tableDbName = tableDbName
	if outsideSqlOption != nil {
		ex.outsideSqlOption = outsideSqlOption
	} else {
		ex.outsideSqlOption = new(OutsideSqlOption)
		ex.outsideSqlOption.TableDbName = tableDbName
	}
	ex.behavior = bhv
	return ex
}
func (d *DefaultOutsideSqlExecutorFactory) createOutsideSqlContextFactory() *OutsideSqlContextFactory {
	factory := new(DefaultOutsideSqlContextFactory)
	var fi OutsideSqlContextFactory = factory
	return &fi
}

type OutsideSqlContextFactory interface {
}
type DefaultOutsideSqlContextFactory struct {
}
type ListHandlingPmb interface {
	GetEntityType() string
	GetOutsideSqlPath() string
}
type OutsideSqlExecuteExecution struct {
	AbstractOutsideSqlExecution
}

func (o *OutsideSqlExecuteExecution) newBasicParameterHandler(sql string) *ParameterHandler {
	bh := new(TnBasicUpdateHandler)
	bh.sql = sql
	//fmt.Printf("StatementFactoty %v\n", o.StatementFactory)
	bh.statementFactory = o.StatementFactory
	var ph ParameterHandler = bh
	return &ph

}

type OutsideSqlSelectExecution struct {
	AbstractOutsideSqlExecution
	ResultType string
}

func (o *OutsideSqlSelectExecution) newBasicParameterHandler(executedSql string) *ParameterHandler {
	handler := new(TnBasicSelectHandler)
	handler.sql = executedSql
	handler.statementFactory = o.StatementFactory
	log.InternalDebug(fmt.Sprintln("o.ResultType:", o.ResultType))
	handler.ResultType = o.ResultType
	var hand ParameterHandler = handler
	return &hand
}

type AbstractOutsideSqlExecution struct {
	AbstractFixedSqlExecution
	removeBlockComment bool
	removeLineComment  bool
	formatSql          bool
	//outsideSqlFilter
}

func (a *AbstractOutsideSqlExecution) filterExecutedSql(executedSql string) string {
	//	        executedSql = super.filterExecutedSql(executedSql);
	//        executedSql = doFilterExecutedSqlByOutsideSqlFilter(executedSql);
	//        if (_removeBlockComment) {
	//            executedSql = Srl.removeBlockComment(executedSql);
	//        }
	//        if (_removeLineComment) {
	//            executedSql = Srl.removeLineComment(executedSql);
	//        }
	//        if (_formatSql) {
	//            executedSql = Srl.removeEmptyLine(executedSql);
	//        }
	//        executedSql = doFilterExecutedSqlByCallbackFilter(executedSql);
	return executedSql
}

type AbstractFixedSqlExecution struct {
	AbstractFixedArgExecution
	rootNode *Node
}

func (a *AbstractFixedSqlExecution) GetRootNode(args []interface{}) (*Node, error) {
	return a.rootNode, nil
}

// DefaultOutsideSqlContextFactory
//
// DefaultOutsideSqlExecutorFactory @CreateBasic @createOutsideSqlContextFactory
//
// OutsideSqlBasicExecutor #behavior #behaviorCommandInvoker #currentDBDef #defaultStatementConfig #outsideSqlContextFactory #outsideSqlOption #tableDbName @SelectList @createSelectListCommand @doSelectList @xsetupCommand
//
// OutsideSqlOption #TableDbName @GenerateUniqueKey
//
// OutsideSqlSelectExecution #ResultType @newBasicParameterHandler
// AbstractOutsideSqlExecution #formatSql #removeBlockComment #removeLineComment @filterExecutedSql
// AbstractFixedSqlExecution #rootNode @GetRootNode
// AbstractFixedArgExecution #ArgNames #ArgTypes @GetArgNames @GetArgTypes
// TnAbstractTwoWaySqlCommand @AnalyzeTwoWaySql @CreateSqlAnalyzer @Execute @apply @createBasicParameterHandler @createCommandContext
// TnAbstractBasicSqlCommand #ResultType #StatementFactory #rc #sqlExecution @GetArgNames @GetArgTypes @GetRootNode @filterExecutedSql @newBasicParameterHandler
//
// PrepareStatement
//
// SelectCBExecution @ExtractTwoWaySql @GetRootNode @filterExecutedSql @newBasicParameterHandler
// AbstractFixedArgExecution #ArgNames #ArgTypes @GetArgNames @GetArgTypes
// TnAbstractTwoWaySqlCommand @AnalyzeTwoWaySql @CreateSqlAnalyzer @Execute @apply @createBasicParameterHandler @createCommandContext
// TnAbstractBasicSqlCommand #ResultType #StatementFactory #rc #sqlExecution @GetArgNames @GetArgTypes @GetRootNode @filterExecutedSql @newBasicParameterHandler
//
// SelectNextValExecution
// SelectSimpleExecution @IsBlockNullParameter @newBasicParameterHandler
// AbstractFixedSqlExecution #rootNode @GetRootNode
// AbstractFixedArgExecution #ArgNames #ArgTypes @GetArgNames @GetArgTypes
// TnAbstractTwoWaySqlCommand @AnalyzeTwoWaySql @CreateSqlAnalyzer @Execute @apply @createBasicParameterHandler @createCommandContext
// TnAbstractBasicSqlCommand #ResultType #StatementFactory #rc #sqlExecution @GetArgNames @GetArgTypes @GetRootNode @filterExecutedSql @newBasicParameterHandler
//
// TnDeleteEntityStaticCommand @setupSql
// TnAbstractEntityStaticCommand #optimisticLockHandling #propertyNames #propertyTypes #sql #targetDBMeta #versionNoAutoIncrementOnMemory @Execute @doExecute @setupDeleteSql @setupDeleteWhere
// TnAbstractBasicSqlCommand #ResultType #StatementFactory #rc #sqlExecution @GetArgNames @GetArgTypes @GetRootNode @filterExecutedSql @newBasicParameterHandler
//
// TnInsertEntityDynamicCommand @Execute @createInsertPropertyTypes @createInsertSql @doExecute @isSpecifiedProperty
// TnAbstractEntityDynamicCommand #propertyNames #targetDBMeta @getModifiedPropertyNames @isOptimisticLockProperty
// TnAbstractBasicSqlCommand #ResultType #StatementFactory #rc #sqlExecution @GetArgNames @GetArgTypes @GetRootNode @filterExecutedSql @newBasicParameterHandler
//
// TnUpdateEntityDynamicCommand #optimisticLockHandling #versionNoAutoIncrementOnMemory @Execute @createUpdateEntityHandler @createUpdatePropertyTypes @createUpdateSql @doExecute @filterExecutedSql @isSpecifiedProperty @isStatementProperty @setupVersionNoAutoIncrementOnQuery
// TnAbstractEntityDynamicCommand #propertyNames #targetDBMeta @getModifiedPropertyNames @isOptimisticLockProperty
// TnAbstractBasicSqlCommand #ResultType #StatementFactory #rc #sqlExecution @GetArgNames @GetArgTypes @GetRootNode @filterExecutedSql @newBasicParameterHandler
//
