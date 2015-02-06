package log

import (
	"github.com/cihub/seelog"
	"time"
	"fmt"
)
var ConfigFile="seelog.xml"
var isEnabled = true

var charCode = "sjis"
var timeFormat = "2006/01/02 15:04:05.000 +0900"

func SetCharCode(code string) {
	charCode = code
}
func SetTimeFormat(format string) {
	timeFormat = format
}
func Init() {
	fmt.Println("Log Init")
	logger, err := seelog.LoggerFromConfigAsFile(ConfigFile)

	if err != nil {
		fmt.Println("df006:fail to load log config")
	}

	seelog.ReplaceLogger(logger)
}
func init() {
	fmt.Println("Log Init")
	logger, err := seelog.LoggerFromConfigAsFile(ConfigFile)

	if err != nil {
		fmt.Println("df006:fail to load log config")
	}

	seelog.ReplaceLogger(logger)
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
func UtfInfo(pack string, log string) {
	seelog.Info(time.Now(), " "+pack+":"+convErr(log))
}
func UtfDebug(pack string, log string) {
	seelog.Debug(time.Now().Format(timeFormat),
		" "+pack+":"+convErr(log))
}
func UtfWarn(pack string, log string) {
	seelog.Warn(time.Now().Format(timeFormat),
		" "+pack+":"+convErr(log))
}
func UtfError(pack string, log string) {
	seelog.Error(time.Now().Format(timeFormat),
		" "+pack+":"+convErr(log))
}
func UtfCritical(pack string, log string) {
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