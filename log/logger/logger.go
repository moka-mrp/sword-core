package logger

import (
	"bufio"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/mattn/go-colorable"
	"github.com/moka-mrp/sword-core/config"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)



const (
	HandlerFile = "file"
	HandlerStdout = "stdout"
	defaultDayTimePattern  = "-%Y-%m-%d"
	defaultHourTimePattern = "-%Y-%m-%d-%H"
)


//返回一个标准输出即可,logrus默认是os.Stderr
//@author sam@2020-04-02 15:40:37
func GetStdOutWriter() (writer io.Writer) {
	return colorable.NewColorableStdout()
}

//初始化日志配置
//@todo 注意 logrus.New()一定要放到函数内，不可以放置文件开头，否则注入多个不同别名的相同资源的时候，会被覆盖，行为出现诡异
//@todo 注意，当采用日志写入文件的时候,本质是通过钩子再写一份的，所以logger.SetOutput(writer)会导致再写入一份
//@link https://www.jianshu.com/p/1de3d4a4e843
//@author sam@2020-04-02 14:31:53
func InitLog(conf config.LogConfig) (*logrus.Logger, error) {
	var formatter *prefixed.TextFormatter   //格式化
	var logger = logrus.New()
	//1.定义并设置formatter
	formatter = &prefixed.TextFormatter{}
	formatter.ForceColors = true
	formatter.DisableColors = false                 //打开颜色
	formatter.ForceFormatting = true                //开启格式化
	formatter.SetColorScheme(&prefixed.ColorScheme{ //针对不同的日志级别设置不同的颜色
		DebugLevelStyle: "blue",
		InfoLevelStyle:  "green",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "red",
		PanicLevelStyle: "red",
		PrefixStyle:     "cyan",
		TimestampStyle:  "37",
	})
	formatter.FullTimestamp = true                    //开启完整时间戳输出和时间戳格式
	formatter.TimestampFormat = "2006-01-02.15:04:05" //设置时间格式  2006-01-02.15:04:05.000000
	logger.SetFormatter(formatter)
	//2.设置日志等级
	level, err := logrus.ParseLevel(conf.Level)
	if err == nil {
		logger.SetLevel(level)
	}

	//3.开启调用函数、文件、代码行信息的输出
	logger.SetReportCaller(true)
	lineHook,err:= NewLineHook(logger,conf.EnableFileName,conf.EnableFuncName)
	if err != nil {
		return nil, err
	}
	logger.AddHook(lineHook) //调用AddHook时, 将Hook加入到LevelHooks map中
	//4.设置日志输出方式 标准输出或文件
	if conf.Handler == HandlerStdout {
		writer := GetStdOutWriter()
		logger.SetOutput(writer)
	} else {
		//设置切割日志文件输出的日志格式
		formatter.DisableColors = true
		//fileFormatter := &logrus.JSONFormatter{}
		//拼接输出日志位置
		logDir, _ := filepath.Abs(conf.Dir) //这里不求绝对值也是可以的额
		logPath := path.Join(logDir,conf.Name)+".log"  //fmt.Sprintf("%s/%s.log",dir,fileName)
		//判断目录是否存在，不存在创建即可
		dir:=filepath.Dir(logPath)
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			err = os.MkdirAll(dir, 0777)
			if err != nil {
				return nil,err
			}
		}
		//todo 这里应该是取消标准输出，用钩子切割到文件中,而不是再次将标准输出写入文件，否则就写两份了
		//创建写入的日志文件
		//writer,err:=OpenNewFile(logPath)
		//logger.SetOutput(writer)
		writer, err :=GetNullWriter()
		if err != nil {
			return  nil,err
		}
		logger.SetOutput(writer)
		//轮转切割开始
		split_writer,err:=GetRotateWriter(logPath,conf.Name,defaultDayTimePattern,7,24)
		if err != nil {
			return  nil,err
		}
		lfHook := lfshook.NewHook(lfshook.WriterMap{
			logrus.DebugLevel: split_writer, // 为不同级别设置不同的输出目的
			logrus.InfoLevel:  split_writer,
			logrus.WarnLevel:  split_writer,
			logrus.ErrorLevel: split_writer,
			logrus.FatalLevel: split_writer,
			logrus.PanicLevel: split_writer,
		}, formatter)
		//增钩
		logger.AddHook(lfHook)
	}
	return logger, nil
}

//打开创建日志文件
//@author sam@2020-04-06 10:30:17
func OpenNewFile(logPath string) (*os.File,error) {
	//判断目录是否存在，不存在创建即可
	dir:=filepath.Dir(logPath)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			return nil,err
		}
	}
	//打开或创建文件
	writer, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777) //如果已经存在，则在尾部添加写
	if err != nil {
		return nil,err
	}
	return writer, nil
}


//设置滚动日志输出writer
//@author sam@2020-04-03 10:43:13
func GetRotateWriter(logPath string,name string,timePattern string,
	maxAge time.Duration,rotationTime time.Duration) (*rotatelogs.RotateLogs, error) {
	//创建切割writer
	writer, err := rotatelogs.New(
		strings.Replace(logPath,name,name+timePattern,1), // 切割后的文件名称
		rotatelogs.WithLinkName(logPath), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge *24*time.Hour), 	// 设置最大保存时间(7天)
		rotatelogs.WithRotationTime(rotationTime *time.Hour), // 设置日志切割时间间隔(1天)
	)
	if err != nil {
		return  nil,err
	}
	return writer, nil
}

//生成一个空垃圾writer
//@author sam@2020-07-30 16:07:53
func GetNullWriter()(*bufio.Writer,error){
	src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return  nil,err
	}
	writer := bufio.NewWriter(src)
	return writer,nil
}
