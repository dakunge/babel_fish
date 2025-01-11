package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	Task struct {
		TaskPath       string `json:"task_path"`
		ResultPath     string `json:"result_path"`
		RetryThreshold int    `json:"retry_threshold"`
		LLMMaxCount    int    `json:"llm_max_count"`
	}
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}
}
