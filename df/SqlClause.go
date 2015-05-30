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
	"bytes"
	//	"container/list"
	"fmt"
	"github.com/mikeshimura/dbflute/log"
	"strconv"
	"strings"
)

const (
	SelectClauseType_COLUMNS        = 0
	SelectClauseType_UNIQUE_COUNT   = 1
	SelectClauseType_PLAIN_COUNT    = 2
	SelectClauseType_COUNT_DISTINCT = 3
	SelectClauseType_MAX            = 4
	SelectClauseType_MIN            = 5
	SelectClauseType_SUM            = 6
	SelectClauseType_AVE            = 7
	RELATION_PATH_DELIMITER         = "_"
	SQIP_BEGIN_MARK_PREFIX          = "--#df:sqbegin#"
	SQIP_END_MARK_PREFIX            = "--#df:sqend#"
	SQIP_IDENTITY_TERMINAL          = "#df:idterm#"
)

type SqlClause interface {
	setTableDbName(tn string)
	GetTableDbName() string
	setSqlSetting(dc *DBCurrent)
	GetSqlSetting() *DBCurrent
	GetBasePorintAliasName() string
	GetOrderByClause() *OrderByClause
	IsOrScopeQueryEffective() bool
	IsOrScopeQueryAndPartEffective() bool
	RegisterWhereClause(crn *ColumnRealName, key *ConditionKey,
		cvalue *ConditionValue, co *ConditionOption, usedAliasName string)
	GetWhereList() *List
	SetDBMeta(dm *DBMeta)
	SetUseSelectIndex(si bool)
	GetClause() string
	RegisterBaseTableInlineWhereClause(columnSqlName *ColumnSqlName,
		key *ConditionKey, cvalue *ConditionValue, option *ConditionOption)
	MakeOrScopeQueryEffective()
	CloseOrScopeQuery()
	BeginOrScopeQueryAndPart()
	EndOrScopeQueryAndPart()
	AllowEmptyStringQuery()
	IsAllowEmptyStringQuery() bool
	ResolveRelationNo(localTableName string, foreignPropertyName string) int
	ResolveJoinAliasName(relationPath string) string
	RegisterOuterJoin(foreignAliasName string, foreignTableDbName string, localAliasName string,
		localTableDbName string, joinOnMap map[*ColumnRealName]*ColumnRealName,
		foreignInfo *ForeignInfo, fixedCondition string, resolver *FixedConditionResolver)
	RegisterSelectedRelation(foreignTableAliasName string, TableDbName string,
		foreignPropertyName string, localRelationPath string,
		foreignRelationPath string, foreignTableDbName string)
	DoFetchPage()
	DoClearFetchPageClause()
	DoFetchFirst()
	FetchFirst(fetchSize int)
	CreateFromHint() string
	CreateSqlSuffix() string
	CreateSelectHint() string
	GetBaseSqlClause() *BaseSqlClause
	GetClauseQueryUpdate(columnParameterKey *StringList, columnParameterValue *StringList) string
	GetClauseQueryDelete() string
}

type BaseSqlClause struct {
	TableDbName string
	//SqlSetting is copy for DBCurrent can chage setting
	SqlSetting                         *DBCurrent
	BasePorintAliasName                string
	OrderByClause                      *OrderByClause
	OrScopeQueryEffective              bool
	WhereList                          *List
	Inline                             bool
	OnClause                           bool
	DBMeta                             *DBMeta
	UseSelectIndex                     bool
	SelectIndexMap                     *StringKeyMap
	SelectIndexReverseMap              *StringKeyMap
	BaseTableInlineWhereList           *List
	OrScopeQueryAndPartEffective       bool
	currentTmpOrScopeQueryInfo         *OrScopeQueryInfo
	orScopeQueryAndPartIdentity        int
	outerJoinMap                       map[string]*LeftOuterJoinInfo
	outerJoinList                      StringList
	emptyStringQueryAllowed            bool
	subQueryLevel                      int
	whereUsedInnerJoinAllowed          bool
	orScopeQueryEffective              bool
	InnerJoinLazyReflector             *List
	pagingAdjustment                   bool
	pagingCountLeastJoin               bool
	unionQueryInfoList                 *List
	selectClauseType                   *SelectClauseType
	structuralPossibleInnerJoinAllowed bool
	checkInvalidQuery                  bool
	disableSelectIndex                 bool
	selectedRelationBasicMap           map[string]string
	selectedRelationColumnMap          map[string]*List
	selectecNextConnectingRelationMap  map[string]string
	selectClauseRealColumnAliasMap     map[string]string
	selectClauseInfo                   StringList
	fetchScopeEffective                bool
	fetchStartIndex                    int
	fetchSize                          int
	fetchPageNumber                    int
	sqlClause                          *SqlClause
	subQueryIndentProcessor            *SubQueryIndentProcessor
}

func (b *BaseSqlClause) GetClauseQueryDelete() string {
	dbmeta := b.DBMeta
	sb := new(bytes.Buffer)
	sb.WriteString("delete")
	useQueryUpdateDirect := b.isUseQueryUpdateDirect(dbmeta)
	whereClause := ""
	if useQueryUpdateDirect { // prepare for direct case
		whereClause = b.processSubQueryIndent(b.GetWhereClause())
		//            if (needsDeleteTableAliasHint(whereClause)) {
		//                sb.append(" ").append(getBasePointAliasName());
		//            }
	}
	sb.WriteString(" from " + (*dbmeta).GetTableSqlName().TableSqlName)
	if useQueryUpdateDirect { // direct (in-scope unsupported or compound primary keys)
		b.buildQueryUpdateDirectClause(nil, nil, whereClause, dbmeta, sb)
	} else { // basically here
		b.buildQueryUpdateInScopeClause(nil, nil, dbmeta, sb)
	}
	return sb.String()
	panic("GetClauseQueryDelete")
	return ""
}
func (b *BaseSqlClause) GetClauseQueryUpdate(
	columnParameterKey *StringList, columnParameterValue *StringList) string {
	if columnParameterKey.Size() == 0 {
		return ""
	}

	dbmeta := b.DBMeta
	sb := new(bytes.Buffer)
	sb.WriteString("update " + (*dbmeta).GetTableSqlName().TableSqlName)

	if b.isUseQueryUpdateDirect(dbmeta) { // direct (in-scope unsupported or compound primary keys)
		whereClause := b.processSubQueryIndent(b.GetWhereClause())
		b.buildQueryUpdateDirectClause(columnParameterKey,
			columnParameterValue, whereClause, dbmeta, sb)
	} else { // basically here
		b.buildQueryUpdateInScopeClause(columnParameterKey,
			columnParameterValue, dbmeta, sb)
	}
	return sb.String()

}
func (b *BaseSqlClause) GetWhereClause() string {
	//
	//	        reflectClauseLazilyIfExists();
	sb := new(bytes.Buffer)
	b.BuildWhereClause(sb)
	return sb.String()
}
func (b *BaseSqlClause) processSubQueryIndent(sql string) string {
	return b.getSubQueryIndentProcessor().processSubQueryIndent(sql, "", sql)
}
func (b *BaseSqlClause) getSubQueryIndentProcessor() *SubQueryIndentProcessor {
	if b.subQueryIndentProcessor == nil {
		b.subQueryIndentProcessor = new(SubQueryIndentProcessor)
	}
	return b.subQueryIndentProcessor
}
func (b *BaseSqlClause) isUseQueryUpdateDirect(dbmeta *DBMeta) bool {
	//return _queryUpdateForcedDirectAllowed || !canUseQueryUpdateInScope(dbmeta);
	return !b.canUseQueryUpdateInScope(dbmeta)
}
func (b *BaseSqlClause) canUseQueryUpdateInScope(dbmeta *DBMeta) bool {
	return b.isUpdateSubQueryUseLocalTableSupported() && !(*dbmeta).HasCompoundPrimaryKey()
}

func (b *BaseSqlClause) isUpdateSubQueryUseLocalTableSupported() bool {
	return true
}
func (b *BaseSqlClause) buildQueryUpdateDirectClause(
	columnParameterKey *StringList, columnParameterValue *StringList,
	whereClause string, dbmeta *DBMeta, sb *bytes.Buffer) {
	//        if (hasUnionQuery()) {
	//            throwQueryUpdateUnavailableFunctionException("union", dbmeta);
	//        }
	useAlias := false
	//if (b.isUpdateTableAliasNameSupported()) {
	//            if (hasQueryUpdateSubQueryPossible(whereClause)) {
	//                useAlias = true;
	//            }
	//} else {
	//            if (hasQueryUpdateSubQueryPossible(whereClause)) {
	//                throwQueryUpdateUnavailableFunctionException("sub-query", dbmeta);
	//            }
	//    }
	directJoin := false
	//        if (isUpdateDirectJoinSupported()) {
	//            if (hasOuterJoin()) {
	//                useAlias = true; // use alias forcedly if direct join
	//                directJoin = true;
	//            }
	//        } else { // direct join unsupported
	if b.hasOuterJoin() {
		panic("QueryUpdateUnavailableFunction")
	}
	//        }
	basePointAliasName := ""
	if useAlias {
		basePointAliasName = b.GetBasePorintAliasName()
	}
	//
	if useAlias {
		sb.WriteString(" " + basePointAliasName)
	}
	directJoin = directJoin
	//        if (directJoin) {
	//            sb.WriteString(b.GetLeftOuterJoinClause());
	//        }
	if columnParameterKey != nil {
		//            final String setClauseAliasName = useAlias ? basePointAliasName + "." : null;
		setClauseAliasName := ""
		b.buildQueryUpdateSetClause(columnParameterKey, columnParameterValue, dbmeta, sb, setClauseAliasName)
	}
	if len(strings.Trim(whereClause, " ")) == 0 {
		return
	}

	//        if (useAlias) {
	//            sb.append(whereClause);
	//        } else {
	sb.WriteString(b.filterQueryUpdateBasePointAliasNameLocalUnsupported(whereClause))
	fmt.Println("sb:" + sb.String())

	panic("buildQueryUpdateDirectClause")

}
func (b *BaseSqlClause) filterQueryUpdateBasePointAliasNameLocalUnsupported(
	subQuery string) string {
	// remove table alias prefix for column
	subQuery = strings.Replace(subQuery, b.GetBasePorintAliasName()+".", "", -1)

	// remove table alias definition
	tableAliasSymbol := " " + b.GetBasePorintAliasName()
	subQuery = strings.Replace(subQuery, tableAliasSymbol+" ", " ", -1)
	subQuery = strings.Replace(subQuery, tableAliasSymbol+Ln, Ln, -1)
	if subQuery[len(subQuery)-len(tableAliasSymbol):] == tableAliasSymbol {
		subQuery = strings.Replace(subQuery, tableAliasSymbol, "", -1)
	}
	return subQuery
}
func (b *BaseSqlClause) hasOuterJoin() bool {
	return b.outerJoinMap != nil && len(b.outerJoinMap) > 0
}
func (b *BaseSqlClause) buildQueryUpdateInScopeClause(
	columnParameterKey *StringList, columnParameterValue *StringList,
	dbmeta *DBMeta, sb *bytes.Buffer) {

	b.buildQueryUpdateSetClause(columnParameterKey, columnParameterValue, dbmeta, sb, "")
	fmt.Println("sb:" + sb.String())
	primaryKeyName := ((*dbmeta).GetPrimaryUniqueInfo().UniqueColumnList.Get(0)).(*ColumnInfo).ColumnSqlName.ColumnSqlName
	columnSqlName := ((*dbmeta).GetPrimaryUniqueInfo().UniqueColumnList.Get(0)).(*ColumnInfo).ColumnSqlName.ColumnSqlName
	selectClause := "select " + b.GetBasePorintAliasName() + "." + columnSqlName
	fromWhereClause := b.buildClauseFromWhereAsTemplate(false)
	// Replace template marks. These are very important!
	fromWhereClause = strings.Replace(fromWhereClause, b.GetUnionSelectClauseMark(), selectClause, -1)
	fromWhereClause = strings.Replace(fromWhereClause, b.GetUnionWhereClauseMark(), "", -1)
	fromWhereClause = strings.Replace(fromWhereClause, b.GetUnionWhereFirstConditionMark(), "", -1)
	subQuery := b.processSubQueryIndent(selectClause + " " + fromWhereClause)
	sb.WriteString(Ln + " where " + primaryKeyName + " in (" + Ln + subQuery)
	if subQuery[len(subQuery)-len(Ln):] != Ln {
		sb.WriteString(Ln)
	}
	sb.WriteString(")")
}
func (b *BaseSqlClause) buildClauseFromWhereAsTemplate(template bool) string {
	sb := new(bytes.Buffer)
	b.BuildFromClause(sb)
	sb.WriteString((*b.sqlClause).CreateFromHint())
	b.BuildWhereClauseSub(sb, template)
	sb.WriteString(b.prepareUnionClause(b.GetUnionSelectClauseMark()))
	return sb.String()
}
func (b *BaseSqlClause) prepareUnionClause(selectClause string) string {
	if !b.hasUnionQuery() {
		return ""
	}
	panic("union not implemented yet")
	//sb := new(bytes.Buffer)
	//        for (UnionQueryInfo unionQueryInfo : _unionQueryInfoList) {
	//            final UnionClauseProvider unionClauseProvider = unionQueryInfo.getUnionClauseProvider();
	//            final String unionQueryClause = unionClauseProvider.provide();
	//            final boolean unionAll = unionQueryInfo.isUnionAll();
	//            sb.append(ln()).append(unionAll ? " union all " : " union ").append(ln());
	//            sb.append(selectClause).append(" ").append(unionQueryClause);
	//        }
	//        return sb.toString();
	return ""
}
func (b *BaseSqlClause) hasUnionQuery() bool {
	return b.unionQueryInfoList != nil && b.unionQueryInfoList.Size() > 0
}
func (b *BaseSqlClause) GetUnionSelectClauseMark() string {
	return "#df:unionSelectClause#"
}
func (b *BaseSqlClause) GetWhereClauseMark() string {
	return "#df:whereClause#"
}
func (b *BaseSqlClause) GetUnionWhereClauseMark() string {
	return "#df:unionWhereClause#"
}
func (b *BaseSqlClause) GetUnionWhereFirstConditionMark() string {
	return "#df:unionWhereFirstCondition#"
}
func (b *BaseSqlClause) buildQueryUpdateSetClause(
	columnParameterKey *StringList, columnParameterValue *StringList,
	dbmeta *DBMeta, sb *bytes.Buffer, aliasName string) {
	if columnParameterKey == nil {
		return
	}
	sb.WriteString(Ln)
	mapSize := columnParameterKey.Size()
	for index, propertyName := range columnParameterKey.data {
		parameter := columnParameterValue.Get(index)
		isVersionColumn, versionParameter := b.checkVersionColumn(parameter)
		fmt.Println("propertyName" + propertyName)
		columnInfo := (*dbmeta).GetColumnInfoByPropertyName(propertyName)
		columnSqlName := columnInfo.ColumnSqlName.ColumnSqlName
		if index == 0 {
			sb.WriteString("   set ")
		} else {
			sb.WriteString("     , ")
		}
		if aliasName != "" {
			sb.WriteString(aliasName)
		}
		sb.WriteString(columnSqlName + " = ")
		valueExp := ""
		if isVersionColumn {
			valueExp = aliasName + versionParameter
		} else {
			valueExp = parameter
		}
		sb.WriteString(valueExp)
		if mapSize-1 > index { // before last loop
			sb.WriteString(Ln)
		}
	}
}
func (b *BaseSqlClause) checkVersionColumn(columnName string) (bool, string) {
	if columnName[0:9] == ":Version:" {
		return true, columnName[9:]
	}
	return false, ""
}
func (b *BaseSqlClause) GetBaseSqlClause() *BaseSqlClause {
	return b
}
func (b *BaseSqlClause) FetchFirst(fetchSize int) {
	b.fetchScopeEffective = true
	if fetchSize < 0 {
		panic("Argument[fetchSize] should be plus")
	}
	b.fetchStartIndex = 0
	b.fetchSize = fetchSize
	b.fetchPageNumber = 1
	(*b.sqlClause).DoClearFetchPageClause()
	(*b.sqlClause).DoFetchFirst()
}
func (b *BaseSqlClause) getPageStartIndex() int {
	if b.fetchPageNumber <= 0 {
		panic("fetchPageNumber must be plus ")
	}
	return b.fetchStartIndex + b.fetchSize*(b.fetchPageNumber-1)
}
func (b *BaseSqlClause) getSelectClauseRealColumnAliasMap() map[string]string {
	if b.selectClauseRealColumnAliasMap == nil {
		b.selectClauseRealColumnAliasMap = make(map[string]string)
	}
	return b.selectClauseRealColumnAliasMap
}
func (b *BaseSqlClause) RegisterSelectedRelation(foreignTableAliasName string, TableDbName string,
	foreignPropertyName string, localRelationPath string,
	foreignRelationPath string, foreignTableDbName string) {
	//        assertObjectNotNull("foreignTableAliasName", foreignTableAliasName);
	//        assertObjectNotNull("localTableDbName", localTableDbName);
	//        assertObjectNotNull("foreignPropertyName", foreignPropertyName);
	//        assertObjectNotNull("foreignRelationPath", foreignRelationPath);
	if b.selectedRelationBasicMap == nil {
		b.selectedRelationBasicMap = make(map[string]string)
	}
	b.selectedRelationBasicMap[foreignRelationPath] = foreignPropertyName
	columnMap := b.createSelectedSelectColumnInfo(foreignTableAliasName,
		TableDbName, foreignPropertyName, localRelationPath, foreignTableDbName)
	if b.selectedRelationColumnMap == nil {
		b.selectedRelationColumnMap = make(map[string]*List)
	}
	b.selectedRelationColumnMap[foreignTableAliasName] = columnMap
	b.analyzeSelectedNextConnectingRelation(foreignRelationPath)
}
func (b *BaseSqlClause) createSelectedSelectColumnInfo(foreignTableAliasName string,
	localTableDbName string, foreignPropertyName string, localRelationPath string,
	foreignTableDbName string) *List {
	dbmeta := DBMetaProvider_I.TableDbNameInstanceMap[localTableDbName]
	foreignInfo := (*dbmeta).FindForeignInfo(foreignPropertyName)
	relationNo := foreignInfo.RelationNo
	nextRelationPath := RELATION_PATH_DELIMITER + strconv.Itoa(relationNo)
	if localRelationPath != "" {
		nextRelationPath = localRelationPath + nextRelationPath
	}
	resultList := new(List)
	foreignDBMeta := DBMetaProvider_I.TableDbNameInstanceMap[foreignTableDbName]
	columnInfoList := (*foreignDBMeta).GetColumnInfoList()
	for _, col := range columnInfoList.data {
		columnInfo := col.(*ColumnInfo)
		columnDbName := columnInfo.ColumnDbName
		selectColumnInfo := new(SelectedRelationColumn)
		selectColumnInfo.tableAliasName = foreignTableAliasName
		selectColumnInfo.columnInfo = columnInfo
		selectColumnInfo.columnAliasName = columnDbName + nextRelationPath
		resultList.Add(selectColumnInfo)
	}
	return resultList
}
func (b *BaseSqlClause) analyzeSelectedNextConnectingRelation(
	foreignRelationPath string) {
	if len(foreignRelationPath) <= 3 { // fast check e.g. _12, _3
		return // three characters cannot make two elements
	}
	delimiter := RELATION_PATH_DELIMITER
	delimiterCount := StringCount(foreignRelationPath, delimiter)
	if delimiterCount < 2 { // has no previous relation
		return
	}
	previousPath := substringLastFront(foreignRelationPath, delimiter)
	if b.selectecNextConnectingRelationMap == nil {
		b.selectecNextConnectingRelationMap = make(map[string]string)
	}
	b.selectecNextConnectingRelationMap[previousPath] = previousPath
}
func (b *BaseSqlClause) RegisterOuterJoin(foreignAliasName string,
	foreignTableDbName string, localAliasName string,
	localTableDbName string, joinOnMap map[*ColumnRealName]*ColumnRealName,
	foreignInfo *ForeignInfo, fixedCondition string,
	fixedConditionResolver *FixedConditionResolver) {
	//        assertAlreadyOuterJoin(foreignAliasName);
	//        assertJoinOnMapNotEmpty(joinOnMap, foreignAliasName);
	outerJoinMap := b.getOuterJoinMap()
	joinInfo := new(LeftOuterJoinInfo)
	joinInfo.foreignAliasName = foreignAliasName
	joinInfo.foreignTableDbName = foreignTableDbName
	joinInfo.localAliasName = localAliasName
	joinInfo.localTableDbName = localTableDbName
	joinInfo.joinOnMap = joinOnMap
	localJoinInfo := outerJoinMap[localAliasName]
	if localJoinInfo != nil { // means local is also joined (not base point)
		joinInfo.localJoinInfo = localJoinInfo
	}
	joinInfo.pureFK = foreignInfo.IsPureFK()
	joinInfo.notNullFKColumn = foreignInfo.IsNotNullFKColumn()
	joinInfo.fixedCondition = fixedCondition
	joinInfo.fixedConditionResolver = fixedConditionResolver
	//        // it should be resolved before registration because
	//        // the process may have Query(Relation) as precondition
	joinInfo.ResolveFixedCondition()
	outerJoinMap[foreignAliasName] = joinInfo
	b.outerJoinList.Add(foreignAliasName)
}

func (b *BaseSqlClause) IsAllowEmptyStringQuery() bool {
	return b.emptyStringQueryAllowed
}
func (b *BaseSqlClause) AllowEmptyStringQuery() {
	b.emptyStringQueryAllowed = true
}
func (b *BaseSqlClause) EndOrScopeQueryAndPart() {
	//assertCurrentTmpOrScopeQueryInfo();
	b.OrScopeQueryAndPartEffective = false
}
func (b *BaseSqlClause) BeginOrScopeQueryAndPart() {
	//assertCurrentTmpOrScopeQueryInfo();
	b.orScopeQueryAndPartIdentity++
	b.OrScopeQueryAndPartEffective = true
}
func (b *BaseSqlClause) CloseOrScopeQuery() {
	//	        assertCurrentTmpOrScopeQueryInfo();
	parentInfo := b.currentTmpOrScopeQueryInfo.parentInfo
	if parentInfo != nil {
		b.currentTmpOrScopeQueryInfo = parentInfo
	} else {
		b.reflectTmpOrClauseToRealObject(b.currentTmpOrScopeQueryInfo)
		b.clearOrScopeQuery()
	}
}
func (b *BaseSqlClause) clearOrScopeQuery() {
	b.currentTmpOrScopeQueryInfo = nil
	b.OrScopeQueryEffective = false
	b.OrScopeQueryAndPartEffective = false
}
func (b *BaseSqlClause) getOuterJoinMap() map[string]*LeftOuterJoinInfo {
	if b.outerJoinMap == nil {
		b.outerJoinMap = make(map[string]*LeftOuterJoinInfo)
	}
	return b.outerJoinMap
}
func (b *BaseSqlClause) reflectTmpOrClauseToRealObject(localInfo *OrScopeQueryInfo) {
	reflector := new(OrScopeQueryReflector)
	if b.WhereList == nil {
		b.WhereList = new(List)
	}
	if b.BaseTableInlineWhereList == nil {
		b.BaseTableInlineWhereList = new(List)
	}
	reflector.whereList = b.WhereList
	reflector.baseTableInlineWhereList = b.BaseTableInlineWhereList
	reflector.outerJoinMap = b.getOuterJoinMap()
	reflector.reflectTmpOrClauseToRealObject(localInfo)
}
func (b *BaseSqlClause) MakeOrScopeQueryEffective() {
	tmpOrScopeQueryInfo := new(OrScopeQueryInfo)
	if b.currentTmpOrScopeQueryInfo != nil {
		b.currentTmpOrScopeQueryInfo.addChildInfo(tmpOrScopeQueryInfo)
	}
	b.currentTmpOrScopeQueryInfo = tmpOrScopeQueryInfo
	b.OrScopeQueryEffective = true
}

func (b *BaseSqlClause) RegisterBaseTableInlineWhereClause(columnSqlName *ColumnSqlName, key *ConditionKey, value *ConditionValue, option *ConditionOption) {
	clauseList := b.GetBaseTableInlineWhereClauseList4Register()
	inlineBaseAlias := "dfinlineloc"
	columnRealName := new(ColumnRealName)
	columnRealName.TableAliasName = inlineBaseAlias
	columnRealName.ColumnSqlName = columnSqlName
	b.DoRegisterWhereClause(clauseList, columnRealName, key, value, option, true, false)

}
func (b *BaseSqlClause) GetBaseTableInlineWhereClauseList4Register() *List {
	if b.OrScopeQueryEffective {
		return b.GetTmpOrBaseTableInlineWhereList()
	} else {
		return b.GetBaseTableInlineWhereList()
	}
}
func (b *BaseSqlClause) GetTmpOrBaseTableInlineWhereList() *List {
	return &b.currentTmpOrScopeQueryInfo.tmpOrBaseTableInlineWhereList
}
func (b *BaseSqlClause) GetBaseTableInlineWhereList() *List {
	if b.BaseTableInlineWhereList == nil {
		b.BaseTableInlineWhereList = new(List)
	}
	return b.BaseTableInlineWhereList
}
func (b *BaseSqlClause) SetUseSelectIndex(si bool) {
	b.UseSelectIndex = si
}

func (b *BaseSqlClause) SetDBMeta(dm *DBMeta) {
	b.DBMeta = dm
}
func (b *BaseSqlClause) GetClause() string {
	log.InternalDebug("GetClause()")
	//        reflectClauseLazilyIfExists();
	sb := new(bytes.Buffer)
	selectClause := b.GetSelectClause()
	sb.WriteString(selectClause)
	b.BuildClauseWithoutMainSelect(sb, selectClause)
	sql := sb.String()
	log.InternalDebug("sql:" + sql)
	//未実装
	//        sql = filterEnclosingClause(sql);
	//        sql = processSubQueryIndent(sql);
	return sql
}
func (b *BaseSqlClause) BuildClauseWithoutMainSelect(sb *bytes.Buffer, selectClause string) {
	b.BuildFromClause(sb)
	sb.WriteString((*b.sqlClause).CreateFromHint())
	b.BuildWhereClause(sb)
	//        sb.append(deleteUnionWhereTemplateMark(prepareUnionClause(selectClause)));
	//	        if (!b.needsUnionNormalSelectEnclosing()) {
	sb.WriteString(b.PrepareClauseOrderBy())
	sb.WriteString((*b.sqlClause).CreateSqlSuffix())
	//	       }
	return
}
func (b *BaseSqlClause) PrepareClauseOrderBy() string {
	//	        if (!_orderByEffective || !hasOrderByClause()) {
	if b.OrderByClause == nil {
		return ""
	}
	sb := new(bytes.Buffer)
	sb.WriteString(" ")
	sb.WriteString(b.getOrderByClause())
	return sb.String()
}
func (b *BaseSqlClause) getOrderByClause() string {
	// reflectClauseLazilyIfExists();
	orderBy := b.OrderByClause
	orderByClause := ""
	//        if (hasUnionQuery()) {
	//            final Map<String, String> selectClauseRealColumnAliasMap = getSelectClauseRealColumnAliasMap();
	//            if (selectClauseRealColumnAliasMap.isEmpty()) {
	//                String msg = "The selectClauseColumnAliasMap should not be empty when union query exists.";
	//                throw new IllegalStateException(msg);
	//            }
	//            orderByClause = orderBy.getOrderByClause(selectClauseRealColumnAliasMap);
	//        } else {
	orderByClause = orderBy.getOrderByClause()
	//      }
	if len(strings.TrimSpace(orderByClause)) > 0 {
		return Ln + " " + orderByClause
	} else {
		return orderByClause
	}
}
func (b *BaseSqlClause) BuildWhereClause(sb *bytes.Buffer) {

	b.BuildWhereClauseSub(sb, false)
	return
}
func (b *BaseSqlClause) BuildWhereClauseSub(sb *bytes.Buffer, template bool) {
	l := b.WhereList
	if l == nil || l.Size() == 0 {
		//		  if (template) {
		//                sb.append(" ").append(getWhereClauseMark());
		//            }
		return
	}
	count := 0
	for _, wl := range l.data {
		var queryClause *QueryClause = (wl).(*QueryClause)
		//fmt.Printf("queryClause %v %T\n", queryClause)
		//     final String clauseElement = filterWhereClauseSimply(whereClause.toString());
		clauseElement := b.FilterWhereClauseSimply((*queryClause).ToString())
		//            if (count == 0) {
		if count == 0 {
			sb.WriteString("\n where " + clauseElement)

			//                sb.append(ln()).append(" ");
			//                sb.append("where ").append(template ? getWhereFirstConditionMark() : "").append(clauseElement);
			//            } else {
		} else {
			sb.WriteString("\n  and " + clauseElement)
		}
		//                sb.append(ln()).append("  ");
		//                sb.append(" and ").append(clauseElement);
		//            }
		//            ++count;
		count++
	}
	return
}
func (b *BaseSqlClause) FilterWhereClauseSimply(clauseElement string) string {
	//        if (_whereClauseSimpleFilterList == null || _whereClauseSimpleFilterList.isEmpty()) {
	//            return clauseElement;
	log.InternalDebug("Where Clause:" + clauseElement)
	return clauseElement
	//        }
	//        for (final Iterator<QueryClauseFilter> ite = _whereClauseSimpleFilterList.iterator(); ite.hasNext();) {
	//            final QueryClauseFilter filter = ite.next();
	//            if (filter == null) {
	//                String msg = "The list of filter should not have null: _whereClauseSimpleFilterList=" + _whereClauseSimpleFilterList;
	//                throw new IllegalStateException(msg);
	//            }
	//            clauseElement = filter.filterClauseElement(clauseElement);
	//        }
	//            return clauseElement;
}
func (b *BaseSqlClause) isJoinInParentheses() bool {
	//defalult implementation
	return false
}
func (b *BaseSqlClause) BuildFromClause(sb *bytes.Buffer) {
	sb.WriteString("\n  from ")
	tablePos := 7 // basically for in-line view indent
	tablePos = tablePos
	if b.isJoinInParentheses() {
		for i := 0; i < len(b.outerJoinMap); i++ {
			sb.WriteString("(")
			tablePos++
		}
	}
	tableSqlName := (*b.DBMeta).GetTableSqlName()
	basePointAliasName := b.BasePorintAliasName
	if b.hasBaseTableInlineWhereClause() {
		baseTableInlineWhereList := b.BaseTableInlineWhereList
		sb.WriteString(b.getInlineViewClause(tableSqlName, baseTableInlineWhereList, tablePos))
		sb.WriteString(" " + basePointAliasName)
	} else {
		sb.WriteString(tableSqlName.TableSqlName + " " + basePointAliasName)
	}
	sb.WriteString(b.createFromBaseTableHint())
	sb.WriteString(b.getLeftOuterJoinClause())
	return
}
func (b *BaseSqlClause) checkFixedConditionLazily() {
	//not implemented yet
}
func (b *BaseSqlClause) reflectInnerJoinAutoDetectLazily() {
	if b.InnerJoinLazyReflector == nil {
		return
	}
	reflectorList := b.InnerJoinLazyReflector
	reflectorList = reflectorList
	for i := 0; i < reflectorList.Size(); i++ {
		reflector := (reflectorList.Get(i)).(*InnerJoinLazyReflector)
		(*reflector).Reflect()
	}
	b.InnerJoinLazyReflector = new(List)
}
func (b *BaseSqlClause) canPagingCountLeastJoin() bool {
	return b.pagingAdjustment && b.pagingCountLeastJoin
}
func (b *BaseSqlClause) isSelectClauseTypeNonUnionCount() bool {
	return !b.hasUnionQuery() && b.selectClauseType.count
}

func (b *BaseSqlClause) hasFixedConditionOverRelationJoin() bool {
	for key := range b.outerJoinMap {
		joinInfo := b.outerJoinMap[key]
		if joinInfo.fixedConditionOverRelation {
			// because over-relation may have references of various relations
			return true
		}
	}
	return false
}
func (b *BaseSqlClause) checkCountLeastJoinAllowed() bool {
	if !b.canPagingCountLeastJoin() {
		return false
	}
	if !b.isSelectClauseTypeNonUnionCount() {
		return false
	}
	return !b.hasFixedConditionOverRelationJoin()
}
func (b *BaseSqlClause) checkStructuralPossibleInnerJoinAllowed() bool {
	if !b.structuralPossibleInnerJoinAllowed {
		return false
	}
	return !b.hasFixedConditionOverRelationJoin()
}
func (b *BaseSqlClause) canBeCountLeastJoin(joinInfo *LeftOuterJoinInfo) bool {
	return !joinInfo.IsCountableJoin()
}
func (b *BaseSqlClause) getLeftOuterJoinClause() string {
	sb := new(bytes.Buffer)
	b.checkFixedConditionLazily()
	b.reflectInnerJoinAutoDetectLazily()
	countLeastJoinAllowed := b.checkCountLeastJoinAllowed()
	structuralPossibleInnerJoinAllowed := b.checkStructuralPossibleInnerJoinAllowed()
	for _, key := range b.outerJoinList.data {
		foreignAliasName := key
		joinInfo := b.outerJoinMap[key]
		if countLeastJoinAllowed && b.canBeCountLeastJoin(joinInfo) {
			continue // means only joined countable
		}
		b.buildLeftOuterJoinClause(sb, foreignAliasName, joinInfo, structuralPossibleInnerJoinAllowed)
	}
	return sb.String()
}
func (b *BaseSqlClause) canBeInnerJoin(joinInfo *LeftOuterJoinInfo,
	structuralPossibleInnerJoinAllowed bool) bool {
	if joinInfo.innerJoin {
		return true
	}

	if structuralPossibleInnerJoinAllowed {
		return joinInfo.isStructuralPossibleInnerJoin()
	}
	return false
}
func (b *BaseSqlClause) buildLeftOuterJoinClause(sb *bytes.Buffer, foreignAliasName string, joinInfo *LeftOuterJoinInfo,
	structuralPossibleInnerJoinAllowed bool) {
	joinOnMap := joinInfo.joinOnMap
	// not implemented yet
	//        assertJoinOnMapNotEmpty(joinOnMap, foreignAliasName);

	sb.WriteString(Ln + "   ")
	joinExp := ""
	canBeInnerJoin := b.canBeInnerJoin(joinInfo, structuralPossibleInnerJoinAllowed)
	if canBeInnerJoin {
		joinExp = " inner join "
	} else {
		joinExp = " left outer join " // is main!
	}
	sb.WriteString(joinExp) // is main!
	b.buildJoinTableClause(sb, joinInfo, joinExp, canBeInnerJoin)
	sb.WriteString(" " + foreignAliasName)
	if joinInfo.hasInlineOrOnClause() || joinInfo.hasFixedCondition() {
		sb.WriteString(Ln + "     ") // only when additional conditions exist
	}
	sb.WriteString(" on ")
	b.buildJoinOnClause(sb, joinInfo, joinOnMap)
	if b.isJoinInParentheses() {
		sb.WriteString(")")
	}
}
func (b *BaseSqlClause) buildJoinOnClause(sb *bytes.Buffer,
	joinInfo *LeftOuterJoinInfo, joinOnMap map[*ColumnRealName]*ColumnRealName) {
	currentConditionCount := 0
	currentConditionCount = b.doBuildJoinOnClauseBasic(sb, joinInfo, joinOnMap, currentConditionCount)
	currentConditionCount = b.doBuildJoinOnClauseFixed(sb, joinInfo, joinOnMap, currentConditionCount)
	currentConditionCount = b.doBuildJoinOnClauseAdditional(sb, joinInfo, joinOnMap, currentConditionCount)

}
func (b *BaseSqlClause) doBuildJoinOnClauseBasic(sb *bytes.Buffer,
	joinInfo *LeftOuterJoinInfo,
	joinOnMap map[*ColumnRealName]*ColumnRealName, currentConditionCount int) int {
	for key := range joinOnMap {
		localRealName := key
		foreignRealName := joinOnMap[key]
		if currentConditionCount > 0 {
			sb.WriteString(" and ")
		}
		sb.WriteString(localRealName.ToString() + " = " +
			foreignRealName.ToString())
		currentConditionCount++
	}
	return currentConditionCount
}
func (b *BaseSqlClause) doBuildJoinOnClauseFixed(sb *bytes.Buffer,
	joinInfo *LeftOuterJoinInfo,
	joinOnMap map[*ColumnRealName]*ColumnRealName, currentConditionCount int) int {
	if joinInfo.hasFixedCondition() {
		fixedCondition := joinInfo.fixedCondition
		if b.isInlineViewOptimizedCondition(fixedCondition) {
			return currentConditionCount
		}
		sb.WriteString(Ln + "    ")
		if currentConditionCount > 0 {
			sb.WriteString(" and ")
		}

		sb.WriteString(fixedCondition)
		currentConditionCount++
	}
	return currentConditionCount
}
func (b *BaseSqlClause) doBuildJoinOnClauseAdditional(sb *bytes.Buffer,
	joinInfo *LeftOuterJoinInfo,
	joinOnMap map[*ColumnRealName]*ColumnRealName, currentConditionCount int) int {
	additionalOnClauseList := joinInfo.additionalOnClauseList
	for _, clause := range additionalOnClauseList.data {
		additionalOnClause := clause.(*QueryClause)
		sb.WriteString(Ln + "    ")
		if currentConditionCount > 0 {
			sb.WriteString(" and ")
		}
		sb.WriteString((*additionalOnClause).ToString())
		currentConditionCount++
	}
	return currentConditionCount
}
func (b *BaseSqlClause) isInlineViewOptimizedCondition(fixedCondition string) bool {
	return OPTIMIZED_MARK == fixedCondition
}
func (b *BaseSqlClause) buildJoinTableClause(sb *bytes.Buffer,
	joinInfo *LeftOuterJoinInfo, joinExp string, canBeInnerJoin bool) {
	foreignTableDbName := joinInfo.foreignTableDbName
	tablePos := 3 + len(joinExp) // basically for in-line view indent
	foreignDBMeta := DBMetaProvider_I.TableDbNameInstanceMap[foreignTableDbName]
	foreignTableSqlName := (*foreignDBMeta).GetTableSqlName()
	inlineWhereClauseList := &joinInfo.inlineWhereClauseList
	tableExp := ""
	if inlineWhereClauseList.Size() == 0 {
		tableExp = foreignTableSqlName.TableSqlName
	} else {
		tableExp = b.getInlineViewClause(foreignTableSqlName, inlineWhereClauseList, tablePos)
	}
	if joinInfo.hasFixedCondition() {
		sb.WriteString(joinInfo.resolveFixedInlineView(tableExp, canBeInnerJoin))
	} else {
		sb.WriteString(tableExp)
	}
}
func (b *BaseSqlClause) createFromBaseTableHint() string {
	//default implementation
	return ""
}
func (b *BaseSqlClause) getInlineViewBasePointAlias() string {
	return "dfinlineloc"
}
func (b *BaseSqlClause) BuildSpaceBar(size int) string {
	sb := new(bytes.Buffer)
	for i := 0; i < size; i++ {
		sb.WriteString(" ")
	}
	return sb.String()
}
func (b *BaseSqlClause) getInlineViewClause(inlineTableSqlName *TableSqlName, inlineWhereClauseList *List, tablePos int) string {
	inlineBaseAlias := b.getInlineViewBasePointAlias()
	sb := new(bytes.Buffer)
	sb.WriteString("(select * from " + inlineTableSqlName.TableSqlName + " " + inlineBaseAlias)
	baseIndent := b.BuildSpaceBar(tablePos + 1)
	sb.WriteString(Ln + baseIndent)
	sb.WriteString(" where ")
	count := 0
	for i := 0; i < inlineWhereClauseList.Size(); i++ {
		whereClause := (inlineWhereClauseList.Get(i)).(*QueryClause)
		clauseElement := b.FilterWhereClauseSimply((*whereClause).ToString())
		if count > 0 {
			sb.WriteString(Ln + baseIndent)
			sb.WriteString("   and ")
		}
		sb.WriteString(clauseElement)
		count++
	}
	sb.WriteString(")")
	return sb.String()
}
func (b *BaseSqlClause) hasBaseTableInlineWhereClause() bool {
	return b.BaseTableInlineWhereList != nil && b.BaseTableInlineWhereList.Size() > 0
}
func (b *BaseSqlClause) isSelectClauseNonUnionScalar() bool {
	return !b.hasUnionQuery() && b.isSelectClauseTypeScalar()
}
func (b *BaseSqlClause) isSelectClauseTypeScalar() bool {
	return b.selectClauseType != nil && b.selectClauseType.scalar
}
func (b *BaseSqlClause) isSelectClauseTypeCount() bool {
	return b.selectClauseType.count
}
func (b *BaseSqlClause) buildSelectClauseCount() string {
	return "select count(*)"
}
func (b *BaseSqlClause) buildSelectClauseScalar(aliasName string) string {
	if b.isSelectClauseTypeCount() {
		return b.buildSelectClauseCount()
	}
	panic("System Error Not Implemented Yet. buildSelectClauseScalar")
}
func (b *BaseSqlClause) GetSelectClause() string {
	//	        reflectClauseLazilyIfExists();
	if b.isSelectClauseNonUnionScalar() {
		return b.buildSelectClauseScalar(b.BasePorintAliasName)
	}
	//        // if it's a scalar-select, it always has union-query since here
	//        final StringBuilder sb = new StringBuilder();
	//
	//        if (_useSelectIndex) {
	//            _selectIndexMap = createSelectIndexMap(); // should be initialized before process
	//        }
	if b.UseSelectIndex {
		b.SelectIndexMap = CreateAsFlexible()
	}
	//
	//        final Integer selectIndex = processSelectClauseLocal(sb);
	//        processSelectClauseRelation(sb, selectIndex);
	//        processSelectClauseDerivedReferrer(sb);
	//
	//        return sb.toString();
	// current 3
	sb := new(bytes.Buffer)
	selectIndex := b.ProcessSelectClauseLocal(sb)
	b.processSelectClauseRelation(sb, selectIndex)
	selectIndex = selectIndex
	log.InternalDebug("Select clause =" + sb.String())
	return sb.String()
}
func (b *BaseSqlClause) processSelectClauseRelation(sb *bytes.Buffer, selectIndex int32) {
	//selectedRelationColumnMap map[string]map[string]*SelectedRelationColumn
	columnMap := b.selectedRelationColumnMap
	for key := range columnMap {
		cmap := columnMap[key]
		//            Map<String, HpSpecifiedColumn> foreginSpecifiedMap = null;
		//            if (_specifiedSelectColumnMap != null) {
		//                foreginSpecifiedMap = _specifiedSelectColumnMap.get(tableAliasName);
		//            }

		validSpecifiedForeign := false
		validSpecifiedForeign = validSpecifiedForeign
		finishedForeignIndent := false
		for _, rc := range cmap.data {
			selectColumnInfo := rc.(*SelectedRelationColumn)
			columnInfo := selectColumnInfo.columnInfo
			columnDbName := columnInfo.ColumnDbName
			columnDbName = columnDbName
			//                if (validSpecifiedForeign && !foreginSpecifiedMap.containsKey(columnDbName)) {
			//                    continue;
			//                }
			realColumnName := selectColumnInfo.BuildRealColumnSqlName()
			b.selectClauseInfo.Add(realColumnName)
			columnAliasName := selectColumnInfo.columnAliasName
			onQueryName := ""
			selectIndex++
			if b.UseSelectIndex {
				onQueryName = b.BuildSelectIndexAlias(columnInfo.ColumnSqlName, columnAliasName, selectIndex)
				b.registerSelectIndex(columnAliasName, onQueryName, selectIndex)

			} else {
				onQueryName = columnAliasName
			}
			if !finishedForeignIndent {
				sb.WriteString(Ln + "     ")
				finishedForeignIndent = true
			}
			sb.WriteString(", ")
			sb.WriteString(realColumnName + " as " + onQueryName)
			b.selectClauseRealColumnAliasMap[realColumnName] = onQueryName
			//                if (validSpecifiedForeign && foreginSpecifiedMap.containsKey(columnDbName)) {
			//                    final HpSpecifiedColumn specifiedColumn = foreginSpecifiedMap.get(columnDbName);
			//                    specifiedColumn.setOnQueryName(onQueryName); // basically for queryInsert()
			//                }
		}
	}
	//        return selectIndex;
	return
}
func (b *BaseSqlClause) ProcessSelectClauseLocal(sb *bytes.Buffer) int32 {
	//        final String basePointAliasName = getBasePointAliasName();
	basePointAliasName := b.BasePorintAliasName
	//        final DBMeta dbmeta = getDBMeta();
	dbmeta := b.DBMeta
	//        final Map<String, HpSpecifiedColumn> localSpecifiedMap;
	//        if (_specifiedSelectColumnMap != null) {
	//            localSpecifiedMap = _specifiedSelectColumnMap.get(basePointAliasName);
	//        } else {
	//            localSpecifiedMap = null;
	//        }
	//        final List<ColumnInfo> columnInfoList;
	//        final boolean validSpecifiedLocal;
	//        if (isSelectClauseTypeUniqueScalar()) {
	//            // it always has union-query because it's handled before this process
	//            if (dbmeta.hasPrimaryKey()) {
	//                columnInfoList = new ArrayList<ColumnInfo>();
	//                columnInfoList.addAll(dbmeta.getPrimaryUniqueInfo().getUniqueColumnList());
	//                if (isSelectClauseTypeSpecifiedScalar()) {
	//                    final ColumnInfo specifiedColumn = getSpecifiedColumnInfoAsOne();
	//                    if (specifiedColumn != null) {
	//                        columnInfoList.add(specifiedColumn);
	//                    }
	//                    // derivingSubQuery is handled after this process
	//                }
	//            } else {
	//                // all columns are target if no-PK and unique-scalar and union-query
	//                columnInfoList = dbmeta.getColumnInfoList();
	/////Current 2
	columnInfoList := (*dbmeta).GetColumnInfoList()
	//            }
	//            validSpecifiedLocal = false; // because specified columns are fixed here
	validSpecifiedLocal := false
	validSpecifiedLocal = validSpecifiedLocal
	//        } else {
	//            columnInfoList = dbmeta.getColumnInfoList();
	//            validSpecifiedLocal = localSpecifiedMap != null && !localSpecifiedMap.isEmpty();
	//        }
	//
	//        Integer selectIndex = 0; // because 1 origin in JDBC
	var selectIndex int32 = 0
	//        boolean needsDelimiter = false;
	needsDelimiter := false
	//
	for _, columnInfo := range columnInfoList.data {
		var cci *ColumnInfo = columnInfo.(*ColumnInfo)
		//            final String columnDbName = columnInfo.getColumnDbName();
		columnDbName := cci.ColumnDbName
		columnDbName = columnDbName
		//            final ColumnSqlName columnSqlName = columnInfo.getColumnSqlName();
		columnSqlName := cci.ColumnSqlName
		log.InternalDebug("columnSqlName := " + columnSqlName.ColumnSqlName)
		//
		//            if (validSpecifiedLocal && !localSpecifiedMap.containsKey(columnDbName)) {
		//                // a case for scalar-select has been already resolved here
		//                continue;
		//            }
		//
		//            if (needsDelimiter) {
		//                sb.append(", ");
		if needsDelimiter {
			sb.WriteString(", ")
		} else {
			//            } else {
			//                sb.append("select");
			//                appendSelectHint(sb);
			//                sb.append(" ");
			//                needsDelimiter = true;
			//            }
			sb.WriteString("select/*$pmb.BaseConditionBean.SelectHint*/ ")
			needsDelimiter = true
		}
		realColumnName := basePointAliasName + "." + columnSqlName.ColumnSqlName
		b.selectClauseInfo.Add(realColumnName)
		var onQueryName string
		selectIndex++
		if b.UseSelectIndex {
			onQueryName = b.BuildSelectIndexAlias(columnSqlName, "", selectIndex)
			b.registerSelectIndex(columnDbName, onQueryName, selectIndex)
		} else {
			onQueryName = columnSqlName.ColumnSqlName
		}
		//            sb.append(decryptSelectColumnIfNeeds(columnInfo, realColumnName)).append(" as ").append(onQueryName);
		sb.WriteString(realColumnName + " as " + onQueryName)
		b.getSelectClauseRealColumnAliasMap()[realColumnName] = onQueryName
		//            if (validSpecifiedLocal && localSpecifiedMap.containsKey(columnDbName)) {
		//                final HpSpecifiedColumn specifiedColumn = localSpecifiedMap.get(columnDbName);
		//                specifiedColumn.setOnQueryName(onQueryName); // basically for queryInsert()
		//            }
		//        }
	}
	return selectIndex
}
func (b *BaseSqlClause) registerSelectIndex(keyName string, onQueryName string, selectIndex int32) {
	if b.SelectIndexMap == nil {
		b.SelectIndexMap = CreateAsFlexible()
	}
	b.SelectIndexMap.Put(keyName, selectIndex)
	if b.SelectIndexReverseMap == nil {
		b.SelectIndexReverseMap = CreateAsFlexible()
	}
	b.SelectIndexReverseMap.Put(onQueryName, keyName)
}
func (b *BaseSqlClause) BuildSelectIndexAlias(columnSqlName *ColumnSqlName, aliasName string, selectIndex int32) string {
	if columnSqlName.IrregularChar {
		return "c" + strconv.Itoa(int(selectIndex)) // use index only for safety
	}
	// regular case only here
	//        final String baseName;
	var baseName string
	//        if (aliasName != null) { // relation column
	if aliasName != "" {
		baseName = aliasName
	} else {
		//            baseName = aliasName;
		//        } else { // local column
		//            baseName = sqlName.toString();
		baseName = columnSqlName.ColumnSqlName
	}
	//        }
	//        final int aliasNameLimitSize = getAliasNameLimitSize();
	//        if (baseName.length() > aliasNameLimitSize) {
	//            final int aliasNameBaseSize = aliasNameLimitSize - 10;
	//            return Srl.substring(baseName, 0, aliasNameBaseSize) + "_c" + selectIndex;
	//        } else {
	//            return baseName;
	//        }
	return baseName
}
func (b *BaseSqlClause) GetOrderByClause() *OrderByClause {
	oc := b.OrderByClause
	if oc != nil {
		return oc
	}
	oc = new(OrderByClause)
	oc.OrderByList = new(List)
	b.OrderByClause = oc
	return oc
}

func (b *BaseSqlClause) setTableDbName(tn string) {
	b.TableDbName = tn
}
func (b *BaseSqlClause) GetTableDbName() string {
	return b.TableDbName
}
func (b *BaseSqlClause) setSqlSetting(dc *DBCurrent) {
	//clone DBCurrent
	var dcx DBCurrent = *dc
	b.SqlSetting = &dcx
	//		(dcx).PagingCountLater=true
	//	fmt.Printf("ok??? %v %v\n",(*dc).PagingCountLater,(dcx).PagingCountLater)
	if dcx.InnerJoinAutoDetect {
		b.whereUsedInnerJoinAllowed = true
		b.structuralPossibleInnerJoinAllowed = true
	}
	if dcx.EmptyStringQueryAllowed {
		b.emptyStringQueryAllowed = true
	}
	if dcx.DisableSelectIndex {
		b.disableSelectIndex = true
	}
}
func (b *BaseSqlClause) GetSqlSetting() *DBCurrent {
	return b.SqlSetting
}
func (b *BaseSqlClause) GetBasePorintAliasName() string {
	return b.BasePorintAliasName
}

func (b *BaseSqlClause) IsOrScopeQueryEffective() bool {
	return b.OrScopeQueryEffective
}

func (b *BaseSqlClause) IsOrScopeQueryAndPartEffective() bool {
	return b.OrScopeQueryAndPartEffective
}
func (b *BaseSqlClause) RegisterWhereClause(crn *ColumnRealName, key *ConditionKey, cvalue *ConditionValue, co *ConditionOption, usedAliasName string) {
	//Assert 省略
	clauseList := b.GetWhereClauseList4Register()
	b.DoRegisterWhereClause(clauseList, crn, key, cvalue, co, false, false)
	cmap := cvalue.Fixed
	log.InternalDebug(fmt.Sprintf(" Sql Clause Cvalue %v cmap %v usedAliasName %s\n", cvalue, cmap, usedAliasName))
	b.ReflectWhereUsedToJoin(usedAliasName)
	if !ConditionKey_IsNullaleConditionKey(key) {
		b.RegisterInnerJoinLazyReflector(usedAliasName)
	}
}
func (b *BaseSqlClause) IsOutOfWhereUsedInnerJoin() bool {
	return !b.whereUsedInnerJoinAllowed || b.orScopeQueryEffective
}
func (b *BaseSqlClause) RegisterInnerJoinLazyReflector(usedAliasName string) {
	if b.IsOutOfWhereUsedInnerJoin() {
		return
	}
	usedAliasInfo := new(QueryUsedAliasInfo)
	usedAliasInfo.usedAliasName = usedAliasName
	b.RegisterInnerJoinLazyReflectorSub(usedAliasInfo)
}
func (b *BaseSqlClause) RegisterInnerJoinLazyReflectorSub(usedAliasInfo *QueryUsedAliasInfo) {
	if b.IsOutOfWhereUsedInnerJoin() {
		return
	}
	reflectorList := b.getInnerJoinLazyReflectorList()
	reflectorList.Add(b.createInnerJoinLazyReflector(usedAliasInfo))
}
func (b *BaseSqlClause) createInnerJoinLazyReflector(usedAliasInfo *QueryUsedAliasInfo) *InnerJoinLazyReflector {
	usedAliasName := usedAliasInfo.usedAliasName
	innerJoinLazyReflectorBase := new(InnerJoinLazyReflectorBase)
	innerJoinLazyReflectorBase.noWaySpeaker = usedAliasInfo.innerJoinNoWaySpeaker
	innerJoinLazyReflectorBase.usedAliasName = usedAliasName
	innerJoinLazyReflectorBase.BaseSqlClause = b
	var ref InnerJoinLazyReflector = innerJoinLazyReflectorBase
	return &ref
}
func (b *BaseSqlClause) DoChangeToInnerJoin(foreignAliasName string, autoDetect bool) {
	outerJoinMap := b.outerJoinMap
	joinInfo := outerJoinMap[foreignAliasName]
	if joinInfo == nil {
		panic("The foreignAliasName was not found:" + " " + foreignAliasName)
	}
	joinInfo.innerJoin = true
	b.reflectUnderInnerJoinToJoin(joinInfo, autoDetect)
}
func (b *BaseSqlClause) reflectUnderInnerJoinToJoin(foreignJoinInfo *LeftOuterJoinInfo, autoDetect bool) {
	currentJoinInfo := foreignJoinInfo.localJoinInfo
	for true {
		if currentJoinInfo == nil { // means base point
			break
		}
		// all join-info are overridden because of complex logic
		if autoDetect {
			currentJoinInfo.innerJoin = true // be inner-join as we can if auto-detect
		} else {
			currentJoinInfo.underInnerJoin = true // manual is pinpoint setting
		}
		currentJoinInfo = currentJoinInfo.localJoinInfo
	}
}
func (b *BaseSqlClause) getInnerJoinLazyReflectorList() *List {
	if b.InnerJoinLazyReflector == nil {
		b.InnerJoinLazyReflector = new(List)
	}
	return b.InnerJoinLazyReflector
}
func (b *BaseSqlClause) ReflectWhereUsedToJoin(usedAliasName string) {
	currentJoinInfo := b.outerJoinMap[usedAliasName]
	for true {
		if currentJoinInfo == nil { // means base point
			break
		}
		if currentJoinInfo.whereUsedJoin { // means already traced
			break
		}
		currentJoinInfo.whereUsedJoin = true
		currentJoinInfo = currentJoinInfo.localJoinInfo // trace back toward base point
	}
}
func (b *BaseSqlClause) DoRegisterWhereClause(clauseList *List, crn *ColumnRealName, key *ConditionKey, cvalue *ConditionValue, co *ConditionOption, inline bool, onClause bool) {
	(*key).AddWhereClause(key, b.XcreateQueryModeProvider(), clauseList, crn, cvalue, co)
	//MarkOrScopeQueryAndPart(clauseList); 未実装
}
func (b *BaseSqlClause) XcreateQueryModeProvider() *QueryModeProvider {
	qm := new(QueryModeProvider)
	qm.SqlClause = b.sqlClause
	qm.IsInline = b.Inline
	qm.IsOnClause = b.OnClause
	return qm
}
func (b *BaseSqlClause) GetWhereClauseList4Register() *List {
	log.InternalDebug(fmt.Sprintf("b.OrScopeQueryEffective %v \n ", b.OrScopeQueryEffective))
	if b.OrScopeQueryEffective {
		return &b.currentTmpOrScopeQueryInfo.tmpOrWhereList
	}
	return b.GetWhereList()
}

func (b *BaseSqlClause) GetWhereList() *List {
	if b.WhereList != nil {
		return b.WhereList
	}
	b.WhereList = new(List)
	return b.WhereList
}
func (b *BaseSqlClause) ResolveRelationNo(localTableName string, foreignPropertyName string) int {
	dbmeta := DBMetaProvider_I.TableDbNameInstanceMap[localTableName]
	foreignInfo := (*dbmeta).FindForeignInfo(foreignPropertyName)
	return foreignInfo.RelationNo
}
func (b *BaseSqlClause) ResolveJoinAliasName(relationPath string) string {
	if b.subQueryLevel > 0 {
		return "sub" + strconv.Itoa(b.subQueryLevel) + "rel" + relationPath
	}
	return "df" + "rel" + relationPath
}

func CreateSqlClause(cb ConditionBean, dc *DBCurrent) *SqlClause {
	return CreateSqlClauseSub(cb.GetBaseConditionBean().AsTableDbName(), dc)
}

func CreateSqlClauseSub(tableDbName string, dc *DBCurrent) *SqlClause {
	code := (*dc.DBDef).Code()
	var sql SqlClause
	if code == "postgresql" {
		sqlp := new(SqlClausePostgres)
		sqlp.BasePorintAliasName = "dfloc"
		sql = sqlp
		sqlp.sqlClause = &sql
	}
	if code == "mysql" {
		sqlp := new(SqlClauseMySql)
		sqlp.BasePorintAliasName = "dfloc"
		sql = sqlp
		sqlp.sqlClause = &sql
	}
	if code == "sqlserver" {
		sqlp := new(SqlClauseSqlServer)
		sqlp.BasePorintAliasName = "dfloc"
		sql = sqlp
		sqlp.sqlClause = &sql
	}
	sql.setTableDbName(tableDbName)
	sql.setSqlSetting(dc)
	return &sql
}

type OrderByClause struct {
	OrderByList *List
}

func (o *OrderByClause) AddElement(orderByElement *OrderByElement) {
	o.OrderByList.Add(orderByElement)
}
func (o *OrderByClause) getOrderByClause() string {
	if o.OrderByList == nil || o.OrderByList.Size() == 0 {
		return ""
	}
	sb := new(bytes.Buffer)
	sb.WriteString("order by ")
	delimiter := ", "
	for i, ele := range o.OrderByList.data {
		element := ele.(*OrderByElement)
		if i > 0 {
			sb.WriteString(delimiter)
		}
		//            if (selectClauseRealColumnAliasMap != null) {
		//                sb.append(element.getElementClause(selectClauseRealColumnAliasMap));
		//            } else {
		sb.WriteString(element.getElementClause())
		//            }
	}
	//        sb.delete(0, delimiter.length()).insert(0, "order by ");
	return sb.String()

}

type OrderByElement struct {
	AliasName  string
	ColumnName string
	AscDesc    string
}

func (o *OrderByElement) getElementClause() string {
	if o.AscDesc == "" {
		panic("The attribute[ascDesc] should not be null.")
	}
	sb := new(bytes.Buffer)
	columnFullName := o.getColumnFullName()
	//        if (_manualOrderOption != null && _manualOrderOption.hasManualOrder()) {
	//            setupManualOrderClause(sb, columnFullName, null);
	//            return sb.toString();
	//        } else {
	sb.WriteString(columnFullName + " " + o.AscDesc)
	clause := sb.String()
	//            if (_orderByNullsSetupper != null) {
	//                return _orderByNullsSetupper.setup(columnFullName, clause, _nullsFirst);
	//            } else {
	return clause
	//            }
	//        }
}
func (o *OrderByElement) getColumnFullName() string {
	sb := new(bytes.Buffer)
	if o.AliasName != "" {
		sb.WriteString(o.AliasName + ".")
	}
	if o.ColumnName == "" {
		panic("The attribute[columnName] should not be null.")
	}
	//        final String derivedMappingAliasPrefix = DerivedMappable.MAPPING_ALIAS_PREFIX;
	//        if (_derivedOrderBy && _columnName.startsWith(derivedMappingAliasPrefix)) {
	//            sb.append(Srl.substringFirstRear(_columnName, derivedMappingAliasPrefix));
	//        } else {
	sb.WriteString(o.ColumnName)
	//        }
	return sb.String()
}

type OrScopeQueryInfo struct {
	tmpOrWhereList                    List
	tmpOrBaseTableInlineWhereList     List
	tmpOrAdditionalOnClauseListMap    map[string]*List
	tmpOrOuterJoinInlineClauseListMap map[string]*List
	parentInfo                        *OrScopeQueryInfo
	childInfoList                     List
}

func (o *OrScopeQueryInfo) addChildInfo(info *OrScopeQueryInfo) {
	info.parentInfo = o
	o.childInfoList.Add(info)
}

type OrScopeQueryClauseListProvider interface {
	provide(tmpOrScopeQueryInfo *OrScopeQueryInfo) interface{}
	provideAlias(tmpOrScopeQueryInfo *OrScopeQueryInfo, aliasName string) interface{}
}
type BaseOrScopeQueryClauseListProvider struct {
}

func (p *BaseOrScopeQueryClauseListProvider) provide(
	tmpOrScopeQueryInfo *OrScopeQueryInfo) interface{} {
	//Dummy implementation
	return nil
}
func (p *BaseOrScopeQueryClauseListProvider) provideAlias(
	tmpOrScopeQueryInfo *OrScopeQueryInfo, aliasName string) interface{} {
	//Dummy implementation
	return nil
}

type OrScopeQueryClauseListProviderWhereList struct {
	BaseOrScopeQueryClauseListProvider
}

func (o *OrScopeQueryClauseListProviderWhereList) provide(tmpOrScopeQueryInfo *OrScopeQueryInfo) interface{} {
	log.InternalDebug(fmt.Sprintf("tmpOrWhereList Size %d \n", tmpOrScopeQueryInfo.tmpOrWhereList.Size()))
	return &tmpOrScopeQueryInfo.tmpOrWhereList
}

type OrScopeQueryClauseListProviderBaseTableInlineWhereList struct {
	BaseOrScopeQueryClauseListProvider
}

func (o *OrScopeQueryClauseListProviderBaseTableInlineWhereList) provide(
	tmpOrScopeQueryInfo *OrScopeQueryInfo) interface{} {
	return &tmpOrScopeQueryInfo.tmpOrBaseTableInlineWhereList
}

type OrScopeQueryClauseListProviderAdditionalOnClauseList struct {
	BaseOrScopeQueryClauseListProvider
}

func (o *OrScopeQueryClauseListProviderAdditionalOnClauseList) provideAlias(
	tmpOrScopeQueryInfo *OrScopeQueryInfo, aliasName string) interface{} {
	return tmpOrScopeQueryInfo.tmpOrAdditionalOnClauseListMap[aliasName]
}

type OrScopeQueryClauseListProviderOuterJoinInlineClauseList struct {
	BaseOrScopeQueryClauseListProvider
}

func (o *OrScopeQueryClauseListProviderOuterJoinInlineClauseList) provideAlias(
	tmpOrScopeQueryInfo *OrScopeQueryInfo, aliasName string) interface{} {
	return tmpOrScopeQueryInfo.tmpOrOuterJoinInlineClauseListMap[aliasName]
}

type OrScopeQueryReflector struct {
	whereList                *List
	baseTableInlineWhereList *List
	outerJoinMap             map[string]*LeftOuterJoinInfo
	setupper                 OrScopeQuerySetupper
}

func (o *OrScopeQueryReflector) reflectTmpOrClauseToRealObject(localInfo *OrScopeQueryInfo) {

	// to Normal Query (where clause)
	wl := new(OrScopeQueryClauseListProviderWhereList)
	var oc OrScopeQueryClauseListProvider = wl
	groupList := (o.setupTmpOrListList(localInfo, &oc)).(*List)
	log.InternalDebug(fmt.Sprintf("groupList %v %d \n", groupList, groupList.Size()))
	o.setupOrScopeQuery(groupList, o.whereList, true)
	// to InlineView for baseTable
	wl2 := new(OrScopeQueryClauseListProviderBaseTableInlineWhereList)
	var oc2 OrScopeQueryClauseListProvider = wl2
	groupList1 := (o.setupTmpOrListList(localInfo, &oc2)).(*List)

	o.setupOrScopeQuery(groupList1, o.baseTableInlineWhereList, false)
	// to OnClause
	for aliasName := range o.outerJoinMap {
		joinInfo := o.outerJoinMap[aliasName]

		wl3 := new(OrScopeQueryClauseListProviderAdditionalOnClauseList)
		var oc3 OrScopeQueryClauseListProvider = wl3
		groupList2 := new(List)
		temp2 := (o.setupTmpOrListListAlias(localInfo, aliasName, &oc3)).(*List)
		for _, temp := range temp2.data {
			//fmt.Printf("groupList2 %v %T\n", temp, temp)
			var tlist *OrScopeQueryClauseGroup = temp.(*OrScopeQueryClauseGroup)
			if tlist.orClauseList != nil {
				groupList2.Add(temp)
			}
		}

		o.setupOrScopeQuery(groupList2, &joinInfo.additionalOnClauseList, false)
	}
	// to InlineView for relation
	for aliasName := range o.outerJoinMap {
		joinInfo := o.outerJoinMap[aliasName]

		wl3 := new(OrScopeQueryClauseListProviderOuterJoinInlineClauseList)
		var oc3 OrScopeQueryClauseListProvider = wl3
		groupList2 := new(List)
		temp2 := (o.setupTmpOrListListAlias(localInfo, aliasName, &oc3)).(*List)
		for _, temp := range temp2.data {
			//fmt.Printf("groupList2 %v %T\n", temp, temp)
			var tlist *OrScopeQueryClauseGroup = temp.(*OrScopeQueryClauseGroup)
			if tlist.orClauseList != nil {
				groupList2.Add(temp)
			}
		}
		o.setupOrScopeQuery(groupList2, &joinInfo.additionalOnClauseList, false)
	}
}
func (o *OrScopeQueryReflector) setupOrScopeQuery(clauseGroupList *List, realList *List, line bool) {
	o.setupper.setupOrScopeQuery(clauseGroupList, realList, line)
}
func (o *OrScopeQueryReflector) setupTmpOrListList(parentInfo *OrScopeQueryInfo, provider *OrScopeQueryClauseListProvider) interface{} {
	resultList := new(List)
	groupInfo := new(OrScopeQueryClauseGroup)
	groupInfo.orClauseList = ((*provider).provide(parentInfo)).(*List)
	resultList.Add(groupInfo)
	if parentInfo.childInfoList.Size() > 0 {
		for _, childInfo := range parentInfo.childInfoList.data {
			var ci *OrScopeQueryInfo = childInfo.(*OrScopeQueryInfo)
			list := (o.setupTmpOrListList(ci, provider)).(*List)
			for _, ele := range list.data {
				resultList.Add(ele)
			}
		}
	}
	return resultList
}
func (o *OrScopeQueryReflector) setupTmpOrListListAlias(
	parentInfo *OrScopeQueryInfo, aliasName string,
	provider *OrScopeQueryClauseListProvider) interface{} {
	resultList := new(List)
	groupInfo := new(OrScopeQueryClauseGroup)
	groupInfo.orClauseList = (*provider).provideAlias(parentInfo, aliasName).(*List)
	resultList.Add(groupInfo)
	if parentInfo.childInfoList.Size() > 0 {
		for _, childInfo := range parentInfo.childInfoList.data {
			var ci *OrScopeQueryInfo = childInfo.(*OrScopeQueryInfo)
			list := (o.setupTmpOrListListAlias(ci, aliasName, provider)).(*List)
			for _, ele := range list.data {
				resultList.Add(ele)
			}
		}
	}
	return resultList
}

type LeftOuterJoinInfo struct {
	foreignAliasName           string
	foreignTableDbName         string
	localAliasName             string
	localTableDbName           string
	joinOnMap                  map[*ColumnRealName]*ColumnRealName
	localJoinInfo              *LeftOuterJoinInfo
	pureFK                     bool
	notNullFKColumn            bool
	inlineWhereClauseList      List
	additionalOnClauseList     List
	fixedCondition             string
	fixedConditionResolver     *FixedConditionResolver
	fixedConditionOverRelation bool
	innerJoin                  bool
	underInnerJoin             bool
	whereUsedJoin              bool
}

func (p *LeftOuterJoinInfo) resolveFixedInlineView(foreignTableSqlName string, canBeInnerJoin bool) string {
	if p.hasFixedCondition() && p.fixedConditionResolver != nil {
		return (*p.fixedConditionResolver).ResolveFixedInlineView(foreignTableSqlName, canBeInnerJoin)
	}
	return foreignTableSqlName
}
func (p *LeftOuterJoinInfo) hasFixedCondition() bool {
	return len(p.fixedCondition) > 0
}
func (p *LeftOuterJoinInfo) hasInlineOrOnClause() bool {
	return p.inlineWhereClauseList.Size() > 0 || p.additionalOnClauseList.Size() > 0
}
func (p *LeftOuterJoinInfo) isPureStructuralPossibleInnerJoin() bool {
	return !p.hasInlineOrOnClause() && p.pureFK && p.notNullFKColumn
}
func (p *LeftOuterJoinInfo) isStructuralPossibleInnerJoin() bool {
	if !p.isPureStructuralPossibleInnerJoin() {
		return false
	}
	// pure structural-possible inner-join here
	// and check all relations from base point are inner-join or not
	// (separated structural-possible should not be inner-join)
	current := p.localJoinInfo
	for true {
		if current == nil { // means first level (not nested) join
			break
		}
		// means nested join here
		// (e.g. SERVICE_RANK if MEMBER is base point)
		if !current.isTraceStructuralPossibleInnerJoin() {
			return false
		}
		current = current.localJoinInfo
	}
	return true
}
func (p *LeftOuterJoinInfo) isTraceStructuralPossibleInnerJoin() bool {
	return !p.hasInlineOrOnClause() && p.pureFK && p.notNullFKColumn
}
func (p *LeftOuterJoinInfo) ResolveFixedCondition() {
	if p.fixedCondition > "" && p.fixedConditionResolver != nil {
		// over-relation should be determined before resolving
		p.fixedConditionOverRelation = (*p.fixedConditionResolver).HasOverRelation(p.fixedCondition)
		p.fixedCondition = (*p.fixedConditionResolver).ResolveVariable(p.fixedCondition, false)
	}
}
func (p *LeftOuterJoinInfo) IsCountableJoin() bool {
	return p.innerJoin || p.underInnerJoin || p.whereUsedJoin
}

type OrScopeQuerySetupper struct {
}

func (o *OrScopeQuerySetupper) setupOrScopeQuery(clauseGroupList *List, realList *List, line bool) {
	if clauseGroupList == nil || clauseGroupList.Size() == 0 {
		return
	}
	or := " or "
	or = or
	and := " and "
	and = and
	var lnIndentOr string
	if line {
		lnIndentOr = Ln + "    "
	} else {
		lnIndentOr = ""
	}
	lnIndentOr = lnIndentOr
	lnIndentAnd := "" // no line separator either way
	lnIndentAnd = lnIndentAnd
	var lnIndentAndLn string
	if line {
		lnIndentAndLn = Ln + "      "
	} else {
		lnIndentAndLn = ""
	}
	lnIndentAndLn = lnIndentAndLn
	sb := new(bytes.Buffer)
	sb = sb
	exists := false
	exists = exists
	validCount := 0
	validCount = validCount
	groupListIndex := 0
	groupListIndex = groupListIndex
	for _, cg := range clauseGroupList.data {
		var clauseGroup *OrScopeQueryClauseGroup = cg.(*OrScopeQueryClauseGroup)
		orClauseList := clauseGroup.orClauseList
		log.InternalDebug(fmt.Sprintf("orClauseList %v %d \n", orClauseList, orClauseList.Size()))
		if orClauseList == nil || orClauseList.Size() == 0 {
			continue // not increment index
		}
		listIndex := 0
		preAndPartIdentity := -1
		for _, oc := range orClauseList.data {
			var clauseElement *QueryClause = oc.(*QueryClause)
			orClause := (*clauseElement).ToString()
			var andPartClause *QueryClause
			log.InternalDebug("clauseElement type" + GetType(clauseElement))
			//next line may have problem.
			if GetType(clauseElement) == "*df.OrScopeQueryAndPartQueryClause" {
				andPartClause = clauseElement
			}
			var beginAndPart bool
			var secondAndPart bool
			if andPartClause != nil {
				identity := (*andPartClause).getIdentity()
				if preAndPartIdentity == -1 { // first of and-part
					preAndPartIdentity = identity
					beginAndPart = true
					secondAndPart = false
				} else if preAndPartIdentity == identity { // same and-part
					beginAndPart = false
					secondAndPart = true
				} else { // other and-part
					sb.WriteString(")") // closing previous and-part
					preAndPartIdentity = identity
					beginAndPart = true
					secondAndPart = false
				}
			} else {
				if preAndPartIdentity != -1 {
					sb.WriteString(")") // closing and-part
					preAndPartIdentity = -1
				}
				beginAndPart = false
				secondAndPart = false
			}
			if groupListIndex == 0 { // first list
				if listIndex == 0 {
					sb.WriteString("(")
				} else {
					containsLn := strings.Contains(orClause, Ln)
					if secondAndPart {
						if containsLn {
							sb.WriteString(lnIndentAndLn)
						} else {
							sb.WriteString(lnIndentAnd)
						}
					} else {
						sb.WriteString(lnIndentOr)
					}
					if secondAndPart {
						sb.WriteString(and)
					} else {
						sb.WriteString(or)
					}

				}
			} else { // second or more list
				if listIndex == 0 {
					// always 'or' here
					sb.WriteString(lnIndentOr)
					sb.WriteString(or)
					sb.WriteString("(")
				} else {
					containsLn := strings.Contains(orClause, Ln)
					if secondAndPart {
						if containsLn {
							sb.WriteString(lnIndentAndLn)
						} else {
							sb.WriteString(lnIndentAnd)
						}
					} else {
						sb.WriteString(lnIndentOr)
					}
					if secondAndPart {
						sb.WriteString(and)
					} else {
						sb.WriteString(or)
					}
				}
			}
			if beginAndPart {
				sb.WriteString("(")
			}

			sb.WriteString(orClause)
			validCount++
			if !exists {
				exists = true
			}
			listIndex++
		}
		if preAndPartIdentity != -1 {
			sb.WriteString(")") // closing and-part
			preAndPartIdentity = -1
		}
		if groupListIndex > 0 { // second or more list
			sb.WriteString(")") // closing or-scope
		}
		groupListIndex++
	}
	if exists {
		if line && validCount > 1 {
			sb.WriteString(Ln + "       )")
		} else {
			sb.WriteString(")")
		}
		sqc := new(StringQueryClause)
		sqc.Clause = sb.String()
		var qc QueryClause = sqc
		realList.Add(&qc)
	}
}

type OrScopeQueryClauseGroup struct {
	orClauseList *List
}

type OrScopeQueryAndPartQueryClause struct {
	BaseQueryClause
	clause   *QueryClause
	identity int
}

func (o *OrScopeQueryAndPartQueryClause) getIdentity() int {
	return o.identity
}

type QueryUsedAliasInfo struct {
	usedAliasName         string
	innerJoinNoWaySpeaker *InnerJoinNoWaySpeaker
}
type InnerJoinNoWaySpeaker struct {
}
type InnerJoinLazyReflector interface {
	Reflect()
}
type InnerJoinLazyReflectorBase struct {
	noWaySpeaker  *InnerJoinNoWaySpeaker
	usedAliasName string
	BaseSqlClause *BaseSqlClause
}

func (p *InnerJoinLazyReflectorBase) Reflect() {
	if p.BaseSqlClause.outerJoinMap[p.usedAliasName] != nil {
		p.BaseSqlClause.DoChangeToInnerJoin(p.usedAliasName, true)
	}
}

type SelectClauseType struct {
	Value           int
	count           bool
	scalar          bool
	uniqueScalar    bool
	specifiedScalar bool
}

func Create_SelectClauseType(value int) *SelectClauseType {
	var s *SelectClauseType = new(SelectClauseType)
	s.Value = value
	switch value {
	case SelectClauseType_COLUMNS:
	case SelectClauseType_UNIQUE_COUNT:
		s.count = true
		s.scalar = true
	case SelectClauseType_PLAIN_COUNT:
		s.count = true
		s.scalar = true
	case SelectClauseType_COUNT_DISTINCT:
		s.scalar = true
		s.uniqueScalar = true
		s.specifiedScalar = true
	case SelectClauseType_MAX:
		s.scalar = true
		s.uniqueScalar = true
		s.specifiedScalar = true
	case SelectClauseType_MIN:
		s.scalar = true
		s.uniqueScalar = true
		s.specifiedScalar = true
	case SelectClauseType_SUM:
		s.scalar = true
		s.uniqueScalar = true
		s.specifiedScalar = true
	case SelectClauseType_AVE:
		s.scalar = true
		s.uniqueScalar = true
		s.specifiedScalar = true
	}
	return s
}

type SelectedRelationColumn struct {
	tableAliasName  string
	columnInfo      *ColumnInfo
	columnAliasName string
}

func (p *SelectedRelationColumn) BuildRealColumnSqlName() string {
	columnSqlName := p.columnInfo.ColumnSqlName
	if p.tableAliasName != "" {
		return p.tableAliasName + "." + columnSqlName.ColumnSqlName
	} else {
		return columnSqlName.ColumnSqlName
	}
}

type SubQueryIndentProcessor struct {
}

func (p *SubQueryIndentProcessor) processSubQueryIndent(
	sql string, preIndent string, originalSql string) string {
	// /= = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = = =
	// it's a super core logic for formatting SQL generated by ConditionBean
	// = = = = = = = = = =/
	beginMarkPrefix := SQIP_BEGIN_MARK_PREFIX
	if !strings.Contains(sql, beginMarkPrefix) {
		return sql
	}
	panic("processSubQueryIndent not implemented yet")
	//	lines := strings.Split(sql, Ln)
	//	endMarkPrefix := SQIP_END_MARK_PREFIX
	//	identityTerminal := SQIP_IDENTITY_TERMINAL
	//	terminalLength := len(identityTerminal)
	//	mainSb := new(bytes.Buffer)
	//	var subSb *bytes.Buffer = nil
	//	throughBegin := false
	//	throughBeginFirst := false
	//	subQueryIdentity := ""
	//	indent := ""
	//	for _, line := range lines {
	//            if (!throughBegin) {
	//                if (line.contains(beginMarkPrefix)) { // begin line
	//                    throughBegin = true;
	//                    subSb = new StringBuilder();
	//                    final int markIndex = line.indexOf(beginMarkPrefix);
	//                    final int terminalIndex = line.indexOf(identityTerminal);
	//                    if (terminalIndex < 0) {
	//                        String msg = "Identity terminal was not found at the begin line: [" + line + "]";
	//                        throw new SubQueryIndentFailureException(msg);
	//                    }
	//                    final String clause = line.substring(0, markIndex) + line.substring(terminalIndex + terminalLength);
	//                    subQueryIdentity = line.substring(markIndex + beginMarkPrefix.length(), terminalIndex);
	//                    subSb.append(clause);
	//                    indent = buildSpaceBar(markIndex - preIndent.length());
	//                } else { // normal line
	//                    if (needsLineConnection(mainSb)) {
	//                        mainSb.append(ln());
	//                    }
	//                    mainSb.append(line).append(ln());
	//                }
	//            } else {
	//                // - - - - - - - -
	//                // In begin to end
	//                // - - - - - - - -
	//                if (line.contains(endMarkPrefix + subQueryIdentity)) { // end line
	//                    final int markIndex = line.indexOf(endMarkPrefix);
	//                    final int terminalIndex = line.indexOf(identityTerminal);
	//                    if (terminalIndex < 0) {
	//                        String msg = "Identity terminal was not found at the begin line: [" + line + "]";
	//                        throw new SubQueryIndentFailureException(msg);
	//                    }
	//                    final String clause = line.substring(0, markIndex);
	//                    // e.g. " + 1" of ColumnQuery calculation for right column
	//                    final String preRemainder = line.substring(terminalIndex + terminalLength);
	//                    subSb.append(clause);
	//                    final String subQuerySql = subSb.toString();
	//                    final String nestedPreIndent = preIndent + indent;
	//                    final String currentSql = processSubQueryIndent(subQuerySql, nestedPreIndent, originalSql);
	//                    if (needsLineConnection(mainSb)) {
	//                        mainSb.append(ln());
	//                    }
	//                    mainSb.append(currentSql);
	//                    if (Srl.is_NotNull_and_NotTrimmedEmpty(preRemainder)) {
	//                        mainSb.append(preRemainder);
	//                    }
	//                    throughBegin = false;
	//                    throughBeginFirst = false;
	//                } else { // scope line
	//                    if (!throughBeginFirst) {
	//                        subSb.append(line.trim()).append(ln());
	//                        throughBeginFirst = true;
	//                    } else {
	//                        subSb.append(indent).append(line).append(ln());
	//                    }
	//                }
	//            }
	//}
	//        final String filteredSql = Srl.rtrim(mainSb.toString()); // removed latest line separator
	//        if (throughBegin) {
	//            throwSubQueryNotFoundEndMarkException(subQueryIdentity, sql, filteredSql, originalSql);
	//        }
	//        if (filteredSql.contains(beginMarkPrefix)) {
	//            throwSubQueryAnyBeginMarkNotHandledException(subQueryIdentity, sql, filteredSql, originalSql);
	//        }
	//        return filteredSql;
	panic("")
	return "processSubQueryInden"
}
