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

type SqlClause interface {
	setTableDbName(tn string)
	GetTableDbName() string
	setSqlSetting(dc *DBCurrent)
	GetSqlSetting() *DBCurrent
	GetBasePorintAliasName() string
	GetOrderByClause() *OrderByClause
	IsOrScopeQueryEffective() bool
	IsOrScopeQueryAndPartEffective() bool
	RegisterWhereClause(crn *ColumnRealName, key *ConditionKey, cvalue *ConditionValue, co *ConditionOption, usedAliasName string)
	GetWhereList() *List
	SetDBMeta(dm *DBMeta)
	SetUseSelectIndex(si bool)
	GetClause() string
	RegisterBaseTableInlineWhereClause(columnSqlName *ColumnSqlName, key *ConditionKey, cvalue *ConditionValue, option *ConditionOption)
	MakeOrScopeQueryEffective()
	CloseOrScopeQuery()
	BeginOrScopeQueryAndPart()
	EndOrScopeQueryAndPart()
	AllowEmptyStringQuery()
	IsAllowEmptyStringQuery() bool
}

type BaseSqlClause struct {
	TableDbName string
	//SqlSetting is copy for DBCurrent can chage setting
	SqlSetting                   *DBCurrent
	BasePorintAliasName          string
	OrderByClause                *OrderByClause
	OrScopeQueryEffective        bool
	WhereList                    *List
	Inline                       bool
	OnClause                     bool
	DBMeta                       *DBMeta
	UseSelectIndex               bool
	SelectIndexMap               *StringKeyMap
	BaseTableInlineWhereList     *List
	orScopeQueryAndPartEffective bool
	currentTmpOrScopeQueryInfo   *OrScopeQueryInfo
	orScopeQueryAndPartIdentity  int
	outerJoinMap                 map[string]*LeftOuterJoinInfo
	emptyStringQueryAllowed bool
}
func (b *BaseSqlClause) IsAllowEmptyStringQuery() bool{
	return b.emptyStringQueryAllowed
}
func (b *BaseSqlClause) AllowEmptyStringQuery(){
	b.emptyStringQueryAllowed = true
}
func (b *BaseSqlClause) EndOrScopeQueryAndPart() {
	//assertCurrentTmpOrScopeQueryInfo();
	b.orScopeQueryAndPartEffective = false
}
func (b *BaseSqlClause) BeginOrScopeQueryAndPart() {
	//assertCurrentTmpOrScopeQueryInfo();
	b.orScopeQueryAndPartIdentity++
	b.orScopeQueryAndPartEffective = true
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
	b.orScopeQueryAndPartEffective = false
}
func (b *BaseSqlClause) reflectTmpOrClauseToRealObject(localInfo *OrScopeQueryInfo) {
	reflector := new(OrScopeQueryReflector)
	if b.WhereList == nil {
		b.WhereList = new(List)
	}
	if b.BaseTableInlineWhereList == nil {
		b.BaseTableInlineWhereList = new(List)
	}
	if b.outerJoinMap == nil {
		b.outerJoinMap = make(map[string]*LeftOuterJoinInfo)
	}
	reflector.whereList = b.WhereList
	reflector.baseTableInlineWhereList = b.BaseTableInlineWhereList
	reflector.outerJoinMap = b.outerJoinMap

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
	//        sb.append(getFromHint());
	b.BuildWhereClause(sb)
	//        sb.append(deleteUnionWhereTemplateMark(prepareUnionClause(selectClause)));
	//	        if (!b.needsUnionNormalSelectEnclosing()) {
	sb.WriteString(b.PrepareClauseOrderBy())
	//            sb.append(prepareClauseSqlSuffix());
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
func (b *BaseSqlClause) BuildFromClause(sb *bytes.Buffer) {
	sb.WriteString("\n  from ")
	//	        sb.append(ln()).append("  ");
	//        sb.append("from ");
	//        int tablePos = 7; // basically for in-line view indent
	tablePos := 7
	tablePos = tablePos
	//        if (isJoinInParentheses()) {
	//            for (int i = 0; i < getOuterJoinMap().size(); i++) {
	//                sb.append("(");
	//                ++tablePos;
	//            }
	//        }
	//        final TableSqlName tableSqlName = getDBMeta().getTableSqlName();
	tableSqlName := (*b.DBMeta).GetTableSqlName()
	//        final String basePointAliasName = getBasePointAliasName();
	basePointAliasName := b.BasePorintAliasName
	//        if (hasBaseTableInlineWhereClause()) {
	//            final List<QueryClause> baseTableInlineWhereList = getBaseTableInlineWhereList();
	//            sb.append(getInlineViewClause(tableSqlName, baseTableInlineWhereList, tablePos));
	//            sb.append(" ").append(basePointAliasName);
	//        } else {
	//            sb.append(tableSqlName).append(" ").append(basePointAliasName);
	sb.WriteString(tableSqlName.TableSqlName + " " + basePointAliasName)
	//        }
	//        sb.append(getFromBaseTableHint());
	//        sb.append(getLeftOuterJoinClause());
	return
}
func (b *BaseSqlClause) GetSelectClause() string {
	//	        reflectClauseLazilyIfExists();
	//        if (isSelectClauseNonUnionScalar()) {
	//            return buildSelectClauseScalar(getBasePointAliasName());
	//        }
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
	sb, selectIndex := b.ProcessSelectClauseLocal()
	selectIndex = selectIndex
	log.InternalDebug("Select clause =" + sb)
	return sb
}
func (b *BaseSqlClause) ProcessSelectClauseLocal() (string, int32) {
	var sb string = ""
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
			sb += ", "
		} else {
			//            } else {
			//                sb.append("select");
			//                appendSelectHint(sb);
			//                sb.append(" ");
			//                needsDelimiter = true;
			//            }
			sb += "select/*$pmb.selectHint*/ "
			needsDelimiter = true
		}

		//            final String realColumnName = basePointAliasName + "." + columnSqlName;
		realColumnName := basePointAliasName + "." + columnSqlName.ColumnSqlName
		var onQueryName string
		//            final String onQueryName;
		//            ++selectIndex;
		selectIndex++
		//            if (_useSelectIndex) {
		//                onQueryName = buildSelectIndexAlias(columnSqlName, null, selectIndex);
		//                registerSelectIndex(columnDbName, onQueryName, selectIndex);
		if b.UseSelectIndex {
			onQueryName = b.BuildSelectIndexAlias(columnSqlName, "", selectIndex)
		} else {
			//            } else {
			//                onQueryName = columnSqlName.toString();
			//            }
			onQueryName = columnSqlName.ColumnSqlName
		}
		//            sb.append(decryptSelectColumnIfNeeds(columnInfo, realColumnName)).append(" as ").append(onQueryName);
		sb += realColumnName + " as " + onQueryName
		//            getSelectClauseRealColumnAliasMap().put(realColumnName, onQueryName);
		//
		//            if (validSpecifiedLocal && localSpecifiedMap.containsKey(columnDbName)) {
		//                final HpSpecifiedColumn specifiedColumn = localSpecifiedMap.get(columnDbName);
		//                specifiedColumn.setOnQueryName(onQueryName); // basically for queryInsert()
		//            }
		//        }
	}
	//        return selectIndex;
	return sb, selectIndex
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
	return b.orScopeQueryAndPartEffective
}
func (b *BaseSqlClause) RegisterWhereClause(crn *ColumnRealName, key *ConditionKey, cvalue *ConditionValue, co *ConditionOption, usedAliasName string) {
	//Assert 省略
	clauseList := b.GetWhereClauseList4Register()
	b.DoRegisterWhereClause(clauseList, crn, key, cvalue, co, false, false)
	cmap := cvalue.Fixed
	log.InternalDebug(fmt.Sprint(" Sql Clause Cvalue %v cmap %v\n", cvalue, cmap))
	//	        reflectWhereUsedToJoin(usedAliasName);
	//        if (!ConditionKey.isNullaleConditionKey(key)) {
	//            registerInnerJoinLazyReflector(usedAliasName);
	//        }
}

func (b *BaseSqlClause) DoRegisterWhereClause(clauseList *List, crn *ColumnRealName, key *ConditionKey, cvalue *ConditionValue, co *ConditionOption, inline bool, onClause bool) {
	(*key).AddWhereClause(key, b.XcreateQueryModeProvider(), clauseList, crn, cvalue, co)
	//MarkOrScopeQueryAndPart(clauseList); 未実装
}
func (b *BaseSqlClause) XcreateQueryModeProvider() *QueryModeProvider {
	qm := new(QueryModeProvider)
	qm.IsInline = b.OrScopeQueryEffective
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
type SqlClauseMySql struct {
	BaseSqlClause
}
type SqlClauseSqlServer struct {
	BaseSqlClause
}
type SqlClausePostgres struct {
	BaseSqlClause
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
	}
		if code == "mysql" {
		sqlp := new(SqlClauseMySql)
		sqlp.BasePorintAliasName = "dfloc"
		sql = sqlp
	}
			if code == "sqlserver" {
		sqlp := new(SqlClauseSqlServer)
		sqlp.BasePorintAliasName = "dfloc"
		sql = sqlp
	}
	sql.setTableDbName(tableDbName)
	sql.setSqlSetting(dc)
	return &sql
	return nil
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
}
type OrScopeQueryClauseListProviderWhereList struct {
}

func (o *OrScopeQueryClauseListProviderWhereList) provide(tmpOrScopeQueryInfo *OrScopeQueryInfo) interface{} {
	log.InternalDebug(fmt.Sprintf("tmpOrWhereList Size %d \n", tmpOrScopeQueryInfo.tmpOrWhereList.Size()))
	return &tmpOrScopeQueryInfo.tmpOrWhereList
}

type OrScopeQueryClauseListProviderBaseTableInlineWhereList struct {
}

func (o *OrScopeQueryClauseListProviderBaseTableInlineWhereList) provide(tmpOrScopeQueryInfo *OrScopeQueryInfo) interface{} {
	return &tmpOrScopeQueryInfo.tmpOrBaseTableInlineWhereList
}

type OrScopeQueryClauseListProviderAdditionalOnClauseList struct {
}

func (o *OrScopeQueryClauseListProviderAdditionalOnClauseList) provide(tmpOrScopeQueryInfo *OrScopeQueryInfo) interface{} {
	return &tmpOrScopeQueryInfo.tmpOrAdditionalOnClauseListMap
}

type OrScopeQueryClauseListProviderOuterJoinInlineClauseList struct {
}

func (o *OrScopeQueryClauseListProviderOuterJoinInlineClauseList) provide(tmpOrScopeQueryInfo *OrScopeQueryInfo) interface{} {
	return &tmpOrScopeQueryInfo.tmpOrOuterJoinInlineClauseListMap
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
		temp2 := (o.setupTmpOrListList(localInfo, &oc3)).(map[string]*List)
		for temp := range temp2 {
			groupList2.Add(temp2[temp])
		}

		o.setupOrScopeQuery(groupList2, &joinInfo.additionalOnClauseList, false)
	}
	// to InlineView for relation
	for aliasName := range o.outerJoinMap {
		joinInfo := o.outerJoinMap[aliasName]

		wl3 := new(OrScopeQueryClauseListProviderOuterJoinInlineClauseList)
		var oc3 OrScopeQueryClauseListProvider = wl3
		groupList2 := new(List)
		temp2 := (o.setupTmpOrListList(localInfo, &oc3)).(map[string]*List)
		for temp := range temp2 {
			groupList2.Add(temp2[temp])
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

type FixedConditionResolver struct {
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
