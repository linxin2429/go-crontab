package master

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
)

// Global_Config 全局配置
var Global_Config *Config

// Config 配置
type Config struct {
	ApiPort         int      `json:"api_port"`
	ApiReadTimeout  int      `json:"api_read_timeout"`
	ApiWriteTimeout int      `json:"api_write_timeout"`
	EtcdEndpoints   []string `json:"etcd_endpoints"`
	EtcdDialTimeout int      `json:"etcd_dial_timeout"`
	EtcdUsername    string   `json:"etcd_username"`
	EtcdPwd         string   `json:"etcd_pwd"`
	LogFilename     string   `json:"log_filename"`
	Webroot         string   `json:"webroot"`
}

// InitConfig 初始化配置
func InitConfig(filename string) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "ioutil readfile error")
	}
	config := &Config{}
	err = json.Unmarshal(content, config)
	if err != nil {
		return errors.Wrap(err, "json")
	}
	Global_Config = config
	return nil
}
