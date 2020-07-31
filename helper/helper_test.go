package helper

import (
	"fmt"
	"testing"
)


//测试简易版的参数检测,仅仅提取第1和第3个参数
//@todo 第2个参数为注入资源对应的配置结构体，不需要提取
//@author sam@2020-07-29 09:50:01
func TestTransformArgs(t *testing.T) {
	diName,lazy,err := TransformArgs("redis","",false)
	fmt.Printf("diName=%s,lazy=%v",diName,lazy)
	if err !=nil{
		t.Error(err)
	}

}

//外部通过注入别名获取资源的时候，获取资源名
//@author sam@2020-07-29 10:13:05
func TestGetDiName(t *testing.T) {
	//只传递一个参数的时候，取第一个参数作为资源名
	dn := "dn"
	a1 := GetDiName(dn)
	if a1 != dn {
		t.Error("must be default")
		return
	}
    //传递多个的时候，取第二个
	a2 := GetDiName(dn, "22")
	if a2 != "22" {
		t.Error("must be args[0]")
		return
	}
}

//测试Map中的key转字符Arr
//todo 注意map是无序的，所以arr[0]可能是k1也可能是k2
//@author sam@2020-07-29 10:18:44
func TestMapToArray(t *testing.T) {
	mp := map[string]interface{}{
		"name": "sam",
		"age": 1,
	}
	arr := MapToArray(mp)
	fmt.Printf("%+v\r\n",arr)

	if len(arr) != 2 {
		t.Error("length of array is not equal 2")
		return
	}

	if arr[0] == "name" {
		if arr[1] != "age" {
			t.Error("part result of array is error")
			return
		}
	} else if arr[0] == "age" {
		if arr[1] != "name" {
			t.Error("part result of array is error")
			return
		}
	} else {
		t.Error("result of array is error")
		return
	}
}
