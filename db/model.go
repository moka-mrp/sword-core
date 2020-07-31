package db

import (
	"errors"
	"github.com/moka-mrp/sword-core/kernel/server"
	"xorm.io/xorm"
)

var (
	ErrIdsEmpty = errors.New("ids is empty")
)

//-------------------基础model----------------------------------------------------------------------------
type Model struct {
	DiName string //依赖注入的别名
}

/**
 * 从容器中获取引擎组
 * @author sam@2020-07-01 14:32:41
 */
func (m *Model) GetDb(args ...string) *xorm.EngineGroup {
	var engineGroup *xorm.EngineGroup

	if len(args) > 0 {
		engineGroup=GetDb(args[0])

	} else if m.DiName != "" {
		engineGroup=GetDb(m.DiName)
	} else {//normal here  本质还是从容器中获取了db资源
		engineGroup=GetDb()
	}
	if server.GetDebug(){
		engineGroup.ShowSQL(true)
	}
	return engineGroup
}

//******************************** 增 ****************************************
/**
 * 插入记录
 * @param beans... 可支持插入连续多个记录 参数可以是一个或多个Struct的指针，一个或多个Struct的Slice的指针。
 * @todo 返回插入成功的函数,如果需要返回插入的id,则传递Struct的指针即可
 * @author  sam@2020-07-01 17:51:40
 * @link https://gobook.io/read/gitea.com/xorm/manual-zh-CN/chapter-04/index.html
 */
func (m *Model) Insert(beans ...interface{}) (int64, error) {
	return m.GetDb().Insert(beans...)
}

//******************************** 删 ****************************************
/**
 * 根据主键删除单个记录 -- 如果有开启delete特性，会触发软删除
 * @param id 主键ID
 * @param bean 数据结构实体(根据它能获知是删除哪张表)
 * @todo  注意自己是否开启了软删除功能额
 * @author sam@2020-07-06 11:09:28
 */
func (m *Model) Delete(id interface{}, bean interface{}) (int64, error) {
	return m.GetDb().ID(id).Delete(bean)
}

/**
 * 根据多个主键批量删除
 * @param ids 主键ID分片
 * @param bean 数据结构实体
 * @todo 建议一次性传递的ids不要太多了，如果真的过多，建议分批次调用该方法
 * @author sam@2020-07-06 11:39:06
 */
func (m *Model) DeleteMulti(ids []interface{}, bean interface{}) (int64, error) {
	if len(ids) == 0 {
		return 0, ErrIdsEmpty
	}
	return m.GetDb().In("id", ids...).Delete(bean)
}

//******************************** 改 ****************************************
/**
 * 更新某个主键ID的数据
 * @param id 主键ID
 * @param bean 数据结构实体
 * @param mustColumns... 因为默认Update只更新非0，非”“，非bool的字段，需要配合此字段
 * @author sam@2020-07-06 11:08:20
 */
func (m *Model) Update(id interface{}, bean interface{}, mustColumns ...string) (int64, error) {
	if len(mustColumns) > 0 {
		return m.GetDb().MustCols(mustColumns...).ID(id).Update(bean)
	} else {
		return m.GetDb().ID(id).Update(bean)
	}
}
//******************************** 查 ****************************************
/**
 * 根据主键获取某行记录
 * @param id 主键ID
 * @param bean 数据结构实体
 * @return has 是否有记录
 * @auth sam@2020-07-01 17:43:57
 */
func (m *Model) GetOneById(id interface{}, bean interface{}) (has bool, err error) {
	//m.GetDb().ShowSQL(true)
	return m.GetDb().ID(id).Get(bean)
}

/**
 * 根据多个主键ID获取多行记录
 * @param ids 主键ID切片
 * @param beans 数据结构实体切片
 * @author sam@2020-07-02 08:57:04
 */
func (m *Model) GetMultiByIds(ids []interface{}, beans interface{}) error {
	if len(ids) == 0 {
		return ErrIdsEmpty
	}
	return m.GetDb().In("id", ids...).Find(beans)
}

/**
 * 根据fastadmin封装的统一的list查询方法
 * @param beans 数据结构实体切片 eg. &students 其中 students := make([]*Student, 0)
 * @params sql  where语句的sql eg. "age > ? or name = ?"
 * @params values where语句中?的替代值  eg. []interfaces{}{30, "sam"}
 * todo 下面的参数是从args中主动提取出来的 ---------------------------------------------------------------------------------
 * @Param []int limit 可选 eg. []int{} 不限量 []int{30} 前30个 []int{30, 20} 从第20个后的前30个
 * @param string order 可选 eg.  "id desc" 单个 "name desc,age asc" 多个
 * @author  sam@2020-07-01 14:53:11
 */
func (m *Model) GetList(beans interface{}, sql string, values []interface{}, args ...interface{}) (err error) {
	//解析args参数
	if len(args) > 0 {
		var (
			order string
			limit int
			start int
		)
		limits, ok := args[0].([]int)
		if ok && len(limits) > 0 {
			limit = limits[0] //取多少条
			if len(limits) > 1 {
				start = limits[1] //从哪个位置开始取
			}
		}
		//排序问题
		if len(args) > 1 {
			order, _ = args[1].(string)
		}
		return m.GetDb().Where(sql, values...).OrderBy(order).Limit(limit, start).Find(beans)
	} else {
		return m.GetDb().Where(sql, values...).Find(beans)
	}
}
