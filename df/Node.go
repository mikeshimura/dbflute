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
	//"errors"
	"fmt"
	"github.com/mikeshimura/dbflute/log"
	"reflect"
	"strings"
)

const (
	COMMENT_TYPE_BIND       = 0
	COMMENT_TYPE_EMBEDDED   = 1
	COMMENT_TYPE_FORCOMMENT = 2
)

type Node interface {
	AddChild(node interface{})
	GetChildSize() int
	GetChild(i int) *Node
	accept(ctx *CommandContext, node *Node)
	getCommentType() *CommentType
	doProcess(ctx *CommandContext, valueAndType *ValueAndType, loopInfo *LoopInfo)
	stype() string
	isImplementedSqlConnectorAdjustable() bool
	getOrgAddress() interface{}
}

type RootNode struct {
	BaseNode
}

func (r *RootNode) stype() string {
	return "RootNode"
}
func (r *RootNode) accept(ctx *CommandContext, node *Node) {
	for i := 0; i < r.GetChildSize(); i++ {
		inter := r.GetChild(i)
		//var node Node = *inter.(*Node)
		stype := GetType(*inter)
		log.InternalDebug(fmt.Sprintln("accept node type :" + stype))
		//node.accept(ctx)
		(*inter).accept(ctx, inter)
	}
	return
}

type SqlPartsNode struct {
	BaseNode
	sqlParts    string
	independent bool
}

func (r *SqlPartsNode) stype() string {
	return "SqlPartsNode"
}
func (r *SqlPartsNode) accept(ctx *CommandContext, node *Node) {
	log.InternalDebug("SqlPartsNode accept")
	(*ctx).addSql(r.sqlParts)
	log.InternalDebug("SQL PART*" + r.sqlParts)
	//以下未実装
	//        if (isMarkAlreadySkipped(ctx)) {
	//            // It does not skipped actually but it has not already needed to skip.
	//            ctx.setAlreadySkippedConnector(true);
	//        }
	return
}

type ScopeNode struct {
	BaseNode
}

func (r *ScopeNode) stype() string {
	return "ScopeNode"
}
func (r *ScopeNode) processAcceptingChildren(ctx *CommandContext, loopInfo *LoopInfo) {
	childSize := r.GetChildSize()
	for i := 0; i < childSize; i++ {
		child := r.GetChild(i)
		if loopInfo != nil { // in loop
			panic("processAcceptingChildren")
			//                if (child instanceof LoopAcceptable) { // accepting loop
			//                    handleLoopElementNullParameter(child, loopInfo);
			//                    ((LoopAcceptable) child).accept(ctx, loopInfo);
			//                } else {
			//                    child.accept(ctx);
			//                }
		} else {
			(*child).accept(ctx, child)
		}
	}
	return
}

type VariableNode struct {
	BaseNode
	expression         string
	testValue          string
	optionDef          string
	specifiedSql       string
	blockNullParameter bool
	nameList           *StringList
}

func (r *VariableNode) stype() string {
	return "VariableNode"
}
func (r *VariableNode) isInScope() bool {
	if r.testValue == "" {
		return false
	}
	return r.testValue[0:1] == "(" && r.testValue[len(r.testValue)-1:len(r.testValue)] == ")"

}
func (r *VariableNode) SetupVariable(expression string, testValue string, specifiedSql string, blockNullParameter bool) {
	if strings.Contains(expression, ":") {
		r.expression = strings.Trim(substringFirstFront(expression, ":"), " ")
		r.optionDef = strings.Trim(substringFirstRear(expression, ":"), " ")
	} else {
		r.expression = expression
		r.optionDef = ""
	}
	r.testValue = testValue
	r.specifiedSql = specifiedSql
	r.blockNullParameter = blockNullParameter
	r.nameList = splitList(expression, ".")
}
func (r *VariableNode) accept(ctx *CommandContext, node *Node) {
	r.doAccept(ctx, nil, node)
}
func (r *VariableNode) doAccept(ctx *CommandContext, loopInfo *LoopInfo, node *Node) {
	firstName := r.nameList.Get(0)
	//未実装
	//assertFirstNameAsNormal(ctx, firstName);

	var firstValue interface{} = (*ctx).GetArgs()[firstName]
	firstType := (*ctx).GetArgTypes()[firstName]
	r.doAcceptSub(ctx, firstValue, firstType, loopInfo, false, node)
	return
}
func (r *VariableNode) doAcceptSub(ctx *CommandContext, firstValue interface{}, firstType string, loopInfo *LoopInfo, inheritLoop bool, node *Node) {
	//        assertInLoopOnlyOptionInLoop(loopInfo);
	valueAndType := new(ValueAndType)
	valueAndType.firstValue = firstValue
	valueAndType.firstType = firstType

	r.setupValueAndType(valueAndType, node)
	//未実装
	//        processLikeSearch(valueAndType, loopInfo, inheritLoop);
	log.InternalDebug(fmt.Sprintf("blockNullParameter %v \n ", r.blockNullParameter))
	if r.blockNullParameter && IsNotNull(valueAndType.targetValue) == false {
		panic("The value of bind variable was null.")
	}
	(*node).doProcess(ctx, valueAndType, loopInfo)
	return
}

func (r *VariableNode) setupValueAndType(valueAndType *ValueAndType, node *Node) {
	ctype := (*node).getCommentType()
	setupper := new(ValueAndTypeSetupper)
	setupper.Setup(r.nameList, r.expression, r.specifiedSql, ctype)
	ctype = ctype
	setupper.setupValueAndType(valueAndType)
}

type BeginNode struct {
	ScopeNode
	nested bool
}

func (b *BeginNode) isNested() bool {
	return b.nested
}
func (b *BeginNode) isImplementedSqlConnectorAdjustable() bool {
	return true
}
func (r *BeginNode) stype() string {
	return "BeginNode"
}
func (r *BeginNode) accept(ctx *CommandContext, node *Node) {
	r.doAccept(ctx, nil)
}
func (r *BeginNode) doAccept(ctx *CommandContext, loopInfo *LoopInfo) {
	childCtx := new(CommandContextImpl)
	childCtx.parent = ctx
	childCtx.beginChild = true
	var child CommandContext = childCtx
	r.processAcceptingChildren(&child, loopInfo)
	if child.isEnabled() {
		(*ctx).addSqlBind(child.getSql(), child.getBindVariables(), child.getBindVariableTypes())
		if (*ctx).isBeginChild() {
			(*ctx).setEnabled(true)
		}
	}
	return
}

//func (r *BeginNode) processAcceptingChildren(ctx *CommandContext, loopInfo *LoopInfo) error {
//	childSize := r.GetChildSize()
//	for i := 0; i < childSize; i++ {
//		child := r.GetChild(i)
//		//            if (loopInfo != null) { // in loop
//		//                if (child instanceof LoopAcceptable) { // accepting loop
//		//                    handleLoopElementNullParameter(child, loopInfo);
//		//                    ((LoopAcceptable) child).accept(ctx, loopInfo);
//		//                } else {
//		//                    child.accept(ctx);
//		//                }
//		//            } else {
//		err := (*child).accept(ctx, child)
//		if err != nil {
//			return err
//		}
//		//            }
//	}
//	return nil
//}

type IfNode struct {
	ScopeNode
	expression   string
	specifiedSql string
	elseNode     *ElseNode
}

func (b *IfNode) isImplementedSqlConnectorAdjustable() bool {
	return true
}
func (r *IfNode) stype() string {
	return "IfNode"
}
func (r *IfNode) accept(ctx *CommandContext, node *Node) {
	r.doAcceptByEvaluator(ctx, nil)
}
func (r *IfNode) doAcceptByEvaluator(ctx *CommandContext, loopInfo *LoopInfo) {
	cmap := (*ctx).GetArgs()
	cmap = cmap
	//fmt.Printf("ctx %v \n", len(cmap))
	evaluator := new(IfCommentEvaluator)
	f := new(ParameterFinderArg)
	var finder ParameterFinder = f
	evaluator.finder = &finder
	evaluator.ctx = ctx
	evaluator.expression = r.expression
	evaluator.specifiedSql = r.specifiedSql
	result := evaluator.evaluate()
	log.InternalDebug(fmt.Sprintf("Else node %v\n", r.elseNode))
	if result {
		r.processAcceptingChildren(ctx, loopInfo)
		(*ctx).setEnabled(true)
	} else if r.elseNode != nil {
		//            if (loopInfo != null) {
		//                _elseNode.accept(ctx, loopInfo);
		//            } else {
		var node Node = r.elseNode
		r.elseNode.accept(ctx, &node)
		//            }
	}
	return
}

type ElseNode struct {
	ScopeNode
}

func (b *ElseNode) isImplementedSqlConnectorAdjustable() bool {
	return true
}
func (r *ElseNode) stype() string {
	return "ElseNode"
}
func (r *ElseNode) accept(ctx *CommandContext, node *Node) {
	r.doAccept(ctx, nil)
}
func (r *ElseNode) doAccept(ctx *CommandContext, loopInfo *LoopInfo) {
	r.processAcceptingChildren(ctx, loopInfo)
	(*ctx).setEnabled(true)
}

type ForNode struct {
	ScopeNode
	expression   string
	specifiedSql string
}

func (b *ForNode) isImplementedSqlConnectorAdjustable() bool {
	return true
}
func (r *ForNode) stype() string {
	return "ForNode"
}
func (r *ForNode) accept(ctx *CommandContext, node *Node) {
	return
}

type EmbeddedVariableNode struct {
	VariableNode
	replaceOnly bool

	terminalDot bool
}

func (r *EmbeddedVariableNode) stype() string {
	return "EmbeddedVariableNode"
}
func (e *EmbeddedVariableNode) doProcess(ctx *CommandContext, valueAndType *ValueAndType, loopInfo *LoopInfo) {
	finalValue := valueAndType.targetValue
	finalType := valueAndType.targetType
	finalType = finalType

	if e.isInScope() {
		if finalValue == nil { // in-scope does not allow null value
			panic("BindOrEmbeddedCommentParameterNullValueException(valueAndType)")
		}
		panic("未実装BindOrEmbeddedCommentInScope")
		//   要確認 finalTypeがJAVAのCollection interfaceをImplementしていればの意味か
		//            if (Collection.class.isAssignableFrom(finalType)) {
		//                embedArray(ctx, ((Collection<?>) finalValue).toArray());
		//            } else if (finalType.isArray()) {
		//                embedArray(ctx, finalValue);
		//            } else {
		//                throwBindOrEmbeddedCommentInScopeNotListException(valueAndType);
		//            }
	} else {
		if finalValue == nil {
			(*ctx).addSql("")
		} else if GetType(finalValue) != "string" {
			embeddedValue := ToStringInterface(finalValue)
			if e.isQuotedScalar() { // basically for condition value
				(*ctx).addSql(e.quote(embeddedValue))
			} else { // basically for cannot-bound condition (for example, paging)
				(*ctx).addSql(embeddedValue)
			}
		} else {
			//                // string type here
			var embeddedStr = finalValue.(string)
			e.assertNotContainBindSymbol(embeddedStr)
			if e.isQuotedScalar() { // basically for condition value
				(*ctx).addSql(e.quote(embeddedStr))
				//未実装
				//if r.isAcceptableLikeSearch(loopInfo) {
				//                        setupRearOption(ctx, valueAndType);
				//}
			} else {
				firstValue := valueAndType.firstValue
				firstType := valueAndType.firstType
				bound := e.processDynamicBinding(ctx, firstValue, firstType, embeddedStr)
				if !bound {
					(*ctx).addSql(embeddedStr)
				}

			}
		}
	}
	//要確認 golang string nullは無いので対応検討
	//if (e.testValue != "") {
	if e.replaceOnly { // e.g. select ... from /*$$pmb.schema*/MEMBER
		// actually the test value is not test value
		// but a part of SQL statement here
		(*ctx).addSql(e.testValue)
	} else if e.terminalDot { // e.g. select ... from /*$$pmb.schema*/dev.MEMBER
		// the real test value is until a dot character
		(*ctx).addSql("." + substringFirstRear(e.testValue, "."))
	}
	return
}

var CMT_0 *CommentType
var CMT_1 *CommentType

func (e *EmbeddedVariableNode) processDynamicBinding(ctx *CommandContext, firstValue interface{}, firstType string, embeddedString string) bool {
	first := extractScopeFirst(embeddedString, "/*", "*/")
	if first == nil {
		return false
	}
	analyzer := new(SqlAnalyzer)
	analyzer.Setup(embeddedString, e.blockNullParameter)
	rootNode := analyzer.Analyze()

	creator := new(CommandContextCreator)
	creator.argNames = []string{"pmb"}
	creator.argTypes = []string{firstType}
	rootCtx := creator.createCommandContext([]interface{}{firstValue})
	(*rootNode).accept(rootCtx, nil)
	sql := (*rootCtx).getSql()
	(*ctx).addSqlBind(sql, (*rootCtx).getBindVariables(), (*rootCtx).getBindVariableTypes())
	return true
}
func (e *EmbeddedVariableNode) quote(value string) string {
	return "'" + value + "'"
}
func (e *EmbeddedVariableNode) isQuotedScalar() bool {
	return strings.Count(e.testValue, "'") > 1 && e.testValue[0:1] == "'" && e.testValue[len(e.testValue)-1:len(e.testValue)] == "'"
}
func (e *EmbeddedVariableNode) assertNotContainBindSymbol(value string) {
	if e.containsBindSymbol(value) {
		br := new(bytes.Buffer)
		br.WriteString("The value of embedded comment contained bind symbols.")
		br.WriteString("Advice")
		br.WriteString("The value of embedded comment should not contain bind symbols.")
		br.WriteString("For example, a question mark '?'.")
		br.WriteString("Comment Expression")
		br.WriteString(e.expression)
		br.WriteString("Embedded Value")
		br.WriteString(value)
		panic(br.String())
	}
	return
}
func (e *EmbeddedVariableNode) containsBindSymbol(value string) bool {
	return strings.Index(value, "?") > -1
}
func (b *EmbeddedVariableNode) getCommentType() *CommentType {
	if CMT_1 == nil {
		CMT_1 = new(CommentType)
		CMT_1.Ctype = 1
		CMT_1.Short = "EMBEDDED"
		CMT_1.TextName = "embedded variable comment"
		CMT_1.TitleName = "Embedded Variable Comment"
	}
	return CMT_1
}

type BindVariableNode struct {
	VariableNode
}

func (r *BindVariableNode) stype() string {
	return "BindVariableNode"
}
func (b *BindVariableNode) getCommentType() *CommentType {
	if CMT_0 == nil {
		CMT_0 = new(CommentType)
		CMT_0.Ctype = 0
		CMT_0.Short = "BIND"
		CMT_0.TextName = "bind variable comment"
		CMT_0.TitleName = "Bind Variable Comment"
	}
	return CMT_0
}

func (b *BindVariableNode) bindList(ctx *CommandContext, list *List) {
	(*ctx).addSql("(")
	for validCount, currentElement := range list.data {
		if currentElement != nil {
			if validCount > 0 {
				(*ctx).addSql(", ")
			}
			(*ctx).addSqlSingle("?", currentElement, GetType(currentElement))
		}
	}
	(*ctx).addSql(")")
}
func (b *BindVariableNode) doProcess(ctx *CommandContext, valueAndType *ValueAndType, loopInfo *LoopInfo) {
	finalValue := valueAndType.targetValue
	finalType := valueAndType.targetType
	if b.isInScope() {
		if finalValue == nil { // in-scope does not allow null value
			panic("BindOrEmbeddedCommentParameterNullValueException(valueAndType)")
		}
		var l *List = finalValue.(*List)
		b.bindList(ctx, l)
	} else {
		(*ctx).addSqlSingle("?", finalValue, finalType) // if null, bind as null
		//未実装
		//	            if (isAcceptableLikeSearch(loopInfo)) {
		//	                setupRearOption(ctx, valueAndType);
		//	            }
	}
	return
}

type LoopBaseNode struct {
	ScopeNode
	expression   string
	replacement  string
	specifiedSql string
}

func (r *LoopBaseNode) stype() string {
	return "LoopBaseNode"
}
func (r *LoopBaseNode) accept(ctx *CommandContext, node *Node) {

}

type LoopFirstNode struct {
	LoopBaseNode
}

func (b *LoopFirstNode) isImplementedSqlConnectorAdjustable() bool {
	return true
}
func (r *LoopFirstNode) stype() string {
	return "LoopFirstNode"
}
func (r *LoopFirstNode) accept(ctx *CommandContext, node *Node) {

}

type LoopNextNode struct {
	LoopBaseNode
}

func (r *LoopNextNode) stype() string {
	return "LoopNextNode"
}
func (r *LoopNextNode) accept(ctx *CommandContext, node *Node) {

}

type LoopLastNode struct {
	LoopBaseNode
}

func (r *LoopLastNode) accept(ctx *CommandContext, node *Node) {

}
func (r *LoopLastNode) stype() string {
	return "LoopLastNode"
}

//FIRST,NEXT,LAST で Factoryが Staticに保管されていてCreateを楽にするだけ
type LoopVariableType_T struct {
	codeValueMap map[string]*Node
}

func (l *LoopVariableType_T) PutCodeNode(code string, node *Node) {
	l.codeValueMap[strings.ToLower(code)] = node
}
func (l *LoopVariableType_T) GetNode(code string) *Node {
	return l.codeValueMap[strings.ToLower(code)]
}
func (l *LoopBaseNode) SetupLoop(expression string, specifiedSql string) {
	l.expression = expression
	l.specifiedSql = specifiedSql
	scope := extractScopeWide(expression, "'", "'")
	if scope != nil {
		l.replacement = scope.content
	}
}

type BaseNode struct {
	list       *List
	orgAddress interface{}
}

func (b *BaseNode) isBeginChildAndValidSql(ctx *CommandContext, sql string) bool {
	return (*ctx).isBeginChild() && len(strings.TrimSpace(sql)) > 0
}
func (b *BaseNode) getOrgAddress() interface{} {
	return b.orgAddress
}

func (b *BaseNode) isImplementedSqlConnectorAdjustable() bool {
	return false
}
func (b *BaseNode) doProcess(ctx *CommandContext, valueAndType *ValueAndType, loopInfo *LoopInfo) {
	//dummy Only EmbenddedVariableNode and BindVariableNode required
	log.InternalDebug("BaseNode Process")
	return
}
func (b *BaseNode) getCommentType() *CommentType {
	return nil
}
func (b *BaseNode) Setup() {
	b.list = new(List)
}
func (b *BaseNode) AddChild(node interface{}) {
	b.list.Add(node)
}
func (b *BaseNode) GetChildSize() int {
	return b.list.Size()
}
func (b *BaseNode) GetChild(i int) *Node {
	inter := b.list.Get(i)
	var node Node = *inter.(*Node)
	return &node
}

type LoopInfo struct {
	expression       string
	specifiedSql     string
	parameterList    *List
	loopSize         int
	likeSearchOption *LikeSearchOption
	loopIndex        int
	parentLoop       *LoopInfo
}

type ValueAndType struct {
	firstValue       interface{}
	firstType        string
	targetValue      interface{}
	targetType       string
	likeSearchOption *LikeSearchOption
}

type CommentType struct {
	Ctype     int
	Short     string
	TextName  string
	TitleName string
}
type ValueAndTypeSetupper struct {
	nameList     *StringList
	expression   string
	specifiedSql string
	commentType  *CommentType
}

func (v *ValueAndTypeSetupper) Setup(nameList *StringList, expression string, specifiedSql string, commentType *CommentType) {
	v.nameList = nameList
	v.expression = expression
	v.specifiedSql = specifiedSql
	v.commentType = commentType
}
func (v *ValueAndTypeSetupper) setupValueAndType(valueAndType *ValueAndType) {
	value := valueAndType.firstValue
	if value == nil { // if null, do nothing
		return
	}
	ctype := valueAndType.firstType
	ctype = ctype

	// LikeSearchOption handling here is for OutsideSql.
	//未実装
	var likeSearchOption LikeSearchOption
	likeSearchOption = likeSearchOption
	//
	for pos := 1; pos < v.nameList.Size(); pos++ {
		//            if (value == null) {
		//                break;
		//            }
		if value == nil {
			break
		}
		currentName := v.nameList.Get(pos)
		ctype, value = getPropertyValue(value, ctype, currentName)
		log.InternalDebug(fmt.Sprintf("new ctype %s\n", ctype))
		log.InternalDebug(fmt.Sprintf("new value %v\n", value))
		if ctype == "map[string]map[string]interface {}" {
			var mp map[string]map[string]interface{} = value.(map[string]map[string]interface{})
			value = mp[v.nameList.Get(pos+1)][v.nameList.Get(pos+2)]
			ctype = GetType(value)
			break
		}
	}
	//未実装
	//adjustLikeSearchDBWay(likeSearchOption)
	//        valueAndType.setTargetValue(value);
	valueAndType.targetValue = value
	log.InternalDebug(fmt.Sprintf("value %v %T\n", value, value))
	//        valueAndType.setTargetType(clazz);
	valueAndType.targetType = ctype
	//未実装
	//        valueAndType.setLikeSearchOption(likeSearchOption);
}
func getPropertyValue(value interface{}, ctype string, currentName string) (string, interface{}) {
	log.InternalDebug(fmt.Sprintf("value  getPropertyValue %v %T %s\n", value, value, currentName))
	xtype := GetType(value)
	log.InternalDebug("xtype:" + xtype)
	v := reflect.ValueOf(value).Elem()
	log.InternalDebug(fmt.Sprintf("new v %v %T\n", v, v))
	newv := v.FieldByName(InitCap(currentName))
	log.InternalDebug(fmt.Sprintf("newv  %v \n", newv))
	if newv.IsValid() == false {
		fmt.Printf("value  %v %T\n", value, value)
		fmt.Println("method" + InitCap(currentName))
		test2 := reflect.ValueOf(value).MethodByName(InitCap(currentName))
		if test2.IsValid() == false {
			test2 = reflect.ValueOf(value).MethodByName("Get" + InitCap(currentName))
			if test2.IsValid() == false {
				return "", nil
			}
		}
		nValuex := test2.Call([]reflect.Value{})
		nValue := nValuex[0].Interface()
		nType := GetType(nValue)
		return nType, nValue
	}
	newValue := newv.Interface()
	newType := GetType(newValue)

	return newType, newValue
}

type SqlConnectorNode struct {
	BaseNode
	connector   string
	sqlParts    string
	independent bool
}

func (s *SqlConnectorNode) stype() string {
	return "SqlConnectorNode"
}
func (s *SqlConnectorNode) accept(ctx *CommandContext, node *Node) {
	if (*ctx).isEnabled() || (*ctx).isAlreadySkippedConnector() {
		(*ctx).addSql(s.connector)
	} else if s.isMarkAlreadySkipped(ctx) {
		// To skip prefix should be done only once
		// so it marks that a prefix already skipped.
		(*ctx).setAlreadySkippedConnector(true)
	}
	(*ctx).addSql(s.sqlParts)

	return
}
func (s *SqlConnectorNode) isMarkAlreadySkipped(ctx *CommandContext) bool {
	return !s.independent && s.isBeginChildAndValidSql(ctx, s.sqlParts)
}

// BeginNode #nested @accept
// ScopeNode
// BaseNode #list @AddChild @GetChild @GetChildSize @Setup @doProcess @getCommentType
//
// BindVariableNode @doProcess @getCommentType
// VariableNode #blockNullParameter #expression #nameList #optionDef #specifiedSql #testValue @SetupVariable @accept @doAccept @doAcceptSub @isInScope @setupValueAndType
// BaseNode #list @AddChild @GetChild @GetChildSize @Setup @doProcess @getCommentType
//
// CommentType #Ctype #Short #TextName #TitleName
//
// ElseNode @accept
// ScopeNode
// BaseNode #list @AddChild @GetChild @GetChildSize @Setup @doProcess @getCommentType
//
// EmbeddedVariableNode #replaceOnly #terminalDot @assertNotContainBindSymbol @containsBindSymbol @doProcess @getCommentType @isQuotedScalar @processDynamicBinding @quote
// VariableNode #blockNullParameter #expression #nameList #optionDef #specifiedSql #testValue @SetupVariable @accept @doAccept @doAcceptSub @isInScope @setupValueAndType
// BaseNode #list @AddChild @GetChild @GetChildSize @Setup @doProcess @getCommentType
//
// ForNode #expression #specifiedSql @accept
// ScopeNode
// BaseNode #list @AddChild @GetChild @GetChildSize @Setup @doProcess @getCommentType
//
// IfNode #elseNode #expression #specifiedSql @accept
// ScopeNode
// BaseNode #list @AddChild @GetChild @GetChildSize @Setup @doProcess @getCommentType
//
// LikeSearchOption
// SimpleStringOption
//
// LoopFirstNode @accept
// LoopBaseNode #expression #replacement #specifiedSql @SetupLoop @accept
// ScopeNode
// BaseNode #list @AddChild @GetChild @GetChildSize @Setup @doProcess @getCommentType
//
// LoopInfo
//
// LoopLastNode @accept
// LoopBaseNode #expression #replacement #specifiedSql @SetupLoop @accept
// ScopeNode
// BaseNode #list @AddChild @GetChild @GetChildSize @Setup @doProcess @getCommentType
//
// LoopNextNode @accept
// LoopBaseNode #expression #replacement #specifiedSql @SetupLoop @accept
// ScopeNode
// BaseNode #list @AddChild @GetChild @GetChildSize @Setup @doProcess @getCommentType
//
// LoopVariableType_T #codeValueMap @GetNode @PutCodeNode
//
// RootNode @accept
// BaseNode #list @AddChild @GetChild @GetChildSize @Setup @doProcess @getCommentType
//
// SqlPartsNode #independent #sqlParts @accept
// BaseNode #list @AddChild @GetChild @GetChildSize @Setup @doProcess @getCommentType
//
// ValueAndType #firstType #firstValue #likeSearchOption #targetType #targetValue
//
// ValueAndTypeSetupper #commentType #expression #nameList #specifiedSql @Setup @setupValueAndType
//
