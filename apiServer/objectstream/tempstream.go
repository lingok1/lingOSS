package objectstream

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type TempPutStream struct {
	Server string
	Uuid   string
}

func NewTempPutStream(server, object string, size int64) (*TempPutStream, error) {
	log.Println("NewTempPutStream1111111")
	request, err := http.NewRequest("POST", "http://"+server+"/temp/"+object, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("size", fmt.Sprintf("%d", size))
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	uuid, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return &TempPutStream{
		server,
		string(uuid),
	}, nil
}

// 以patch方法访问数据服务的temp接口，将需要写入的数据上传
// 用在io.TeeReadr 处 TempPutStream实现了Write([]byte)(int,err) 是一个io.Reader
func (w *TempPutStream) Write(p []byte) (n int, err error) {
	log.Println("Write222222222")
	request, err := http.NewRequest("PATCH", "http://"+w.Server+"/temp/"+w.Uuid, strings.NewReader(string(p)))
	if err != nil {
		return 0, err
	}
	client := http.Client{}

	response, err := client.Do(request) //向数据节点发送patch请求 得到回复
	if err != nil {
		return 0, err
	}
	if response.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("dataServer return http code %d", response.StatusCode)
	}
	return len(p), nil
}

// commit 根据传入参数决定项数据节点发送put还是delete
func (w *TempPutStream) Commit(good bool) {
	log.Println("Commit3333333")
	method := "DELETE"
	if good {
		method = "PUT"
	}
	request, _ := http.NewRequest(method, "http://"+w.Server+"/temp/"+w.Uuid, nil)
	client := http.Client{}
	client.Do(request)
}

func NewTempGetStream(server, uuid string) (*GetStream, error) {
	return newGetStream("http://" + server + "/temp/" + uuid)
}
