package dftest

import (
	"fmt"
	"testing"
	"github.com/mikeshimura/dbflute/df"
)

func TestConditionKey(t *testing.T) {
	var c1 df.ConditionKey
	var c2 df.ConditionKey
	var c3 df.ConditionKey
	c1 = df.CK_EQ
	c2 = df.CK_EQ
	c3 = df.CK_GT

	if c1 != c2 {
		t.Error("Expected c1 == c2 but not equal ")
	}
	if c2 == c3 {
		t.Error("Expected c2 != c2 but equal ")
	}
	if c1.GetConditionKeyS() != df.C_EQ {
		t.Error("Expected ==C_EQ but:" + c1.GetConditionKeyS())
	}
	if c3.GetConditionKeyS() != df.C_GT {
		t.Error("Expected ==C_GT but:" + c3.GetConditionKeyS())
	}
	sname := new(df.ColumnSqlName)
	sname.ColumnSqlName = "tTest-"
	sname.AnalyzeIrregularChar()
	if sname.IrregularChar == false {
		t.Error("Expected IrregulaChar is true but false ")
	}
	sname.ColumnSqlName = "tTest"
	sname.AnalyzeIrregularChar()
	if sname.IrregularChar {
		t.Error("Expected IrregulaChar is false but true ")
	}
	cv := new(df.ConditionValue)
	var v2 int64 = 12
	cv.SetupFixedValue(&c2, v2)
	cv.SetupFixedValue(&c3, "vv")

	if cv.Fixed["query"]["equal"].(int64) != v2 {
		t.Error("Expected 12 but " + 
			fmt.Sprintf("%v", cv.Fixed["query"]["equal"]))
	}
	cv.SetupEqual("TEST","LOC")
	cv.SetupGreaterThan(123,"gtloc")
	ev:=cv.EqualValueHandler.GetValue()
	fmt.Printf("%v\n",ev)
	gtv:=cv.GreaterThanValueHandler.GetValue()
	fmt.Printf("%v\n",gtv)

}
