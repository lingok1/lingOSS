package objects

import (
	"OSS/dataServer/conf"
	"OSS/dataServer/locate"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodGet {
		log.Println("Handler11111")
		log.Println(strings.Split(r.URL.EscapedPath(), "/")[2])
		get(w, r)
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	log.Println("Handler 222222222")
	log.Println(strings.Split(r.URL.EscapedPath(), "/")[2])
	file := getFile(strings.Split(r.URL.EscapedPath(), "/")[2])
	log.Println("Handler file",file)
	if file == "" {
		w.WriteHeader(http.StatusNotFound)
		return

	}
	sendFile(w, file)
	log.Println("w",w)

}

//func getFile(hash string) string {
//	file := conf.Conf.Dir + "/objects/" + hash
//	f, _ := os.Open(file)
//	d := url.PathEscape(utils.CalculateHash(f))
//	f.Close()
//	if d != hash {
//		log.Println("object hash mismatch, remove", file)
//		locate.Del(hash)
//		os.Remove(file)
//		return ""
//
//	}
//	return file
//
//}

func getFile(name string) string {
	log.Println("getFile111")
	log.Println("name: ", name)
	files, _ := filepath.Glob(conf.Conf.Dir + "/objects/" + name + ".*")
	// files, _ := filepath.Glob(conf.Conf.Dir + "/temp/" + name + ".*")
	if len(files) != 1 {
		return ""

	}
	file := files[0]
	h := sha256.New()
	sendFile(h, file)
	d := url.PathEscape(base64.StdEncoding.EncodeToString(h.Sum(nil)))
	hash := strings.Split(file, ".")[2]
	if d != hash {
		log.Println("object hash mismatch, remove", file)
		locate.Del(hash)
		os.Remove(file)
		return ""

	}
	return file

}

func sendFile(w io.Writer, file string) {
	f, _ := os.Open(file)
	defer f.Close()
	io.Copy(w, f)
}
