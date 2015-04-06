package df

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mikeshimura/dbflute/log"
	"reflect"
	"strings"
)

type ConditionQuery interface {
	GetBaseConditionQuery() *BaseConditionQuery
	//	GetCQProerty() string
}

type BaseConditionQuery struct {
	TableDbName    string
	ReferrerQuery  *ConditionQuery
	SqlClause      *SqlClause
	AliasName      string
	NestLevel      int8
	DBMetaProvider *DBMetaProvider
	Inline         bool
	OnClause       bool
	CQ_PROPERTY    string
	ConditionQuery interface{}
}

func (b *BaseConditionQuery) FRES(value interface{}) interface{} {
	if (*b.SqlClause).IsAllowEmptyStringQuery() {
		return value
	}
	switch value.(type) {
	case sql.NullString:
		var nstr sql.NullString = value.(sql.NullString)
		if nstr.Valid && nstr.String == "" {
			nstr.Valid = false
		}
		return nstr
	case *sql.NullString:
		var nstr *sql.NullString = value.(*sql.NullString)
		if nstr.Valid && nstr.String == "" {
			nstr.Valid = false
		}
		return nstr
	case string:
		var str string = value.(string)
		if str == "" {
			var null sql.NullString
			null.Valid = false
			return null
		}
	case *string:
		var strx string = *value.(*string)
		if strx == "" {
			var null sql.NullString
			null.Valid = false
			return null
		}
	default:
		panic("This type not supported :" + GetType(value))
	}
	return value
}
func (b *BaseConditionQuery) CLSOP() *LikeSearchOption {
	lso := new(LikeSearchOption)
	lso.LikePrefix()
	return lso
}
func (b *BaseConditionQuery) RegROO(minNumber interface{}, maxNumber interface{}, cvalue *ConditionValue, col string, option *RangeOfOption) error {
	if option == nil {
		return errors.New("RangeOfOption is nil")
	}
	//not implemented yet
	//   if (option.hasCalculationRange()) {
	//            final ConditionBean dreamCruiseCB = _baseCB.xcreateDreamCruiseCB();
	//            //dreamCruiseCB.x
	//            //dreamCruiseCB.overTheWaves(xcreateManualOrderSpecifiedColumn(dreamCruiseCB));
	//            option.xinitCalculationRange(_baseCB, dreamCruiseCB);
	//        }

	minKey := option.getMinNumberConditionKey()
	minValidQuery := true
	//        final boolean minValidQuery = isValidQueryNoCheck(minKey, minNumber, cvalue, columnDbName);
	//
	maxKey := option.getMaxNumberConditionKey()
	maxValidQuery := true
	//        final boolean maxValidQuery = isValidQueryNoCheck(maxKey, maxNumber, cvalue, columnDbName);
	//
	needsAndPart := b.isOrScopeQueryDirectlyUnder() && minValidQuery && maxValidQuery
	if needsAndPart {
		(*b.SqlClause).BeginOrScopeQueryAndPart()
	}
	var co ConditionOption = option
	//        try {
	//            if (minValidQuery) {
	b.SetupConditionValueAndRegisterWhereClauseSub(minKey, minNumber, cvalue, col, &co)
	//            }
	//            if (maxValidQuery) {
	b.SetupConditionValueAndRegisterWhereClauseSub(maxKey, maxNumber, cvalue, col, &co)
	//            } else {
	//                if (!minValidQuery) { // means both queries are invalid
	//                    final List<ConditionKey> keyList = newArrayList(minKey, maxKey);
	//                    final List<Number> valueList = newArrayList(minNumber, maxNumber);
	//                    handleInvalidQueryList(keyList, valueList, columnDbName);
	//                }
	//            }
	//        } finally {
	if needsAndPart {
		(*b.SqlClause).EndOrScopeQueryAndPart()
	}
	return nil
}
func (b *BaseConditionQuery) RegLSQ(key *ConditionKey, value string, cvalue *ConditionValue, col string, option *LikeSearchOption) error {
	if option == nil {
		panic("LikeSearchOption nil")
	}
	var co ConditionOption = option
	if !b.IsValidQueryChecked(key, value, cvalue, col) {
		return errors.New("Invalid Query")
	}
	//not implemented
	//        if l.xsuppressEscape() {
	//            option.notEscape();
	//        }
	// basically for DBMS that has original wild-cards
	//not implemented
	//        b.SqlClause.adjustLikeSearchDBWay(option);

	if !option.isSplit() {
		//            if (option.canOptimizeCompoundColumnLikePrefix()) {
		//                // - - - - - - - - - -
		//                // optimized compound
		//                // - - - - - - - - - -
		//                doRegisterLikeSearchQueryCompoundOptimized(value, cvalue, columnDbName, option);
		//            } else {
		//                // - - - - - - - - - - - - -
		//                // normal or normal compound
		//                // - - - - - - - - - - - - -
		b.SetupConditionValueAndRegisterWhereClauseSub(key, value, cvalue, col, &co)
		//            }
		return nil
	}
	// - - - - - - -
	// splitByXxx()
	// - - - - - - -
	return b.doRegisterLikeSearchQuerySplitBy(key, value, cvalue, col, option)

}
func (b *BaseConditionQuery) doRegisterLikeSearchQuerySplitBy(key *ConditionKey, value string, cvalue *ConditionValue, col string, option *LikeSearchOption) error {
	//      assertObjectNotNull("option(LikeSearchOption)", option);
	// these values should be valid only (already filtered before)
	// and invalid values are ignored even at the check mode
	// but if all elements are invalid, it is an exception

	strArray := option.GenerateSplitValueArray(value)
	//fmt.Printf("%v\n", strArray)
	if len(strArray) == 0 {
		return errors.New("Parameter Empty")
	}

	if !option.asOrSplit {
		// as 'and' condition
		needsAndPart := b.isOrScopeQueryDirectlyUnder()
		if needsAndPart {
			(*b.SqlClause).BeginOrScopeQueryAndPart()
		}
		//            try {
		for i := 0; i < len(strArray); i++ {
			currentValue := strArray[i]
			var co ConditionOption = option
			b.SetupConditionValueAndRegisterWhereClauseSub(key, currentValue, cvalue, col, &co)
		}
		//            } finally {
		if needsAndPart {
			(*b.SqlClause).EndOrScopeQueryAndPart()
		}
		//            }
	} else {
		//            // as 'or' condition
		//            if (isOrScopeQueryAndPartEffective()) {
		//                // limit because of so complex
		//                String msg = "The AsOrSplit in and-part is unsupported: " + getTableDbName();
		//                throw new OrScopeQueryAndPartUnsupportedOperationException(msg);
		//            }
		needsNewOrScope := !(*b.SqlClause).IsOrScopeQueryEffective()
		log.InternalDebug(fmt.Sprintf("needsNewOrScope %v \n", needsNewOrScope))
		if needsNewOrScope {

			(*b.SqlClause).MakeOrScopeQueryEffective()
		}
		//            try {
		for i := 0; i < len(strArray); i++ {
			currentValue := strArray[i]
			var co ConditionOption = option
			if i == 0 {

				b.SetupConditionValueAndRegisterWhereClauseSub(key, currentValue, cvalue, col, &co)

			} else {
				b.invokeQueryLikeSearch(col, currentValue, option)
			}
		}
		//            } finally {
		log.InternalDebug(fmt.Sprintf("needsNewOrScope end %v \n", needsNewOrScope))
		if needsNewOrScope {
			(*b.SqlClause).CloseOrScopeQuery()
		}
		//            }
	}

	return nil
}
func (b *BaseConditionQuery) invokeQueryLikeSearch(col string, value interface{}, option interface{}) {
	b.doInvokeQuery(col, "likeSearch", value, option)
}
func (b *BaseConditionQuery) doInvokeQuery(col string, ckey string, value interface{}, option interface{}) {
	//	      assertStringNotNullAndNotTrimmedEmpty("columnFlexibleName", colName);
	//        assertStringNotNullAndNotTrimmedEmpty("conditionKeyName", ckey);
	if value == nil {
		return // do nothing if the value is null when the key has arguments
	}
	container := b.xhelpExtractingPropertyNameCQContainer(col)
	flexibleName := container.flexibleName
	cq := container.cq
	var dbmeta *DBMeta
	bcq := reflect.ValueOf(cq).MethodByName("GetBaseConditionQuery").Call([]reflect.Value{})
	bcqx := (bcq[0].Interface()).(*BaseConditionQuery)
	//fmt.Printf("bcq %v %T \n", bcqx, bcqx)
	dbmeta = DBMetaProvider_I.TableDbNameInstanceMap[bcqx.TableDbName]
	//fmt.Printf("dbmeta %v %T \n", dbmeta, dbmeta)
	cino := (*dbmeta).GetColumnInfoMap()[flexibleName]
	//fmt.Printf("cino %d \n", cino)
	ci := (*dbmeta).GetColumnInfoList().Get(cino)
	if ci == nil {
		panic("ColumnFindFailure :" + flexibleName)
	}
	var columnInfo *ColumnInfo = ci.(*ColumnInfo)
	columnCapPropName := InitCap(columnInfo.PropertyName)
	var noArg bool
	var rangeOf bool
	rangeOf = rangeOf
	var fromTo bool
	fromTo = fromTo
	ckeyl := strings.ToLower(ckey)
	if ckeyl == "isnull" || ckeyl == "isnotnull" || ckeyl == "isnullorempty" || ckeyl == "emptystring" {
		noArg = true
	}
	if ckeyl == "rangeof" {
		rangeOf = true
	}
	if ckeyl == "fromto" || ckeyl == "datefromto" {
		fromTo = true
	}
	if !noArg {
		//not implemented
		//            try {
		//                value = columnInfo.toPropretyType(value); // convert type
		//            } catch (RuntimeException e) {
		//                throwConditionInvokingValueConvertFailureException(colName, ckey, value, option, e);
		//            }
	}

	methodName := "Set" + columnCapPropName + "_" + InitCap(ckey)
	//fmt.Printf("cino %s \n", methodName)
	var param []reflect.Value = make([]reflect.Value, 2)
	param[0] = reflect.ValueOf(value)
	param[1] = reflect.ValueOf(option)
	reflect.ValueOf(cq).MethodByName(methodName).Call(param)

	//        final List<Class<?>> typeList = newArrayList();
	//        if (fromTo) {
	//            typeList.add(Date.class);
	//            typeList.add(Date.class);
	//        } else if (rangeOf) {
	//            final Class<?> propertyType = columnInfo.getPropertyType();
	//            typeList.add(propertyType);
	//            typeList.add(propertyType);
	//        } else {
	//            if (!noArg) {
	//                typeList.add(value.getClass());
	//            }
	//        }
	//        if (option != null) {
	//            typeList.add(option.getClass());
	//        }
	//        final Class<?>[] parameterTypes = typeList.toArray(new Class<?>[] {});
	//        final Method method = xhelpGettingCQMethod(cq, methodName, parameterTypes);
	//        if (method == null) {
	//            throwConditionInvokingSetMethodNotFoundException(colName, ckey, value, option, methodName, parameterTypes);
	//        }
	//        try {
	//            final List<Object> argList = newArrayList();
	//            if (fromTo || rangeOf) {
	//                if (!(value instanceof List<?>)) { // check type
	//                    throwConditionInvokingDateFromToValueInvalidException(colName, ckey, value, option, methodName,
	//                            parameterTypes);
	//                }
	//                argList.addAll((List<?>) value);
	//            } else {
	//                if (!noArg) {
	//                    argList.add(value);
	//                }
	//            }
	//            if (option != null) {
	//                argList.add(option);
	//            }
	//            xhelpInvokingCQMethod(cq, method, argList.toArray());
	//        } catch (ReflectionFailureException e) {
	//            throwConditionInvokingSetReflectionFailureException(colName, ckey, value, option, methodName,
	//                    parameterTypes, e);
	//        }

}
func (b *BaseConditionQuery) xhelpExtractingPropertyNameCQContainer(col string) *PropertyNameCQContainer {
	//	        final String[] strings = name.split("\\.");
	//        final int length = strings.length;
	//        String propertyName = null;
	//        ConditionQuery cq = this;
	//        int index = 0;
	//        for (String element : strings) {
	//            if (length == (index + 1)) { // at last loop!
	//                propertyName = element;
	//                break;
	//            }
	//            cq = cq.invokeForeignCQ(element);
	//            ++index;
	//        }
	//        return new PropertyNameCQContainer(propertyName, cq);
	//temporary implimation
	container := new(PropertyNameCQContainer)
	container.cq = b.ConditionQuery
	container.flexibleName = col
	return container
}
func (b *BaseConditionQuery) isOrScopeQueryDirectlyUnder() bool {
	orScopeQuery := (*b.SqlClause).IsOrScopeQueryEffective()
	orScopeQueryAndPart := (*b.SqlClause).IsOrScopeQueryAndPartEffective()
	return orScopeQuery && !orScopeQueryAndPart
}

//func (b *BaseConditionQuery) isValidQueryChecked(key *ConditionKey, value string, cvalue *ConditionValue, col string) bool {
//	//mot implemented
//	return true
//}
func (b *BaseConditionQuery) GetCQProerty() string {
	return b.CQ_PROPERTY
}
func (b *BaseConditionQuery) IsBaseQuery() bool {
	return b.ReferrerQuery == nil
}

func (b *BaseConditionQuery) RegQ(key *ConditionKey, value interface{}, cvalue *ConditionValue, col string) {
	//fmt.Printf("RegQ col %s\n", col)
	if b.IsValidQueryChecked(key, value, cvalue, col) == false {
		return
	}
	b.SetupConditionValueAndRegisterWhereClause(key, value, cvalue, col)
}

func (b *BaseConditionQuery) SetupConditionValueAndRegisterWhereClause(key *ConditionKey, value interface{}, cvalue *ConditionValue, col string) {
	eo := b.CreateEmbeddedOption(key, value, cvalue, col)
	b.SetupConditionValueAndRegisterWhereClauseSub(key, value, cvalue, col, eo)

}

func (b *BaseConditionQuery) SetupConditionValueAndRegisterWhereClauseSub(key *ConditionKey, value interface{}, cvalue *ConditionValue, col string, co *ConditionOption) {
	dm := b.DBMetaProvider.TableDbNameInstanceMap[b.TableDbName]
	cinfo := (*dm).GetColumnInfoByPropertyName(col)
	pn := cinfo.PropertyName
	un := InitUnCap(pn)
	loc := b.XgetLocation(un)
	log.InternalDebug(fmt.Sprintln("Location :" + loc))
	(*key).SetupConditionValue(key, b.XcreateQueryModeProvider(), cvalue, value, loc, co)
	crn := b.ToColumnRealName(col, cinfo.ColumnSqlName)
	usedAliasName := b.AliasName
	(*b.SqlClause).RegisterWhereClause(crn, key, cvalue, co, usedAliasName)

}
func (b *BaseConditionQuery) ToColumnRealName(col string, csn *ColumnSqlName) *ColumnRealName {
	var crn *ColumnRealName
	if csn != nil {
		crn = CreateColumnRealName(b.AliasName, csn)
	} else {
		dbmeta := b.DBMetaProvider.TableDbNameInstanceMap[b.TableDbName]
		log.InternalDebug(fmt.Sprintf("ToColumnRealName dbmeta %v \n", dbmeta))
		cno := (*dbmeta).GetColumnInfoMap()[col]
		ci := (*dbmeta).GetColumnInfoList().Get(cno)
		columnInfo := ci.(*ColumnInfo)
		crn = CreateColumnRealName(b.AliasName, columnInfo.ColumnSqlName)
	}
	return crn
}
func (b *BaseConditionQuery) XgetLocation(propertyName string) string {
	return b.XgetLocationBase() + InitCap(propertyName)
}
func (b *BaseConditionQuery) XgetLocationBase() string {
	res := ""

	for {
		if b.IsBaseQuery() {
			res = b.GetCQProerty() + "." + res
			break
		}
		//xgetForeignPropertyName
		//////未実装
		break

	}
	return res
}
func (b *BaseConditionQuery) IsValidQueryChecked(key *ConditionKey, value interface{}, cvalue *ConditionValue, col string) bool {
	return b.XdoIsValidQuery(key, value, cvalue, col, true)
}
func (b *BaseConditionQuery) XdoIsValidQuery(key *ConditionKey, value interface{}, cvalue *ConditionValue, col string, checked bool) bool {
	callerName := b.ToColumnRealName(col, nil) // logging only
	if (*key).IsValidRegistration(b.XcreateQueryModeProvider(), cvalue, value, callerName) {
		return true
	} else {
		if checked {
			b.handleInvalidQuery(key, value, col)
		}
		return false
	}

	return true
}
func (b *BaseConditionQuery) handleInvalidQuery(key *ConditionKey, value interface{}, col string) {
	//not implemented yet
}
func (b *BaseConditionQuery) CreateEmbeddedOption(key *ConditionKey, value interface{}, cvalue *ConditionValue, col string) *ConditionOption {
	return nil
}
func (b *BaseConditionQuery) RegOBA(col string) {
	// 	if b.SqlClause.GetSqlSetting().NoOrderBy はError
	//	if b.SqlClause.GetSqlSetting(). {
	//		log.Error("BaseConditionBean", "df007:BaseCondition Bean Purpose Error:col="+col)
	//		return
	//	}
	b.RegisterOrderBy(col, true)
}
func (b *BaseConditionQuery) RegOBD(col string) {
	b.RegisterOrderBy(col, false)
}
func (b *BaseConditionQuery) FindDbMeta() *DBMeta {
	dbp := b.DBMetaProvider
	tn := dbp.TableDbNameFlexibleMap.Get(b.TableDbName).(string)
	return dbp.TableDbNameInstanceMap[tn]
}
func (b *BaseConditionQuery) RegisterOrderBy(col string, ascDesc bool) {
	dbm := b.FindDbMeta()
	ci := (*dbm).GetColumnInfoByPropertyName(col)
	//	rn:=CreateColumnRealName (b.AliasName,ci.ColumnSqlName)
	//	rn=rn
	oe := new(OrderByElement)
	oe.AliasName = b.AliasName
	oe.ColumnName = ci.ColumnSqlName.ColumnSqlName
	if ascDesc {
		oe.AscDesc = "asc"
	} else {
		oe.AscDesc = "desc"
	}
	(*b.SqlClause).GetOrderByClause().AddElement(oe)

}
func (b *BaseConditionQuery) XcreateQueryModeProvider() *QueryModeProvider {
	qm := new(QueryModeProvider)
	qm.IsInline = (*b.SqlClause).IsOrScopeQueryEffective()
	qm.IsInline = b.Inline
	qm.IsOnClause = b.OnClause
	return qm
}

type ColumnRealName struct {
	TableAliasName string
	ColumnSqlName  *ColumnSqlName
}

func (c *ColumnRealName) ToString() string {
	return c.TableAliasName + "." + c.ColumnSqlName.ColumnSqlName
}
func CreateColumnRealName(aliasName string, columnSqlName *ColumnSqlName) *ColumnRealName {
	rn := new(ColumnRealName)
	rn.TableAliasName = aliasName
	rn.ColumnSqlName = columnSqlName
	return rn
}

type QueryModeProvider struct {
	IsOrScopeQuery bool
	IsInline       bool
	IsOnClause     bool
}
type PropertyNameCQContainer struct {
	flexibleName string
	cq           interface{}
}
