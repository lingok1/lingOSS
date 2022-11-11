package es

import (
	"OSS/apiServer/conf"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Metadata struct {
	Name    string
	Version int
	Size    int64
	Hash    string
}

//getMetadata用于获取对象的名字和版本号来获取对象的元数据 URL = ESserver + index:metadata + type:_doc
func getMetadata(name string, versionId int) (meta Metadata, err error) {
	url := fmt.Sprintf("http://%s/metadata/_doc/%s_%d/_source",
		conf.Conf.EsAddr, name, versionId)
	r, err := http.Get(url)
	if err != nil {
		return
	}
	if r.StatusCode != http.StatusOK {
		err = fmt.Errorf("fail to get %s_%d: %d", name, versionId, r.StatusCode)
		return
	}
	result, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(result, &meta)
	return
}

type hit struct {
	Source Metadata `json:"_source"`
}

type searchResult struct {
	Hits struct {
		Total int
		Hits  []hit
	}
}

//以对象的名字作为参数，调用ES搜索API 在URL中指定了对象的名字
func SearchLatestVersion(name string) (meta Metadata, err error) {
	url := fmt.Sprintf("http://%s/metadata/_search?q=name:%s&size=1&sort=version:desc",
		conf.Conf.EsAddr, url.PathEscape(name))
	r, e := http.Get(url)
	if e != nil {
		return

	}
	if r.StatusCode != http.StatusOK {
		e = fmt.Errorf("fail to search latest metadata: %d", r.StatusCode)
		return
	}

	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	if len(sr.Hits.Hits) != 0 {
		meta = sr.Hits.Hits[0].Source

	}
	//fmt.Println(sr)
	return
}

//getMetadata 的封装接口
func GetMetadata(name string, version int) (Metadata, error) {
	if version == 0 {
		return SearchLatestVersion(name)
	}
	return getMetadata(name, version)
}

func PutMetadata(name string, version int, size int64, hash string) error {
	doc := fmt.Sprintf(`{"name":"%s","version":"%d","size":%d,"hash":"%s"}`,
		name, version, size, hash)
	fmt.Println(doc)
	client := http.Client{}
	url := fmt.Sprintf("http://%s/metadata/_doc/%s_%d?op_type=create",
		conf.Conf.EsAddr, name, version)
	fmt.Println(url)
	request, _ := http.NewRequest("PUT", url, strings.NewReader(doc))
	request.Header.Set("Content-Type", "application/json")
	r, err := client.Do(request)
	if err != nil {
		return err
	}
	if r.StatusCode == http.StatusConflict { //conflict 冲突 多个客户端上传同一个元数据,只上传一份 版本号+1
		return PutMetadata(name, version+1, size, hash)
	}
	if r.StatusCode != http.StatusCreated {
		result, _ := ioutil.ReadAll(r.Body)
		return fmt.Errorf("fail to put metadata: %d %s", r.StatusCode, string(result))
	}
	return nil
}

func AddVersion(name, hash string, size int64) error {
	version, e := SearchLatestVersion(name)
	if e != nil {
		return e

	}
	return PutMetadata(name, version.Version+1, size, hash)

}

func SearchAllVersions(name string, from, size int) ([]Metadata, error) {
	url := fmt.Sprintf("http://%s/metadata/_search?sort=name,version&from=%d&size=%d",
		conf.Conf.EsAddr, from, size)
	if name != "" {
		url += "&q=name:" + name

	}
	r, e := http.Get(url)
	if e != nil {
		return nil, e

	}
	metas := make([]Metadata, 0)
	result, _ := ioutil.ReadAll(r.Body)
	var sr searchResult
	json.Unmarshal(result, &sr)
	for i := range sr.Hits.Hits {
		metas = append(metas, sr.Hits.Hits[i].Source)

	}
	return metas, nil

}
