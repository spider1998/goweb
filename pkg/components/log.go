package components

type Logger interface {
	Run(...Logger)
	Debugf(format string, a ...interface{})
	Infof(format string, a ...interface{})
	Errorf(format string, a ...interface{})
	Fatalf(format string, a ...interface{})
	Warnf(format string, a ...interface{})
	WriteStorage(msg string)
}
