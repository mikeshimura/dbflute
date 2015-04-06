package df

import (
	"testing"
	"fmt"
	"strings"
)

func TestString(t *testing.T) {
	s1 := InitUnCap("")
	if s1 != "" {
		t.Error("expected null string but :" + s1)
	}
	s2 := InitUnCap("X")
	if s2 != "x" {
		t.Error("expected x string but :" + s2)
	}
		s3 := InitUnCap("XxxAAA")
	if s3 != "xxxAAA" {
		t.Error("expected xxxAAA string but :" + s3)
	}
	fmt.Println(substringFirstRear("AAcA bb-c D"," ","C"))
	
	fmt.Println(strings.Replace("aabbcc","a","x",1))
	list:=splitListTrimmed("aa  : bbb:cccc  ",":")
	fmt.Printf("list %v \n",list)
}
