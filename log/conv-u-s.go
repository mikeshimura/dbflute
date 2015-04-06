package log

import (
 iconv "github.com/djimenez/iconv-go"
)

func ConvertUtf(in string,code string)(output string, err error){
	output, err = iconv.ConvertString(in, "utf-8",code)
	return
}
