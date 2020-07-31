package logger



func Trace(args ...interface{}) {
	GetLogger().Trace(args...)
}
func Debug(args ...interface{}) {
	GetLogger().Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	GetLogger().Debugf(format,args...)
}

func Info(args ...interface{}) {
	GetLogger().Info(args...)
}
func Warn(args ...interface{}) {
	GetLogger().Warn(args...)
}

func Warning(args ...interface{}) {
	GetLogger().Warning(args...)
}

func Error(args ...interface{}) {
	GetLogger().Error(args...)
}

func Fatal(args ...interface{}) {
	GetLogger().Fatal(args...)
}

func Panic(args ...interface{}) {
	GetLogger().Panic(args...)
}

