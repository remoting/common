package log

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"runtime"

	"github.com/natefinch/lumberjack"
)

var errorLog, debugLog, warnLog *log.Logger
var _level int = -999

type LogConfig struct {
	Prefix     string
	LogDir     string
	Level      int
	MaxSize    int
	MaxBackups int
}

// InitConfig ...
// 5=debug,10=warn,15=error
func InitConfig(conf LogConfig) {
	if len(conf.LogDir) <= 0 {
		conf.LogDir = "./logs/"
	}
	if conf.Level <= 0 {
		conf.Level = 5
	}
	if conf.MaxBackups <= 0 {
		conf.MaxBackups = 3
	}
	_level = conf.Level
	errorLogFile := &lumberjack.Logger{
		Filename:   conf.LogDir + "error.log",
		MaxSize:    conf.MaxSize, // megabytes
		MaxBackups: conf.MaxBackups,
		MaxAge:     28, // days
	}
	errorLog = log.New(io.MultiWriter(os.Stdout, errorLogFile), "["+conf.Prefix+"-ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)

	debugLogFile := &lumberjack.Logger{
		Filename:   conf.LogDir + "debug.log",
		MaxSize:    conf.MaxSize, // megabytes
		MaxBackups: conf.MaxBackups,
		MaxAge:     28, // days
	}
	debugLog = log.New(io.MultiWriter(os.Stdout, debugLogFile), "["+conf.Prefix+"-DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)

	warnLogFile := &lumberjack.Logger{
		Filename:   conf.LogDir + "warn.log",
		MaxSize:    conf.MaxSize, // megabytes
		MaxBackups: conf.MaxBackups,
		MaxAge:     28, // days
	}
	warnLog = log.New(io.MultiWriter(os.Stdout, warnLogFile), "["+conf.Prefix+"-WARN] ", log.Ldate|log.Ltime|log.Llongfile)
}

func Error(format string, v ...interface{}) {
	if _level == -999 {
		fmt.Printf("Log Uninitialized:"+format, v...)
	} else {
		if _level <= 15 {
			errorLog.Output(2, fmt.Sprintf(format, v...))
			errorLog.Output(2, fmt.Sprintf("%s", Stack()))
		}
	}
}
func Debug(format string, v ...interface{}) {
	if _level == -999 {
		fmt.Printf("Log Uninitialized:"+format, v...)
	} else {
		if _level <= 5 {
			debugLog.Output(2, fmt.Sprintf(format, v...))
		}
	}
}
func Warn(format string, v ...interface{}) {
	if _level == -999 {
		fmt.Printf("Log Uninitialized:"+format, v...)
	} else {
		if _level <= 10 {
			warnLog.Output(2, fmt.Sprintf(format, v...))
		}
	}
}

func Stack() []byte {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			buf = buf[:n]
			break
		}
		buf = make([]byte, 2*len(buf))
	}
	line := []byte("\n")
	data := bytes.Split(buf, line)
	data = append(data[:1], data[5:]...)
	return bytes.Join(data, line)
}
