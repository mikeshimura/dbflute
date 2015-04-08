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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

func TestUtil(t *testing.T) {
	var s Stack
	s.Push(1)
	s.Push(2)
	s.Push(1)
	s.Push(2)
	if len(s.data) != 4 {
		t.Error("expected 4 but :" + strconv.Itoa(len(s.data)))
	}
	v := (s.Pop()).(int)
	//fmt.Println(v)
	v = (s.Pop()).(int)
	//fmt.Println(v)
	v = (s.Pop()).(int)
	//fmt.Println(v)
	v = (s.Pop()).(int)
	//fmt.Println(v)
	if v != 1 {
		t.Error("expected 1 but :" + strconv.Itoa(v))
	}
	v2, ok := (s.Pop()).(int)
	if ok {
		fmt.Println(v2)
	}
	var l List
	l.Add(1)
	l.Add(2)
	l.Add(3)
	l.Add(4)
	if l.Size() != 4 {
		t.Error("expected 4 but :" + strconv.Itoa(l.Size()))
	}
	ll, ok := (l.Get(2)).(int)
	if ll != 3 {
		t.Error("expected 3 but :" + strconv.Itoa(ll))
	}
	dir, err := filepath.Abs(".")
	if err != nil {
		fmt.Println("PATH ERROR ")
	} else {
		fmt.Println(dir)
	}
	gopath := os.Getenv("GOPATH")
	p := filepath.Join(gopath, "src/dbflute/adf/bhv/sql", "LiTblBhv_test1Name.sql")
	fmt.Println("gopath:=", gopath, " p:=", p)
	buf,_ := ioutil.ReadFile(p)
	fmt.Println(string(buf))
}
