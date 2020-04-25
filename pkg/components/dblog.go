package components

import (
	"fmt"
	"github.com/go-xorm/core"
)

type DbLogger struct {
	logger  Logger
	level   core.LogLevel
	showSQL bool
}

func makeDbLogger(logger Logger, l core.LogLevel) *DbLogger {
	return &DbLogger{
		logger: logger,
		level:  l,
	}
}

func (s *DbLogger) Error(v ...interface{}) {
	if s.level <= core.LOG_ERR {
		s.logger.Errorf(fmt.Sprint(s.format(v)...))
	}
}

func (s *DbLogger) Errorf(format string, v ...interface{}) {
	if s.level <= core.LOG_ERR {
		s.logger.Errorf(fmt.Sprintf(format, s.format(v)...))
	}
}

func (s *DbLogger) Debug(v ...interface{}) {
	if s.level <= core.LOG_DEBUG {
		s.logger.Debugf(fmt.Sprint(s.format(v)...))
	}
}

func (s *DbLogger) Debugf(format string, v ...interface{}) {
	if s.level <= core.LOG_DEBUG {
		s.logger.Debugf(fmt.Sprintf(format, s.format(v)...))
	}
}

func (s *DbLogger) Info(v ...interface{}) {
	if s.level <= core.LOG_INFO {
		s.logger.Infof(fmt.Sprint(s.format(v)...))
	}
	return
}

func (s *DbLogger) Infof(format string, v ...interface{}) {
	if s.level <= core.LOG_INFO {
		s.logger.Infof(fmt.Sprintf(format, s.format(v)...))
	}
	return
}

func (s *DbLogger) Warn(v ...interface{}) {
	if s.level <= core.LOG_WARNING {
		s.logger.Warnf(fmt.Sprint(s.format(v)...))
	}
	return
}

func (s *DbLogger) Warnf(format string, v ...interface{}) {
	if s.level <= core.LOG_WARNING {
		s.logger.Warnf(fmt.Sprintf(format, s.format(v)...))
	}
	return
}

func (s *DbLogger) Level() core.LogLevel {
	return s.level
}

func (s *DbLogger) SetLevel(l core.LogLevel) {
	s.level = l
	return
}

func (s *DbLogger) ShowSQL(show ...bool) {
	if len(show) == 0 {
		s.showSQL = true
		return
	}
	s.showSQL = show[0]
}

func (s *DbLogger) IsShowSQL() bool {
	return s.showSQL
}

func (s *DbLogger) format(v []interface{}) []interface{} {
	tmpV := make([]interface{}, len(v))
	copy(tmpV, v)
	if len(tmpV) >= 2 {
		if slice, ok := tmpV[1].([]interface{}); ok {
			tmpSlice := make([]interface{}, len(slice))
			copy(tmpSlice, slice)
			for i, item := range tmpSlice {
				switch raw := item.(type) {
				case []byte:
					tmpSlice[i] = fmt.Sprintf("<%d bytes>", len(raw))
					tmpV[1] = tmpSlice
				}
			}
		}
	}
	return tmpV
}
