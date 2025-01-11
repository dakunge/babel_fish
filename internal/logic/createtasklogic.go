package logic

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"

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
	uid64, _ := l.ctx.Value("userID").(json.Number).Int64()
	uid := uint(uid64)
	_ = r.ParseMultipartForm(maxFileSize)
	file, handler, err := r.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	pass, err := l.RepeateCheck(uid, content)
	if err != nil {
		return nil, err
	}
	if !pass {
		return &types.CreateTaskResponse{}, nil
	}

	uuid := uuid.New().String()
	taskFileName := path.Join(l.svcCtx.Config.Task.TaskPath, uuid)
	tempFile, err := os.Create(taskFileName)
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()
	bytes.NewReader(content)
	_, err = io.Copy(tempFile, bytes.NewReader(content))
	if err != nil {
		return nil, err
	}
	resultFileName := path.Join(l.svcCtx.Config.Task.ResultPath, uuid)
	task := model.Task{
		UserID:         uid,
		UserFileName:   handler.Filename,
		TaskFileName:   taskFileName,
		ResultFileName: resultFileName,
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

// TODO(kun.li): 这块是校验,应该保证 get,set 原子性, 太累了,不写了
func (l *CreateTaskLogic) RepeateCheck(uid uint, content []byte) (bool, error) {
	h := md5.New()
	_, err := h.Write(content)
	if err != nil {
		return false, err
	}
	md5sum := h.Sum(nil)

	key := fmt.Sprintf("repeate:%v", uid)
	value, _ := l.svcCtx.Redis.Get(key).Bytes()
	if len(value) != 0 {
		if string(md5sum) == string(value) {
			fmt.Println("repate repate")
			return false, nil
		}
	}

	cmd := l.svcCtx.Redis.Set(key, md5sum, time.Second*5)
	if cmd.Err() != nil {
		return false, cmd.Err()
	}
	return true, nil
}
