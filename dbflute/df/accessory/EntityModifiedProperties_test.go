package accessory

import (
"testing"
	"fmt"
)

func TestEntityModifiedProperties(t *testing.T) {
	ep := new(EntityModifiedProperties)
	err:=ep.AddPropertyName("")
	if err==nil{
		t.Error("Expected Error for id length 0 but no error ")
	}
	err=ep.AddPropertyName("id")
	if err!=nil{
		t.Error("Expected No Error but error ")
	}
	res:=ep.GetPropertyNamesArray()
	fmt.Printf("%v\n",res)
	found:=ep.IsModifiedProperty("xx")
	if found{
		t.Error("Expected Notfound but found ")
	}
	found=ep.IsModifiedProperty("id")
	if !found{
		t.Error("Expected found but not found ")
	}

}