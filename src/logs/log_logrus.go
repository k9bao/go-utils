package logs

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

var (
	flagColor    = flag.Bool("Log.color", true, "logger with color")
	flagDebug    = flag.Bool("Log.debug", false, "logger with debug")
	flagLongTime = flag.Bool("Log.longtime", false, "logger with longtime")
	Log          *logrus.Logger
)

func NewLoggerWithRotate() *logrus.Logger {
	if Log != nil {
		return Log
	}

	path := "log/test.log"
	writer, _ := rotatelogs.New(
		path+"%Y%m%d%H%M",
		rotatelogs.WithLinkName(path),               // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(60*time.Second),       // 文件最大保存时间
		rotatelogs.WithRotationTime(20*time.Second), // 日志切割时间间隔
	)
	writerTrack, _ := os.Create("log/track.log")

	pathMap := lfshook.WriterMap{
		logrus.PanicLevel: writer,
		logrus.FatalLevel: writer,
		logrus.ErrorLevel: writer,
		logrus.WarnLevel:  writer,
		logrus.InfoLevel:  writer,
		logrus.DebugLevel: writer,
		logrus.TraceLevel: writerTrack,
	}

	Log = logrus.New()
	Log.Hooks.Add(lfshook.NewHook(
		pathMap,
		&logrus.TextFormatter{},
	))

	return Log
}

type MyFormatter struct {
	Fmt   logrus.TextFormatter
	Color bool
}

func (f *MyFormatter) GetSource(entry *logrus.Entry) []byte {
	file := ""
	if entry.Caller != nil {
		file = entry.Caller.File
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
		file = fmt.Sprintf("[%v:%v]", file, entry.Caller.Line)
	}
	return []byte(file)
}

func (f *MyFormatter) GetColor(entry *logrus.Entry) string {
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = 32 // gray
	case logrus.WarnLevel:
		levelColor = 33 // yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 31 // red
	default:
		levelColor = 36 // blue
	}
	return fmt.Sprintf("\x1b[%dm", levelColor)
}

func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	//[INFO] 2020-12-15 18:39:33 - message [common/log_logrus_test.go:16] , INFO 和 logText 有颜色显示
	var buffer bytes.Buffer
	buffer.Grow(100)
	if f.Color {
		buffer.WriteString(f.GetColor(entry))
	}
	buffer.WriteByte('[')                                          //[]
	buffer.WriteString(strings.ToUpper(entry.Level.String()[0:4])) //INFO
	buffer.WriteByte(']')                                          //]
	if f.Color {
		buffer.WriteString("\x1b[0m")
	}
	buffer.WriteByte(' ')                                        //space
	buffer.WriteString(entry.Time.Format(f.Fmt.TimestampFormat)) //2020-12-15 18:39:33
	buffer.WriteString(" - ")                                    //space-space
	if f.Color {
		buffer.WriteString(f.GetColor(entry))
	}
	buffer.WriteString(entry.Message) //message
	if len(entry.Data) > 0 {
		needVertical := true
		for k, v := range entry.Data {
			if k == "" {
				continue
			}
			if needVertical {
				buffer.WriteString(" | ")
				needVertical = false
			}
			buffer.WriteString(k)
			buffer.WriteByte(':')
			buffer.WriteString(fmt.Sprintf("%v", v))
			buffer.WriteByte(' ')
		}
	}
	if f.Color {
		buffer.WriteString("\x1b[0m")
	}
	buffer.Write(f.GetSource(entry)) //source:line
	buffer.WriteByte('\n')
	return buffer.Bytes(), nil
}

func SetLogInfo() {
	Log.SetReportCaller(true)
	if *flagDebug {
		Log.SetLevel(logrus.DebugLevel)
	} else {
		Log.SetLevel(logrus.InfoLevel)
	}
	timeFormat := "15:04:05.000"
	if *flagLongTime {
		timeFormat = "2006-01-02 15:04:05"
	}
	Log.SetFormatter(&MyFormatter{Fmt: logrus.TextFormatter{
		FullTimestamp:          true,
		TimestampFormat:        timeFormat,
		DisableLevelTruncation: false,
	}, Color: *flagColor})
	Log.Warnf("Log.debug = %v,Log.color = %v,Log.longtime=%v", flagDebug, flagColor, flagLongTime)
}

func init() {
	Log = logrus.New()
	// Log = NewLoggerWithRotate()
	SetLogInfo() //default set，here flagDebug,flagColor is use default
}
