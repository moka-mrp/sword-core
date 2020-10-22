## sword
sword是一套简单易用的Go语言业务框架，整体逻辑设计简洁，支持HTTP服务、分布式调度和和脚本任务等常用业务场景模式。

## Goals
让PHPer更加容易切换角色成为Gopher。




## Quick start

### Requirements
- Go version >= 1.13
- Global environment configure (Linux/Mac)  

```
export  GOPROXY=https://goproxy.cn
export  GO111MODULE=on
```

### Installation
```shell
go get github.com/moka-mrp/sword
sword new sword-demo -p /tmp  -m moye #sword-demo为创建的项目名称, -p指明项目放置目录，未指明取当面目录,-m指明go mod模块名，未指明取项目名称
```

### Build & Run
```shell
cd /tmp/sword-demo
go run main.go  api
```

### Test demo
```
curl "http://127.0.0.1:9999"
```

## Documents

- [项目脚手架](https://github.com/moka-mrp/sword)
- [中文文档](https://github.com/moka-mrp/sword-core/wiki)
- [changelog](https://github.com/moka-mrp/sword-core/blob/master/CHANGELOG.md)

## Contributors

- Sam 
- Rick 
- Lucifer




