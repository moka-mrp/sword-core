package server

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"github.com/moka-mrp/sword-core/kernel/close"

)

const (
	Version     = "1.0.0"
	BuildCommit = "cf53ea9c98b42e937bba08d21b52667c6ae1c9c4"
	BuildDate   = "2020-03-26 17:22:03"
)

type serverInfo struct {
	stop  chan bool  //这个通道中有值，则服务结束
	debug bool  //是否测试模式启动
}

var srv *serverInfo

func init() {
	srv = new(serverInfo)
	srv.stop = make(chan bool, 0)
}

//将进程号写入文件,未指明进程号的，则直接获取当前的进程号写入即可
//@reviser sam@2020-04-14 14:39:43
func WritePidFile(path string, pidArgs ...int) error {
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	var pid int
	if len(pidArgs) > 0 {
		pid = pidArgs[0]
	} else {
		pid = os.Getpid()
	}
	_, err = fd.WriteString(fmt.Sprintf("%d\n", pid))
	return err
}

//读取文件的进程号
func ReadPidFile(path string) (int, error) {
	fd, err := os.Open(path)
	if err != nil {
		return -1, err
	}
	defer fd.Close()

	buf := bufio.NewReader(fd)
	line, err := buf.ReadString('\n')
	if err != nil {
		return -1, err
	}
	line = strings.TrimSpace(line)
	return strconv.Atoi(line)
}

//阻塞等待程序内部的Stop通道信号
//@reviser sam@2020-04-14 15:19:06
func WaitStop() {
	<-srv.stop
}

//关闭整个服务启动过程加载的资源
//@reviser sam@2020-04-14 15:12:59
func CloseService() {
	if srv.debug {
		fmt.Println("close service")
	}
	close.Free()
}

//处理进程的信号量
//@reviser sam@2020-04-14 15:13:38
func HandleSignal(sig os.Signal) {
	switch sig {
	case syscall.SIGINT:
		fallthrough
	case syscall.SIGTERM:
		Stop()
	default:
	}
}

//监听信号量
//@reviser sam@2020-04-14 15:01:19
func RegisterSignal() {
	go func() {
		var sigs = []os.Signal{
			syscall.SIGHUP,
			syscall.SIGUSR1,
			syscall.SIGUSR2,
			syscall.SIGINT,
			syscall.SIGTERM,
		}
		//创建一个接收通道c
		c := make(chan os.Signal)
		//只接收上述定义的那几个信号
		signal.Notify(c, sigs...)
		//阻塞等待有信号的来临
		for {
			sig := <-c //blocked
			HandleSignal(sig)
		}
	}()
}

// HandleUserCmd use to stop/reload the proxy service
func HandleUserCmd(cmd string, pidFile string) error {
	var sig os.Signal

	switch cmd {
	case "stop":
		sig = syscall.SIGTERM
	case "restart":
		//目前api使用endless平滑重启，需要传递此信号，其他只需要平滑关闭就可以了
		sig = syscall.SIGHUP
	default:
		return fmt.Errorf("unknown user command %s", cmd)
	}

	pid, err := ReadPidFile(pidFile)
	if err != nil {
		return err
	}

	if srv.debug {
		fmt.Printf("send %v to pid %d \n", sig, pid)
	}

	proc := new(os.Process)
	proc.Pid = pid
	return proc.Signal(sig)
}

//收到要停止的信号了，我们往stop通道中发送一个元素即可
//@reviser sam@2020-04-14 15:18:19
func Stop() {
	srv.stop <- true
}


//--------------------------------------------------Debug-----------------------------------------------------
//给外部设置启动服务的debug模式
//@author sam@2020-04-13 17:37:46
func SetDebug(debug bool) {
	srv.debug = debug
	return
}
//给外部获取启动服务的debug模式
//@author  sam@2020-04-13 17:38:26
func GetDebug() bool {
	return srv.debug
}
