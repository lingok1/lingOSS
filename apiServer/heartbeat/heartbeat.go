package heartbeat

import (
	"OSS/apiServer/conf"
	"OSS/apiServer/rabbitmq"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

/*
监听集群服务心跳
移除过期的集群服务
获取全部集群服务
随机选择一个服务节点
*/

var dataServers = make(map[string]time.Time)
var mutex sync.Mutex

func ListenHeartbeat() {
	q := rabbitmq.New(conf.Conf.RabbitmqAddr)
	defer q.Close()
	q.Bind("apiServers")
	c := q.Consume()
	go removeExpiredDataServer()
	for msg := range c {
		dataServer, e := strconv.Unquote(string(msg.Body))
		// color.Green("来自数据节点的心跳%v\n", dataServer)
		if e != nil {
			panic(e)
		}
		mutex.Lock()
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}

func removeExpiredDataServer() {
	for {
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers {
			if t.Add(10 * time.Second).Before(time.Now()) {
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}

}

func GetDataServers() []string {
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)
	for s, _ := range dataServers {
		ds = append(ds, s)
	}
	return ds

}

//func ChooseRandomDataServer() string {
//	ds := GetDataServers()
//	n := len(ds)
//	if n == 0 {
//		return ""
//	}
//	return ds[rand.Intn(n)]
//
//}

func ChooseRandomDataServers(n int, exclude map[int]string) (ds []string) {
	log.Println("n", n, "exclude", exclude)
	candidates := make([]string, 0)
	reverseExcludeMap := make(map[string]int)
	for id, addr := range exclude {
		reverseExcludeMap[addr] = id

	}
	servers := GetDataServers()
	for i := range servers {
		s := servers[i]
		_, excluded := reverseExcludeMap[s]
		if !excluded {
			candidates = append(candidates, s)

		}

	}
	length := len(candidates)
	if length < n {
		return

	}
	p := rand.Perm(length)
	for i := 0; i < n; i++ {
		ds = append(ds, candidates[p[i]])

	}
	log.Println("ds",ds)
	return

}
