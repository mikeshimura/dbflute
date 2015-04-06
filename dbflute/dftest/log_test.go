package dftest

import (
	"fmt"
	"testing"
	"github.com/mikeshimura/dbflute/df"
	"database/sql"
	"time"
)

func TestLog(t *testing.T) {
	db:=new(df.DisplaySqlBuilder)
	var it df.Date
	it.Date=time.Now()
	var ft df.Timestamp
	ft.Timestamp=time.Now()
	var st string="123456"
	var ns sql.NullString
	ns.Valid=true
	ns.String="NULLTEST"
	var testi *df.List=new(df.List)
	testi.Add(it)
	testi.Add(ft)
	testi.Add(st)
	testi.Add(ns)
	fmt.Println(db.BuildDisplaySql("test ? ? ? ?",testi))
}