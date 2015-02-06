package conditionKey

import (
	"container/list"
	//"fmt"
	"strings"
)

const (
	FIXED_KEY_QUERY = "query" //取り敢えず使わず
	C_EQ            = "equal"
	C_GT            = "greaterThan"
)

type ConditionKey interface {
	DoAddWhereClause(conditionList *list.List, columnRealName ColumnRealName, cvalue ConditionValue)
	GetConditionKeyS() string
	GetOperand() string
	DoSetupConditionValue(cvalue ConditionValue, value interface{}, location string)
}

type CK_EQ_T struct {
	ConditionKeyS string
	Operand       string
}

func (w *CK_EQ_T) DoAddWhereClause(conditionList *list.List, columnRealName ColumnRealName, cvalue ConditionValue) {

}
func (w *CK_EQ_T) GetConditionKeyS() string {
	return w.ConditionKeyS
}

func (w *CK_EQ_T) GetOperand() string {
	return w.Operand
}
func (w *CK_EQ_T) DoSetupConditionValue(cvalue ConditionValue, value interface{}, location string) {
	cvalue.SetupEqual(value, location)
}

var CK_EQ *CK_EQ_T

type CK_GT_T struct {
	ConditionKeyS string
	Operand       string
}

func (w *CK_GT_T) DoAddWhereClause(conditionList *list.List, columnRealName ColumnRealName, cvalue ConditionValue) {

}
func (w *CK_GT_T) GetConditionKeyS() string {
	return w.ConditionKeyS
}

func (w *CK_GT_T) GetOperand() string {
	return w.Operand
}

func (w *CK_GT_T) DoSetupConditionValue(cvalue ConditionValue, value interface{}, location string) {
	cvalue.SetupGreaterThan(value, location)
}

var CK_GT *CK_GT_T

type ColumnRealName struct {
	TableAliasName string
	ColumnSqlName  ColumnSqlName
}

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
	EqualLatestLocation       string
	GreaterThanLatestLocation string
	EqualValueHandler         *StandardValueHandler
	GreaterThanValueHandler   *StandardValueHandler
	FixedValueMap             map[string]interface{}
}

func (w *ConditionValue) SetupFixedValue(cond ConditionKey, value interface{}) string {
	if w.FixedValueMap == nil {
		w.FixedValueMap = make(map[string]interface{})
	}
	key := cond.GetConditionKeyS()
	w.FixedValueMap[key] = value
	return "fixed.query." + key
}

func CreateStandardValueHandler(cv *ConditionValue,  ck ConditionKey) *StandardValueHandler{
	sv := new(StandardValueHandler)
	sv.ConditionValue = cv
	sv.ConditionKey = ck
	return sv
}
func (w *ConditionValue) SetupEqual(value interface{}, location string) {
	if w.EqualValueHandler == nil {
		w.EqualValueHandler=CreateStandardValueHandler(w, CK_EQ)
	}
	w.EqualLatestLocation = location + "." + w.EqualValueHandler.SetValue(value)
}

func (w *ConditionValue) SetupGreaterThan(value interface{}, location string) {
	if w.GreaterThanValueHandler == nil {
		w.GreaterThanValueHandler = CreateStandardValueHandler(w, CK_GT)
	}
	w.GreaterThanLatestLocation = location + "." + w.GreaterThanValueHandler.SetValue(value)
}
func (w *ConditionValue) GetFixedValue(ck ConditionKey) interface{} {
	if w.FixedValueMap == nil {
		return nil
	}
	return w.FixedValueMap[ck.GetConditionKeyS()]
}

type StandardValueHandler struct {
	ConditionValue *ConditionValue
	ConditionKey   ConditionKey
}

func (w *StandardValueHandler) SetValue(value interface{}) string {
	//Or Query非対応
	return w.ConditionValue.SetupFixedValue(w.ConditionKey, value)
}
func (w *StandardValueHandler) GetValue() interface{} {
	if w.ConditionValue.FixedValueMap == nil {
		return nil
	}
	return w.ConditionValue.FixedValueMap[w.ConditionKey.GetConditionKeyS()]
}
