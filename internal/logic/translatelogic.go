package logic

import (
	"context"
	"encoding/json"

	"tuwei/babel_fish/internal/model"
	"tuwei/babel_fish/internal/svc"
	"tuwei/babel_fish/internal/types"
	"tuwei/babel_fish/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type TranslateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTranslateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TranslateLogic {
	return &TranslateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TranslateLogic) Translate(req *types.TranslateRequest) (resp *types.TranslateResponse, err error) {
	uid64, _ := l.ctx.Value("userID").(json.Number).Int64()
	uid := uint(uid64)
	task, err := l.svcCtx.TaskModel.GetTask(l.ctx, uid, uint(req.ID))
	if err != nil {
		return nil, err
	}
	if task.LLMCallCount >= l.svcCtx.Config.Task.LLMMaxCount {
		_, err := l.svcCtx.TaskModel.UpdateState(l.ctx, task.ID, task.State, model.FinalFailedTaskState, task.LLMCallCount)
		if err != nil {
			return nil, err
		}
		return &types.TranslateResponse{}, nil
	}
	/*
		采用乐观锁的方式,避免并发调用 llm
		可能会出现 llm 调用失败, state 更新成功的问题, 低概率时间,使用 check 人物进行兜底
	*/
	switch task.State {
	case model.FailedTaskState:
		fallthrough
	case model.WaitTaskState:
		err := l.TranslateImpl(task, model.DoingTaskState)
		if err != nil {
			return nil, err
		}
		return &types.TranslateResponse{}, nil
	default:
		return &types.TranslateResponse{}, nil
	}
}

func (l *TranslateLogic) TranslateImpl(task model.Task, targetState model.TaskState) error {
	ok, err := l.svcCtx.TaskModel.UpdateState(l.ctx, task.ID, task.State, targetState, task.LLMCallCount)
	if err != nil {
		return err
	}
	if !ok {
		return err
	}
	sourceLang, destLangs, contents, err := util.ReadTask(l.ctx, task.TaskFileName)
	if err != nil {
		return err
	}
	results, err := l.svcCtx.LLM.Translate(l.ctx, sourceLang, destLangs, contents)
	if err != nil {
		ok, err2 := l.svcCtx.TaskModel.UpdateState(l.ctx, task.ID, targetState, model.FailedTaskState, task.LLMCallCount+1)
		if err2 != nil {
			l.Logger.Errorf("update task state failed err %v", err)
			return err
		}
		l.Infof("update state failed success %v, %v", ok, task)
		return err
	}

	err = util.WriteResult(l.ctx, task.ResultFileName, results)
	if err != nil {
		return err
	}
	ok, err = l.svcCtx.TaskModel.UpdateState(l.ctx, task.ID, targetState, model.SuccessTaskState, task.LLMCallCount+1)
	if err != nil {
		l.Logger.Errorf("update task state success state failed")
		return err
	}
	return nil
}
