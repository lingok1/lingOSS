package main

import (
	"OSS/apiServer/conf"
	"OSS/apiServer/heartbeat"
	"OSS/apiServer/locate"
	"OSS/apiServer/objects"
	"OSS/apiServer/versions"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func init() {
	confile := "conf/conf.json"
	conf.Conf.Parse(confile)
}

func main() {

	var url string //监听地址:端口

	url = conf.Conf.ListenAddr + ":" + conf.Conf.ListenPort
	fmt.Println("run apiServer start...")
	fmt.Println("ListenPort:6060")
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/OSS/", indexPage)
	http.HandleFunc("/OSS/objects/", objects.Handler) 
	http.HandleFunc("/OSS/locate/", locate.Handler)
	http.HandleFunc("/OSS/versions/", versions.Handler)
	log.Fatalln(http.ListenAndServe(url, nil))
	fmt.Println("run apiServer end...")

}

func indexPage(w http.ResponseWriter, r *http.Request) {
	b, _ := ioutil.ReadFile("../index.html")
	_, _ = fmt.Fprintln(w, string(b))
}
