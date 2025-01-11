package logic

import (
	"context"
	"net/http"
	"os"

	"tuwei/babel_fish/internal/model"
	"tuwei/babel_fish/internal/svc"
	"tuwei/babel_fish/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DownloadTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDownloadTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DownloadTaskLogic {
	return &DownloadTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DownloadTaskLogic) DownloadTask(w http.ResponseWriter, req *types.DownloadTaskRequest) (resp *types.DownloadTaskResponse, err error) {
	uid := uint(0)
	task, err := l.svcCtx.TaskModel.GetTask(l.ctx, uid, uint(req.ID))
	if err != nil {
		return nil, err
	}
	if task.State == model.SuccessTaskState {
		content, err := os.ReadFile(task.ResultFileName)
		if err != nil {
			return nil, err
		}
		_, err = w.Write(content)
		if err != nil {
			return nil, err
		}

	}

	resp = &types.DownloadTaskResponse{
		State: int(task.State),
	}
	return resp, nil
}
