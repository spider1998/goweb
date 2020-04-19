package log

import (
	"fmt"
	"goweb/pkg/config"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

type Logger interface {
	Run(...Logger)
	Debugf(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Errorf(format string, a ...interface{})
	Fatalf(format string, a ...interface{})
	Warnf(format string, a ...interface{})
}

type WebLogger struct {
	Debug          bool
	Level          int
	LogChan        chan string
	LogMsg         []string
	SingleCapacity int
	RuntimePath    string
	Truncate       Truncate
}

type Truncate struct {
	Mode      TruncateMode
	TimingDay int
	FrontDay  int
	LimitSize int
}

type priority int

type TruncateMode int8

const (
	errorFormat = "[%s]【%s】 %s:%d %s\n"
	logFormat   = "[%s]【%s】 %s\n"
	timeFormat  = "2006-01-02 15:04:05"

	priorityFatal priority = iota
	priorityError
	priorityWarn
	priorityInfo
	priorityDebug

	TruncateModeDay  TruncateMode = 1
	TruncateModeSize TruncateMode = 2
)

func (p priority) String() string {
	switch p {
	case priorityFatal:
		return "fatal"
	case priorityError:
		return "error"
	case priorityWarn:
		return "warn"
	case priorityInfo:
		return "info"
	case priorityDebug:
		return "debug"
	}

	return ""
}

func (l *WebLogger) truncate() {

}

func (l WebLogger) Run(logger ...Logger) {
	WriteLogFunc := func() {
		var logs []string
		logs = l.LogMsg
		l.cleanMsg()
		for _, logRecord := range logs {
			writeLog(l.RuntimePath, logRecord)
		}
	}
	for {
		select {
		case msg := <-l.LogChan:
			l.LogMsg = append(l.LogMsg, msg)
			if len(l.LogMsg) < l.SingleCapacity {
			} else {
				go WriteLogFunc()
			}
		case <-time.After(time.Minute * 60 * 24):
			go WriteLogFunc()
		}
	}
}

func (l *WebLogger) cleanMsg() {
	l.LogMsg = []string{}
}

func (l WebLogger) Debugf(format string, a ...interface{}) {
	if os.Getenv("DEBUG") != "" {
		logf(os.Stderr, priorityDebug, format, a...)
	}
}

func (l WebLogger) Infof(format string, a ...interface{}) {
	l.LogChan <- priorityInfo.String() + ":" + logf(os.Stdout, priorityInfo, format, a...)

}

func (l WebLogger) Errorf(format string, a ...interface{}) {
	l.LogChan <- priorityError.String() + ":" + logf(os.Stderr, priorityError, format, a...)
}

func (l WebLogger) Fatalf(format string, a ...interface{}) {
	l.LogChan <- priorityFatal.String() + ":" + logf(os.Stderr, priorityFatal, format, a...)
	os.Exit(1)
}

func (l WebLogger) Warnf(format string, a ...interface{}) {
	l.LogChan <- priorityWarn.String() + ":" + logf(os.Stderr, priorityWarn, format, a...)
}

func logf(stream io.Writer, level priority, format string, a ...interface{}) string {
	var prefix string

	if level <= priorityError || level == priorityDebug {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "<unknown file>"
			line = -1
		} else {
			file = file[strings.LastIndex(file, "/")+1:]
		}
		prefix = fmt.Sprintf(errorFormat, time.Now().Format(timeFormat), level.String(), file, line, format)
	} else {
		prefix = fmt.Sprintf(logFormat, time.Now().Format(timeFormat), level.String(), format)
	}
	return fmt.Sprintf(prefix, a...)
}

func NewLogger() Logger {
	var l WebLogger
	l.SingleCapacity = config.GlobalConfig.Log.SingleCapacity
	l.RuntimePath = config.GlobalConfig.Log.RuntimePath
	l.LogChan = make(chan string)
	return l
}
