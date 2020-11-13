## v0.1.3(2020-11-13)
### New Features
- 钉钉机器人增加 文本类型、link 类型、markdown 类型、整体跳转 ActionCard 类型

## v0.1.2(2020-11-04)
### New Features
- 增加中台接口平台接入签名包sign

## v0.1.1(2020-09-09)
### New Features
- 增加守护进程优雅退出事件

## v0.1.0(2020-08-26)
### New Features
- 增加命令行参数、命令行选项智能解析
- 增加xorm时间格式友好解析

## v0.0.9(2020-08-24)
### Changes
- ctxkit增加jwt claims上下文


## v0.0.8(2020-08-22)
### New Features
- Etcd组件服务

## v0.0.7(2020-08-18)
### Changes
-通用上下文kit增加http相关标识量
-通用上下文kit增加jwt相关标识量
-通用上下文kit增加app交互时序相关标识量

## v0.0.6(2020-08-17)
### Changes
-增加随机字符串产生库
-增加md5校验函数


## v0.0.5(2020-08-11)
### Bug Fix
- 修复log配置名称与路径匹配替换冲突的问题


## v0.0.4(2020-08-09)
### Changes
-增加默认的http中间件


## v0.0.3(2020-08-08)
### Changes
- redis方法糖优化。支持多连接池使用


## v0.0.2(2020-08-07)
### New Features
- 做个有趣的程序员图库支持
### Bug Fix
- 修复rds包配置秒级未进行转换变成纳秒的问题


## v0.0.1(2020-07-28)
### New Features
- Redis组件服务
- Log组件服务
- DB组件服务
- Config通用配置结构
- Http的通用中间件和通用上下文kit
- Kernel内核包
    - close服务注册
    - provider组件注册
    - container容器注入
    - server通用服务启动
- utils工具包
    - HTTP请求工具包
    - 其他常用函数工具包