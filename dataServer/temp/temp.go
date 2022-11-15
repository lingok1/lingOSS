package temp

import (
	"OSS/dataServer/conf"
	"OSS/dataServer/locate"
	"OSS/dataServer/utils"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		put(w, r)
		return

	}
	if m == http.MethodPatch {
		patch(w, r)
		return

	}
	if m == http.MethodPost {
		post(w, r)
		return

	}
	if m == http.MethodDelete {
		del(w, r)
		return

	}
	w.WriteHeader(http.StatusMethodNotAllowed)

}

//func commitTempObject(datFile string, tempinfo *tempInfo) {
//	os.Rename(datFile, conf.Conf.Dir+"/objects/"+tempinfo.Name)
//	locate.Add(tempinfo.Name)
//}

func (t *tempInfo) hash() string {
	s := strings.Split(t.Name, ".")
	return s[0]

}

func (t *tempInfo) id() int {
	s := strings.Split(t.Name, ".")
	id, _ := strconv.Atoi(s[1])
	return id

}

// commitTempObject 提交临时对象,os.Rename()将临时对象文件重命名为正式对象文件
func commitTempObject(datFile string, tempinfo *tempInfo) {
	log.Println("commitTempObject11111")
	f, _ := os.Open(datFile)
	d := url.PathEscape(utils.CalculateHash(f))
	log.Println("d 1111111", d)
	f.Close()
	os.Rename(datFile, conf.Conf.Dir+"/objects/"+tempinfo.Name+"."+d)
	locate.Add(tempinfo.hash(), tempinfo.id())

}

// put 接口服务数据校验一致，将该临时对象转正
func put(w http.ResponseWriter, r *http.Request) {
	log.Println("put 111111")
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	log.Println("uuid1111", uuid)
	tempinfo, e := readFromFile(uuid)
	log.Println("tempinfo2222", tempinfo)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return

	}
	infoFile := conf.Conf.Dir + "/temp/" + uuid
	datFile := infoFile + ".dat"
	log.Println("datFile 333333", datFile)
	f, e := os.Open(datFile)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	actual := info.Size()
	os.Remove(infoFile)
	if actual != tempinfo.Size {
		os.Remove(datFile)
		log.Println("actual size mismatch, expect", tempinfo.Size, "actual", actual)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	log.Println("actual4444", actual, "expected", tempinfo.Size)
	commitTempObject(datFile, tempinfo)

}

// patch 访问数据服务节点上的临时数据，http请求的正文会被写入该临时对象
func patch(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	tempinfo, e := readFromFile(uuid)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return

	}
	infoFile := conf.Conf.Dir + "/temp/" + uuid
	datFile := infoFile + ".dat"
	f, e := os.OpenFile(datFile, os.O_WRONLY|os.O_APPEND, 0)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	defer f.Close()
	// 写入分片数据
	_, e = io.Copy(f, r.Body)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	info, e := f.Stat()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	actual := info.Size()
	if actual > tempinfo.Size {
		os.Remove(datFile)
		os.Remove(infoFile)
		log.Println("actual size", actual, "exceeds", tempinfo.Size)
		w.WriteHeader(http.StatusInternalServerError)

	}
}

func readFromFile(uuid string) (*tempInfo, error) {
	f, e := os.Open(conf.Conf.Dir + "/temp/" + uuid)
	if e != nil {
		return nil, e

	}
	defer f.Close()
	b, _ := ioutil.ReadAll(f)
	var info tempInfo
	json.Unmarshal(b, &info)
	return &info, nil

}

type tempInfo struct {
	Uuid string
	Name string
	Size int64
}

// post 数据服务节点创建临时对象，写入tempInfo信息
func post(w http.ResponseWriter, r *http.Request) {
	output, _ := exec.Command("uuidgen").Output()
	uuid := strings.TrimSuffix(string(output), "\n")
	name := strings.Split(r.URL.EscapedPath(), "/")[2]
	size, e := strconv.ParseInt(r.Header.Get("size"), 0, 64)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	t := tempInfo{uuid, name, size}
	e = t.writeToFile()
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}
	os.Create(conf.Conf.Dir + "/temp/" + t.Uuid + ".dat")
	w.Write([]byte(uuid))

}

func (t *tempInfo) writeToFile() error {
	f, e := os.Create(conf.Conf.Dir + "/temp/" + t.Uuid)
	if e != nil {
		return e

	}
	defer f.Close()
	b, _ := json.Marshal(t)
	f.Write(b)
	return nil

}

// del 接口服务数据校验不一致，删除该临时对象
func del(w http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.EscapedPath(), "/")[2]
	infoFile := conf.Conf.Dir + "/temp/" + uuid
	datFile := infoFile + ".dat"
	os.Remove(infoFile)
	os.Remove(datFile)

}
