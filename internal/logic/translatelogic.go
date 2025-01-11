package logic

import (
	"context"

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
	task, err := l.svcCtx.TaskModel.GetTask(l.ctx, uint(req.ID))
	if err != nil {
		return nil, err
	}
	/*
		采用乐观锁的方式,避免并发调用 llm
		可能会出现 llm 调用失败, state 更新成功的问题, 低概率时间,使用 check 人物进行兜底
	*/
	if task.State == model.WaitTaskState {
		ok, err := l.svcCtx.TaskModel.UpdateState(l.ctx, task.ID, task.State, model.DoingTaskState)
		if err != nil {
			return nil, err
		}
		if !ok {
			return &types.TranslateResponse{}, nil
		}
		sourceLang, destLangs, contents, err := util.ReadTask(l.ctx, task.TaskFileName)
		if err != nil {
			return nil, err
		}
		results, err := l.svcCtx.LLM.Translate(l.ctx, sourceLang, destLangs, contents)
		if err != nil {
			ok, err = l.svcCtx.TaskModel.UpdateState(l.ctx, task.ID, model.DoingTaskState, model.FailedTaskState)
			if err != nil {
				l.Logger.Errorf("update task state failed")
				return nil, err
			}
			return nil, err
		}

		err = util.WriteResult(l.ctx, task.ResultFileName, results)
		if err != nil {
			return nil, err
		}
		ok, err = l.svcCtx.TaskModel.UpdateState(l.ctx, task.ID, model.DoingTaskState, model.SuccessTaskState)
		if err != nil {
			ok, err = l.svcCtx.TaskModel.UpdateState(l.ctx, task.ID, model.DoingTaskState, model.FailedTaskState)
			if err != nil {
				l.Logger.Errorf("update task state failed")
				return nil, err
			}
			return nil, err
		}

		return &types.TranslateResponse{}, nil
	}

	return &types.TranslateResponse{}, nil
}
