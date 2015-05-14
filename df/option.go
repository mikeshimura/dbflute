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
	//"fmt"
	"strings"
)

const (
	LIKE_PREFIX  = "prefix"
	LIKE_SUFFIX  = "suffix"
	LIKE_CONTAIN = "contain"
)

type InsertOption struct {
}
//tentative implementation
func (i * InsertOption)isPrimaryKeyIdentityDisabled()bool{
	return false
}
type UpdateOption struct {
}
type DeleteOption struct {
}

type ConditionOption interface {
	HasCompoundColumn() bool
	GenerateRealValue(s string) string
	GetRearOption() string
}

type DummyOption struct {
}

func (d *DummyOption) HasCompoundColumn() bool {
	return false
}
func (d *DummyOption) GenerateRealValue(s string) string {
	return ""
}
func (d *DummyOption) GetRearOption() string {
	return ""
}

type LikeSearchOption struct {
	SimpleStringOption SimpleStringOption
	like                   string
	escapechar             string
	asOrSplit              bool
	originalWildCardList   *StringList
	compoundColumnList     *List
	compoundColumnSizeList *List
	stringConnector        *StringConnector
}

func (l *LikeSearchOption) HasCompoundColumn() bool {
	return l.compoundColumnList != nil && l.compoundColumnList.Size() > 0
}

func (l *LikeSearchOption) LikePrefix() *LikeSearchOption{
	l.like = LIKE_PREFIX
	l.EscapeByPipeLine()
	return l
}

func (l *LikeSearchOption) LikeSuffix() *LikeSearchOption{
	l.like = LIKE_SUFFIX
	l.EscapeByPipeLine()
	return l
}

func (l *LikeSearchOption) LikeContain() *LikeSearchOption{
	l.like = LIKE_CONTAIN
	l.EscapeByPipeLine()
	return l
}
func (l *LikeSearchOption) Escape() *LikeSearchOption {
	l.EscapeByPipeLine()
	return l
}

func (l *LikeSearchOption) EscapeByPipeLine() *LikeSearchOption {
	l.escapechar = "|"
	return l
}
func (l *LikeSearchOption) EscapeByAtMark() *LikeSearchOption {
	l.escapechar = "@"
	return l
}
func (l *LikeSearchOption) EscapeBySlash() *LikeSearchOption {
	l.escapechar = "/"
	return l
}

func (l *LikeSearchOption) EscapeByBackSlash() *LikeSearchOption {
	l.escapechar = "\\"
	return l
}
func (l *LikeSearchOption) NotEscape() *LikeSearchOption {
	l.escapechar = ""
	return l
}
func (l *LikeSearchOption) AsOrSplit() *LikeSearchOption {
	l.asOrSplit = true
	return l
}
func (l *LikeSearchOption) GetRearOption() string {
	if len(strings.Trim(l.escapechar, " ")) == 0 {
		return ""
	}
	return " escape '" + l.escapechar + "'"
}

func (l *LikeSearchOption) GenerateRealValue(value string) string {
	if l.originalWildCardList == nil {
		l.originalWildCardList = new(StringList)
	}
	// escape wild-cards  change from Java version for add l.like != ""
	if len(strings.Trim(l.escapechar, " ")) > 0 && l.like != "" {
		tmp := strings.Replace(value, l.escapechar, l.escapechar+l.escapechar, -1)
		// basic wild-cards
		tmp = l.filterEscape(tmp, "%")
		tmp = l.filterEscape(tmp, "_")

		if l.originalWildCardList.Size() > 0 {
			for _, wildCard := range l.originalWildCardList.data {
				tmp = l.filterEscape(tmp, wildCard)
			}
		}

		value = tmp
	}
	wildCard := "%"
	if len(strings.Trim(l.like, " ")) == 0 {
		return value
	} else if l.like == LIKE_PREFIX {
		return value + wildCard
	} else if l.like == LIKE_SUFFIX {
		return wildCard + value
	} else if l.like == LIKE_CONTAIN {
		return wildCard + value + wildCard
	} else {
		msg := "The like was wrong string: " + l.like
		panic(msg)
	}

}
func (l *LikeSearchOption) filterEscape(target string, wildCard string) string {
	return strings.Replace(target, wildCard, l.escapechar+wildCard, -1)
}
func (s *LikeSearchOption) SplitByBlank() *LikeSearchOption{
	s.SimpleStringOption.splitByBlank()
	return s
}
func (s *LikeSearchOption) SplitBySpace() *LikeSearchOption {
	s.SimpleStringOption.splitBySpace()
	return s
}
func (s *LikeSearchOption) SplitBySpaceContainsDoubleByte() *LikeSearchOption {
	s.SimpleStringOption.splitBySpaceContainsDoubleByte()
	return s
}
func (s *LikeSearchOption) SplitByPipeLine() *LikeSearchOption {
	s.SimpleStringOption.splitByPipeLine()
	return s
}
func (s *LikeSearchOption) LimitSplit(splitLimitCount int) *LikeSearchOption {
	s.SimpleStringOption.limitSplit(splitLimitCount)
	return s
}
func (s *LikeSearchOption) isSplit() bool {
	//fmt.Printf("delimiter %v\n", s.delimiter)
	return s.SimpleStringOption.delimiter != ""
}
func (s *LikeSearchOption) GenerateSplitValueArray(value string) []string {
	return s.SimpleStringOption.GenerateSplitValueArray(value)
}
type StringConnector interface {
	connect(element ...string) string
}
type StandardStringConnector struct {
}

func (s *StandardStringConnector) connect(element ...string) string {
	sb := new(bytes.Buffer)
	for i, ele := range element {
		if i > 0 {
			sb.WriteString(" || ")
		}
		sb.WriteString(ele)
	}
	return sb.String()
}
type PlusStringConnector struct {
}

func (s *PlusStringConnector) connect(element ...string) string {
	sb := new(bytes.Buffer)
	for i, ele := range element {
		if i > 0 {
			sb.WriteString(" + ")
		}
		sb.WriteString(ele)
	}
	return sb.String()
}
type MysqlStringConnector struct {
}

func (s *MysqlStringConnector) connect(element ...string) string {
	sb := new(bytes.Buffer)
	sb.WriteString("concat(")
	for i, ele := range element {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(ele)
	}
	sb.WriteString(")")
	return sb.String()
}

type SimpleStringOption struct {
	SplitOptionParts
}

func (s *SimpleStringOption) SplitByBlank() *SimpleStringOption {
	s.splitByBlank()
	return s
}
func (s *SimpleStringOption) SplitBySpace() *SimpleStringOption {
	s.splitBySpace()
	return s
}
func (s *SimpleStringOption) SplitBySpaceContainsDoubleByte() *SimpleStringOption {
	s.splitBySpaceContainsDoubleByte()
	return s
}
func (s *SimpleStringOption) SplitByPipeLine() *SimpleStringOption {
	s.splitByPipeLine()
	return s
}
func (s *SimpleStringOption) LimitSplit(splitLimitCount int) *SimpleStringOption {
	s.limitSplit(splitLimitCount)
	return s
}

type SplitOptionParts struct {
	delimiter        string
	subDelimiterList *StringList
	splitLimitCount  int
}

func (s *SplitOptionParts) isSplit() bool {
	//fmt.Printf("delimiter %v\n", s.delimiter)
	return s.delimiter != ""
}
func (s *SplitOptionParts) splitByBlank() {
	s.delimiter = " "
	s.subDelimiterList = new(StringList)
	s.subDelimiterList.Add("\u3000")
	s.subDelimiterList.Add("\t")
	s.subDelimiterList.Add("\r")
	s.subDelimiterList.Add("\n")
}
func (s *SplitOptionParts) splitBySpace() {
	s.delimiter = " "
	s.subDelimiterList = new(StringList)
}
func (s *SplitOptionParts) splitBySpaceContainsDoubleByte() {
	s.splitBySpace()
	s.subDelimiterList = new(StringList)
	s.subDelimiterList.Add("\u3000")
}
func (s *SplitOptionParts) splitByPipeLine() {
	s.delimiter = "|"
	s.subDelimiterList = new(StringList)
}
func (s *SplitOptionParts) limitSplit(splitLimitCount int) {
	s.splitLimitCount = splitLimitCount
}
func (s *SplitOptionParts) GenerateSplitValueArray(value string) []string {
	sl := new(StringList)
	if strings.Index(value, s.delimiter) > -1 {
		sl.Add(s.delimiter)
	}
	for _, delim := range s.subDelimiterList.data {
		if strings.Index(value, delim) > -1 {
			sl.Add(delim)
		}
	}
	pos := -1
	result := new(StringList)
	for len(value) > 0 {
		for _, delim := range sl.data {
			tpos := strings.Index(value, delim)
			if tpos > -1 {
				if pos == -1 {
					pos = tpos
				} else {
					if tpos < pos {
						pos = tpos
					}
				}
			}
		}
		if pos > -1 {
			result.Add(value[0:pos])
			value = value[pos+1:]
			pos = -1
		} else {
			break
		}
		if result.Size() == s.splitLimitCount {
			break
		}
	}
	if (s.splitLimitCount == 0 || result.Size() < s.splitLimitCount) && (len(value) > 0) {
		result.Add(value)
	}
	var sres []string = make([]string, result.Size())
	for i, str := range result.data {
		sres[i] = str
	}
	return sres
}

type RangeOfOption struct {
	DummyOption
	greaterThan bool
	lessThan    bool
	orIsNull    bool
}

func (r *RangeOfOption) GreaterThan() *RangeOfOption {
	r.greaterThan = true
	return r
}
func (r *RangeOfOption) LessThan() *RangeOfOption {
	r.lessThan = true
	return r
}
func (r *RangeOfOption) OrIsNull() *RangeOfOption {
	r.orIsNull = true
	return r
}
func (r *RangeOfOption) getMinNumberConditionKey() *ConditionKey {
	var ck ConditionKey
	if r.greaterThan {
		if r.orIsNull {
			ck = CK_GTISN

		} else {
			ck = CK_GT
		}

	} else {
		if r.orIsNull {
			ck = CK_GEISN
		} else {
			ck = CK_GE
		}
	}
	return &ck
}
func (r *RangeOfOption) getMaxNumberConditionKey() *ConditionKey {
	var ck ConditionKey
	if r.lessThan {
		if r.orIsNull {
			ck = CK_LTISN

		} else {
			ck = CK_LT
		}

	} else {
		if r.orIsNull {
			ck = CK_LEISN
		} else {
			ck = CK_LE
		}
	}
	return &ck
}
