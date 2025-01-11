package logic

import (
	"context"
	"errors"

	"tuwei/babel_fish/internal/model"
	"tuwei/babel_fish/internal/svc"
	"tuwei/babel_fish/internal/types"
	"tuwei/babel_fish/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLogic) User(req *types.UserRequest) (resp *types.UserResponse, err error) {
	exist, err := l.svcCtx.UserModel.Exist(l.ctx, req.UserName)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errors.New("user exist")
	}
	hashPwd, err := util.HashPassword(l.ctx, req.UserPwd)
	if err != nil {
		return nil, err
	}
	u := model.User{
		UserName: req.UserName,
		UserPwd:  hashPwd,
	}
	_, err = l.svcCtx.UserModel.Create(l.ctx, u)
	if err != nil {
		return nil, err
	}
	return &types.UserResponse{}, nil
}
