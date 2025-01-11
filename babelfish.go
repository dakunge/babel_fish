package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"tuwei/babel_fish/internal/config"
	"tuwei/babel_fish/internal/handler"
	"tuwei/babel_fish/internal/logic"
	"tuwei/babel_fish/internal/model"
	"tuwei/babel_fish/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logc"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/babelfish-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)
	go monitor(ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

func monitor(s *svc.ServiceContext) {
	ctx := context.Background()
	states := []model.TaskState{model.DoingTaskState, model.FailedTaskState}
	begin := uint(0)
	for {
		time.Sleep(time.Second * 3)
		tasks, err := s.TaskModel.GetTasks(ctx, begin, states)
		if err != nil {
			logc.Errorf(ctx, "monitor get tasks err %v", err)
			continue
		}
		if len(tasks) == 0 {
			logc.Infof(ctx, "monitor all task is normal")
			continue
		}

		// 正常不需要排序, task 查出来就是有序的, 加上排序防止 GetTasks 以后修改
		sort.Slice(tasks, func(i, j int) bool {
			return tasks[i].ID < tasks[j].ID
		})
		// 提高性能
		begin = tasks[0].ID
		for _, t := range tasks {
			// 实际翻译任务已经完成,直接修改状态为 success
			_, err := os.Stat(t.ResultFileName)
			if os.IsExist(err) {
				s.TaskModel.UpdateState(ctx, t.ID, t.State, model.SuccessTaskState, t.LLMCallCount+1)
				logc.Infof(ctx, "task actually finsh, so set state to success")
				continue
			}

			maxCount := s.Config.Task.LLMMaxCount
			if t.LLMCallCount >= maxCount {
				s.TaskModel.UpdateState(ctx, t.ID, t.State, model.FinalFailedTaskState, t.LLMCallCount)
				continue
			}
			if t.State == model.FailedTaskState {
				logc.Infof(ctx, "retry task %v, %v, %v", t.ID, t.State, t.UpdatedAt)
				err := logic.NewTranslateLogic(ctx, s).TranslateImpl(*t, model.DoingTaskState)
				if err != nil {
					logc.Errorf(ctx, "retry task %v, err %v", t, err)
					continue
				}

			}
			if t.State == model.DoingTaskState {
				threshold := s.Config.Task.RetryThreshold
				if int(time.Now().Sub(t.UpdatedAt).Seconds()) > threshold {
					logc.Infof(ctx, "retry task %v, %v, %v", t.ID, t.State, t.UpdatedAt)
					err := logic.NewTranslateLogic(ctx, s).TranslateImpl(*t, model.DoingTaskState)
					if err != nil {
						logc.Errorf(ctx, "retry task %v, err %v", t, err)
						continue
					}
				}
			}
		}
	}
}
