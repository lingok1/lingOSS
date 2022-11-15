package objects

import (
	"OSS/apiServer/es"
	"OSS/apiServer/heartbeat"
	"OSS/apiServer/locate"

	//"OSS/apiServer/objectstream"
	"OSS/apiServer/rs"
	"OSS/apiServer/utils"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

/*
http.HandleFunc:put get del
*/

//	func get(w http.ResponseWriter, r *http.Request) {
//		name := strings.Split(r.URL.EscapedPath(), "/")[3]
//		versionId := r.URL.Query()["version"]
//		version := 0
//		var e error
//		if len(versionId) != 0 {
//			version, e = strconv.Atoi(versionId[0])
//			if e != nil {
//				log.Println(e)
//				w.WriteHeader(http.StatusBadRequest)
//				return
//			}
//		}
//
//		meta, e := es.GetMetadata(name, version)
//		if e != nil {
//			log.Println(e)
//			w.WriteHeader(http.StatusInternalServerError)
//			return
//		}
//
//		if meta.Hash == "" {
//			w.WriteHeader(http.StatusNotFound)
//			return
//		}
//
//		object := url.PathEscape(meta.Hash)
//		stream, e := getStream(object)
//		if e != nil {
//			log.Println(e)
//			w.WriteHeader(http.StatusNotFound)
//			return
//		}
//
//		io.Copy(w, stream)
//	}
//

func get(w http.ResponseWriter, r *http.Request) {
	log.Println("r.URL.EscapedPath()", r.URL.EscapedPath())
	log.Println("11", strings.Split(r.URL.EscapedPath(), "/"))
	name := strings.Split(r.URL.EscapedPath(), "/")[3]
	log.Println("name", name)
	versionId := r.URL.Query()["version"]
	log.Println("version 111", versionId)
	// log.Println("version 111", versionId[0])
	version := 0
	var e error
	if len(versionId) != 0 {
		version, e = strconv.Atoi(versionId[0])
		log.Println("version 2222", versionId)
		log.Println("version 2222", versionId[0])
		if e != nil {
			log.Println(e)
			w.WriteHeader(http.StatusBadRequest)
			return

		}

	}
	meta, e := es.GetMetadata(name, version)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	if meta.Hash == "" {
		w.WriteHeader(http.StatusNotFound)
		return

	}
	log.Println("meta", meta)
	log.Println("meta.Size", meta.Size)
	log.Println("meta.Name", meta.Name)
	hash := url.PathEscape(meta.Hash)
	log.Println("hash", hash)
	stream, e := GetStream(hash, meta.Size)
	log.Println("stream", stream)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return

	}
	_, e = io.Copy(w, stream)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return

	}
	log.Println("e", e)
	log.Println("w", w)
	stream.Close()

}

//func getStream(object string) (io.Reader, error) {
//	server := locate.Locate(object)
//	if server == "" {
//		return nil, fmt.Errorf("object %s locate fail", object)
//	}
//
//	return objectstream.NewGetStream(server, object)
//}

func GetStream(hash string, size int64) (*rs.RSGetStream, error) {
	log.Println("hash", hash, "size", size)
	locateInfo := locate.Locate(hash)
	log.Println("locateInfo", locateInfo)
	if len(locateInfo) < rs.DATA_SHARDS {
		return nil, fmt.Errorf("object %s locate fail, result %v", hash, locateInfo)

	}
	dataServers := make([]string, 0)
	if len(locateInfo) == rs.ALL_SHARDS {
		dataServers = heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS-len(locateInfo), locateInfo)

	}
	log.Println("locateInfo 1111", locateInfo)
	log.Println("dataServers2222", dataServers)
	log.Println("hash333", hash)
	log.Println("size4444", size)
	return rs.NewRSGetStream(locateInfo, dataServers, hash, size)

}

// put hash是前端算好传过来的，先hash后base64
func put(w http.ResponseWriter, r *http.Request) {
	hash := utils.GetHashFromHeader(r.Header)
	if hash == "" {
		log.Println("missing object hash in digest header")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	size := utils.GetSizeFromHeader(r.Header)
	log.Println("来自客户端的PUT信息:")
	color.Yellow("Hash : %v \nSize : %v\n", url.PathEscape(hash), strconv.FormatInt(size, 10))
	log.Println("r.Body", r.Body)
	log.Println("r.Body", &r.Body)
	log.Println("hash", hash)
	log.Println("size", size)
	c, e := storeObject(r.Body, hash, size) //向dataserver写入时附带hash信息 等元数据
	if e != nil {
		log.Println(e)
		w.WriteHeader(c)
		return
	}

	if c != http.StatusOK {
		w.WriteHeader(c)
		w.Write([]byte("111"))
		return
	}

	name := strings.Split(r.URL.EscapedPath(), "/")[3]
	e = es.AddVersion(name, hash, size)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//	func putStream(hash string, size int64) (*objectstream.TempPutStream, error) {
//		server := heartbeat.ChooseRandomDataServer()
//		if server == "" {
//			return nil, fmt.Errorf("cannot find any dataServer")
//		}
//
//		return objectstream.NewTempPutStream(server, hash, size)
//	}

func putStream(hash string, size int64) (*rs.RSPutStream, error) {
	servers := heartbeat.ChooseRandomDataServers(rs.ALL_SHARDS, nil)
	log.Println("servers", servers)
	if len(servers) != rs.ALL_SHARDS {
		return nil, fmt.Errorf("cannot find enough dataServer")

	}

	return rs.NewRSPutStream(servers, hash, size)

}

func storeObject(r io.Reader, hash string, size int64) (int, error) {
	if locate.Exist(url.PathEscape(hash)) {
		return http.StatusOK, nil
	}
	//调用封装的http post方式写入流向发给detaserver的请求写入数据
	stream, e := putStream(url.PathEscape(hash), size)
	if e != nil {
		return http.StatusServiceUnavailable, e
	}

	reader := io.TeeReader(r, stream)
	d := utils.CalculateHash(reader)
	color.Red("api hash = %v\n", d)
	color.Red("header hash = %v\n", hash)
	if d != hash {
		stream.Commit(false)
		return http.StatusBadRequest, fmt.Errorf("object hash mismatch,Calculated=%s,requested=%s", d, hash)
	}
	stream.Commit(true)
	return http.StatusOK, nil

}

func del(w http.ResponseWriter, r *http.Request) {
	name := strings.Split(r.URL.EscapedPath(), "/")[3]
	version, e := es.SearchLatestVersion(name)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	e = es.PutMetadata(name, version.Version+1, 0, "")
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

}
func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		put(w, r)
		return
	} else if m == http.MethodGet {
		get(w, r)
		return
	} else if m == http.MethodDelete {
		del(w, r)
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
