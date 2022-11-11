package objects

import (
	"OSS/apiServer/rs"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	fileName := strings.Split(r.URL.EscapedPath(), "/")[2]
	log.Println("fileName", fileName)
	body, err := GetObject(fileName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Write(body)
	// f, e := os.Open("/home/ling/go/src/go-implement-your-object-storage/c1" + "/objects/" +
	// 	strings.Split(r.URL.EscapedPath(), "/")[2])
	// if e != nil {
	// 	log.Println(e)
	// 	w.WriteHeader(http.StatusNotFound)
	// 	return
	// }
	// defer f.Close()
	// io.Copy(w, f)
}
// main.txt 976
func GetObject(fileName string) ([]byte, error) {
	var err error
	fileSize := 976
	writers := make([]io.Writer, rs.ALL_SHARDS)
	readers := make([]io.Reader, rs.ALL_SHARDS)
	for i := 0; i < rs.ALL_SHARDS; i++ {
		readers[i], err = os.Open(fmt.Sprintf(FileDIr+fileName+".%d", i))
		if err != nil {
			return nil, err
		}
	}
	log.Println("readers111",readers)
	log.Println("writers2222",writers)
	dec := rs.NewDecoder(readers, writers, int64(fileSize))
	body := make([]byte, fileSize+10)
	count := 0
	for {
		n, err := dec.Read(body[count:])
		count += n
		if err != nil {
			log.Println("err1111111",err)
			break
		}
	}
	log.Println("body1213123", body)
	return body, nil
}
