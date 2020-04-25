package components

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
	"time"

	"goweb/pkg/util"
)

const (
	_LineFeed    = "\r\n"
	_UnknownFile = "<unknown file>"
	_DirSymbol   = "/"
	_LogSuffix   = ".log"
)

func NewLogger() Logger {
	var l WebLogger
	l.SingleCapacity = GlobalConfig.Log.SingleCapacity
	l.RuntimePath = GlobalConfig.Log.RuntimePath
	l.LogChan = make(chan string)
	return l
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
			l.WriteStorage(logRecord)
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
			file = _UnknownFile
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

func (l WebLogger) WriteStorage(msg string) {
	r, _ := regexp.Compile(`\d{4}-\d{2}-\d{2}`)
	times := r.FindString(msg)

	var (
		err error
		f   *os.File
	)
	path := l.RuntimePath + strings.Split(msg, ":")[0]
	if !util.IsExist(path) {
		if err = util.CreateDir(path); err != nil {
			logf(os.Stderr, priorityFatal, err.Error())
		}
	}

	f, err = os.OpenFile(path+_DirSymbol+times+_LogSuffix, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	_, err = io.WriteString(f, _LineFeed+msg)

	defer f.Close()
	return
}
