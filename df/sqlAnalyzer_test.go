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
	"fmt"
	"strconv"
	"testing"
	"github.com/mikeshimura/dbflute/log"
)

func TestSqlAnalyzer(t *testing.T) {
	pos := IndexAfter("test", "t", 1)
	if pos != 3 {
		t.Error("Expected 3 but " + strconv.Itoa(pos))
	}
	pos = IndexAfter("test", "x", 1)
	if pos != -1 {
		t.Error("Expected -1 but " + strconv.Itoa(pos))
	}
	pos = IndexAfter("test", "t", 0)
	if pos != 0 {
		t.Error("Expected 0 but " + strconv.Itoa(pos))
	}
	pos = IndexAfter("test", "s", 2)
	if pos != 2 {
		t.Error("Expected 2 but " + strconv.Itoa(pos))
	}
	tn := new(SqlTokenizer)
	tn.Setup(`select/*$pmb.selectHint*/ dfloc.id as id, dfloc.login_id as login_id, dfloc.name as name, dfloc.version_no as version_no
  from li_tbl dfloc
 where dfloc.id = /*pmb.conditionQuery.id.fixed.query.equal*/null
 `)
	log.InternalDebug(fmt.Sprintf("%v \n", tn))
	tn.parseSql()
	log.InternalDebug(fmt.Sprintf("%v \n", tn))
	log.InternalDebug(fmt.Sprintln("token:" + tn.token))
	log.InternalDebug(fmt.Sprintln("nexttokentype:" + strconv.Itoa(tn.nextTokenType)))
	log.InternalDebug(fmt.Sprintln("position:" + strconv.Itoa(tn.position)))
	tn.parseComment()
	log.InternalDebug(fmt.Sprintf("%v \n", tn))
	log.InternalDebug(fmt.Sprintln("token:" + tn.token))
	log.InternalDebug(fmt.Sprintln("nexttokentype:" + strconv.Itoa(tn.nextTokenType)))
	log.InternalDebug(fmt.Sprintln("position:" + strconv.Itoa(tn.position)))
	tn.parseSql()
	log.InternalDebug(fmt.Sprintf("%v \n", tn))
	log.InternalDebug(fmt.Sprintln("token:" + tn.token))
	log.InternalDebug(fmt.Sprintln("nexttokentype:" + strconv.Itoa(tn.nextTokenType)))
	log.InternalDebug(fmt.Sprintln("position:" + strconv.Itoa(tn.position)))
	tn.parseComment()
	log.InternalDebug(fmt.Sprintf("%v \n", tn))
	log.InternalDebug(fmt.Sprintln("token:" + tn.token))
	log.InternalDebug(fmt.Sprintln("nexttokentype:" + strconv.Itoa(tn.nextTokenType)))
	log.InternalDebug(fmt.Sprintln("position:" + strconv.Itoa(tn.position)))
	tn.parseSql()
	log.InternalDebug(fmt.Sprint("%v \n", tn))
	log.InternalDebug(fmt.Sprintln("token:" + tn.token))
	log.InternalDebug(fmt.Sprintln("nexttokentype:" + strconv.Itoa(tn.nextTokenType)))
	log.InternalDebug(fmt.Sprintln("position:" + strconv.Itoa(tn.position)))
	sn := new(SqlPartsNode)
	var node Node = sn
	stype := fmt.Sprintf("%T", node)
	log.InternalDebug(fmt.Sprintln(stype))
	sa := new(SqlAnalyzer)
	sa.Setup(`select/*$pmb.selectHint*/ dfloc.id as id, dfloc.login_id as login_id, dfloc.name as name, dfloc.version_no as version_no
  from li_tbl dfloc
 where dfloc.id = /*pmb.conditionQuery.id.fixed.query.equal*/null`, false)
	nd,_ := sa.Analyze()
	//fmt.Printf("Node %v %T\n",nd,nd)
	var sanode Node = *nd
	//fmt.Println("len",strconv.Itoa(sanode.GetChildSize()))
	if sanode.GetChildSize() != 4 {
		t.Error("len 4 expected")
	}
	for i := 0; i < 4; i++ {
		//fmt.Printf("type %T \n",sanode.GetChild(i))
		var x interface{} = sanode.GetChild(i)
		//fmt.Printf("type %T \n",x)
		var xx Node = *x.(*Node)
		//fmt.Printf("xx %d %v %T\n",i,xx,xx)
		if i == 0 || i == 2 {
			var ib0 *SqlPartsNode = xx.(*SqlPartsNode)
			ib0=ib0
		}
		if i == 1 {
			var ib1 *EmbeddedVariableNode = xx.(*EmbeddedVariableNode)
			if ib1.expression != "pmb.selectHint" {
				t.Error("pmb.selectHint expected but :" + ib1.expression)
			}
		}
		if i == 3 {
			var ib3 *BindVariableNode = xx.(*BindVariableNode)
			if ib3.expression != "pmb.conditionQuery.id.fixed.query.equal" {
				t.Error("pmb.conditionQuery.id.fixed.query.equal expected but :" + ib3.expression)
			}
			if ib3.testValue != "null" {
				t.Error("null expected but :" + ib3.testValue)
			}
		}
	}
}
