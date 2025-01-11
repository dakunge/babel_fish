package svc

import (
	"log"
	"tuwei/babel_fish/internal/config"
	"tuwei/babel_fish/internal/llm"
	"tuwei/babel_fish/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	LLM       llm.LLM
	TaskModel model.TaskModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	dsn := "root:@tcp(127.0.0.1:3306)/babel_fish?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("open db err ")
	}

	llm := llm.NewLLM()
	taskModel := model.NewTaskModel(db)
	return &ServiceContext{
		Config:    c,
		LLM:       llm,
		TaskModel: taskModel,
	}
}
