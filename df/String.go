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
	"strings"
)
func StringCount(s string, ele string)int{
	count:=0
	for (true){
		pos:=strings.Index(s,ele)
		if pos > -1 {
			count ++
			s=s[pos+len(ele):]
		} else{
			break
		}
	}
	return count
}
func IndexAfter(s string, find string, spos int) int {
	pos := strings.Index(s[spos:], find)
	if pos == -1 {
		return -1
	}
	return spos + pos
}
var Ln string
func InitUnCap(str string) string {
	if str == "" {
		return str
	}
	if len(str) == 1 {
		res := strings.ToLower(str)
		return res
	}
	res := strings.ToLower(string(str[0])) + string(str[1:len(str)])
	return res
}
func InitCap(str string) string {
	if str == "" {
		return str
	}
	if len(str) == 1 {
		res := strings.ToUpper(str)
		return res
	}
	res := strings.ToUpper(string(str[0])) + string(str[1:len(str)])
	return res
}
func substringFirstFront(str string, delimiters ...string) string {

	return doSubstringFirstRear(false, false, false, str, delimiters...)

}
func substringFirstFrontIgnoreCase(str string, delimiters ...string) string {

	return doSubstringFirstRear(false, false, true, str, delimiters...)

}
func substringFirstRear(str string, delimiters ...string) string {

	return doSubstringFirstRear(false, true, false, str, delimiters...)

}
func substringFirstRearIgnoreCase(str string, delimiters ...string) string {

	return doSubstringFirstRear(false, true, true, str, delimiters...)

}
func substringLastFront(str string, delimiters ...string) string {

	return doSubstringFirstRear(true, false, false, str, delimiters...)

}
func substringLastFrontIgnoreCase(str string, delimiters ...string) string {

	return doSubstringFirstRear(true, false, true, str, delimiters...)

}
func substringLastRear(str string, delimiters ...string) string {

	return doSubstringFirstRear(true, true, false, str, delimiters...)

}
func substringLastRearIgnoreCase(str string, delimiters ...string) string {

	return doSubstringFirstRear(true, true, true, str, delimiters...)

}
func doSubstringFirstRear(last bool, rear bool, ignoreCase bool, str string, delimiters ...string) string {
	var info *IndexOfInfo
	if ignoreCase {
		if last {
			info = indexOfLastIgnoreCase(str, delimiters...)

		} else {
			info = indexOfFirstIgnoreCase(str, delimiters...)
		}
	} else {
		if last {
			info = indexOfLast(str, delimiters...)
		} else {
			info = indexOfFirst(str, delimiters...)
		}
	}
	if info == nil {
		return str
	}
	if rear {
		return str[info.index+len(info.delimiter):]
	} else {
		return str[0:info.index]
	}

}
func indexOfFirst(str string, delimiters ...string) *IndexOfInfo {
	return doIndexOfFirst(false, str, delimiters...)
}
func indexOfLast(str string, delimiters ...string) *IndexOfInfo {
	return doIndexOfLast(false, str, delimiters...)
}
func indexOfLastIgnoreCase(str string, delimiters ...string) *IndexOfInfo {
	return doIndexOfLast(true, str, delimiters...)
}
func indexOfFirstIgnoreCase(str string, delimiters ...string) *IndexOfInfo {
	return doIndexOfFirst(true, str, delimiters...)
}
func doIndexOfLast(ignoreCase bool, str string, delimiters ...string) *IndexOfInfo {
	return doIndexOf(ignoreCase, true, str, delimiters...)
}
func doIndexOfFirst(ignoreCase bool, str string, delimiters ...string) *IndexOfInfo {
	return doIndexOf(ignoreCase, false, str, delimiters...)
}
func doIndexOf(ignoreCase bool, last bool, str string, delimiters ...string) *IndexOfInfo {
	var filteredStr string
	if ignoreCase {
		filteredStr = strings.ToLower(str)
	} else {
		filteredStr = str
	}
	targetIndex := -1
	targetDelimiter := ""

	for _, delimiter := range delimiters {
		var filteredDelimiter string
		if ignoreCase {
			filteredDelimiter = strings.ToLower(delimiter)
		} else {
			filteredDelimiter = delimiter
		}
		var index int
		if last {
			index = strings.LastIndex(filteredStr, filteredDelimiter)
		} else {
			index = strings.Index(filteredStr, filteredDelimiter)
		}
		if index < 0 {
			continue
		}
		if targetIndex < 0 || (last && targetIndex < index) || (!last && targetIndex > index) {
			targetIndex = index
			targetDelimiter = delimiter
		}
	}
	if targetIndex < 0 {
		return nil
	}
	info := new(IndexOfInfo)
	info.baseString = str
	info.index = targetIndex
	info.delimiter = targetDelimiter
	return info

}

type IndexOfInfo struct {
	baseString string
	index      int
	delimiter  string
}

func (i *IndexOfInfo) substringFront() string {
	return i.baseString[0:i.index]
}
func (i *IndexOfInfo) substringFrontTrimmed() string {
	return strings.Trim(i.substringFront(), " ")
}
func (i *IndexOfInfo) substringRear() string {
	return i.baseString[i.index:]
}
func (i *IndexOfInfo) substringRearTrimmed() string {
	return strings.Trim(i.substringRear(), " ")
}

func (i *IndexOfInfo) getRearIndex() int {
	return i.index + len(i.delimiter)
}
func extractScopeWide(str string, beginMark string, endMark string) *ScopeInfo {
	first := indexOfFirst(str, beginMark)
	if first == nil {
		return nil
	}
	last := indexOfLast(str, endMark)
	if last == nil {
		return nil
	}
	content := str[first.index+len(first.delimiter) : last.index]
	info := new(ScopeInfo)
	info.baseString = str
	info.beginIndex = first.index
	info.endIndex = last.index
	info.beginMark = beginMark
	info.endMark = endMark
	info.content = content
	info.scope = beginMark + content + endMark
	return info
}

type ScopeInfo struct {
	baseString string
	beginIndex int
	endIndex   int
	beginMark  string
	endMark    string
	content    string
	scope      string
	previous   *ScopeInfo
	next       *ScopeInfo
}

func (s *ScopeInfo) isBeforeScope(index int) bool {
	return index < s.beginIndex
}
func (s *ScopeInfo) isInScope(index int) bool {
	return index >= s.beginIndex && index <= s.endIndex
}
func (s *ScopeInfo) replaceContentOnBaseString(toStr string) string {
	scopeList := s.takeScopeList()

	//            final StringBuilder sb = new StringBuilder();
	sb := new(bytes.Buffer)
	//            for (ScopeInfo scope : scopeList) {
	for _, scopei := range scopeList.data {
		//                sb.append(scope.substringInterspaceToPrevious());
		scope, ok := scopei.(*ScopeInfo)
		if ok == false {
			panic("Scope type error")
		}
		sb.WriteString(scope.beginMark)
		sb.WriteString(toStr)
		sb.WriteString(scope.endMark)
		if scope.next == nil {
			sb.WriteString(scope.substringInterspaceToNext())
		}
	}
	return sb.String()
}
func (s *ScopeInfo) replaceContentOnBaseStringFrom(fromStr string, toStr string) string {
	scopeList := s.takeScopeList()
	sb := new(bytes.Buffer)
	for _, scopei := range scopeList.data {
		scope, ok := scopei.(*ScopeInfo)
		if ok == false {
			panic("Scope type error")
		}
		sb.WriteString(scope.substringInterspaceToPrevious())
		sb.WriteString(scope.beginMark)
		sb.WriteString(strings.Replace(scope.content, fromStr, toStr, 1))
		sb.WriteString(scope.endMark)
		if scope.next == nil { // last
			sb.WriteString(scope.substringInterspaceToNext())
		}
	}

	return sb.String()
}
func (s *ScopeInfo) replaceScopeInterspace(str string, fromStr string, toStr string, beginMark string, endMark string) string {
	scopeList := extractScopeList(str, beginMark, endMark)
	if len(scopeList.data) == 0 {
		return str
	}
	sl, ok := scopeList.data[0].(*ScopeInfo)
	if ok == false {
		panic("Not ScopeInfo")
	}
	return sl.replaceInterspaceOnBaseString(fromStr, toStr)
}
func (s *ScopeInfo) replaceInterspaceOnBaseString(fromStr string, toStr string) string {
	scopeList := s.takeScopeList()
	sb := new(bytes.Buffer)
	for _, scopei := range scopeList.data {
		scope, ok := scopei.(*ScopeInfo)
		if ok == false {
			panic("Not ScopeInfo")
		}
		sb.WriteString(strings.Replace(scope.substringInterspaceToPrevious(), fromStr, toStr, 1))
		sb.WriteString(scope.scope)
		if scope.next == nil { // last
			sb.WriteString(strings.Replace(scope.substringInterspaceToNext(), fromStr, toStr, 1))
		}
	}
	return sb.String()
}
func extractScopeList(str string, beginMark string, endMark string) *List {
	scopeList := doExtractScopeList(str, beginMark, endMark, false)
	if scopeList == nil {
		return new(List)
	} else {
		return scopeList
	}
}
func doExtractScopeList(str string, beginMark string, endMark string, firstOnly bool) *List {
	var resultList *List
	var previous *ScopeInfo = nil
	rear := str
	for {
		beginIndex := strings.Index(rear, beginMark)
		if beginIndex < 0 {
			break
		}
		rear = rear[beginIndex:] // scope begins
		if len(rear) <= len(beginMark) {
			break
		}
		rear = rear[len(beginMark):] // skip begin-mark
		endIndex := strings.Index(rear, endMark)
		if endIndex < 0 {
			break
		}
		scope := beginMark + rear[0:endIndex+len(endMark)]
		info := new(ScopeInfo)
		info.baseString = str
		var absoluteIndex int
		if previous != nil {
			absoluteIndex = previous.endIndex + beginIndex
		} else {
			absoluteIndex = 0 + beginIndex
		}
		info.beginIndex = absoluteIndex
		info.endIndex = absoluteIndex + len(scope)
		info.beginMark = beginMark
		info.endMark = endMark
		info.content = strings.TrimRight(strings.TrimLeft(scope, beginMark), endMark)
		info.scope = scope
		if previous != nil {
			info.previous = previous
			previous.next = info
		}
		if resultList == nil {
			resultList = new(List) // lazy load
		}
		resultList.Add(info)
		if previous == nil && firstOnly {
			break
		}
		previous = info
		rear = str[info.endIndex:]
	}
	return resultList // nullable if not found to suppress unneeded ArrayList creation

}
func (s *ScopeInfo) substringInterspaceToPrevious() string {
	previousEndIndex := -1
	if s.previous != nil {
		previousEndIndex = s.previous.endIndex
	}
	if previousEndIndex >= 0 {
		return s.baseString[previousEndIndex:s.beginIndex]
	} else {
		return s.baseString[0:s.beginIndex]
	}
}
func (s *ScopeInfo) substringInterspaceToNext() string {
	nextBeginIndex := -1
	if s.next != nil {
		nextBeginIndex = s.next.beginIndex
	}
	if nextBeginIndex >= 0 {
		return s.baseString[s.endIndex:nextBeginIndex]
	} else {
		return s.baseString[s.endIndex:]
	}
}
func (s *ScopeInfo) takeScopeList() *List {
	scope := s
	for {
		previous := scope.previous
		if previous == nil {
			break
		}
		scope = previous
	}
	scopeList := new(List)
	scopeList.Add(scope)
	for {
		next := scope.next
		if next == nil {
			break
		}
		scope = next
		scopeList.Add(next)
	}
	return scopeList
}
func splitList(str string, delimitter string) *StringList {
	return doSplitList(str, delimitter, false)
}
func splitListTrimmed(str string, delimitter string) *StringList {
	return doSplitList(str, delimitter, true)
}
func doSplitList(str string, delimiter string, trim bool) *StringList {
	list := new(StringList)
	elementIndex := 0
	delimiterIndex := strings.Index(str, delimiter)
	var element string
	for delimiterIndex >= 0 {
		element = str[elementIndex:delimiterIndex]
		if trim {
			list.Add(strings.Trim(element, " "))
		} else {
			list.Add(element)
		}
		elementIndex = delimiterIndex + len(delimiter)
		temp := strings.Index(str[elementIndex:], delimiter)
		if temp < 0 {
			delimiterIndex = temp
		} else {
			delimiterIndex = temp + elementIndex
		}
	}
	element = str[elementIndex:]
	if trim {
		list.Add(strings.Trim(element, " "))
	} else {
		list.Add(element)
	}
	return list
}
func extractScopeFirst(str string,beginMark string,endMark string)*ScopeInfo{
	     scopeList := doExtractScopeList(str, beginMark, endMark, false);
        if (scopeList == nil || scopeList.Size()==0) {
            return nil
        }
        sl:= scopeList.Get(scopeList.Size() - 1)
        var res *ScopeInfo=sl.(*ScopeInfo)
        return res
}
