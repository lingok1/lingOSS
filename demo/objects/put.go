package objects

import (
	"OSS/apiServer/rs"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	log.Println("r.Body", r.Body)
	fileName := strings.Split(r.URL.EscapedPath(), "/")[2]
	log.Println("fileName", fileName)
	// ----------------------------------------------
	const maxUploadSize = 20 * 1024 * 1024 // 20 mb
	// parse and validate file and post parameters
	file, fileHeader, err := r.FormFile(fileName)
	if err != nil {
		return
	}
	defer file.Close()
	// Get and print out file size
	fileSize := fileHeader.Size
	fmt.Printf("File size (bytes): %v\n", fileSize)
	// validate file size
	if fileSize > maxUploadSize {
		return
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	log.Println("fileBytes", fileBytes)
	log.Println("len(fileBytes)1111", len(fileBytes))
	detectedFileType := http.DetectContentType(fileBytes)
	log.Println("detectedFileType", detectedFileType)

	if PutObject(fileBytes, fileName) != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// ----------------------------------------------
	// log.Println("w", w)
	// f, e := os.Create(FileDIr + strings.Split(r.URL.EscapedPath(), "/")[2])
	// if e != nil {
	// 	log.Println(e)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }
	// defer f.Close()
	// io.Copy(f, r.Body)
}

var FileDIr string = "/home/ling/go/src/Orginone/pkg/minio/demo/data/"

func PutObject(body []byte, fileName string) error {
	var err error
	writers := make([]io.Writer, rs.ALL_SHARDS)
	for i := range writers {
		writers[i], err = os.Create(fmt.Sprintf(FileDIr+fileName+".%d", i))
		if err != nil {
			return err
		}
	}
	enc := rs.NewEncoder(writers)
	length := len(body)
	for count := 0; count != length; {
		n, err := enc.Write(body[count:])
		if err != nil {
			log.Println(err)
			return err
		}
		count += n
	}
	enc.Flush()
	for i := range writers {
		writers[i].(*os.File).Close()
		writers[i] = nil
	}
	return nil
}
