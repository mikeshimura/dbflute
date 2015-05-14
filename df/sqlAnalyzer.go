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
	//	"strconv"
	//"errors"
	"fmt"
	"github.com/mikeshimura/dbflute/log"
	//"reflect"
	"strings"
)

const (
	TOK_SQL                                     = 1
	TOK_COMMENT                                 = 2
	TOK_ELSE                                    = 3
	TOK_BIND_VARIABLE                           = 4
	TOK_EOF                                     = 99
	FOR_NODE_PREFIX                             = "FOR "
	FOR_NODE_CURRENT_VARIABLE                   = "#current"
	JAVA_FIRST_CHAR                             = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz$_"
	BEGIN_NODE_MARK                             = "BEGIN"
	IF_NODE_PREFIX                              = "IF "
	LOOP_FIRSTNODE_MARK                         = "FIRST"
	LOOP_NEXTNODE_MARK                          = "NEXT"
	LOOP_LASTNODE_MARK                          = "LAST"
	EMBEDDINT_VARIABLE_NODE_PREFIX_NORMAL       = "$"
	EMBEDDINT_VARIABLE_NODE_PREFIX_REPLACE_ONLY = "$$"
	EMBEDDINT_VARIABLE_NODE_PREFIX_TERMINAL_DOT = "$."
)

type SqlAnalyzer struct {
	specifiedSql       string
	blockNullParameter bool
	stack              *Stack
	tokenizer          *SqlTokenizer
	inBeginScope       bool
}

func (s *SqlAnalyzer) Setup(sql string, blockNullParameter bool) {
	//fmt.Println("analyzer sql :"+sql)
	s.specifiedSql = sql
	s.stack = new(Stack)
	s.tokenizer = new(SqlTokenizer)
	s.tokenizer.Setup(sql)
	s.blockNullParameter = blockNullParameter
}
func (s *SqlAnalyzer) Analyze() *Node {

	//	     push(createRootNode()); // root node of all
	rn := s.CreateRootNode()
	log.InternalDebug(fmt.Sprintf("AnalyzeRoot Node %v\n", rn))
	var node Node = rn
	s.Push(&node)
	//        while (SqlTokenizer.EOF != _tokenizer.next()) {
	//for TOK_EOF != s.tokenizer.Next() {
	for {
		res:= s.tokenizer.Next()
		if res == TOK_EOF {
			break
		}
		s.parseToken()
	}
	//            parseToken();
	//        }
	//        return pop();
	res := s.Pop()
	log.InternalDebug(fmt.Sprintf("AnalyzeRoot End Node %v\n", res))
	return res
}
func (s *SqlAnalyzer) parseToken() {
	switch s.tokenizer.tokenType {
	case TOK_SQL:
		s.parseSql()
		break
	case TOK_COMMENT:
		s.parseComment()
		break
	case TOK_ELSE:
		s.parseElse()
		break
	case TOK_BIND_VARIABLE:
		s.parseBindVariable()
		break
	}
	return
}
func (s *SqlAnalyzer) parseSql() {
	sql := ""

	token := s.tokenizer.token
	if s.isElseMode() {
		token = strings.Replace(token, "--", "", -1)
	}
	sql = token

	node := s.Peek()
	if s.isSqlConnectorAdjustable(node) {
		s.processSqlConnectorAdjustable(node, sql)
	} else {
		log.InternalDebug(fmt.Sprintf("conn %v \n", s.createSqlPartsNodeOutOfConnector(node, sql)))
		(*node).AddChild(s.createSqlPartsNodeOutOfConnector(node, sql))
	}
	return
}
func (s *SqlAnalyzer) processSqlConnectorAdjustable(node *Node, sql string) {
	st := new(SqlTokenizer)
	st.sql = sql
	st.skipWhitespaceNoPos()
	skippedToken := st.skipToken(false)
	st.skipWhitespaceNoPos()

	if s.processSqlConnectorMark(node, sql) { // comma, ...
		return
	}
	if s.processSqlConnectorCondition(node, st, skippedToken) { // and/or
		return
	}
	//is not connector
	(*node).AddChild(s.createSqlPartsNodeThroughConnector(node, sql))
}
func (s *SqlAnalyzer) createSqlPartsNodeThroughConnector(node *Node, sql string) *Node {
	if s.isNestedBegin(node) { // basically nested if BEGIN node because checked before
		// connector adjustment of BEGIN is independent
		sqn := new(SqlConnectorNode)
		sqn.sqlParts = sql
		sqn.independent = true
		var n Node=sqn
		return &n
	} else {
		sqn := new(SqlConnectorNode)
		sqn.sqlParts = sql
		var n Node=sqn
		return &n
	}
	return nil
}
func (s *SqlAnalyzer) isNestedBegin(node *Node) bool {
	if (*node).stype() != "BeginNode" {
		return false
	}
	bgn:=((*node).getOrgAddress()).(*BeginNode)
	return bgn.isNested()
}
func (s *SqlAnalyzer) processSqlConnectorCondition(node *Node, st *SqlTokenizer, skippedToken string) bool {
	if strings.ToLower(skippedToken) == "and" || strings.ToLower(skippedToken) == "or" { // is connector
		(*node).AddChild(s.createSqlConnectorNode(node, st.getBefore(), st.getAfter()))
		return true
	}
	return false
}
func (s *SqlAnalyzer) createSqlConnectorNode(node *Node, connector string, sqlParts string) *Node {
	if s.isNestedBegin(node) { // basically nested if BEGIN node because checked before
		// connector adjustment of BEGIN is independent
		sqn := new(SqlConnectorNode)
		sqn.connector = connector
		sqn.sqlParts = sqlParts
		sqn.independent = true
		var n Node=sqn
		return &n
	} else {
		sqn := new(SqlConnectorNode)
		sqn.connector = connector
		sqn.sqlParts = sqlParts
			var n Node=sqn
		return &n
	}
}
func (s *SqlAnalyzer) processSqlConnectorMark(node *Node, sql string) bool {
	if s.doProcessSqlConnectorMark(node, sql, ",") { // comma
		return true
	}
	return false
}
func (s *SqlAnalyzer) doProcessSqlConnectorMark(node *Node, sql string, mark string) bool {
	ltrimmedSql := strings.TrimLeft(sql, " ")  // for mark
	if strings.Index(ltrimmedSql, mark) == 0 { // is connector
		markSpace := mark + " "
		realMark := ""
		if strings.Index(ltrimmedSql, markSpace) == 0 {
			realMark = markSpace
		} else {
			realMark = mark
		}
		(*node).AddChild(s.createSqlConnectorNode(node, realMark, ltrimmedSql[len(realMark):]))
		return true
	}
	return false
}
func (s *SqlAnalyzer) isSqlConnectorAdjustable(node *Node) bool {
	if (*node).GetChildSize() > 0 {
		return false
	}
	return (*node).isImplementedSqlConnectorAdjustable() && !s.isTopBegin(node)
}
func (s *SqlAnalyzer) isElseMode() bool {
	//	未実装
	//	        for (int i = 0; i < _nodeStack.size(); ++i) {
	//            if (_nodeStack.get(i) instanceof ElseNode) {
	//                return true;
	//            }
	//        }
	for _, node := range s.stack.data {
		n := node.(*Node)
		log.InternalDebug(fmt.Sprintf("nodetype %v %s\n", node, (*n).stype()))
	}
	return false
}
func (s *SqlAnalyzer) createSqlPartsNodeOutOfConnector(node *Node, sqlParts string) *Node {
	if s.isTopBegin(node) { // top BEGIN only (nested goes 'else' statement)
		pn := new(SqlPartsNode)
		pn.Setup()
		pn.sqlParts = sqlParts
		pn.independent = true
		var node Node = pn
		return &node
	} else {
		return s.createSqlPartsNode(sqlParts)
	}
}
func (s *SqlAnalyzer) createSqlPartsNode(sqlParts string) *Node {
	//return SqlPartsNode.createSqlPartsNode(sqlParts)
	spn := new(SqlPartsNode)
	spn.Setup()
	spn.sqlParts = sqlParts
	var node Node = spn
	return &node
}
func (s *SqlAnalyzer) isTopBegin(node *Node) bool {
	if (*node).stype() != "BeginNode" {
		return false
	}
	bgn:=((*node).getOrgAddress()).(*BeginNode)
	return !bgn.isNested()
}

func (s *SqlAnalyzer) isTargetComment(comment string) bool {
	//if (Srl.is_Null_or_TrimmedEmpty(comment)) {
	if len(strings.Trim(comment, " ")) == 0 {
		return false
	}
	//if (!comment.startsWith(ForNode.CURRENT_VARIABLE)) {
	if strings.Index(comment, FOR_NODE_CURRENT_VARIABLE) != 0 { // except current variable from check
		//            if (!Character.isJavaIdentifierStart(comment.charAt(0))) {
		firstchar := string(comment[0])
		if strings.Index(JAVA_FIRST_CHAR, firstchar) == -1 {
			return false
		}
	}
	return true
}

func (s *SqlAnalyzer) isBeginComment(comment string) bool {
	return comment == BEGIN_NODE_MARK
}
func (s *SqlAnalyzer) isIfComment(comment string) bool {
	return strings.Index(comment, IF_NODE_PREFIX) == 0
}
func (s *SqlAnalyzer) isForComment(comment string) bool {
	return strings.Index(comment, FOR_NODE_PREFIX) == 0
}
func (s *SqlAnalyzer) isLoopVariableComment(comment string) bool {
	return strings.Index(comment, LOOP_FIRSTNODE_MARK) == 0 ||
		strings.Index(comment, LOOP_NEXTNODE_MARK) == 0 ||
		strings.Index(comment, LOOP_LASTNODE_MARK) == 0
}
func (s *SqlAnalyzer) isEndComment(content string) bool {
	return content == "END"
}
func (s *SqlAnalyzer) parseBegin() {
	//final BeginNode beginNode = createBeginNode();
	beginNode := new(BeginNode)
	beginNode.Setup()
	beginNode.nested = s.inBeginScope
beginNode.orgAddress=beginNode
	s.inBeginScope = true
	var node Node = beginNode
	(*s.Peek()).AddChild(&node)
	s.Push(&node)
	s.parseEnd()

	//finally
	s.inBeginScope = false
	return 

}
func (s *SqlAnalyzer) parseCommentBindVariable() error {
	expr := s.tokenizer.token
	testValue := s.tokenizer.skipToken(true)
	if strings.Index(expr, EMBEDDINT_VARIABLE_NODE_PREFIX_NORMAL) == 0 {
		if strings.Index(expr, EMBEDDINT_VARIABLE_NODE_PREFIX_REPLACE_ONLY) == 0 { // replaceOnly
			realExpr := expr[len(EMBEDDINT_VARIABLE_NODE_PREFIX_REPLACE_ONLY):]
			(*s.Peek()).AddChild(s.createEmbeddedVariableNode(realExpr, testValue, true, false))
		} else if strings.Index(expr, EMBEDDINT_VARIABLE_NODE_PREFIX_TERMINAL_DOT) == 0 { // terminalDot
			realExpr := expr[len(EMBEDDINT_VARIABLE_NODE_PREFIX_TERMINAL_DOT):]
			(*s.Peek()).AddChild(s.createEmbeddedVariableNode(realExpr, testValue, false, true))
		} else { // normal
			realExpr := expr[len(EMBEDDINT_VARIABLE_NODE_PREFIX_NORMAL):]
			(*s.Peek()).AddChild(s.createEmbeddedVariableNode(realExpr, testValue, false, false))
		}
	} else {
		(*s.Peek()).AddChild(s.createBindVariableNode(expr, testValue))
	}
	return nil
}
func (s *SqlAnalyzer) createEmbeddedVariableNode(expr string, testValue string, replaceOnly bool, terminalDot bool) *Node {
	//research 未実装
	//        researchIfNeeds(_researchEmbeddedVariableCommentList, expr); // for research
	//return new EmbeddedVariableNode(expr, testValue, _specifiedSql, _blockNullParameter, replaceOnly, terminalDot);
	en := new(EmbeddedVariableNode)
	en.Setup()
	en.expression = expr
	en.testValue = testValue
	en.specifiedSql = s.specifiedSql
	en.blockNullParameter = s.blockNullParameter
	en.SetupVariable(expr, testValue, s.specifiedSql, s.blockNullParameter)
	var node Node = en
	return &node
}
func (s *SqlAnalyzer) createBindVariableNode(expr string, testValue string) *Node {
	//research 未実装
	//	        researchIfNeeds(_researchBindVariableCommentList, expr); // for research
	//       return new BindVariableNode(expr, testValue, _specifiedSql, _blockNullParameter);
	bn := new(BindVariableNode)
	bn.Setup()
	bn.expression = expr
	bn.testValue = testValue
	bn.specifiedSql = s.specifiedSql
	bn.blockNullParameter = s.blockNullParameter
	bn.SetupVariable(expr, testValue, s.specifiedSql, s.blockNullParameter)
	var node Node = bn
	return &node
}
func (s *SqlAnalyzer) parseEnd()  {
	commentType := TOK_COMMENT
	//for TOK_EOF != s.tokenizer.Next() {
	for {
		res:= s.tokenizer.Next()
		if res == TOK_EOF {
			break
		}
		if s.tokenizer.tokenType == commentType && s.isEndComment(s.tokenizer.token) {
			s.Pop()
			return
		}
		s.parseToken()
	}
	//throwEndCommentNotFoundException();
	panic("EndCommentNotFound")
}
func (s *SqlAnalyzer) parseIf()  {
	comment := s.tokenizer.token
	condition := strings.Trim(comment[len(IF_NODE_PREFIX):], " ")
	if condition == "" {
		//throwIfCommentConditionEmptyException();
		panic("IfCommentConditionEmpty")
	}
	//final IfNode ifNode = createIfNode(condition);
	in := new(IfNode)
	in.Setup()
	in.expression = condition
	in.specifiedSql = s.specifiedSql
	var node Node = in
	(*s.Peek()).AddChild(&node)
	s.Push(&node)
	s.parseEnd()
	return
}
func (s *SqlAnalyzer) parseFor()  {
	comment := s.tokenizer.token
	condition := strings.Trim(comment[len(FOR_NODE_PREFIX):], " ")
	if condition == "" {
		//throwForCommentExpressionEmptyException();
		panic("ForCommentExpressionEmpty")
	}

	fn := new(ForNode)
	fn.Setup()
	fn.expression = condition
	fn.specifiedSql = s.specifiedSql
	var node Node = fn
	(*s.Peek()).AddChild(&node)
	s.Push(&node)
	s.parseEnd()
}
func (s *SqlAnalyzer) parseLoopVariable()  {
	//未実装 Priority Low
	comment := s.tokenizer.token
	code := substringFirstFront(comment, " ")
	code = code
	//        if (Srl.is_Null_or_TrimmedEmpty(code)) { // no way
	//            String msg = "Unknown loop variable comment: " + comment;
	//            throw new IllegalStateException(msg);
	//        }
	if strings.Trim(code, " ") == "" {
		panic("Unknown loop variable comment: " + comment)
	}
	return
	//        final LoopVariableType type = LoopVariableType.codeOf(code);
	//        if (type == null) { // no way
	//            String msg = "Unknown loop variable comment: " + comment;
	//            throw new IllegalStateException(msg);
	//        }
	//        final String condition = comment.substring(type.name().length()).trim();
	//        final LoopAbstractNode loopFirstNode = createLoopFirstNode(condition, type);
	//        peek().addChild(loopFirstNode);
	//        if (Srl.count(condition, "'") < 2) {
	//            push(loopFirstNode);
	//            parseEnd();
	//        }
}
func (s *SqlAnalyzer) parseElse()  {
	parent := s.Peek()
	stype := (*parent).stype()
	stype = stype
	if stype != "IfNode" {
		return
	}
	ifNode, ok := (*s.Pop()).(*IfNode)
	if ok == false {
		panic("not *IfNode")
	}
	elseNode := s.newElseNode()
	ifNode.elseNode = elseNode
	var node Node = elseNode
	s.Push(&node)
	s.tokenizer.skipWhitespaceNoPos()
	return
}
func (s *SqlAnalyzer) newElseNode() *ElseNode {
	en := new(ElseNode)
	en.Setup()
	return en
}
func (s *SqlAnalyzer) parseBindVariable() error {
	expr := s.tokenizer.token
	bn := s.createBindVariableNode(expr, "")
	var node Node = *bn
	(*s.Peek()).AddChild(&node)
	return nil
}

func (s *SqlAnalyzer) Push(node *Node) {
	s.stack.Push(node)
}
func (s *SqlAnalyzer) Pop() *Node {
	node, ok := (s.stack.Pop()).(*Node)
	if ok == false {
		panic("Type Unmatch expected Node")
	}
	var n Node = *node
	return &n
}
func (s *SqlAnalyzer) Peek() *Node {
	if len(s.stack.data) == 0 {
		panic("Cant't Peek 0 stack")
	}
	node, ok := (s.stack.data[len(s.stack.data)-1]).(*Node)
	if ok == false {
		panic("Type Unmatch expected Node")
	}
	var n Node = *node
	return &n
}

func (s *SqlAnalyzer) CreateRootNode() *RootNode {
	node := new(RootNode)
	node.Setup()
	return node
}

func (s *SqlAnalyzer) parseComment()  {
	comment := s.tokenizer.token
	if s.isTargetComment(comment) { // parameter comment
		if s.isBeginComment(comment) {
			s.parseBegin()
		} else if s.isIfComment(comment) {
			s.parseIf()
		} else if s.isForComment(comment) {
			s.parseFor()
		} else if s.isLoopVariableComment(comment) {
			s.parseLoopVariable()
		} else if s.isEndComment(comment) {
			return
		} else {
			s.parseCommentBindVariable()
		}
		//       } else if (Srl.is_NotNull_and_NotTrimmedEmpty(comment)) { // plain comment
	} else if len(strings.Trim(comment, " ")) > 0 {
		before := s.tokenizer.getBefore()
		content := before[strings.LastIndex(before, "/*"):]
		(*s.Peek()).AddChild(s.createSqlPartsNode(content))
	}
	return

}
