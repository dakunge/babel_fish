package svc

import (
	"fmt"
	"log"
	"tuwei/babel_fish/internal/config"
	"tuwei/babel_fish/internal/llm"
	"tuwei/babel_fish/internal/model"

	"github.com/go-redis/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config    config.Config
	LLM       llm.LLM
	TaskModel model.TaskModel
	UserModel model.UserModel
	Redis     *redis.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	//dsn := "root:@tcp(127.0.0.1:3306)/babel_fish?charset=utf8mb4&parseTime=True&loc=Local"
	user := c.Mysql.User
	pwd := c.Mysql.Pwd
	host := c.Mysql.Host
	port := c.Mysql.Port
	database := c.Mysql.Database
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local", user, pwd, host, port, database)
	fmt.Println("dsn", dsn)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("open db err ")
	}

	host = c.Redis.Host
	port = c.Redis.Port
	rdb := redis.NewClient(&redis.Options{
		//Addr: "localhost:6379",
		Addr: fmt.Sprintf("%v:%v", host, port),
	})
	if rdb.Ping().Err() != nil {
		log.Fatalf("open redis err %v", rdb.Ping().Err())
	}

	err = db.AutoMigrate(&model.Task{})
	if err != nil {
		log.Fatal("automigrate err")
	}
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		log.Fatal("automigrate err")
	}

	llm := llm.NewLLM()
	taskModel := model.NewTaskModel(db)
	userModel := model.NewUserModel(db)
	return &ServiceContext{
		Config:    c,
		LLM:       llm,
		TaskModel: taskModel,
		UserModel: userModel,
		Redis:     rdb,
	}
}
