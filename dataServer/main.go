package main

import (
	"OSS/dataServer/conf"
	"OSS/dataServer/heartbeat"
	"OSS/dataServer/locate"
	"OSS/dataServer/objects"
	"OSS/dataServer/temp"
	"flag"
	"fmt"
	"log"
	"net/http"
)

var configFile = flag.String("f", "config/conf.json", "the config file")

func init() {

}

func main() {

	flag.Parse()

	cfg := *configFile
	conf.Conf.Parse(cfg)

	var url string //监听地址:端口
	url = conf.Conf.ListenAddr + ":" + conf.Conf.ListenPort
	log.Println(url)
	fmt.Println("run dataServer start...")
	locate.CollectObjects()

	go heartbeat.StartHeartbeat(url)
	go locate.StartLocate(url)
	http.HandleFunc("/objects/", objects.Handler)//对象获取
	http.HandleFunc("/temp/", temp.Handler)//对象上传
	log.Fatal(http.ListenAndServe(url, nil))
	fmt.Println("run dataServer end...")

}
