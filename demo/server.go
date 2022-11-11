package main

import (
	"OSS/demo/objects"
	"log"
	"net/http"
)

func main() {
	log.Println("server listening 1212~~~~")

	http.HandleFunc("/objects/", objects.Handler)
	http.ListenAndServe(":1212", nil)
	log.Println("server listening 1212```````")
}
