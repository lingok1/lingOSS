# 存储桶改造

<div align=center>
<img src="https://www.helloimg.com/images/2022/11/11/ZfU1ST.png" width="200" height="200"/>
</div>
<br/>

书籍：《分布式对象存储原理架构及Go语言实现》

书籍包括以下内容:
对象存储简介;数据冗余和即时修复;分布式系统原理;断点续传;元数据以及元数据服务;数据压缩;数据校验和去重;数据维护。

中间件：rabbitMQ（集群心跳）和ElasticSearch（全文索引）
## apiServer目录结构
```go
.
├── conf
│   ├── conf.json
│   └── config.go
├── es
│   └── es.go
├── heartbeat
│   └── heartbeat.go
├── locate
│   └── locate.go
├── main.go
├── objects
│   └── objects.go
├── objectstream
│   ├── objectstream.go
│   └── tempstream.go
├── rabbitmq
│   └── rabbitmq.go
├── rs
│   ├── common.go
│   ├── decoder.go
│   ├── decoder_test.go
│   ├── encoder.go
│   ├── get.go
│   ├── put.go
│   ├── resumable_get.go
│   └── resumable_put.go
├── start.sh
├── types
│   └── types.go
├── utils
│   └── utils.go
└── versions
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
├── heartbeat
│   └── heartbeat.go
├── kill.sh
├── locate
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
├── objects
│   └── objects.go
├── rabbitmq
│   └── rabbitmq.go
├── start.sh
├── temp
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
api层
data层


