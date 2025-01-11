package logic

import (
	"context"
	"encoding/json"

	"tuwei/babel_fish/internal/svc"
	"tuwei/babel_fish/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTaskLogic {
	return &GetTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTaskLogic) GetTask(req *types.GetTaskRequest) (resp *types.GetTaskResponse, err error) {
	// todo: add your logic here and delete this line
	uid64, _ := l.ctx.Value("userID").(json.Number).Int64()
	uid := uint(uid64)
	task, err := l.svcCtx.TaskModel.GetTask(l.ctx, uid, uint(req.ID))
	if err != nil {
		return nil, err
	}

	//task.UserFileName
	resp = &types.GetTaskResponse{
		ID:       int(task.ID),
		FileName: task.UserFileName,
		State:    int(task.State),
	}

	return resp, nil
}
