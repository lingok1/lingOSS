# 存储桶改造
<div align=center>
<img src="https://www.helloimg.com/images/2022/11/15/ZhxpK0.jpg" width="400" height="200"/>

<br/>

![Go version](https://img.shields.io/badge/go-v1.18-9cf)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/lling1234/go-workflow/blob/master/LICENSE)

</div>

参考书籍：《分布式对象存储原理架构及Go语言实现》

书籍包括以下内容:
对象存储简介;数据冗余和即时修复;分布式系统原理;断点续传;元数据以及元数据服务;数据压缩;数据校验和去重;数据维护。

中间件：rabbitMQ（集群心跳）和ElasticSearch（全文索引）

## 功能
### 已实现功能
* 文件上传
* 文件下载
* 文件删除
* 减少数据冗余
* 即时修复
* 元数据服务（对象的名字、版本、大小和散列值）
* 数据校验和去重
* 集群心跳
* 全文索引
* 分布式实现
### 未实现功能
* 多级文件夹增删改
* 断点续传
* 数据压缩
* 数据维护
* 大小文件上传

## apiServer目录结构
```go
.
├── conf
│   ├── conf.json
│   └── config.go
├── es elasticsearch GetMetadata和PuttMetadata http接口封装
│   └── es.go 
├── heartbeat
│   └── heartbeat.go 数据集群服务心跳，5s监测一次
├── locate 用于实际定位文件对象
│   └── locate.go 
├── main.go
├── objects 增查删对象路由处理
│   └── objects.go
├── objectstream 临时对象流和对象流http接口封装，调用dataServer的实现
│   ├── objectstream.go
│   └── tempstream.go
├── rabbitmq 数据集群服务调用
│   └── rabbitmq.go 
├── rs Reed Solomon算法编码解码封装
│   ├── common.go
│   ├── decoder.go
│   ├── decoder_test.go
│   ├── encoder.go
│   ├── get.go
│   ├── put.go
│   ├── resumable_get.go
│   └── resumable_put.go
├── start.sh 启动脚本
├── types
│   └── types.go
├── utils
│   └── utils.go
└── versions 文件版本
    └── version.go
```

## dataServer目录结构
```go
.
├── common.sh
├── conf
│   ├── 1.json
│   ├── 2.json
│   ├── 3.json
│   ├── 4.json
│   ├── 5.json
│   ├── 6.json
│   ├── conf.json
│   └── config.go
├── delObjects.sh
├── heartbeat 当前数据服务节点注册到数据服务集群
│   └── heartbeat.go
├── kill.sh
├── locate 先扫描一下当前数据服务器所有对象objects存到内存中
│   └── locate.go
├── logs
│   ├── 1.log
│   ├── 2.log
│   ├── 3.log
│   ├── 4.log
│   ├── 5.log
│   ├── 6.log
│   └── conf.log
├── main.go
├── objects 对象获取路由和具体实现
│   └── objects.go
├── rabbitmq
│   └── rabbitmq.go
├── start.sh
├── temp 对象上传路由和具体实现
│   └── temp.go
├── types
│   └── types.go
└── utils
    └── utils.go
```

## objectsData对象文件存储目录结构
```go
.
├── 1
│   ├── objects
│   └── temp
├── 2
│   ├── objects
│   └── temp
├── 3
│   ├── objects
│   └── temp
├── 4
│   ├── objects
│   └── temp
├── 5
│   ├── objects
│   └── temp
└── 6
    ├── objects
    └── temp
```

## 代码逻辑
### apiServer层
```go
1.对象上传
从http Header获取的hash值与对象流hash值作对比（hash值先sha256后bash64）
对象分片，先调用dataServer http post路由再调用patch和put路由
post路由：数据服务节点创建临时对象，写入tempInfo信息
patch路由：访问数据服务节点上的临时数据，http请求的正文会被写入该临时对象
put路由：接口服务数据校验一致，将该临时对象转正
2.对象获取
从es中获取对象元数据信息，根据元数据信息请求dataServer http get路由获取文件流返回
3.对象删除
获取对象最新版本的元数据信息，在es中创建一条新的元数据信息记录（name不变，version+1，size=0,uuid=""）
```

### dataServer层
```go
1.文件上传先调用http的post patch del路由，调用put路由 接口服务数据校验一致，将该临时对象转正
临时对象temp文件夹存储 uuid和uuid.data两个文件文件，
uuid存储结构体信息
type tempInfo struct {
	Uuid string
	Name string
	Size int64
}
uuid.data存储对象分片数据
put路由 接口服务数据校验一致，将该临时对象转正
读取uuid文件结构体信息校验一致，uuid.data移动到objects文件夹，
文件重命名uuid.1.uuid
D4lNxkL5lsqhLH+llEIrzRN%2Ft5qBjXeUHP1TSf%2FDbeY=.0.a%2FQAT3VLYdnPZLTxYXClze%2Fx6Vuk0hhR6pK78dlV4Xg=
删除temp文件夹里的uuid和uuid.data两个文件

2.对象下载
根据文件名在对象存储节点查找，查找到通过http响应
```

