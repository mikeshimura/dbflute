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
	//	"strings"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"github.com/mikeshimura/dbflute/log"
	"reflect"
	"strconv"
	"time"
)

const (
	C_DISP_SQL_DEFAULT_DATE_FORMAT      = "2006-01-02"
	C_DISP_SQL_DEFAULT_TIME_FORMAT      = "15:04:05"
	C_DISP_SQL_DEFAULT_TIMESTAMP_FORMAT = "2006-01-02 15:04:05.000"
	DISP_SQL_NULL                       = "null"
)

var DISP_SQL_DEFAULT_DATE_FORMAT string = C_DISP_SQL_DEFAULT_DATE_FORMAT
var DISP_SQL_DEFAULT_TIME_FORMAT string = C_DISP_SQL_DEFAULT_TIME_FORMAT
var DISP_SQL_DEFAULT_TIMESTAMP_FORMAT string = C_DISP_SQL_DEFAULT_TIMESTAMP_FORMAT

type DisplaySqlBuilder struct {
	//format 対応未実装
}

func (d *DisplaySqlBuilder) BuildDisplaySql(sql string, args *List) string {
	ep := new(EmbeddingProcessor)
	return ep.Embed(sql, args)
}

type EmbeddingProcessor struct {
	processPointer           int
	loopIndex                int
	questionMarkIndex        int
	quotationScopeBeginIndex int
	quotationScopeEndIndex   int
	blockCommentBeginIndex   int
	blockCommentEndIndex     int
}

func (e *EmbeddingProcessor) Embed(sql string, args *List) string {
	if args == nil || args.Size() == 0 {
		return sql
	}
	sb := new(bytes.Buffer)
	sb.WriteString(Ln)
	for {
		e.questionMarkIndex = IndexAfter(sql, "?", e.processPointer)
		//next line actually do nothing
		//                setupQuestionMarkIndex(sql);
		if e.questionMarkIndex < 0 {
			e.processLastPart(sb, sql, args)
			break
		}
		if e.questionMarkIndex == 0 {
			e.processBindVariable(sb, sql, args)
			continue
		}
		e.setupBlockCommentIndex(sql)
		if e.hasBlockComment() {
			if e.isBeforeBlockComment() {
				e.processQuotationScope(sb, sql, args)
			} else { // in or after the block comment
				e.processBlockComment(sb, sql, args)
			}
		} else { // means no more block comment
			e.processQuotationScope(sb, sql, args)
		}

	}
	return sb.String()
}
func (e *EmbeddingProcessor) processBlockComment(sb *bytes.Buffer, sql string, args *List) {
	nextPointer := e.blockCommentEndIndex + 1
	beforeCommentEnd := sql[e.processPointer:nextPointer]
	sb.WriteString(beforeCommentEnd)
	e.processPointer = nextPointer
}
func (e *EmbeddingProcessor) processQuotationScope(sb *bytes.Buffer, sql string, args *List) {
	e.setupQuotationScopeIndex(sql)
	if e.isQuotationScopeOverBlockComment() {
		//                // means the quotation end in or after the comment (invalid scope)
		e.quotationScopeBeginIndex = -1
		e.quotationScopeEndIndex = -1
	}
	if e.hasQuotationScope() {
		if e.isInQuotationScope() {
			nextPointer := e.quotationScopeEndIndex + 1
			beforeScopeEnd := sql[e.processPointer:nextPointer]
			sb.WriteString(beforeScopeEnd)
			e.processPointer = nextPointer
		} else {
			e.processBindVariable(sb, sql, args)
		}
	} else {
		e.processBindVariable(sb, sql, args)
	}
}
func (e *EmbeddingProcessor) isInQuotationScope() bool {
	if !e.hasQuotationScope() {
		return false
	}
	return e.quotationScopeBeginIndex < e.questionMarkIndex && e.questionMarkIndex < e.quotationScopeEndIndex

}
func (e *EmbeddingProcessor) isQuotationScopeOverBlockComment() bool {
	return e.hasBlockComment() && e.hasQuotationScope() && e.quotationScopeEndIndex > e.blockCommentBeginIndex

}
func (e *EmbeddingProcessor) hasQuotationScope() bool {
	return e.quotationScopeBeginIndex >= 0 && e.quotationScopeEndIndex >= 0
}
func (e *EmbeddingProcessor) setupQuotationScopeIndex(sql string) {
	e.quotationScopeBeginIndex = IndexAfter(sql, "'", e.processPointer)
	e.quotationScopeEndIndex = IndexAfter(sql, "'", e.quotationScopeBeginIndex+1)
}
func (e *EmbeddingProcessor) isBeforeBlockComment() bool {
	if !e.hasBlockComment() {
		return false
	}
	return e.questionMarkIndex < e.blockCommentBeginIndex
}
func (e *EmbeddingProcessor) hasBlockComment() bool {
	return e.blockCommentBeginIndex >= 0 && e.blockCommentEndIndex >= 0
}
func (e *EmbeddingProcessor) setupBlockCommentIndex(sql string) {
	e.blockCommentBeginIndex = IndexAfter(sql, "/*", e.processPointer)
	e.blockCommentEndIndex = IndexAfter(sql, "*/", e.blockCommentBeginIndex+1)
}
func (e *EmbeddingProcessor) processLastPart(sb *bytes.Buffer, sql string, args *List) {
	sb.WriteString(sql[e.processPointer:])
	e.processPointer = 0
	e.loopIndex = 0
}
func (e *EmbeddingProcessor) processBindVariable(sb *bytes.Buffer, sql string, args *List) {
	//未実装
	//            assertArgumentSize(args, sql);
	beforeParameter := sql[e.processPointer:e.questionMarkIndex]
	bindVariableText := InterfaceToStringQuote(args.Get(e.loopIndex))
	sb.WriteString(beforeParameter + bindVariableText)
	e.processPointer = e.questionMarkIndex + 1
	e.loopIndex++
}

func Quote(str string) string {
	return "'" + str + "'"
}
func InterfaceToStringQuote(arg interface{}) string {
	switch arg.(type) {
	case string:
		return Quote(arg.(string))
	case *string:
		return Quote(*arg.(*string))
	case int:
		intv := arg.(int)
		return strconv.Itoa(intv)
	case int64:
		int64v := arg.(int64)
		return strconv.Itoa(int(int64v))
	case *int64:
		int64v := arg.(*int64)
		return strconv.Itoa(int(*int64v))
	case float64:
		float64v := arg.(float64)
		return strconv.FormatFloat(float64v, 'f', -1, 64)
	case float32:
		floatv := arg.(float32)
		return strconv.FormatFloat(float64(floatv), 'f', -1, 64)
	case *float64:
		float64v := arg.(*float64)
		return strconv.FormatFloat(*float64v, 'f', -1, 64)
	case bool:
		boolv := arg.(bool)
		return strconv.FormatBool(boolv)
	case *bool:
		boolv := arg.(*bool)
		return strconv.FormatBool(*boolv)
	case time.Time:
		tv := arg.(time.Time)
		return Quote(fmt.Sprint(tv.Format(DISP_SQL_DEFAULT_TIME_FORMAT)))
	case *time.Time:
		tv := arg.(*time.Time)
		return Quote(fmt.Sprint((*tv).Format(DISP_SQL_DEFAULT_TIME_FORMAT)))
	case Date:
		dv := arg.(Date)
		return Quote(fmt.Sprint(dv.Date.Format(DISP_SQL_DEFAULT_DATE_FORMAT)))
	case *Date:
		dv := arg.(*Date)
		return Quote(fmt.Sprint((*dv).Date.Format(DISP_SQL_DEFAULT_DATE_FORMAT)))
	case Timestamp:
		dv := arg.(Timestamp)
		return Quote(fmt.Sprint(dv.Timestamp.Format(DISP_SQL_DEFAULT_TIMESTAMP_FORMAT)))
	case *Timestamp:
		dv := arg.(*Timestamp)
		return Quote(fmt.Sprint((*dv).Timestamp.Format(DISP_SQL_DEFAULT_TIMESTAMP_FORMAT)))
	case Numeric:
		nv := arg.(Numeric)
		return nv.String()
	case *Numeric:
		nv := arg.(*Numeric)
		return (*nv).String()
	case MysqlDate:
		nv := arg.(MysqlDate)
		return Quote(nv.String())
	case *MysqlDate:
		nv := arg.(*MysqlDate)
		return Quote((*nv).String())
	case MysqlTime:
		nv := arg.(MysqlTime)
		return Quote(nv.String())
	case *MysqlTime:
		nv := arg.(*MysqlTime)
		return Quote((*nv).String())
	case MysqlTimestamp:
		nv := arg.(MysqlTimestamp)
		return Quote(nv.String())
	case *MysqlTimestamp:
		nv := arg.(*MysqlTimestamp)
		return Quote((*nv).String())
	case sql.NullString:
		nsv := arg.(sql.NullString)
		if nsv.Valid {
			return Quote(nsv.String)
		}
		return "null"
	case *sql.NullString:
		nsv := arg.(*sql.NullString)
		if nsv.Valid {
			return Quote((*nsv).String)
		}
		return "null"
		//	case NullString:
		//		nsv := arg.(NullString)
		//		if nsv.Valid {
		//			return Quote(nsv.String)
		//		}
		//		return "null"
		//	case *NullString:
		//		nsv := arg.(*NullString)
		//		if nsv.Valid {
		//			return Quote((*nsv).String)
		//		}
		//		return "null"
	case pq.NullTime:
		ntv := arg.(pq.NullTime)
		if ntv.Valid {
			return Quote(fmt.Sprint(ntv.Time.Format(DISP_SQL_DEFAULT_TIME_FORMAT)))
		}
		return "null"
	case *pq.NullTime:
		ntv := arg.(*pq.NullTime)
		if ntv.Valid {
			return Quote(fmt.Sprint((*ntv).Time.Format(DISP_SQL_DEFAULT_TIME_FORMAT)))
		}
		return "null"
	case NullDate:
		ndv := arg.(NullDate)
		if ndv.Valid {
			return Quote(fmt.Sprint(ndv.Date.Format(DISP_SQL_DEFAULT_DATE_FORMAT)))
		}
		return "null"
	case *NullDate:
		ndv := arg.(*NullDate)
		if ndv.Valid {
			return Quote(fmt.Sprint((*ndv).Date.Format(DISP_SQL_DEFAULT_DATE_FORMAT)))
		}
		return "null"
	case NullTimestamp:
		ntsv := arg.(NullTimestamp)
		if ntsv.Valid {
			return Quote(fmt.Sprint(ntsv.Timestamp.Format(DISP_SQL_DEFAULT_TIMESTAMP_FORMAT)))
		}
		return "null"
	case *NullTimestamp:
		ntsv := arg.(*NullTimestamp)
		if ntsv.Valid {
			return Quote(fmt.Sprint((*ntsv).Timestamp.Format(DISP_SQL_DEFAULT_TIMESTAMP_FORMAT)))
		}
		return "null"
	case NullNumeric:
		nnv := arg.(NullNumeric)
		if nnv.Valid {
			return nnv.String()
		}
		return "null"
	case *NullNumeric:
		nnv := arg.(*NullNumeric)
		if nnv.Valid {
			return (*nnv).String()
		}
		return "null"
	case MysqlNullDate:
		nnv := arg.(MysqlNullDate)
		if nnv.Valid {
			return nnv.String()
		}
		return "null"
	case *MysqlNullDate:
		nnv := arg.(*MysqlNullDate)
		if nnv.Valid {
			return (*nnv).String()
		}
		return "null"
	case MysqlNullTime:
		nnv := arg.(MysqlNullTime)
		if nnv.Valid {
			return nnv.String()
		}
		return "null"
	case *MysqlNullTime:
		nnv := arg.(*MysqlNullTime)
		if nnv.Valid {
			return (*nnv).String()
		}
		return "null"
	case MysqlNullTimestamp:
		nnv := arg.(MysqlNullTimestamp)
		if nnv.Valid {
			return nnv.String()
		}
		return "null"
	case *MysqlNullTimestamp:
		nnv := arg.(*MysqlNullTimestamp)
		if nnv.Valid {
			return (*nnv).String()
		}
		return "null"
	case sql.NullInt64:
		nnv := arg.(sql.NullInt64)
		if nnv.Valid {
			return strconv.Itoa(int(nnv.Int64))
		}
		return "null"
	case *sql.NullInt64:
		nnv := arg.(*sql.NullInt64)
		if nnv.Valid {
			return strconv.Itoa(int(nnv.Int64))
		}
		return "null"
	case sql.NullFloat64:
		nnv := arg.(sql.NullFloat64)
		if nnv.Valid {
			return strconv.FormatFloat(nnv.Float64, 'f', -1, 64)
		}
		return "null"
	case *sql.NullFloat64:
		nnv := arg.(*sql.NullFloat64)
		if nnv.Valid {
			return strconv.FormatFloat(nnv.Float64, 'f', -1, 64)
		}
		return "null"
	case sql.NullBool:
		nnv := arg.(sql.NullBool)
		if nnv.Valid {
			return strconv.FormatBool(nnv.Bool)
		}
		return "null"
	case *sql.NullBool:
		nnv := arg.(*sql.NullBool)
		if nnv.Valid {
			return strconv.FormatBool(nnv.Bool)
		}
		return "null"
	case []byte:
		return "[]byte"
	case *[]byte:
		return "[]byte"
	default:
		fmt.Printf("dump %v %T\n", arg, arg)
		panic("This type not supported :" + GetType(arg))

	}
	return ""
}

func IsNotNull(arg interface{}) bool {
	switch arg.(type) {
	case sql.NullString:
		nsv := arg.(sql.NullString)
		return nsv.Valid
	case *sql.NullString:
		nsv := arg.(*sql.NullString)
		return nsv.Valid
		//	case NullString:
		//		nsv := arg.(NullString)
		//		return nsv.Valid
		//	case *NullString:
		//		nsv := arg.(*NullString)
		//		return nsv.Valid
	case pq.NullTime:
		ntv := arg.(pq.NullTime)
		return ntv.Valid
	case *pq.NullTime:
		ntv := arg.(*pq.NullTime)
		return ntv.Valid
	case NullDate:
		ndv := arg.(NullDate)
		return ndv.Valid
	case *NullDate:
		ndv := arg.(*NullDate)
		return ndv.Valid
	case NullTimestamp:
		ntsv := arg.(NullTimestamp)
		return ntsv.Valid
	case *NullTimestamp:
		ntsv := arg.(*NullTimestamp)
		return ntsv.Valid
	case NullNumeric:
		nnv := arg.(NullNumeric)
		return nnv.Valid
	case *NullNumeric:
		nnv := arg.(*NullNumeric)
		return nnv.Valid
	case MysqlNullDate:
		nnv := arg.(MysqlNullDate)
		return nnv.Valid
	case *MysqlNullDate:
		nnv := arg.(*MysqlNullDate)
		return nnv.Valid
	case MysqlNullTime:
		nnv := arg.(MysqlNullTime)
		return nnv.Valid
	case *MysqlNullTime:
		nnv := arg.(*MysqlNullTime)
		return nnv.Valid
	case MysqlNullTimestamp:
		nnv := arg.(MysqlNullTimestamp)
		return nnv.Valid
	case *MysqlNullTimestamp:
		nnv := arg.(*MysqlNullTimestamp)
		return nnv.Valid
	case sql.NullInt64:
		nnv := arg.(sql.NullInt64)
		return nnv.Valid
	case *sql.NullInt64:
		nnv := arg.(*sql.NullInt64)
		return nnv.Valid
	case sql.NullFloat64:
		nnv := arg.(sql.NullFloat64)
		return nnv.Valid
	case *sql.NullFloat64:
		nnv := arg.(*sql.NullFloat64)
		return nnv.Valid
	case sql.NullBool:
		nnv := arg.(sql.NullBool)
		return nnv.Valid
	case *sql.NullBool:
		nnv := arg.(*sql.NullBool)
		return nnv.Valid
	default:
		return true

	}
}
func InterfaceToString(arg interface{}) string {
	switch arg.(type) {
	case string:
		return arg.(string)
	case *string:
		return *arg.(*string)
	case int64:
		int64v := arg.(int64)
		return strconv.Itoa(int(int64v))
	case *int64:
		int64v := arg.(*int64)
		return strconv.Itoa(int(*int64v))
	case float64:
		float64v := arg.(float64)
		return strconv.FormatFloat(float64v, 'f', -1, 64)
	case *float64:
		float64v := arg.(*float64)
		return strconv.FormatFloat(*float64v, 'f', -1, 64)
	case bool:
		boolv := arg.(bool)
		return strconv.FormatBool(boolv)
	case *bool:
		boolv := arg.(*bool)
		return strconv.FormatBool(*boolv)
	case time.Time:
		tv := arg.(time.Time)
		return fmt.Sprint(tv.Format(DISP_SQL_DEFAULT_TIME_FORMAT))
	case *time.Time:
		tv := arg.(*time.Time)
		return fmt.Sprint((*tv).Format(DISP_SQL_DEFAULT_TIME_FORMAT))
	case Date:
		dv := arg.(Date)
		return fmt.Sprint(dv.Date.Format(DISP_SQL_DEFAULT_DATE_FORMAT))
	case *Date:
		dv := arg.(*Date)
		return fmt.Sprint((*dv).Date.Format(DISP_SQL_DEFAULT_DATE_FORMAT))
	case Timestamp:
		dv := arg.(Timestamp)
		return fmt.Sprint(dv.Timestamp.Format(DISP_SQL_DEFAULT_TIMESTAMP_FORMAT))
	case *Timestamp:
		dv := arg.(*Timestamp)
		return fmt.Sprint((*dv).Timestamp.Format(DISP_SQL_DEFAULT_TIMESTAMP_FORMAT))
	case Numeric:
		nv := arg.(Numeric)
		return nv.String()
	case *Numeric:
		nv := arg.(*Numeric)
		return (*nv).String()
	case MysqlDate:
		nv := arg.(MysqlDate)
		return nv.String()
	case *MysqlDate:
		nv := arg.(*MysqlDate)
		return (*nv).String()
	case MysqlTime:
		nv := arg.(MysqlTime)
		return nv.String()
	case *MysqlTime:
		nv := arg.(*MysqlTime)
		return (*nv).String()
	case MysqlTimestamp:
		nv := arg.(MysqlTimestamp)
		return nv.String()
	case *MysqlTimestamp:
		nv := arg.(*MysqlTimestamp)
		return (*nv).String()
	case sql.NullString:
		nsv := arg.(sql.NullString)
		if nsv.Valid {
			return nsv.String
		}
		return "null"
	case *sql.NullString:
		nsv := arg.(*sql.NullString)
		if nsv.Valid {
			return (*nsv).String
		}
		return "null"
		//	case NullString:
		//		nsv := arg.(NullString)
		//		if nsv.Valid {
		//			return nsv.String
		//		}
		//		return "null"
		//	case *NullString:
		//		nsv := arg.(*NullString)
		//		if nsv.Valid {
		//			return (*nsv).String
		//		}
		//		return "null"
	case pq.NullTime:
		ntv := arg.(pq.NullTime)
		if ntv.Valid {
			return fmt.Sprint(ntv.Time.Format(DISP_SQL_DEFAULT_TIME_FORMAT))
		}
		return "null"
	case *pq.NullTime:
		ntv := arg.(*pq.NullTime)
		if ntv.Valid {
			return fmt.Sprint((*ntv).Time.Format(DISP_SQL_DEFAULT_TIME_FORMAT))
		}
		return "null"
	case NullDate:
		ndv := arg.(NullDate)
		if ndv.Valid {
			return fmt.Sprint(ndv.Date.Format(DISP_SQL_DEFAULT_DATE_FORMAT))
		}
		return "null"
	case *NullDate:
		ndv := arg.(*NullDate)
		if ndv.Valid {
			return fmt.Sprint((*ndv).Date.Format(DISP_SQL_DEFAULT_DATE_FORMAT))
		}
		return "null"
	case NullTimestamp:
		ntsv := arg.(NullTimestamp)
		if ntsv.Valid {
			return fmt.Sprint(ntsv.Timestamp.Format(DISP_SQL_DEFAULT_TIMESTAMP_FORMAT))
		}
		return "null"
	case *NullTimestamp:
		ntsv := arg.(*NullTimestamp)
		if ntsv.Valid {
			return fmt.Sprint((*ntsv).Timestamp.Format(DISP_SQL_DEFAULT_TIMESTAMP_FORMAT))
		}
		return "null"
	case NullNumeric:
		nnv := arg.(NullNumeric)
		if nnv.Valid {
			return nnv.String()
		}
		return "null"
	case *NullNumeric:
		nnv := arg.(*NullNumeric)
		if nnv.Valid {
			return (*nnv).String()
		}
		return "null"
	case MysqlNullDate:
		nnv := arg.(MysqlNullDate)
		if nnv.Valid {
			return nnv.String()
		}
		return "null"
	case *MysqlNullDate:
		nnv := arg.(*MysqlNullDate)
		if nnv.Valid {
			return (*nnv).String()
		}
		return "null"
	case MysqlNullTime:
		nnv := arg.(MysqlNullTime)
		if nnv.Valid {
			return nnv.String()
		}
		return "null"
	case *MysqlNullTime:
		nnv := arg.(*MysqlNullTime)
		if nnv.Valid {
			return (*nnv).String()
		}
		return "null"
	case MysqlNullTimestamp:
		nnv := arg.(MysqlNullTimestamp)
		if nnv.Valid {
			return nnv.String()
		}
		return "null"
	case *MysqlNullTimestamp:
		nnv := arg.(*MysqlNullTimestamp)
		if nnv.Valid {
			return (*nnv).String()
		}
		return "null"
	case sql.NullInt64:
		nnv := arg.(sql.NullInt64)
		if nnv.Valid {
			return strconv.Itoa(int(nnv.Int64))
		}
		return "null"
	case *sql.NullInt64:
		nnv := arg.(*sql.NullInt64)
		if nnv.Valid {
			return strconv.Itoa(int(nnv.Int64))
		}
		return "null"
	case sql.NullFloat64:
		nnv := arg.(sql.NullFloat64)
		if nnv.Valid {
			return strconv.FormatFloat(nnv.Float64, 'f', -1, 64)
		}
		return "null"
	case *sql.NullFloat64:
		nnv := arg.(*sql.NullFloat64)
		if nnv.Valid {
			return strconv.FormatFloat(nnv.Float64, 'f', -1, 64)
		}
		return "null"
	case sql.NullBool:
		nnv := arg.(sql.NullBool)
		if nnv.Valid {
			return strconv.FormatBool(nnv.Bool)
		}
		return "null"
	case *sql.NullBool:
		nnv := arg.(*sql.NullBool)
		if nnv.Valid {
			return strconv.FormatBool(nnv.Bool)
		}
		return "null"
	case []byte:
		return "[]byte"
	case *[]byte:
		return "[]byte"
	default:
		fmt.Printf("dump %v %T\n", arg, arg)
		panic("This type not supported :" + GetType(arg))

	}
	return ""
}

type BehaviorResultBuilder struct {
}

func (b *BehaviorResultBuilder) buildResultExp(elapse time.Duration, entityType string, ret interface{}) string {
	//	        try {
	//            return doBuildResultExp(retType, ret, before, after);
	//        } catch (RuntimeException e) {
	//            String msg = "Failed to build result expression for logging: retType=" + retType + " ret=" + ret;
	//            throw new IllegalStateException(msg, e);
	//        }
	return b.doBuildResultExp(elapse, entityType, ret)
}
func (b *BehaviorResultBuilder) getElapse(elapse time.Duration) string {
	m := elapse / time.Minute
	rem2 := elapse - m*time.Minute
	s := rem2 / time.Second
	rem3 := rem2 - s*time.Second
	ms := rem3 / time.Millisecond
	return fmt.Sprintf("%02dm%02ds%03dms ", m, s, ms)

}
func (b *BehaviorResultBuilder) doBuildResultExp(elapse time.Duration, entityType string, ret interface{}) string {
	var resultExp string
	//        final String prefix = "===========/ [" + DfTraceViewUtil.convertToPerformanceView(after - before) + " ";
	prefix := "===========/ [" + b.getElapse(elapse)
	stype := GetType(ret)
	log.InternalDebug("return Type:" + stype)
	log.InternalDebug("entity Type:" + entityType)
	if stype == "*df.ListResultBean" {
		var lrb *ListResultBean = ret.(*ListResultBean)
		if entityType == "D_Int64" {
			if lrb.AllRecordCount == 0 {
				resultExp = prefix + "(0)]"
			} else if lrb.AllRecordCount == 1 {
				var cnt int = int((lrb.List.Get(0)).(*D_Int64).GetValue())
				resultExp = prefix + "result=" + strconv.Itoa(cnt) + "]"
			}
		} else {
			if lrb.AllRecordCount == 0 {
				resultExp = prefix + "(0)]"
			} else if lrb.AllRecordCount == 1 {
				resultExp = prefix + "(1) result=" + b.buildEntityExp(lrb.List.Get(0), entityType) + "]"
			} else {
				resultExp = prefix + "(" + strconv.Itoa(lrb.AllRecordCount) + ") first=" + b.buildEntityExp(lrb.List.Get(0), entityType) + "]"
			}
		}

	}
	if stype == "int64" {
		var iv int64 = ret.(int64)
		resultExp = prefix + "result=" + strconv.Itoa(int(iv)) + "]"
	}
	//        if (List.class.isAssignableFrom(retType)) {
	//            if (ret == null) {
	//                resultExp = prefix + "(null)]";
	//            } else {
	//                final List<?> ls = (java.util.List<?>) ret;
	//                if (ls.isEmpty()) {
	//                    resultExp = prefix + "(0)]";
	//                } else if (ls.size() == 1) {
	//                    resultExp = prefix + "(1) result=" + buildEntityExp(ls.get(0)) + "]";
	//                } else {
	//                    resultExp = prefix + "(" + ls.size() + ") first=" + buildEntityExp(ls.get(0)) + "]";
	//                }
	//            }
	//        } else if (Entity.class.isAssignableFrom(retType)) {
	//            if (ret == null) {
	//                resultExp = prefix + "(null)" + "]";
	//            } else {
	//                final Entity entity = (Entity) ret;
	//                resultExp = prefix + "(1) result=" + buildEntityExp(entity) + "]";
	//            }
	//        } else if (int[].class.isAssignableFrom(retType)) {
	//            if (ret == null) { // basically not come here
	//                resultExp = prefix + "(null)" + "]";
	//            } else {
	//                final int[] resultArray = (int[]) ret;
	//                resultExp = buildBatchUpdateResultExp(prefix, resultArray);
	//            }
	//        } else {
	//            resultExp = prefix + "result=" + ret + "]";
	//        }
	//        return resultExp;
	return resultExp
}

func (b *BehaviorResultBuilder) buildEntityExp(data interface{}, entityType string) string {
	meta := DBMetaProvider_I.TableDbNameInstanceMap[entityType]
	if meta == nil {
		return ""
	}
	sb := new(bytes.Buffer)
	sflag := true
	v := reflect.ValueOf(data)
	list := (*meta).GetColumnInfoList()
	for _, colInf := range list.data {
		var ci *ColumnInfo = colInf.(*ColumnInfo)
		res := v.MethodByName("Get" + InitCap(ci.PropertyName)).Call([]reflect.Value{})
		if sflag == false {
			sb.WriteString(", ")
		}

		sb.WriteString(InterfaceToString(res[0].Interface()))
		//fmt.Printf("%v\n", res[0].Interface())
		sflag = false
	}
	return "{" + sb.String() + "}"
}
