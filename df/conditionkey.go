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
	"fmt"
	"github.com/mikeshimura/dbflute/log"
	"strconv"
	"strings"
)

const (
	FIXED_KEY_QUERY    = "query"
	FIXED_KEY_INLINE   = "inline"
	FIXED_KEY_ONCLAUSE = "onClause"
	C_EQ               = "equal"
	C_GT               = "greaterThan"
	C_NE               = "notEqual"
	C_LT               = "lessThan"
	C_GE               = "greaterEqual"
	C_LE               = "lessEqual"
	C_ISN              = "isNull"
	C_ISNN             = "isNotNull"
	C_ISNOE            = "isNullOrEmpty"
	C_LS               = "likeSearch"
	C_NLS              = "notLikeSearch"
	C_GTISN            = "greaterThanOrIsNull"
	C_GEISN            = "greaterEqualOrIsNull"
	C_LTISN            = "lessThanOrIsNull"
	C_LEISN            = "lessEqualOrIsNull"
	C_INS              = "inScope"
)

type ConditionKey interface {
	AddWhereClause(ck *ConditionKey, qmp *QueryModeProvider, conditionList *List,
		columnRealName *ColumnRealName, cvalue *ConditionValue, co *ConditionOption)
	DoAddWhereClause(conditionList *List, columnRealName *ColumnRealName,
		cvalue *ConditionValue, co *ConditionOption)
	GetConditionKeyS() string
	GetOperand() string
	DoSetupConditionValue(cvalue *ConditionValue, value interface{},
		location string, co *ConditionOption)
	SetupConditionValue(ck *ConditionKey, qmp *QueryModeProvider,
		cvalue *ConditionValue, value interface{}, location string, co *ConditionOption)
	BuildBindClause(columnRealName *ColumnRealName, loc string,
		co *ConditionOption) *QueryClause
	IsValidRegistration(provider *QueryModeProvider, cvalue *ConditionValue,
		value interface{}, callerName *ColumnRealName) bool
	GetBindVariableDummyValue() string
}

func ConditionKey_IsNullaleConditionKey(key *ConditionKey) bool {
	return CK_GEISN.conditionKey == key || CK_GTISN.conditionKey == key ||
		CK_LEISN.conditionKey == key || CK_LTISN.conditionKey == key ||
		CK_ISN.conditionKey == key || CK_ISNOE.conditionKey == key
}

type BaseConditionKey struct {
	ConditionKeyS string
	Operand       string
	conditionKey  *ConditionKey
}

func (w *BaseConditionKey) IsValidRegistration(provider *QueryModeProvider,
	cvalue *ConditionValue, value interface{}, callerName *ColumnRealName) bool {
	//only null check implemented
	return IsNotNull(value)
}
func (w *BaseConditionKey) BuildBindClauseOrIsNull(columnRealName *ColumnRealName,
	loc string, co *ConditionOption) *QueryClause {
	mainQuery := w.DoBuildBindClause(columnRealName, loc, co)
	clause := "(" + mainQuery + " or " + columnRealName.ToString() + " is null)"
	sqc := new(StringQueryClause)
	sqc.Clause = clause
	var qc QueryClause = sqc
	log.InternalDebug(fmt.Sprintf(" QueryClause %v \n", qc))
	return &qc
}

func (w *BaseConditionKey) BuildBindClause(columnRealName *ColumnRealName,
	loc string, co *ConditionOption) *QueryClause {
	sqc := CreateStringQueryClause(w.DoBuildBindClause(columnRealName, loc, co))
	var qc QueryClause = sqc
	log.InternalDebug(fmt.Sprintf(" QueryClause %v \n", qc))
	return &qc
}
func (w *BaseConditionKey) DoBuildBindClause(columnRealName *ColumnRealName,
	loc string, co *ConditionOption) string {
	result := w.ResolveBindClause(columnRealName, loc, co)
	log.InternalDebug(fmt.Sprintf(" result %v \n", result))
	return result.ToBindClause()
}
func (w *BaseConditionKey) ResolveBindClause(columnRealName *ColumnRealName,
	loc string, co *ConditionOption) *BindClauseResult {
	var basicBindExp string
	if loc != "" {
		basicBindExp = w.BuildBindVariableExp(loc, co)
		log.InternalDebug("basicBindExp :" + basicBindExp)
	}
	resolvedColumn := w.ResolveOptionalColumn(columnRealName, co)
	return w.CreateBindClauseResult(resolvedColumn, basicBindExp, co)

}
func (w *BaseConditionKey) ResolveOptionalColumn(columnExp *ColumnRealName,
	co *ConditionOption) *ColumnRealName {
	return w.ResolveCalculationColumn(w.ResolveCompoundColumn(columnExp, co), co)
}
func (w *BaseConditionKey) ResolveCalculationColumn(columnRealName *ColumnRealName,
	co *ConditionOption) *ColumnRealName {
	if co == nil {
		return columnRealName
	}
	//未実装
	return columnRealName
}

func (w *BaseConditionKey) ResolveCompoundColumn(baseRealName *ColumnRealName,
	co *ConditionOption) *ColumnRealName {
	if co == nil || !(*co).HasCompoundColumn() {
		return baseRealName
	}
	//未実装
	return baseRealName
}
func (w *BaseConditionKey) CreateBindClauseResult(columnExp *ColumnRealName,
	bindExp string, co *ConditionOption) *BindClauseResult {
	op := w.ResolveOperand(co)
	rearOption := w.ResolveRearOption(co)
	//fmt.Printf("op %v rearOption %v \n", op, rearOption)
	result := NewBindClauseResult(columnExp, op, bindExp, rearOption)
	result.Arranger = w.ResolveWhereClauseArranger(co)
	return result
}
func (w *BaseConditionKey) ResolveWhereClauseArranger(
	co *ConditionOption) *QueryClauseArranger {
	//未実装
	return nil
}

func (w *BaseConditionKey) ResolveRearOption(co *ConditionOption) string {
	if co != nil {
		return (*co).GetRearOption()
	}
	return ""
}
func (w *BaseConditionKey) ResolveOperand(co *ConditionOption) string {
	op := w.ExtractExtOperand(co)
	if op == "" {
		return w.GetOperand()
	} else {
		return op
	}
}
func (w *BaseConditionKey) ExtractExtOperand(co *ConditionOption) string {
	//未実装
	return ""
}
func (w *BaseConditionKey) BuildBindVariableExp(
	loc string, co *ConditionOption) string {
	return "/*pmb." + loc + "*/" + (*w.conditionKey).GetBindVariableDummyValue()
}
func (w *BaseConditionKey) GetBindVariableDummyValue() string {
	return ""
}

type BindClauseResult struct {
	ColumnExp  *ColumnRealName
	Operand    string
	BindExp    string
	RearOption string
	Arranger   *QueryClauseArranger
}

func NewBindClauseResult(ColumnExp *ColumnRealName, Operand string,
	BindExp string, RearOption string) *BindClauseResult {
	bcr := new(BindClauseResult)
	bcr.ColumnExp = ColumnExp
	bcr.Operand = Operand
	bcr.BindExp = BindExp
	bcr.RearOption = RearOption
	return bcr
}

func (b *BindClauseResult) ToBindClause() string {
	var clause string
	//Temporary for T/S
	//	if b.Arranger != nil {
	//		clause = (*b.Arranger).Arrange(b.ColumnExp, b.Operand, b.BindExp, b.RearOption)
	//	} else {
	clause = b.ColumnExp.ToString() + " " + b.Operand + " " +
		b.BindExp + b.RearOption
	//	}
	log.InternalDebug("operand :" + b.Operand)
	log.InternalDebug("clause:" + clause)
	log.InternalDebug("bindExp :" + b.BindExp)
	log.InternalDebug("clause :" + clause)
	return clause
}

type QueryClauseArranger interface {
	Arrange(columnRealName *ColumnRealName, operand string,
		bindExpression string, rearOption string) string
}

func (w *BaseConditionKey) AddWhereClause(ck *ConditionKey,
	qmp *QueryModeProvider, conditionList *List, columnRealName *ColumnRealName,
	cvalue *ConditionValue, co *ConditionOption) {
	cvalue.OrScopeQuery = qmp.IsOrScopeQuery()
	cvalue.Inline = qmp.IsInline
	cvalue.OnClause = qmp.IsOnClause
	(*ck).DoAddWhereClause(conditionList, columnRealName, cvalue, co)
//	cvalue.OrScopeQuery = false
//	cvalue.Inline = false
//	cvalue.OnClause = false
}

func (w *BaseConditionKey) DoAddWhereClause(conditionList *List,
	columnRealName *ColumnRealName, cvalue *ConditionValue, co *ConditionOption) {
}
func (w *BaseConditionKey) GetConditionKeyS() string {
	return w.ConditionKeyS
}

func (w *BaseConditionKey) GetOperand() string {
	return w.Operand
}
func (w *BaseConditionKey) DoSetupConditionValue(cvalue *ConditionValue,
	value interface{}, location string, co *ConditionOption) {
}

func (w *BaseConditionKey) SetupConditionValue(ck *ConditionKey,
	qmp *QueryModeProvider, cvalue *ConditionValue, value interface{},
	location string, co *ConditionOption) {
	cvalue.OrScopeQuery = qmp.IsOrScopeQuery()
	cvalue.Inline = qmp.IsInline
	cvalue.OnClause = qmp.IsOnClause
	(*ck).DoSetupConditionValue(cvalue, value, location, co)
//	cvalue.OrScopeQuery = false
//	cvalue.Inline = false
//	cvalue.OnClause = false
	//fmt.Printf("CK EQ lastlocation %v \n", cvalue.EqualLatestLocation)
}

type QueryClause interface {
	ToString() string
	getIdentity() int
}
type BaseQueryClause struct {
}

func (b *BaseQueryClause) ToString() string {
	return ""
}
func (b *BaseQueryClause) getIdentity() int {
	return -1
}

type StringQueryClause struct {
	BaseQueryClause
	Clause string
}

func (s *StringQueryClause) ToString() string {
	return s.Clause
}
func CreateStringQueryClause(str string) *StringQueryClause {
	sqc := new(StringQueryClause)
	sqc.Clause = str
	return sqc
}

type CK_EQ_T struct {
	BaseConditionKey
}

func (w *CK_EQ_T) DoSetupConditionValue(cvalue *ConditionValue, value interface{},
	location string, co *ConditionOption) {
	//fmt.Println("EQl DoSetup")
	cvalue.SetupEqual(value, location)
}
func (w *CK_EQ_T) DoAddWhereClause(conditionList *List, columnRealName *ColumnRealName,
	cvalue *ConditionValue, co *ConditionOption) {
	conditionList.Add(w.BuildBindClause(columnRealName, cvalue.EqualLatestLocation, co))
}

var CK_EQ *CK_EQ_T
var CK_EQ_C *ConditionKey

type CK_GT_T struct {
	BaseConditionKey
}

func (w *CK_GT_T) DoSetupConditionValue(cvalue *ConditionValue, value interface{},
	location string, co *ConditionOption) {
	cvalue.SetupGreaterThan(value, location)
}
func (w *CK_GT_T) DoAddWhereClause(conditionList *List, columnRealName *ColumnRealName,
	cvalue *ConditionValue, co *ConditionOption) {
	conditionList.Add((*w.conditionKey).BuildBindClause(columnRealName,
		cvalue.GreaterThanLatestLocation, co))
}

var CK_GT *CK_GT_T
var CK_GT_C *ConditionKey

type CK_NE_T struct {
	BaseConditionKey
}

func (w *CK_NE_T) DoSetupConditionValue(cvalue *ConditionValue,
	value interface{}, location string, co *ConditionOption) {
	cvalue.SetupNotEqual(value, location)
}
func (w *CK_NE_T) DoAddWhereClause(conditionList *List,
	columnRealName *ColumnRealName, cvalue *ConditionValue, co *ConditionOption) {
	conditionList.Add((*w.conditionKey).BuildBindClause(
		columnRealName, cvalue.NotEqualLatestLocation, co))
}

var CK_NE *CK_NE_T
var CK_NE_C *ConditionKey

type CK_LT_T struct {
	BaseConditionKey
}

func (w *CK_LT_T) DoSetupConditionValue(cvalue *ConditionValue, value interface{},
	location string, co *ConditionOption) {
	cvalue.SetupLessThan(value, location)
}
func (w *CK_LT_T) DoAddWhereClause(conditionList *List,
	columnRealName *ColumnRealName, cvalue *ConditionValue, co *ConditionOption) {
	conditionList.Add((*w.conditionKey).BuildBindClause(
		columnRealName, cvalue.LessThanLatestLocation, co))

}

var CK_LT *CK_LT_T
var CK_LT_C *ConditionKey

type CK_GE_T struct {
	BaseConditionKey
}

func (w *CK_GE_T) DoSetupConditionValue(cvalue *ConditionValue,
	value interface{}, location string, co *ConditionOption) {
	cvalue.SetupGreaterEqual(value, location)
}
func (w *CK_GE_T) DoAddWhereClause(conditionList *List, columnRealName *ColumnRealName,
	cvalue *ConditionValue, co *ConditionOption) {
	conditionList.Add((*w.conditionKey).BuildBindClause(
		columnRealName, cvalue.GreaterEqualLatestLocation, co))

}

var CK_GE *CK_GE_T
var CK_GE_C *ConditionKey

type CK_LE_T struct {
	BaseConditionKey
}

func (w *CK_LE_T) DoSetupConditionValue(cvalue *ConditionValue, value interface{},
	location string, co *ConditionOption) {
	cvalue.SetupLessEqual(value, location)
}
func (w *CK_LE_T) DoAddWhereClause(conditionList *List, columnRealName *ColumnRealName,
	cvalue *ConditionValue, co *ConditionOption) {
	conditionList.Add((*w.conditionKey).BuildBindClause(columnRealName,
		cvalue.LessEqualLatestLocation, co))

}

var CK_LE *CK_LE_T
var CK_LE_C *ConditionKey

type CK_ISN_T struct {
	BaseConditionKey
}

func (w *CK_ISN_T) DoSetupConditionValue(cvalue *ConditionValue,
	value interface{}, location string, co *ConditionOption) {
	cvalue.SetupIsNull()
}
func (w *CK_ISN_T) DoAddWhereClause(conditionList *List,
	columnRealName *ColumnRealName, cvalue *ConditionValue, co *ConditionOption) {
	conditionList.Add((*w.conditionKey).BuildBindClause(columnRealName, "", co))

}

var CK_ISN *CK_ISN_T
var CK_ISN_C *ConditionKey

type CK_ISNN_T struct {
	BaseConditionKey
}

func (w *CK_ISNN_T) DoSetupConditionValue(
	cvalue *ConditionValue, value interface{}, location string, co *ConditionOption) {
	cvalue.SetupIsNotNull()
}
func (w *CK_ISNN_T) DoAddWhereClause(conditionList *List, columnRealName *ColumnRealName,
	cvalue *ConditionValue, co *ConditionOption) {
	conditionList.Add((*w.conditionKey).BuildBindClause(columnRealName, "", co))

}

var CK_ISNN *CK_ISNN_T
var CK_ISNN_C *ConditionKey

type CK_ISNOE_T struct {
	BaseConditionKey
}

func (w *CK_ISNOE_T) DoSetupConditionValue(cvalue *ConditionValue, value interface{},
	location string, co *ConditionOption) {
	cvalue.SetupIsNullOrEmpty()
}
func (w *CK_ISNOE_T) DoAddWhereClause(conditionList *List, columnRealName *ColumnRealName,
	cvalue *ConditionValue, co *ConditionOption) {
	sql := "(" + columnRealName.ToString() + " " + w.Operand + " or " +
		columnRealName.ToString() + " = '')"
	sqc := CreateStringQueryClause(sql)
	var qc QueryClause = sqc
	conditionList.Add(&qc)
}

var CK_ISNOE *CK_ISNOE_T
var CK_ISNOE_C *ConditionKey

type CK_LS_T struct {
	BaseConditionKey
}

func (w *CK_LS_T) DoSetupConditionValue(cvalue *ConditionValue, value interface{},
	location string, co *ConditionOption) {
	cvalue.SetupLikeSearch(value, location, co)
}
func (w *CK_LS_T) DoAddWhereClause(conditionList *List, columnRealName *ColumnRealName,
	cvalue *ConditionValue, co *ConditionOption) {
	conditionList.Add((*w.conditionKey).BuildBindClause(
		columnRealName, cvalue.LikeSearchLatestLocation, co))

}

var CK_LS *CK_LS_T
var CK_LS_C *ConditionKey

type CK_NLS_T struct {
	BaseConditionKey
}

func (w *CK_NLS_T) DoSetupConditionValue(cvalue *ConditionValue, value interface{},
	location string, co *ConditionOption) {
	cvalue.SetupNotLikeSearch(value, location, co)
}
func (w *CK_NLS_T) DoAddWhereClause(conditionList *List, columnRealName *ColumnRealName,
	cvalue *ConditionValue, co *ConditionOption) {
	conditionList.Add((*w.conditionKey).BuildBindClause(columnRealName,
		cvalue.NotLikeSearchLatestLocation, co))

}

var CK_NLS *CK_NLS_T
var CK_NLS_C *ConditionKey

type CK_GTISN_T struct {
	CK_GT_T
}

func (w *CK_GTISN_T) BuildBindClause(columnRealName *ColumnRealName, location string,
	co *ConditionOption) *QueryClause {
	return w.BuildBindClauseOrIsNull(columnRealName, location, co)
}

var CK_GTISN *CK_GTISN_T
var CK_GTISN_C *ConditionKey

type CK_GEISN_T struct {
	CK_GE_T
}

func (w *CK_GEISN_T) BuildBindClause(columnRealName *ColumnRealName, location string,
	co *ConditionOption) *QueryClause {
	return w.BuildBindClauseOrIsNull(columnRealName, location, co)
}

var CK_GEISN *CK_GEISN_T
var CK_GEISN_C *ConditionKey

type CK_LTISN_T struct {
	CK_LT_T
}

func (w *CK_LTISN_T) BuildBindClause(columnRealName *ColumnRealName, location string,
	co *ConditionOption) *QueryClause {
	return w.BuildBindClauseOrIsNull(columnRealName, location, co)
}

var CK_LTISN *CK_LTISN_T
var CK_LTISN_C *ConditionKey

type CK_LEISN_T struct {
	CK_LE_T
}

func (w *CK_LEISN_T) BuildBindClause(columnRealName *ColumnRealName, location string,
	co *ConditionOption) *QueryClause {
	return w.BuildBindClauseOrIsNull(columnRealName, location, co)
}

var CK_LEISN *CK_LEISN_T
var CK_LEISN_C *ConditionKey

type CK_INS_T struct {
	BaseConditionKey
}

func (w *CK_INS_T) DoSetupConditionValue(cvalue *ConditionValue, value interface{},
	location string, co *ConditionOption) {
	cvalue.SetupInScope(value, location)
}
func (w *CK_INS_T) DoAddWhereClause(conditionList *List, columnRealName *ColumnRealName,
	cvalue *ConditionValue, co *ConditionOption) {
	conditionList.Add((*w.conditionKey).BuildBindClause(columnRealName,
		cvalue.InScopeLatestLocation, co))

}
func (w *CK_INS_T) GetBindVariableDummyValue() string {
	return "('a1', 'a2')" // to indicate inScope
}

var CK_INS *CK_INS_T
var CK_INS_C *ConditionKey

type ColumnSqlName struct {
	ColumnSqlName string
	IrregularChar bool
}

func (w *ColumnSqlName) AnalyzeIrregularChar() {
	const okchars = "_0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	s := w.ColumnSqlName
	ok := true
	for i := 0; i < len(s); i++ {
		ss := string(s[i])
		ck := strings.Index(okchars, ss)
		if ck == -1 {
			ok = false
		}
	}
	w.IrregularChar = !ok
}

type ConditionValue struct {
	EqualLatestLocation         string
	GreaterThanLatestLocation   string
	NotEqualLatestLocation      string
	LessThanLatestLocation      string
	GreaterEqualLatestLocation  string
	LessEqualLatestLocation     string
	LikeSearchLatestLocation    string
	NotLikeSearchLatestLocation string
	InScopeLatestLocation       string
	EqualValueHandler           *StandardValueHandler
	GreaterThanValueHandler     *StandardValueHandler
	NotEqualValueHandler        *StandardValueHandler
	LessThanValueHandler        *StandardValueHandler
	GreaterEqualValueHandler    *StandardValueHandler
	LessEqualValueHandler       *StandardValueHandler
	IsNullValueHandler          *StandardValueHandler
	IsNotNullValueHandler       *StandardValueHandler
	IsNullOrEmptyValueHandler   *StandardValueHandler
	LikeSearchValueHandler      *VaryingValueHandler
	NotLikeSearchValueHandler   *VaryingValueHandler
	InScopeValueHandler         *VaryingValueHandler
	Fixed                       map[string]map[string]interface{}
	Varying                     map[string]map[string]interface{}
	OrScopeQuery                bool
	Inline                      bool
	OnClause                    bool
}

func (w *ConditionValue) SetupFixedValue(cond *ConditionKey, value interface{}) string {
	if w.Fixed == nil {
		w.Fixed = make(map[string]map[string]interface{})
	}
	fixedValueKey := w.getFixedValueKey()
	var elementMap map[string]interface{} = w.Fixed[fixedValueKey]
	if elementMap == nil {
		elementMap = make(map[string]interface{})
		w.Fixed[fixedValueKey] = elementMap
	}
	key := (*cond).GetConditionKeyS()

	elementMap[key] = value
	return "Fixed." + fixedValueKey + "." + key
}
func (w *ConditionValue) SetupVaryingValue(cond *ConditionKey, value interface{}) string {
	if w.Varying == nil {
		log.InternalDebug("Varying created:")
		w.Varying = make(map[string]map[string]interface{})
	}
	key := (*cond).GetConditionKeyS()
	log.InternalDebug("GetConditionKeyS  :" + key)
	var elementMap map[string]interface{} = w.Varying[key]
	if elementMap == nil {
		elementMap = make(map[string]interface{})
		w.Varying[key] = elementMap
	}
	elementkey := key + strconv.Itoa(len(elementMap))
	elementMap[elementkey] = value
	return "Varying." + key + "." + elementkey
}
func (w *ConditionValue) getFixedValueKey() string {
	if w.Inline {
		if w.OnClause {
			return FIXED_KEY_ONCLAUSE
		} else {
			return FIXED_KEY_INLINE
		}
	} else { // normal query
		return FIXED_KEY_QUERY
	}
}
func CreateVaryingValueHandler(cv *ConditionValue, ck *ConditionKey) *VaryingValueHandler {
	sv := new(VaryingValueHandler)
	sv.ConditionValue = cv
	sv.ConditionKey = ck
	return sv
}
func CreateStandardValueHandler(cv *ConditionValue, ck *ConditionKey) *StandardValueHandler {
	sv := new(StandardValueHandler)
	sv.ConditionValue = cv
	sv.ConditionKey = ck
	return sv
}
func (w *ConditionValue) SetupEqual(value interface{}, location string) {
	if w.EqualValueHandler == nil {
		var ck ConditionKey = CK_EQ
		w.EqualValueHandler = CreateStandardValueHandler(w, &ck)
	}
	w.EqualLatestLocation = location + "." + w.EqualValueHandler.SetValue(value)
}

func (w *ConditionValue) SetupGreaterThan(value interface{}, location string) {
	if w.GreaterThanValueHandler == nil {
		var ck ConditionKey = CK_GT
		w.GreaterThanValueHandler = CreateStandardValueHandler(w, &ck)
	}
	w.GreaterThanLatestLocation = location + "." + w.GreaterThanValueHandler.SetValue(value)
}
func (w *ConditionValue) SetupNotEqual(value interface{}, location string) {
	if w.NotEqualValueHandler == nil {
		var ck ConditionKey = CK_NE
		w.NotEqualValueHandler = CreateStandardValueHandler(w, &ck)
	}
	w.NotEqualLatestLocation = location + "." + w.NotEqualValueHandler.SetValue(value)
}
func (w *ConditionValue) SetupLessThan(value interface{}, location string) {
	if w.LessThanValueHandler == nil {
		var ck ConditionKey = CK_LT
		w.LessThanValueHandler = CreateStandardValueHandler(w, &ck)
	}
	w.LessThanLatestLocation = location + "." + w.LessThanValueHandler.SetValue(value)
}
func (w *ConditionValue) SetupGreaterEqual(value interface{}, location string) {
	if w.GreaterEqualValueHandler == nil {
		var ck ConditionKey = CK_GE
		w.GreaterEqualValueHandler = CreateStandardValueHandler(w, &ck)
	}
	w.GreaterEqualLatestLocation = location + "." + w.GreaterEqualValueHandler.SetValue(value)
}
func (w *ConditionValue) SetupLessEqual(value interface{}, location string) {
	if w.LessEqualValueHandler == nil {
		var ck ConditionKey = CK_LE
		w.LessEqualValueHandler = CreateStandardValueHandler(w, &ck)
	}
	w.LessEqualLatestLocation = location + "." + w.LessEqualValueHandler.SetValue(value)
}
func (w *ConditionValue) SetupIsNull() {
	if w.IsNullValueHandler == nil {
		var ck ConditionKey = CK_ISN
		w.IsNullValueHandler = CreateStandardValueHandler(w, &ck)
	}
}
func (w *ConditionValue) SetupIsNotNull() {
	if w.IsNotNullValueHandler == nil {
		var ck ConditionKey = CK_ISNN
		w.IsNotNullValueHandler = CreateStandardValueHandler(w, &ck)
	}
}

func (w *ConditionValue) SetupIsNullOrEmpty() {
	if w.IsNullOrEmptyValueHandler == nil {
		var ck ConditionKey = CK_ISNOE
		w.IsNullOrEmptyValueHandler = CreateStandardValueHandler(w, &ck)
	}
}
func (w *ConditionValue) SetupLikeSearch(
	value interface{}, location string, co *ConditionOption) {
	if w.LikeSearchValueHandler == nil {
		var ck ConditionKey = CK_LS
		w.LikeSearchValueHandler = CreateVaryingValueHandler(w, &ck)
	}
	var v string = value.(string)
	log.InternalDebug("LikeSearch value :" + (*co).GenerateRealValue(v))
	w.LikeSearchLatestLocation = location + "." +
		w.LikeSearchValueHandler.SetValue((*co).GenerateRealValue(v))

}
func (w *ConditionValue) SetupNotLikeSearch(value interface{}, location string,
	co *ConditionOption) {
	if w.NotLikeSearchValueHandler == nil {
		var ck ConditionKey = CK_LS
		w.NotLikeSearchValueHandler = CreateVaryingValueHandler(w, &ck)
	}
	var v string = value.(string)
	log.InternalDebug("LikeSearch value :" + (*co).GenerateRealValue(v))
	w.NotLikeSearchLatestLocation = location + "." +
		w.NotLikeSearchValueHandler.SetValue((*co).GenerateRealValue(v))

}
func (w *ConditionValue) GetFixedValue(ck *ConditionKey) interface{} {
	if w.Fixed == nil {
		return nil
	}
	return w.Fixed[w.getFixedValueKey()][(*ck).GetConditionKeyS()]
}
func (w *ConditionValue) SetupInScope(value interface{}, location string) {
	if w.InScopeValueHandler == nil {
		var ck ConditionKey = CK_INS
		w.InScopeValueHandler = CreateVaryingValueHandler(w, &ck)
	}
	w.InScopeLatestLocation = location + "." +
		w.InScopeValueHandler.SetValue(value)
}

type VaryingValueHandler struct {
	ConditionValue *ConditionValue
	ConditionKey   *ConditionKey
}

func (w *VaryingValueHandler) SetValue(value interface{}) string {

	return w.ConditionValue.SetupVaryingValue(w.ConditionKey, value)
}
func (w *VaryingValueHandler) GetValue() interface{} {
	if w.ConditionValue.Fixed == nil {
		return nil
	}
	if w.ConditionValue.OrScopeQuery {
		//return w.getVaryingValue(w.ConditionKey)
		panic("Get Value for VaryingValue 未実装")
	} else {
		return w.ConditionValue.GetFixedValue(w.ConditionKey)
	}
	//return w.ConditionValue.fixed[(*w.ConditionKey).GetConditionKeyS()]
}

type StandardValueHandler struct {
	ConditionValue *ConditionValue
	ConditionKey   *ConditionKey
}

func (w *StandardValueHandler) SetValue(value interface{}) string {
	if w.ConditionValue.OrScopeQuery{
		return w.ConditionValue.SetupVaryingValue(w.ConditionKey, value)
	}

	return w.ConditionValue.SetupFixedValue(w.ConditionKey, value)
}
func (w *StandardValueHandler) GetValue() interface{} {
	if w.ConditionValue.Fixed == nil {
		return nil
	}
	if w.ConditionValue.OrScopeQuery {
		panic("Get Value for VaryingValue 未実装")
	} else {
		return w.ConditionValue.GetFixedValue(w.ConditionKey)
	}
	//return w.ConditionValue.fixed[(*w.ConditionKey).GetConditionKeyS()]
}
