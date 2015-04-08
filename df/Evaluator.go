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
	"fmt"
	"github.com/mikeshimura/dbflute/log"
	"reflect"
	"strconv"
	"strings"
)

const (
	Evaluator_AND           = " && "
	Evaluator_OR            = " || "
	Evaluator_EQUAL         = " == "
	Evaluator_NOT_EQUAL     = " != "
	Evaluator_GREATER_THAN  = " > "
	Evaluator_LESS_THAN     = " < "
	Evaluator_GREATER_EQUAL = " >= "
	Evaluator_LESS_EQUAL    = " <= "
	Evaluator_BOOLEAN_NOT   = "!"
	Evaluator_METHOD_SUFFIX = "()"
)

type IfCommentEvaluator struct {
	finder       *ParameterFinder
	expression   string
	specifiedSql string
	loopInfo     *LoopInfo
	ctx          *CommandContext
}

func (i *IfCommentEvaluator) evaluate() bool {
	i.assertExpression()
	if strings.Contains(i.expression, Evaluator_AND) {
		splitList := strings.Split(i.expression, Evaluator_AND)
		for _, booleanClause := range splitList {
			result := i.evaluateBooleanClause(booleanClause)
			if !result {
				return false
			}
		}
		return true
	} else if strings.Contains(i.expression, Evaluator_OR) {
		splitList := strings.Split(i.expression, Evaluator_OR)
		for _, booleanClause := range splitList {
			result := i.evaluateBooleanClause(booleanClause)
			if result {
				return true
			}
		}
		return false
	} else {
		return i.evaluateBooleanClause(i.expression)
	}
}
func (i *IfCommentEvaluator) assertExpression() {
	if len(strings.TrimSpace(i.expression)) == 0 {
		panic("IfCommentEmpty")
	}

	//            filtered := strings.Replace(i.expression, "()", "",-1);
	//            filtered = Srl.replace(filtered, ".get(", "");
	if strings.Contains(i.expression, "(") {
		panic("IfCommentUnsupported")
	}
	if strings.Contains(i.expression, Evaluator_AND) && strings.Contains(i.expression, Evaluator_OR) {
		panic("IfCommentUnsupportedExpression")
	}
	if strings.Contains(i.expression, " = ") || strings.Contains(i.expression, " <> ") {
		panic("IfCommentUnsupportedExpression")
	}
	if strings.Contains(i.expression, "\"") {
		panic("IfCommentUnsupportedExpression")
	}

}
func (i *IfCommentEvaluator) evaluateBooleanClause(booleanClause string) bool {
	if strings.Contains(booleanClause, Evaluator_EQUAL) {
		return i.evaluateCompareClause(booleanClause, Evaluator_EQUAL)
	}
	if strings.Contains(booleanClause, Evaluator_NOT_EQUAL) {
		return i.evaluateCompareClause(booleanClause, Evaluator_NOT_EQUAL)
	}
	if strings.Contains(booleanClause, Evaluator_GREATER_THAN) {
		return i.evaluateCompareClause(booleanClause, Evaluator_GREATER_THAN)
	}
	if strings.Contains(booleanClause, Evaluator_LESS_THAN) {
		return i.evaluateCompareClause(booleanClause, Evaluator_LESS_THAN)
	}
	if strings.Contains(booleanClause, Evaluator_GREATER_EQUAL) {
		return i.evaluateCompareClause(booleanClause, Evaluator_GREATER_EQUAL)
	}
	if strings.Contains(booleanClause, Evaluator_LESS_EQUAL) {
		return i.evaluateCompareClause(booleanClause, Evaluator_LESS_EQUAL)
	}
	panic("Operator missing")
	return true
}
func (i *IfCommentEvaluator) evaluateCompareClause(booleanClause string, operand string) bool {
	pos := strings.Index(booleanClause, operand)
	left := strings.TrimSpace(booleanClause[0:pos])
	right := strings.TrimSpace(booleanClause[pos+len(operand):])
	leftResult := i.evaluateComparePiece(left, nil)
	rightResult := i.evaluateComparePiece(right, nil)
	leftResult = leftResult
	rightResult = rightResult
	//fmt.Printf("leftResult %v %T\n", leftResult, leftResult)
	//fmt.Printf("rightResult %v %T\n", rightResult, rightResult)
	evaluator := new(OperandEvaluator)
	return evaluator.evaluate(operand, leftResult, rightResult)
}

//func (i *IfCommentEvaluator) getEvaluator(operand string) *OperandEvaluator {
//	var evaluator OperandEvaluator
//	if operand == Evaluator_EQUAL {
//		eval := new(OperandEvaluatorEqual)
//		evaluator = eval
//	}
//	if operand == Evaluator_NOT_EQUAL {
//		eval := new(OperandEvaluatorNotEqual)
//		evaluator = eval
//	}
//	return &evaluator
//}
func (i *IfCommentEvaluator) evaluateComparePiece(piece string, leftResult interface{}) interface{} {
	piece = strings.TrimSpace(piece)
	lower := strings.ToLower(piece)
	if strings.Index(piece, "pmb") != 0 {
		if lower == "null" {
			var n sql.NullString
			return n
		}
		if lower == "true" {
			return true
		}
		if lower == "false" {
			return false
		}
		quote := "'"
		qlen := len(quote)
		if strings.Index(piece, quote) == 0 && strings.LastIndex(piece, quote) == len(piece)-1 {
			return piece[qlen : len(piece)-qlen]
		}
		big, err := strconv.ParseFloat(piece, 64)
		if err != nil {
			panic("can't parse :" + piece)
		}
		return big
	}
	propertyList := new(StringList)
	preProperty := i.setupPropertyList(piece, propertyList)
	baseObject := i.findBaseObject(preProperty)
	//fmt.Printf("baseObject %v %T\n", baseObject, baseObject)
	for _, property := range propertyList.data {
		baseObject = i.processOneProperty(baseObject, preProperty, property)
	}
	//fmt.Printf("baseObject %v %T\n", baseObject, baseObject)
	return baseObject
}
func (i *IfCommentEvaluator) processOneProperty(baseObject interface{}, preProperty string, property string) interface{} {
	if baseObject == nil {
		panic("Null Oblect found:" + preProperty)
	}
	v := reflect.ValueOf(baseObject).Elem()
	log.InternalDebug(fmt.Sprintf("pmb new v %v %T\n", v, v))
	newv := v.FieldByName(InitCap(property))
	log.InternalDebug(fmt.Sprintf("newv  %v \n", newv))
	if newv.IsValid() == false {
		//どちらが良いか要検討
		//return "string", ""
		return nil
	}
	newValue := newv.Interface()
	return newValue
}
func (i *IfCommentEvaluator) findBaseObject(firstName string) interface{} {
	//loop not implemented
	return (*i.finder).find(i.ctx, firstName)
}
func (i *IfCommentEvaluator) setupPropertyList(piece string, propertyList *StringList) string {
	splitList := strings.Split(piece, ".")
	firstName := ""
	for i := 0; i < len(splitList); i++ {
		token := splitList[i]
		if i == 0 {
			firstName = token
			continue
		}
		propertyList.Add(token)
	}
	return firstName
}

type OperandEvaluator struct {
}

func (o *OperandEvaluator) getFloat(arg interface{}) (float64, bool) {
	switch arg.(type) {

	case string:
		res, _ := strconv.ParseFloat(arg.(string), 64)
		return res, false
	case sql.NullString:
		ns := arg.(sql.NullString)
		if ns.Valid {
			res, _ := strconv.ParseFloat(ns.String, 64)
			return res, false
		} else {
			return 0.0, true
		}
	case int64:
		var f float64 = (float64)(arg.(int64))
		return f, false
	case sql.NullInt64:
		ns := arg.(sql.NullInt64)
		if ns.Valid {
			var f float64 = (float64)(ns.Int64)
			return f, false
		} else {
			return 0.0, true
		}
	case float64:
		return arg.(float64), false
	case sql.NullFloat64:
		ns := arg.(sql.NullFloat64)
		if ns.Valid {
			return ns.Float64, false
		} else {
			return 0.0, true
		}

	case Numeric:
		num := arg.(Numeric)
		res, _ := strconv.ParseFloat(num.String(), 64)
		return res, false

	case NullNumeric:
		num := arg.(NullNumeric)
		if num.Valid {
			res, _ := strconv.ParseFloat(num.String(), 64)
			return res, false
		} else {
			return 0.0, true
		}
	}
	return 0.0, true
}

func (o *OperandEvaluator) getString(arg interface{}) (string, bool) {
	switch arg.(type) {
	case string:
		return arg.(string), false
	case sql.NullString:
		ns := arg.(sql.NullString)
		if ns.Valid {
			return ns.String, false
		} else {
			return "", true
		}
	}
	return "", true
}

func (o *OperandEvaluator) evaluate(operand string, leftResult interface{}, rightResult interface{}) bool {
	if operand == Evaluator_EQUAL {
		if IsNotNull(rightResult) == false {
			if IsNotNull(leftResult) == false {
				return true
			} else {
				return false
			}
		}
		if IsNotNull(leftResult) == false {
			if IsNotNull(rightResult) == false {
				return true
			} else {
				return false
			}
		}
	}
	if operand == Evaluator_NOT_EQUAL {
		if IsNotNull(rightResult) == false {
			if IsNotNull(leftResult) == false {
				return false
			} else {
				return true
			}
		}
		if IsNotNull(leftResult) == false {
			if IsNotNull(rightResult) == false {
				return false
			} else {
				return true
			}
		}
	}
	leftType := GetType(leftResult)

	if leftType == "string" || leftType == "sql.NullString" {
		leftString, nulll := o.getString(leftResult)
		rightString, nullr := o.getString(rightResult)
		if nulll || nullr {
			panic("system logic error")
		}
		switch operand {
		case Evaluator_EQUAL:
			return leftString == rightString
		case Evaluator_NOT_EQUAL:
			return leftString != rightString
		case Evaluator_GREATER_THAN:
			return leftString > rightString
		case Evaluator_LESS_THAN:
			return leftString < rightString
		case Evaluator_GREATER_EQUAL:
			return leftString >= rightString
		case Evaluator_LESS_EQUAL:
			return leftString <= rightString
		}
	}
	if leftType == "int64" || leftType == "sql.NullInt64" || leftType == "float64" || leftType == "sql.NullFloat64" || leftType == "df.Numeric" || leftType == "df.NullNumeric" {
		log.InternalDebug(fmt.Sprintf("leftResult %v %T rightResult %v %T\n", leftResult, leftResult, rightResult, rightResult))
		leftFloat, nulll2 := o.getFloat(leftResult)
		rightFloat, nullr2 := o.getFloat(rightResult)
		if nulll2 || nullr2 {
			panic("system logic error")
		}
		switch operand {
		case Evaluator_EQUAL:
			return leftFloat == rightFloat
		case Evaluator_NOT_EQUAL:
			return leftFloat != rightFloat
		case Evaluator_GREATER_THAN:
			return leftFloat > rightFloat
		case Evaluator_LESS_THAN:
			return leftFloat < rightFloat
		case Evaluator_GREATER_EQUAL:
			return leftFloat >= rightFloat
		case Evaluator_LESS_EQUAL:
			return leftFloat <= rightFloat
		}
	}
	panic("not supported this type yet:" + leftType)
	return false
}

//type OperandEvaluatorNotEqual struct {
//	BaseOperandEvaluator
//}
//
//func (o *OperandEvaluatorNotEqual) evaluate(leftResult interface{}, rightResult interface{}) bool {
//	if IsNotNull(rightResult) == false {
//		if IsNotNull(leftResult) == false {
//			return false
//		} else {
//			return true
//		}
//	}
//	if IsNotNull(leftResult) == false {
//		if IsNotNull(rightResult) == false {
//			return false
//		} else {
//			return true
//		}
//	}
//	leftType := GetType(leftResult)
//
//	if leftType == "string" || leftType == "sql.NullString" {
//		leftString, nulll := o.getString(leftResult)
//		rightString, nullr := o.getString(rightResult)
//		if nulll || nullr {
//			panic("logic error")
//		}
//		return leftString != rightString
//	}
//	if leftType == "int64" || leftType == "sql.NullInt64" || leftType == "float64" || leftType == "sql.NullFloat64" || leftType == "df.Numeric" || leftType == "df.NullNumeric" {
//		log.InternalDebug(fmt.Sprintf("leftResult %v %T rightResult %v %T\n", leftResult, leftResult, rightResult, rightResult))
//		leftBig, nulll2 := o.getBig(leftResult)
//		rightBig, nullr2 := o.getBig(rightResult)
//		fmt.Printf("leftBig %v %T\n", leftBig, leftBig)
//		fmt.Printf("rightBig %v %T\n", rightBig, rightBig)
//		if nulll2 || nullr2 {
//			panic("logic error")
//		}
//		return leftBig != rightBig
//	}
//	panic("not supported this type yet:" + leftType)
//	return false
//}

type ParameterFinder interface {
	find(ctx *CommandContext, name string) interface{}
}

type ParameterFinderArg struct {
}

func (p *ParameterFinderArg) find(ctx *CommandContext, name string) interface{} {
	return (*ctx).getArg(name)
}
