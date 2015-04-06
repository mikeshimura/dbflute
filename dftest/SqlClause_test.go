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
