package converter

import (
 iconv "github.com/djimenez/iconv-go"
)

func ConvertUtoS(in string)(output string, err error){
	output, err = iconv.ConvertString(in, "utf-8","sjis")
	return
}
