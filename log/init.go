package log

import (
	"fmt"
	"github.com/cihub/seelog"
	"os"
	"path/filepath"
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
	path := ""
	logPath := os.Getenv("LOGPATH")
	goPath := os.Getenv("GOPATH")
	if len(logPath) > 0 {
		path = logPath
	} else {
		path = goPath
	}
	fullpath := filepath.Join(path, "seelog.xml")
	files, _ := filepath.Glob(fullpath)
	var logger seelog.LoggerInterface
	var err error
	if len(files) > 0 {
		logger, err = seelog.LoggerFromConfigAsFile(fullpath)
	} else {
		logger, err = seelog.LoggerFromConfigAsBytes([]byte(appConfig))
	}

	if err != nil {
		fmt.Println("fail to load log config")
	}

	seelog.ReplaceLogger(logger)

}
