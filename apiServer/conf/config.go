package conf

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

//全局配置变量
var Conf *Config = &Config{}

type Config struct {
	RabbitmqAddr string //rabbitmq服务器地址
	EsAddr       string //es服务器地址
	ListenAddr   string //监听地址
	ListenPort   string //监听端口
}

func (conf *Config) Parse(confile string) {
	if _, err := os.Stat(confile); os.IsNotExist(err) {
		panic(err)
	} else {
		f, err := os.Open(confile)
		if err != nil {
			panic(err)
		}
		data, _ := ioutil.ReadAll(f)
		if err := json.Unmarshal(data, conf); err != nil {
			panic(err)
		}
	}
}
