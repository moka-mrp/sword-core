package container

import (
	"fmt"
	"testing"
)

//测试资源注入容器以及从容器中取出来
//@author sam@2020-07-31 10:58:34
func TestContainerSetSingleton(t *testing.T) {
	//注入资源di1
	App.SetSingleton("di1", "1")
	//注入资源di2
	App.SetSingleton("di2", 2)
	//-------------------------------------------------------
	//获取资源di1
	a1 := App.GetSingleton("di1")
	fmt.Printf("%v\r\n",a1)
	//获取资源di3
	//@todo 获取不存在的资源的时候返回nil
	a3 := App.GetSingleton("di3")
	fmt.Printf("%v",a3)
}

//---------------------创建一个启动资源-------------------------------------------------------
type Single01 struct {
	Name  string `json:"name"`
	Age   int
}

func (*Single01)t001(){
	fmt.Println("t001")
}

func (*Single01)t002(){
	fmt.Println("t002")
}

//-------------------- 创建依赖注入的实例---------------------

type Inject01 struct {
	Name  *Single01 `di:"di1"`
}

func (*Inject01)in01(){
	fmt.Println("in001")
}

func (*Inject01)in02(){
	fmt.Println("in002")
}


//测试容器中的demo
//@author sam@2020-07-31 11:34:17
func TestContainerDemo(t *testing.T) {
	App.SetSingleton("di1",new(Single01))
	App.SetSingleton("di2", "2")
	App.SetPrototype("f01", factoryDemo)
	ret, err := App.GetPrototype("f01")
	fmt.Println(ret,err)
	//打印容器内部
	fmt.Println("----------------------------App.String()---------------------------------------------")
	nameStr := App.String()
	fmt.Println(nameStr)
	//是否单例依赖
	fmt.Println("----------------------------App.isSingleton()------------------------------------------")
	bool1 := App.isSingleton("di1")
	fmt.Printf("%v\r\n",bool1)
	//是否实例依赖
	fmt.Println("----------------------------App.isPrototype()------------------------------------------")
    bool2 := App.isPrototype("di1")
	fmt.Printf("%v\r\n",bool2)
	fmt.Println("----------------------------App.injectName()------------------------------------------")
	strTest := App.injectName("in01,in02")
	fmt.Println(strTest)
	fmt.Println("----------------------------App.Ensure()------------------------------------------")
	in:=new(Inject01)
	//App.Ensure(in)

	in.Name.t001()

}

func factoryDemo() (i interface{}, err error) {
	fmt.Println("call factoryDemo")
	return
}