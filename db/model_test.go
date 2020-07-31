package db

import (
	"fmt"
	"testing"
	"time"

	//go test时需要开启
	_ "github.com/go-sql-driver/mysql"
)
//----------------------------- init -------------------------------------------------------------------------
func init() {
	dbInit(false)
	//todo 这里就不需要在这里取资源了,model中封装了公共的方法
}

//---------- studentModel ---------------------------------------------------------------------------------------------

type studentModel struct {
	Model
}


//-----------------------  测试部分 -------------------------------


//**********增******************
//todo 传递指针，否则您拿不到插入的id返回值的额，别怪我等没提醒你
//@author sam@2020-07-08 11:24:00
func TestInsert(t *testing.T) {
	stu01:=&Student{Age:18,Name:"rick01",ImageUrl:"",Url:"",Status:1,CreatedAt:time.Now(),UpdatedAt:time.Now()}
	stu02:=&Student{Age:19,Name:"rick02",ImageUrl:"",Url:"",Status:1,CreatedAt:time.Now(),UpdatedAt:time.Now()}
	stu03:=&Student{Age:20,Name:"rick03",ImageUrl:"",Url:"",Status:1,CreatedAt:time.Now(),UpdatedAt:time.Now()}
	model := new(studentModel)
    //插入单个
	rs, err := model.Insert(stu01)
	fmt.Println(rs,err)
	fmt.Println(stu01.Id)
	//插入多个
	rs,err=model.Insert(stu02,stu03)
	fmt.Println(rs,err)
	fmt.Println(stu02.Id,stu03.Id)
}
//**********删******************

//@todo 注意是软删除额
func TestDelete(t *testing.T) {
	model := new(studentModel)
	student := new(Student)
	id := 3
	ret, err := model.Delete(id,student)

	if err != nil {
		t.Errorf("Delete error: %v", err)
		return
	}
	fmt.Println("Delete.ret", ret)
}

func TestDeleteMulti(t *testing.T) {
	model := new(studentModel)
	student := new(Student)
	var id = []interface{}{5, 6}
	// 批量删除id 为5，6的  数据来源参考TestInsert
	ret, err := model.DeleteMulti(id,student)

	if err != nil {
		t.Errorf("DeleteMulti error: %v", err)
		return
	}
	fmt.Println("DeleteMulti.ret", ret)

	// 测试参数为空的异常分支
	var idErr []interface{}
	_, err = model.DeleteMulti(idErr,student)
	fmt.Println("DeleteMulti.CheckExceptionBranch.ret", err)
}



//**********改******************
/**
 *演示修改
 *@author sam@2020-07-06 10:20:40
 */
func TestUpdate(t *testing.T) {
	model := new(studentModel)
	student := new(Student)
	student.ImageUrl = ""
	student.Status=0
	student.Age=18

	var id = 7
	// 注意：直接用默认的update对上面的ImageUrl字段/Status字段不会更新的，因为这两个都设置了特殊的零值，必须结合mustColumns强制更改
	_, err := model.Update(id,student)
	if err != nil {
		t.Errorf("Update error: %v", err)
		return
	}
	fmt.Println("Update.success")
	//===============================================强制

	student.ImageUrl = ""
	student.Status=0
	student.Age=19
	// xorm默认对更新字段数据为""的不会执行，需要加mustColumns，这样保证为空的数据字段也能更新，详情搜索xorm手册
	_, err = model.Update(id,student, "img_url", "status")

	if err != nil {
		t.Errorf("Update mustColumns error: %v", err)
		return
	}
	fmt.Println("Update mustColumns.success")
}
//**********查******************

//根据主键获取一条记录
//@author sam@2020-07-01 17:45:52
func TestGetOneById(t *testing.T) {
	model := new(studentModel)
	ret := new(Student)
	id := 1
	_, err := model.GetOneById(id, ret)
	if err != nil {
		t.Errorf("getOne error: %v", err)
		return
	}
	fmt.Println("getOne.Ret", ret)
}

//根据多个主键ids获取多行记录
//@author  sam@2020-07-02 09:20:29
func TestGetMultiByIds(t *testing.T) {
	model := new(studentModel)
	ret := make([]*Student, 0)
	var idList = []interface{}{1, 2}
	err := model.GetMultiByIds(idList, &ret)
	if err != nil {
		t.Errorf("getMulti error: %v", err)
		return
	}
	for _, v := range ret {
		fmt.Println("getMulti.ItemRet", v)
	}
}

/**
 * 分页测试
 * @author sam@2020-07-06 10:24:38
 */
func TestGetList(t *testing.T) {
	model := new(studentModel)
	students := make([]*Student, 0)

	//设置where语句
	sql := "age > ? and age < ? and status = ?"
	//填充where语句的?
	var values = []interface{}{"1", "20",1}


	err := model.GetList(&students, sql, values)
	if err != nil {
		t.Errorf("Getlist error: %v", err)
		return
	}
	for _, v := range students {
		fmt.Println("GetList.ret", v)
	}
	fmt.Println("-----------------------------------------------------------------------------------------------------")
	//测试其他if分支 覆盖getList所有代码
	students1 := make([]*Student, 0)
	//where
	sql = "age >= ? and age <= ?"
	var valuesTest = []interface{}{"1", "20"}

	err = model.GetList(&students1, sql, valuesTest, []int{3, 3}, "name desc")
	if err != nil {
		t.Errorf("GetlistLimitAndOrderBranch error: %v", err)
		return
	}
	for _, v := range students1 {
		fmt.Println("GetlistLimitAndOrderBranch.ret", v)
	}

}


//测试关闭链接
//@author sam@2020-07-06 11:44:55
func TestProviderClose(t *testing.T) {
	// 关闭链接，此时再执行sql都无法执行会报 sql: database is closed， 所以在sql执行完之后做close操作
	err := Pr.Close()
	if err != nil {
		t.Error("Close Fail")
		return
	}
}