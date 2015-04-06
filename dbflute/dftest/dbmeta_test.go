package dftest

import (
"fmt"
	"testing"
	"github.com/mikeshimura/dbflute/df"
//	"encoding/json"
//	"reflect"
)
func TestDbmeta(t *testing.T) {
	sm:=df.CreateAsFlexible()
	sm.Put("test-a","AAA")
	sm.Put("tes_AA",123)
	r:=sm.Get("tesAA")
	fmt.Printf("%v\n",sm.SearchMap)
	if r.(int)!=123 {
		t.Error("Expected Error for id length 0 but no error ")
	}

}