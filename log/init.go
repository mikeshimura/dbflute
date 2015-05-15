package log

import (
	"fmt"
	"github.com/cihub/seelog"
)

func init() {
	appConfig := `	
<seelog>
    <outputs>
        <console formatid="nocolor"/>
    </outputs>
    <formats>
        <format id="nocolor"  format="%Level %Msg%n"/>
    </formats>
</seelog>
`
	fmt.Println("Log Init")
	logger, err := seelog.LoggerFromConfigAsBytes([]byte(appConfig))

	if err != nil {
		fmt.Println("df006:fail to load log config")
	}

	seelog.ReplaceLogger(logger)

}
