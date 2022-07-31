package logs

import (
	"flag"
	"fmt"
	"runtime"
	"testing"

	"github.com/sirupsen/logrus"
)

//go test -run TestLog$ -timeout 200s
func TestLog(t *testing.T) {
	flag.Parse()
	fmt.Println("test")
	logs.Log.Debugln("test")
	logs.Log.Infoln("test")
	logs.Log.Warnln("test")
	logs.Log.Errorln("test")
}

func MyReport(fr *runtime.Frame) (function string, file string) {
	function = ""
	file = fr.File
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n += 1
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	file = fmt.Sprintf("%v:%v -", file, fr.Line)
	return
}

//go test -run TestLogrus -timeout 200s
func TestLogrus(t *testing.T) {
	log := logrus.New()
	logs.Log.SetReportCaller(true)
	logs.Log.Level = logrus.DebugLevel
	logs.Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		ForceColors:            true,
		DisableLevelTruncation: false,
		CallerPrettyfier:       MyReport,
	})

	logs.Log.Debugln("github.com/sirupsen/logrus")
	logs.Log.Infoln("github.com/sirupsen/logrus")
	logs.Log.Warnln("github.com/sirupsen/logrus")
	logs.Log.Errorln("github.com/sirupsen/logrus")
}

func TestLogrus2(t *testing.T) {
	log := logrus.New()
	logs.Log.SetReportCaller(true)
	logs.Log.Level = logrus.DebugLevel
	logs.Log.SetFormatter(&MyFormatter{logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        "2006-01-02 15:04:05",
		ForceColors:            true,
		DisableLevelTruncation: true,
		// CallerPrettyfier:       MyReport,
	}, true})

	logs.Log.Debugln("github.com/sirupsen/logrus")
	logs.Log.Infoln("github.com/sirupsen/logrus")
	logs.Log.Warnln("github.com/sirupsen/logrus")
	logs.Log.Errorln("github.com/sirupsen/logrus")
}
