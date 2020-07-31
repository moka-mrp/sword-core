package close

import (
	"fmt"
	"testing"
)


//-----------------------定义一个实现Closer接口的结构体实例------------------------------
type samClose struct {
}

func (m *samClose) Close() error {
	fmt.Println("sam is closed.")
	return nil
}



//测试资源关闭
//@author sam@2020-07-31 15:04:04
func TestRegister(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			t.Error(e)
		}

	}()

	cl := new(samClose)
	Register(cl)
	//MultiRegister(cl)
	Free()
}
