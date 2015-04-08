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