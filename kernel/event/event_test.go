package event

import (
	"fmt"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)


//event.On("exit",f1,f2,f3)
//event.On("wait",f1)
func TestEvent(t *testing.T) {
	i := []int{}
	f := func(s interface{}) {
		fmt.Println("f ...")
		i = append(i, 1)
	}
	f2 := func(s interface{}) {
		i = append(i, 2)
		i = append(i, 3)
	}

	Convey("events package test", t, func() {
		//1.初始化事件包应该成功
		//Convey("init events package should be success", func() {
		//	fmt.Printf("%+v\r\n",i)
		//	fmt.Printf("%+v\r\n",Events)
		//
		//	So(len(i), ShouldEqual, 0)
		//	So(len(Events[EXIT]), ShouldEqual, 0)
		//	fmt.Println("-------------------------------------------------")
		//})
		//2.空事件执行关闭功能不应成功
		//Convey("empty events execute Off function should not be success", func() {
		//	err:=Off(EXIT, f)
		//	if err !=nil{
		//		fmt.Println(err)
		//	}
		//	So(Off(EXIT, f), ShouldNotBeNil)
		//	fmt.Println("-------------------------------------------------")
		//})
		//3.一个函数的多执行On功能不应该成功
		//Convey("multi execute On function for a function should not be success", func() {
		//	So(On(EXIT, f), ShouldBeNil)
		//	So(On(EXIT, f), ShouldNotBeNil)
		//
		//	fmt.Printf("%+v\r\n",i)
		//	fmt.Printf("%+v\r\n",Events)
		//
		//	fmt.Println("-------------------------------------------------")
		//})

		//4.执行Emit功能应该成功
		//Convey("execute Emit function should be success", func() {
		//	Emit(EXIT, nil)
		//	So(len(i), ShouldEqual, 1)
		//})


		//5.全面测试
		Convey("events package should be work", func() {
			fmt.Println("on ...")
			So(On(EXIT, f), ShouldBeNil)
			So(On(EXIT, f2), ShouldBeNil)
			fmt.Printf("%+v\r\n",i)
			fmt.Printf("%+v\r\n",Events)
			fmt.Println("off ...")
			//So(Off(EXIT, f), ShouldBeNil)
			//fmt.Printf("%+v\r\n",i)
			//fmt.Printf("%+v\r\n",Events)

			Wait()
			Emit(EXIT,nil)


		})


	})
}
