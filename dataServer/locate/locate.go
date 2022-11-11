package locate

import (
	"OSS/dataServer/conf"
	"OSS/dataServer/rabbitmq"
	"OSS/dataServer/types"
	"log"

	//"github.com/fatih/color"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

var objects = make(map[string]int)
var mutex sync.Mutex

func Locate(hash string) int {
	mutex.Lock()
	id, ok := objects[hash]
	mutex.Unlock()
	if !ok {
		return -1
	}
	return id
}

func Add(hash string, id int) {
	mutex.Lock()
	objects[hash] = id
	mutex.Unlock()
}

func Del(hash string) {
	mutex.Lock()
	delete(objects, hash)
	mutex.Unlock()
}

// func StartLocate(url string) {
// 	q := rabbitmq.New(conf.Conf.RabbitmqAddr)
// 	defer q.Close()
// 	q.Bind("dataServers")
// 	c := q.Consume()
// 	for msg := range c {
// 		hash, e := strconv.Unquote(string(msg.Body))
// 		if e != nil {
// 			panic(e)
// 		}

//			exist := Locate(hash)
//			color.Yellow("服务节点请求数据:%v\t状态: %v\n", hash, exist)
//			if exist {
//				q.Send(msg.ReplyTo, url)
//			}
//		}
//	}
//
// StartLocate 用于监听定位消息
func StartLocate(url string) {
	q := rabbitmq.New(conf.Conf.RabbitmqAddr)
	defer q.Close()
	q.Bind("dataServers")
	c := q.Consume()
	for msg := range c {
		hash, e := strconv.Unquote(string(msg.Body))
		if e != nil {
			panic(e)

		}
		id := Locate(hash)
		if id != -1 {
			q.Send(msg.ReplyTo, types.LocateMessage{Addr: url, Id: id})

		}

	}

}

//	func CollectObjects() {
//		files, _ := filepath.Glob(conf.Conf.Dir + "/objects/*")
//		for i := range files {
//			hash := filepath.Base(files[i])
//			objects[hash] = 1
//		}
//	}
func CollectObjects() {
	files, _ := filepath.Glob(conf.Conf.Dir + "/objects/*")
	// files, _ := filepath.Glob(conf.Conf.Dir + "/temp/*")
	log.Println("files", files)
	log.Println("files2222", len(files))
	for i := range files {
		log.Println("i", i)
		// file := strings.Split(filepath.Base(files[i]), ".")
		file := strings.Split(filepath.Base(files[i]), ".")
		log.Println("file", file)
		log.Println("len(file)", len(file))
		if len(file) != 3 {
			panic(files[i])
		}
		hash := file[0]
		id, e := strconv.Atoi(file[1])
		if e != nil {
			panic(e)

		}
		objects[hash] = id

	}
	log.Println("CollectObjects objects", objects)
}
