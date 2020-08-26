package command

import (
	"errors"
)

type fmsState int

const (
	stateArgumentOutside fmsState = iota //0
	stateArgumentStart  //1
	stateArgumentEnd    //2
)

//表示到了解析的命令行的行尾了
var errEndOfLine = errors.New("End of line")

type cmdArgumentParser struct {
	s            string
	i            int
	length       int
	state        fmsState
	startToken   byte
	shouldEscape bool
	currArgument []byte
	err          error
}

//返回参数解析结构体指针
//@author sam@2020-08-25 10:36:19
func newCmdArgumentParser(s string) *cmdArgumentParser {
	return &cmdArgumentParser{
		s:            s,
		i:            -1,
		length:       len(s),
		currArgument: make([]byte, 0, 16),
	}
}

//分析命令行参数
//@author  sam@2020-08-25 10:38:30
func (cap *cmdArgumentParser) parse() (arguments []string) {
	for {
		//fmt.Println(cap.i,"state",cap.state )
		//next方法，每次让i的值随着每次的循环递增1  -1 0 1
		cap.next()
		if cap.err != nil {//当遍历到了行末尾的时候会进入这里
		    //fmt.Println("shouldEscape",cap.shouldEscape,"currArgument",cap.currArgument)
			if cap.shouldEscape {
				cap.currArgument = append(cap.currArgument, '\\')
			}
			if len(cap.currArgument) > 0 {
				arguments = append(arguments, string(cap.currArgument))
			}
			return
		}

		switch cap.state {
		case stateArgumentOutside://0 表示需要我们侦测startToken的值
			cap.detectStartToken()
		case stateArgumentStart://1
			if !cap.detectEnd() {
				cap.detectContent()
			}
		case stateArgumentEnd://2
			cap.state = stateArgumentOutside
			arguments = append(arguments, string(cap.currArgument))
			cap.currArgument = cap.currArgument[:0]
		}
	}
}

//该方法每调用一次就会让字段i递增1
//@author sam@2020-08-25 11:22:38
func (cap *cmdArgumentParser) next() {
	//fmt.Println("length",cap.length,"i",cap.i)
	if cap.length-cap.i == 1 { //判断是否到行尾了
		cap.err = errEndOfLine
		return
	}
	cap.i++
}

//与next相反，即将i减少1(i>0否则不做任何处理)
//@author sam@2020-08-25 11:22:28
func (cap *cmdArgumentParser) previous() {
	if cap.i >= 0 {
		cap.i--
	}
}



//根据当前遍历到的字符串侦测startToken的值，同时判断是否需要回滚一次，当侦测到了startToken值之后就将state值置为1
//侦测的时候我们会根据如下几类字符做出不同的响应：
//1.空字符 ===>飘过它
//2.\ ====>将startToken变为0   shouldEscape变为true
//3."和' ===>将startToken变为该值
//4.其余的字符 ===>将startToken变为0  同时  previous
//@author sam@2020-08-25 11:22:15
func (cap *cmdArgumentParser) detectStartToken() {
	//获取当前遍历到的字符
	c := cap.s[cap.i]
	//fmt.Printf("当前遍历的字符为 %c\r\n",c)
	//如果是空，则侦测不到任何结果
	if c == ' ' {
		return
	}
	switch c {
	case '\\'://是否是反斜线\
		cap.startToken = 0
		cap.shouldEscape = true
	case '"', '\'': //是否是双引号或者单引号
		cap.startToken = c
	default://只要不是上述说的那3个字符都会i都会回去一次
		cap.startToken = 0
		cap.previous()
	}
	cap.state = stateArgumentStart //state为1
}



//大部分情况都是返回false,只有两种情况会返回true
//@author sam@2020-08-25 11:22:01
func (cap *cmdArgumentParser) detectEnd() (detected bool) {
	c := cap.s[cap.i]
    //该字符不是"和'
	if cap.startToken == 0 {
		if c == ' ' && !cap.shouldEscape {
			cap.state = stateArgumentEnd //2
			cap.previous()
			return true
		}
		return false
	}
	//当遍历的字符为单引号或者双引号的时候会进入到这里
	if c == cap.startToken && !cap.shouldEscape {
		cap.state = stateArgumentEnd //2
		return true
	}

	return false
}



func (cap *cmdArgumentParser) detectContent() {
	c := cap.s[cap.i]

	if cap.shouldEscape {
		switch c {
		case ' ', '\\', cap.startToken:
			cap.currArgument = append(cap.currArgument, c)
		default:
			cap.currArgument = append(cap.currArgument, '\\', c)
		}
		cap.shouldEscape = false
		return
	}

	if c == '\\' {
		cap.shouldEscape = true
	} else {
		cap.currArgument = append(cap.currArgument, c)
	}
}



//使用我们自定义的 cmdArgumentParser结构体来分析处理传递的字符串s
//@author sam@2020-08-25 10:22:09
func ParseCmdArguments(s string) (arguments []string) {
	return newCmdArgumentParser(s).parse()
}
