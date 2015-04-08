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
	"testing"
	//"fmt"
	"github.com/mikeshimura/dbflute/df"
	//"github.com/mikeshimura/dbflute/log"
)

func TestBaseEntity(t *testing.T) {
	
//	ep := new(df.BaseEntity)
//	err:=ep.AddPropertyName("")
//	if err==nil{
//		t.Error("Expected Error for id length 0 but no error ")
//	}
//	err=ep.AddPropertyName("id")
//	if err!=nil{
//		t.Error("Expected No Error but error ")
//	}
//	res:=ep.GetPropertyNamesArray()
//	fmt.Printf("%v\n",res)
//	found:=ep.IsModifiedProperty("xx")
//	if found{
//		t.Error("Expected Notfound but found ")
//	}
//	found=ep.IsModifiedProperty("id")
//	if !found{
//		t.Error("Expected found but not found ")
//	}
}

func TestEntityModifiedProperties(t *testing.T) {
	ep := new(df.EntityModifiedProperties)
	err:=ep.AddPropertyName("")
	if err==nil{
		t.Error("Expected Error for id length 0 but no error ")
	}
	err=ep.AddPropertyName("id")
	if err!=nil{
		t.Error("Expected No Error but error ")
	}
//	res:=ep.GetPropertyNamesArray()
//	fmt.Printf("%v\n",res)
	found:=ep.IsModifiedProperty("xx")
	if found{
		t.Error("Expected Notfound but found ")
	}
	found=ep.IsModifiedProperty("id")
	if !found{
		t.Error("Expected found but not found ")
	}

}