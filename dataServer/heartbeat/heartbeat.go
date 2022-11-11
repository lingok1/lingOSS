package heartbeat

import (
	"OSS/dataServer/conf"
	"OSS/dataServer/rabbitmq"
	"time"
)

func StartHeartbeat(url string) {
	q := rabbitmq.New(conf.Conf.RabbitmqAddr)
	defer q.Close()
	for {
		q.Publish("apiServers", url)
		time.Sleep(5 * time.Second)
	}
}
