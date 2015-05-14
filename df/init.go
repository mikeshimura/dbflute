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
	//"fmt"
	"os"
)

func init() {
	//fmt.Println("df init:")
	BhvUtil_I = new(BhvUtil)
	BhvUtil_I.SetUp()
	DBMetaProvider_I = CreateDBMetaProvider()
	CK_EQ = new(CK_EQ_T)
	CK_EQ.ConditionKeyS = C_EQ
	CK_EQ.Operand = "="
	var CK_EQ_ck ConditionKey = CK_EQ
	CK_EQ.conditionKey = &CK_EQ_ck
	CK_EQ_C=&CK_EQ_ck
	CK_GT = new(CK_GT_T)
	CK_GT.ConditionKeyS = C_GT
	CK_GT.Operand = ">"
	var CK_GT_ck ConditionKey = CK_GT
	CK_GT.conditionKey = &CK_GT_ck
	CK_GT_C =&CK_GT_ck
	CK_NE = new(CK_NE_T)
	CK_NE.ConditionKeyS = C_NE
	CK_NE.Operand = "<>"
	var CK_NE_ck ConditionKey = CK_NE
	CK_NE.conditionKey = &CK_NE_ck
	CK_NE_C = &CK_NE_ck
	CK_LT = new(CK_LT_T)
	CK_LT.ConditionKeyS = C_LT
	CK_LT.Operand = "<"
	var CK_LT_ck ConditionKey = CK_LT
	CK_LT.conditionKey = &CK_LT_ck
	CK_LT_C=&CK_LT_ck
	CK_GE = new(CK_GE_T)
	CK_GE.ConditionKeyS = C_GE
	CK_GE.Operand = ">="
	var CK_GE_ck ConditionKey = CK_GE
	CK_GE.conditionKey = &CK_GE_ck
	CK_GE_C = &CK_GE_ck
	CK_LE = new(CK_LE_T)
	CK_LE.ConditionKeyS = C_LE
	CK_LE.Operand = "<="
	var CK_LE_ck ConditionKey = CK_LE
	CK_LE.conditionKey = &CK_LE_ck
	CK_LE_C = &CK_LE_ck
	CK_ISN = new(CK_ISN_T)
	CK_ISN.ConditionKeyS = C_ISN
	CK_ISN.Operand = "is null"
	var CK_ISN_ck ConditionKey = CK_ISN
	CK_ISN.conditionKey = &CK_ISN_ck
	CK_ISN_C = &CK_ISN_ck
	CK_ISNN = new(CK_ISNN_T)
	CK_ISNN.ConditionKeyS = C_ISNN
	CK_ISNN.Operand = "is not null"
	var CK_ISNN_ck ConditionKey = CK_ISNN
	CK_ISNN.conditionKey = &CK_ISNN_ck
	CK_ISNN_C =  &CK_ISNN_ck
	CK_ISNOE = new(CK_ISNOE_T)
	CK_ISNOE.ConditionKeyS = C_ISNOE
	CK_ISNOE.Operand = "is null"
	var CK_ISNOE_ck ConditionKey = CK_ISNOE
	CK_ISNOE.conditionKey = &CK_ISNOE_ck
	CK_ISNOE_C = &CK_ISNOE_ck
	CK_LS = new(CK_LS_T)
	CK_LS.ConditionKeyS = C_LS
	CK_LS.Operand = "like"
	var CK_LS_ck ConditionKey = CK_LS
	CK_LS.conditionKey = &CK_LS_ck
	CK_LS_C = &CK_LS_ck
	CK_NLS = new(CK_NLS_T)
	CK_NLS.ConditionKeyS = C_NLS
	CK_NLS.Operand = "not like"
	var CK_NLS_ck ConditionKey = CK_NLS
	CK_NLS.conditionKey = &CK_NLS_ck
	CK_NLS_C =  &CK_NLS_ck
	CK_GTISN = new(CK_GTISN_T)
	CK_GTISN.ConditionKeyS = C_GTISN
	CK_GTISN.Operand = ">"
	var CK_GTISN_ck ConditionKey = CK_GTISN
	CK_GTISN.conditionKey = &CK_GTISN_ck
	CK_GTISN_C = &CK_GTISN_ck
	CK_GEISN = new(CK_GEISN_T)
	CK_GEISN.ConditionKeyS = C_GEISN
	CK_GEISN.Operand = ">="
	var CK_GEISN_ck ConditionKey = CK_GEISN
	CK_GEISN.conditionKey = &CK_GEISN_ck
	CK_GEISN_C = &CK_GEISN_ck
	CK_LTISN = new(CK_LTISN_T)
	CK_LTISN.ConditionKeyS = C_LTISN
	CK_LTISN.Operand = "<"
	var CK_LTISN_ck ConditionKey = CK_LTISN
	CK_LTISN.conditionKey = &CK_LTISN_ck
	CK_LTISN_C = &CK_LTISN_ck
	CK_LEISN = new(CK_LEISN_T)
	CK_LEISN.ConditionKeyS = C_LEISN
	CK_LEISN.Operand = "<="
	var CK_LEISN_ck ConditionKey = CK_LEISN
	CK_LEISN.conditionKey = &CK_LEISN_ck
	CK_LEISN_C =  &CK_LEISN_ck
		CK_INS = new(CK_INS_T)
	CK_INS.ConditionKeyS = C_INS
	CK_INS.Operand = "in"
	var CK_INS_ck ConditionKey = CK_INS
	CK_INS.conditionKey = &CK_INS_ck
	CK_INS_C =  &CK_INS_ck
	LoopVariableType := new(LoopVariableType_T)
	LoopVariableType.codeValueMap = make(map[string]*Node)
	CreateDBMetaInstanceHandle()
	Ln = "\n"
	Gopath = os.Getenv("GOPATH")

	d_int64 := func() *Entity {
		//tbl := new(O_Name)
		var te Entity = new(D_Int64)
		return &te
	}

	BhvUtil_I.AddEntity("D_Int64", d_int64)
	var dm DBMeta
	Create_D_Int64Dbm()
	dm = D_Int64Dbm
	DBMetaProvider_I.TableDbNameInstanceMap["D_Int64"] = &dm
	DBMetaProvider_I.TableDbNameFlexibleMap.Put("D_Int64", "D_Int64")

}
