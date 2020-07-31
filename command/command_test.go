package command

import (
	"testing"
	"fmt"
)

//测试命令行，绑定一个d001，然后再调用这个命令行
//@author sam@2020-07-28 17:30:23
func TestNew(t *testing.T) {
	cmd := New()
	cmd.AddFunc("d001", test)
	cmd.Execute("d001")
	cmd.Execute("test1")
}

func test() {
	fmt.Println("run test")
}
