package logic

import (
	"context"

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

func (l *DownloadTaskLogic) DownloadTask(req *types.DownloadTaskRequest) (resp *types.DownloadTaskResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
