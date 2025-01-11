package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	TaskPath   string `json:"task_path"`
	ResultPath string `json:"result_path"`
}
