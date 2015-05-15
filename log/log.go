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
package log

import (
	"github.com/cihub/seelog"
	"time"
	"fmt"
)
var ConfigFile="seelog.xml"
var InternalDebugFlag = false
var InternalDebugTest = false
var isEnabled = true

var charCode = "sjis"
var timeFormat = "2006/01/02 15:04:05.000 +0900"

func SetCharCode(code string) {
	charCode = code
}
func SetTimeFormat(format string) {
	timeFormat = format
}

func InternalDebug(log string) {
	if InternalDebugFlag==false{
		return
	}
	if InternalDebugTest {
		fmt.Println("InternalDebug"+":"+convErr(log))
	} else {
	seelog.Debug(time.Now().Format(timeFormat),
		"InternalDebug"+":"+convErr(log))
	}
}
func Flush() {
	seelog.Flush()
}
func Info(pack string, log string) {
	seelog.Info(time.Now().Format(timeFormat),
		" "+pack+":"+log)
}
func Debug(pack string, log string) {
	seelog.Debug(time.Now().Format(timeFormat),
		" "+pack+":"+log)
}
func Warn(pack string, log string) {
	seelog.Warn(time.Now().Format(timeFormat),
		" "+pack+":"+log)
}
func Error(pack string, log string) {
	seelog.Error(time.Now().Format(timeFormat),
		" "+pack+":"+log)
}
func Critical(pack string, log string) {
	seelog.Critical(time.Now().Format(timeFormat),
		" "+pack+":"+log)
}
func InfoConv(pack string, log string) {
	seelog.Info(time.Now(), " "+pack+":"+convErr(log))
}
func DebugConv(pack string, log string) {
	seelog.Debug(time.Now().Format(timeFormat),
		" "+pack+":"+convErr(log))
}
func WarnConv(pack string, log string) {
	seelog.Warn(time.Now().Format(timeFormat),
		" "+pack+":"+convErr(log))
}
func ErrorConv(pack string, log string) {
	seelog.Error(time.Now().Format(timeFormat),
		" "+pack+":"+convErr(log))
}
func CriticalConv(pack string, log string) {
	seelog.Critical(time.Now().Format(timeFormat),
		" "+pack+":"+convErr(log))
}
func convErr(log string) string {
	if charCode != "utf8" {
		res, _ := ConvertUtf(log,charCode)
		return res
	}
	return log
}

func IsEnabled () bool {
	return isEnabled
}

func DisAble(){
	isEnabled = false
}

func Enable(){
	isEnabled = true
}