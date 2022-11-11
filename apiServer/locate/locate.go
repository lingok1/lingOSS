package locate

import (
	"OSS/apiServer/conf"
	"OSS/apiServer/rabbitmq"
	"OSS/apiServer/rs"
	"OSS/apiServer/types"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return

	}
	info := Locate(strings.Split(r.URL.EscapedPath(), "/")[3])
	if len(info) == 0 {
		w.WriteHeader(http.StatusNotFound)
		return

	}
	b, _ := json.Marshal(info)
	w.Write(b)

}

//func Locate(name string) string {
//	q := rabbitmq.New(conf.Conf.RabbitmqAddr)
//	q.Publish("dataServers", name)
//	c := q.Consume()
//	go func() {
//		time.Sleep(time.Second)
//		q.Close()
//	}()
//	msg := <-c
//	s, _ := strconv.Unquote(string(msg.Body))
//	return s
//
//}
//
//func Exist(name string) bool {
//	return Locate(name) != ""
//
//}

// Locate 用于实际定位对象
func Locate(name string) (locateInfo map[int]string) {
	q := rabbitmq.New(conf.Conf.RabbitmqAddr)
	q.Publish("dataServers", name)
	c := q.Consume()
	go func() {
		time.Sleep(time.Second)
		q.Close()

	}()
	locateInfo = make(map[int]string)
	for i := 0; i < rs.ALL_SHARDS; i++ {
		msg := <-c
		if len(msg.Body) == 0 {
			return

		}
		var info types.LocateMessage
		json.Unmarshal(msg.Body, &info)
		locateInfo[info.Id] = info.Addr

	}
	return

}

func Exist(name string) bool {
	return len(Locate(name)) >= rs.DATA_SHARDS

}
