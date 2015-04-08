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
