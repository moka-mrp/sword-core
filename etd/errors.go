package etd

import "errors"

var (
	ErrNotFound        = errors.New("Record not found.")//当前获取的记录不存在，比如任务详情
	ErrValueMayChanged = errors.New("The value has been changed by others on this time.") //持乐观锁执行input操作
)

