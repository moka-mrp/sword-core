package db

import (
	"fmt"
	"github.com/moka-mrp/sword-core/config"
	"time"
	"xorm.io/xorm"
	"testing"
	//go test时需要开启
	_ "github.com/go-sql-driver/mysql"
)

//todo  测试db资源的连接问题
//todo  直接测试engineGroup的可以在这里测试

var engineGroup *xorm.EngineGroup

//--------------------- Student struct -------------------------------------------------------------
type Student struct {
	Id        int64 `xorm:"pk autoincr"` // 注：使用getOne 或者ID() 需要设置主键
	Age      int
	Name     string
	ImageUrl  string `xorm:"'img_url'"`
	Url       string
	Status    uint8
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `xorm:"deleted"` // 此特性会激发软删除
}

func (m *Student) TableName() string {
	return "student"
}

//-------------------------------- init -----------------------------------------------
func init() {
	//初始化的时候会注入的
	dbInit(true)
	//从容器中获取资源
	engineGroup = GetDb()
}

//注入容器以及从容器中快速取出来
func dbInit(lazyBool bool) {
	//基本配置项
	base:= config.DbBaseConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "root",
		Name:   "gwp",
	}
	//连接超时相关配置项
	options:=config.DbOptionConfig{
		MaxIdleConns:    32,
		MaxOpenConns:    128,
		ConnMaxLifetime: 180,//秒
		ConnectTimeout:  3,
		Charset:         "utf8mb4",
	}
	dbConf := config.DbConfig{
		Driver: "mysql",
		Master: base,
		Option:options,
	}
	//测试容器注入功能(容器本身已经自动在kernel/container/app.go中初始化好了)
	err := Pr.Register("db", dbConf, lazyBool)
	if err != nil {
		fmt.Println(err)
	}
}
//----------------------下面的单元测试依赖上面engineGroup的初始化-----------------------------------------------------------

//测试一下直接获取个数据
//@author sam@2020-07-29 14:07:42
func TestGet(t *testing.T) {
	student := new(Student)
	// sql是否打印开关
	engineGroup.ShowSQL(true)
	_, err := engineGroup.ID(1).Get(student)

	if err != nil {
		t.Errorf("get error: %v", err)
		return
	}
    fmt.Printf("%+v\r\n",student)
}


//测试一下注入的别名有哪些
//@author sam@2020-07-29 14:07:52
func TestPrProvides(t *testing.T) {
	retList := Pr.Provides()
	if len(retList) == 0 {
		t.Error("Provides empty")
		return
	}
	fmt.Printf("%+v\r\n",retList)
	//for k, v := range retList {
	//	fmt.Println("Provides list", k, v)
	//}
}
