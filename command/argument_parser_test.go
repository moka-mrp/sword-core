package command

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

//@link https://github.com/smartystreets/goconvey
//@author sam@2020-08-25 10:51:36
func TestCmdArgumentParser(t *testing.T) {
	var args []string
	var str string

	fmt.Println("----------------------------------")
	Convey("Parse Cmd Arguments ["+str+"]", t, func() {
		args = ParseCmdArguments(str)
		fmt.Printf("%#v\r\n",args)
		So(len(args), ShouldEqual, 0)
	})

	fmt.Println("----------------------------------")
	str = "    "
	Convey("Parse Cmd Arguments ["+str+"]", t, func() {
		args = ParseCmdArguments(str)
		fmt.Printf("%#v\r\n",args)
		So(len(args), ShouldEqual, 0)
	})

	fmt.Println("----------------------------------")

	str = "aa bbb  ccc "
	Convey("Parse Cmd Arguments ["+str+"]", t, func() {
		args = ParseCmdArguments(str)
		fmt.Printf("%#v\r\n",args)
		So(len(args), ShouldEqual, 3)
		So(args[0], ShouldEqual, "aa")
		So(args[1], ShouldEqual, "bbb")
		So(args[2], ShouldEqual, "ccc")
	})
	fmt.Println("----------------------------------")
	str = "' \\\"" //[]string{" \\\""}
	Convey("Parse Cmd Arguments ["+str+"]", t, func() {
		args = ParseCmdArguments(str)
		fmt.Printf("%#v\r\n",args)
		So(len(args), ShouldEqual, 1)
		So(args[0], ShouldEqual, " \\\"")
	})
	fmt.Println("----------------------------------")
	str = `a "b c"` //[]string{"a", "b c"}  注意命令行中双引号括起来的算一个参数额
	Convey("Parse Cmd Arguments ["+str+"]", t, func() {
		args = ParseCmdArguments(str)
		fmt.Printf("%#v\r\n",args)
		So(len(args), ShouldEqual, 2)
		So(args[0], ShouldEqual, "a")
		So(args[1], ShouldEqual, "b c")
	})
	fmt.Println("----------------------------------")
	str = `a '\''"` //[]string{"a", "'"}
	Convey("Parse Cmd Arguments ["+str+"]", t, func() {
		args = ParseCmdArguments(str)
		fmt.Printf("%#v\r\n",args)
		So(len(args), ShouldEqual, 2)
		So(args[0], ShouldEqual, "a")
		So(args[1], ShouldEqual, "'")
	})


	fmt.Println("----------------------------------")
	str = `   \\a   'b c'   c\ d\  ` //[]string{"\\a", "b c", "c d "}
	Convey("Parse Cmd Arguments ["+str+"]", t, func() {
		args = ParseCmdArguments(str)
		fmt.Printf("%#v\r\n",args)
		So(len(args), ShouldEqual, 3)
		So(args[0], ShouldEqual, "\\a")
		So(args[1], ShouldEqual, "b c")
		So(args[2], ShouldEqual, "c d ")
	})

	fmt.Println("----------------------------------")
	str = `\` //[]string{"\\"}
	Convey("Parse Cmd Arguments ["+str+"]", t, func() {
		args = ParseCmdArguments(str)
	   fmt.Printf("%#v\r\n",args)
		So(len(args), ShouldEqual, 1)
		So(args[0], ShouldEqual, "\\")
	})
	fmt.Println("----------------------------------")
	str = `  \   ` //[]string{" "}    \SPACE
	Convey("Parse Cmd Arguments ["+str+"]", t, func() {
		args = ParseCmdArguments(str)
		fmt.Printf("%#v\r\n",args)
		So(len(args), ShouldEqual, 1)
		So(args[0], ShouldEqual, " ")
	})



}
