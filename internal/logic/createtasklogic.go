package logic

import (
	"context"
	"io"
	"net/http"
	"os"
	"path"

	"tuwei/babel_fish/internal/model"
	"tuwei/babel_fish/internal/svc"
	"tuwei/babel_fish/internal/types"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

const maxFileSize = 10 << 20 // 10 MB

type CreateTaskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateTaskLogic {
	return &CreateTaskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateTaskLogic) CreateTask(r *http.Request, req *types.CreateTaskRequest) (resp *types.CreateTaskResponse, err error) {
	_ = r.ParseMultipartForm(maxFileSize)
	file, handler, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	internalName := path.Join(l.svcCtx.Config.TaskPath, uuid.New().String())
	tempFile, err := os.Create(internalName)
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()
	_, err = io.Copy(tempFile, file)
	if err != nil {
		return nil, err
	}

	task := model.Task{
		UserFileName:     handler.Filename,
		InternalFileName: internalName,
	}
	id, err := l.svcCtx.TaskModel.Create(l.ctx, task)
	if err != nil {
		return nil, err
	}

	resp = &types.CreateTaskResponse{
		ID: int(id),
	}
	return resp, nil
}
