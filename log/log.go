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
	"fmt"
	"github.com/cihub/seelog"
	"time"
)

var ConfigFile = "seelog.xml"
var InternalDebugFlag = false
var InternalDebugTest = false
var isEnabled = true
var LogInfoFunc = seeInfo
var LogDebugFunc = seeDebug
var LogErrorFunc = seeError
var LogWarnFunc = seeWarn
var LogCriticalFunc = seeCritical
var LogFlushFunc = seeFlush

var timeFormat = "2006/01/02 15:04:05.000 +0900"

func seeInfo(param string) {
	seelog.Info(time.Now().Format(timeFormat),
		" "+param)
}
func seeDebug(param string) {
	seelog.Debug(time.Now().Format(timeFormat),
		" "+param)
}
func seeWarn(param string) {
	seelog.Warn(time.Now().Format(timeFormat),
		" "+param)
}
func seeError(param string) {
	seelog.Error(time.Now().Format(timeFormat),
		" "+param)
}
func seeCritical(param string) {
	seelog.Critical(time.Now().Format(timeFormat),
		" "+param)
}

func seeFlush() {
	seelog.Flush()
}

func SetTimeFormat(format string) {
	timeFormat = format
}

func InternalDebug(log string) {
	if InternalDebugFlag == false {
		return
	}
	if InternalDebugTest {
		fmt.Println("InternalDebug" + ":" + log)
	} else {
		LogDebugFunc("InternalDebug" + ":" + log)
	}
}
func Flush() {
	LogFlushFunc()
}
func Info(pack string, log string) {
	LogInfoFunc(pack + ":" + log)
}
func Debug(pack string, log string) {
	LogDebugFunc(pack + ":" + log)
}
func Warn(pack string, log string) {
	LogWarnFunc(pack + ":" + log)
}
func Error(pack string, log string) {
	LogErrorFunc(pack + ":" + log)
}
func Critical(pack string, log string) {
	LogCriticalFunc(pack + ":" + log)
}

func IsEnabled() bool {
	return isEnabled
}

func DisAble() {
	isEnabled = false
}

func Enable() {
	isEnabled = true
}
