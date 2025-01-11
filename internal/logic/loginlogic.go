package logic

import (
	"context"
	"errors"
	"time"

	"tuwei/babel_fish/internal/svc"
	"tuwei/babel_fish/internal/types"
	"tuwei/babel_fish/internal/util"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// 不暴露错误信息,返回同样的错误
func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.LoginResponse, err error) {
	user, err := l.svcCtx.UserModel.GetPwd(l.ctx, req.UserName)
	if err != nil {
		return nil, errors.New("user or password error")
	}
	if !util.CheckPasswordHash(req.UserPwd, user.UserPwd) {
		return nil, errors.New("user or password error")
	}
	secret := l.svcCtx.Config.Auth.AccessSecret
	expires := l.svcCtx.Config.Auth.AccessExpire
	iat := time.Now().Unix()
	// 生成 JWT 令牌
	tokenString, err := util.GenerateJWT(l.ctx, secret, user.ID, iat, expires)
	if err != nil {
		return nil, err
	}

	return &types.LoginResponse{
		Token: tokenString,
	}, nil

}
