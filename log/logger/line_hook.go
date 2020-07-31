package logger

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"path"
)
//logrus 代码日志文件，函数和代码行位置输出hook
//@author sam@2020-04-03 10:13:23
type LineHook struct {
	EnableFileNameLog bool //启用文件名称log
	EnableFuncNameLog bool 	//启用函数名称log
}

func NewLineHook(logger *log.Logger,enableFileNameLog, enableFuncNameLog bool) (*LineHook, error) {
	lh := new(LineHook)
	lh.EnableFileNameLog=enableFileNameLog
	lh.EnableFuncNameLog=enableFuncNameLog
	return lh, nil
}

func (hooks LineHook) Levels() []log.Level {
	return log.AllLevels
}
func (hook *LineHook) Fire(entry *log.Entry) error {
	var (
		file, function string
		line           int
	)
	if entry.HasCaller() {
		frame := entry.Caller
		line = frame.Line
		function = frame.Function
		dir, filename := path.Split(frame.File)
		f := path.Base(dir)
		file = fmt.Sprintf("%s/%s", f, filename)
	}
	if hook.EnableFileNameLog && hook.EnableFuncNameLog {
		entry.Message = fmt.Sprintf("[%s(%s:%d)] %s", function, file, line, entry.Message)
	}
	if hook.EnableFileNameLog && !hook.EnableFuncNameLog {
		entry.Message = fmt.Sprintf("[%s(%d)] %s", file, line, entry.Message)
	}
	if !hook.EnableFileNameLog && hook.EnableFuncNameLog {
		entry.Message = fmt.Sprintf("[%s(%d)] %s", function, line, entry.Message)
	}
	return nil
}
