package db

import (
	"errors"
	"fmt"
	"github.com/moka-mrp/sword-core/config"
	"time"
	_ "github.com/go-sql-driver/mysql"
	"xorm.io/xorm"
	"xorm.io/core"
)

const (
	defaultTimeout = 10 //默认连接超时时间
	defaultCharset = "utf8mb4" //默认连接字符编码
)

//创建Engine组
//@author sam@2020-07-01 10:15:11
func NewEngineGroup(dbConf config.DbConfig) (*xorm.EngineGroup, error) {
	//创建主engine
	master, err := newConn(dbConf.Driver, dbConf.Master, dbConf.Option)
	if err != nil {
		panicConnectionErr(dbConf.Driver, dbConf.Master.Host, dbConf.Master.Port, err)
	}
	//创建从engine
	slaves := make([]*xorm.Engine, len(dbConf.Slaves))
	for k, slaveConf := range dbConf.Slaves {
		slave, err := newConn(dbConf.Driver, slaveConf, dbConf.Option)
		if err != nil {
			panicConnectionErr(dbConf.Driver, slaveConf.Host, slaveConf.Port, err)
		}
		slaves[k] = slave
	}
    //创建Engine组返回
	return xorm.NewEngineGroup(master, slaves)
}

//新创建Engine并设置相关引擎属性
//@author sam@2020-07-01 10:24:06
func newConn(driver string, base config.DbBaseConfig, option config.DbOptionConfig) (db *xorm.Engine, err error) {
	//1.组装Data Source Name(dsn)格式
	dsn := formatDSN(driver, base, option)
	if dsn == "" {
		return nil, errors.New(fmt.Sprintf("missing db driver %s or db config", driver))
	}
	//2.创建引擎
	db, err = xorm.NewEngine(driver,dsn)
	if err != nil {
		return
	}
	//3.设置表名和字段的映射规则：驼峰转下划线
	db.SetMapper(core.SnakeMapper{})
	//4.设置资源池等配置
	if option.MaxIdleConns > 0 { //32
		db.SetMaxIdleConns(option.MaxIdleConns)
	}
	if option.MaxOpenConns > 0 { //128
		db.SetMaxOpenConns(option.MaxOpenConns)
	}
	if option.ConnMaxLifetime > 0 { //180
		db.SetConnMaxLifetime(time.Second * option.ConnMaxLifetime)
	}
	return
}

/**
 * 各驱动的dsn
 * @author sam@2020-07-01 10:31:02
 */
func formatDSN(driver string, base config.DbBaseConfig, option config.DbOptionConfig) string {
	switch driver {
	case "mysql":
		return formatMysqlDSN(base, option)
	}
	return ""
}

//Mysql DSN
//@todo  用户名:密码@tcp(IP:端口)/数据库?timeout=3&charset=utf8mb4
//@author sam@2020-07-01 10:32:19
func formatMysqlDSN(base config.DbBaseConfig, option config.DbOptionConfig) string {
	//端口号
	port := getPortOrDefault(base.Port, 3306)
	charset := option.Charset
	if charset == "" {
		charset = defaultCharset
	}
	timeout := option.ConnectTimeout
	if timeout <= 0 {
		timeout = defaultTimeout
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?timeout=%ds&charset=%s&parseTime=true&loc=Local",
		base.User, base.Password, base.Host, port, base.Name, timeout, charset)
}

//获取数据库的端口号
//@author sam@2020-07-01 10:36:23
func getPortOrDefault(port int, defaultPort int) int {
	if port == 0 {
		return defaultPort
	}
	return port
}

//统一封装的创建引擎报错时候的panic输出格式
//@author sam@2020-07-01 10:49:04
func panicConnectionErr(driver string, host string, port int, err error) {
	panic(fmt.Sprintf("%s connect error %s:%d, error:%v", driver, host, port, err))
}
