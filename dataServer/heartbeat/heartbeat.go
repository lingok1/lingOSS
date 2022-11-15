package heartbeat

import (
	"OSS/dataServer/conf"
	"OSS/dataServer/rabbitmq"
	"time"
)

// StartHeartbeat 当前数据服务节点注册到数据服务集群
func StartHeartbeat(url string) {
	q := rabbitmq.New(conf.Conf.RabbitmqAddr)
	defer q.Close()
	for {
		q.Publish("apiServers", url)
		time.Sleep(5 * time.Second)
	}
}
