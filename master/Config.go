package master

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	ApiPort int 	`json:"apiPort"`
	ApiReadTimeout int 	`json:"apiReadTimeout"`
	ApiWriteTimeout int 	`json:"apiWriteTimeout"`
	EtcdEndpoints	[]string	`json:"etcdEndpoints"`
	EtcdDialTimeout	int		`json:"etcdDialTimeout"`
	WebRoot string `json:"webroot"`
}
var (
	G_config *Config
)

//加载配置
func InitConfig(filename string) (err error) {
	var (
		content []byte
		config Config
	)
	//把配置文件读取进来
	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}
	//json反序列化
	if err = json.Unmarshal(content, &config); err != nil {
		return
	}
	G_config = &config
	return
}
