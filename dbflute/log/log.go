package log

import (
	"github.com/cihub/seelog"
	 "time"
	 "github.com/mikeshimura/dbflute/converter"
)

func init(){
	    logger, err := seelog.LoggerFromConfigAsFile("seelog.xml")

    if err != nil {
        panic("fail to load config")
    }

    seelog.ReplaceLogger(logger)
}
func Flush(){
	seelog.Flush()
}
func Info(pack string,log string){
	seelog.Info(time.Now()," "+pack+":"+log)
}
func Debug(pack string,log string){
	seelog.Debug(time.Now()," "+pack+":"+log)
}
func Warn(pack string,log string){
	seelog.Warn(time.Now()," "+pack+":"+log)
}
func Error(pack string,log string){
	seelog.Error(time.Now()," "+pack+":"+log)
}
func Critical(pack string,log string){
	seelog.Critical(time.Now()," "+pack+":"+log)
}
func UtfInfo(pack string,log string){
	seelog.Info(time.Now()," "+pack+":"+convErrToSjis(log))
}
func UtfDebug(pack string,log string){
	seelog.Debug(time.Now()," "+pack+":"+convErrToSjis(log))
}
func UtfWarn(pack string,log string){
	seelog.Warn(time.Now()," "+pack+":"+convErrToSjis(log))
}
func UtfError(pack string,log string){
	seelog.Error(time.Now()," "+pack+":"+convErrToSjis(log))
}
func UtfCritical(pack string,log string){
	seelog.Critical(time.Now()," "+pack+":"+convErrToSjis(log))
}
func convErrToSjis(log string)(string){
	res,_ :=converter.ConvertUtoS(log)
	return res
}