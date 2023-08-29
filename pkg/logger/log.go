package logger

// Log log interface
type Log interface {
	Print(v ...any)
	Printf(format string, v ...any)
	Println(v ...any)
	Debug(v ...any)
	Debugf(format string, v ...any)
	Info(v ...any)
	Infof(format string, v ...any)
	Warn(v ...any)
	Warnf(format string, v ...any)
	Error(err error, v ...any)
	Errorf(err error, format string, v ...any)
	Fatal(err error, v ...any)
	Fatalf(err error, format string, v ...any)
}
