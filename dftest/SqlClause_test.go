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
package dftest

import (
	"bytes"
	"testing"
	"github.com/mikeshimura/dbflute/df"
)

func TestSqlClause(t *testing.T) {
	di := new(df.DBCurrent)
	di.ProjectName = "postgresql"
	var dd df.DBDef = new(df.PostgreSQL)
	di.DBDef = &dd

	sql := df.CreateSqlClauseSub("LiTbl", di)
	tn:=(*sql).GetTableDbName()
	if tn != "LiTbl" {
		t.Error("expected LiTbl but :"+tn)
	}
	
	
	sb:=new(bytes.Buffer)
	sb.WriteString("TEST")
	sb.WriteString("DAYO")
	res:=sb.String()
	if res !="TESTDAYO" {
		t.Error("expect TESTDAYO but:"+res)
	}
}
