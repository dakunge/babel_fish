package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `json:"user_name"`
	UserPwd  string `json:"user_pwd"`
}

type UserModel interface {
	Exist(ctx context.Context, userName string) (bool, error)
	Create(ctx context.Context, user User) (uint, error)
	Get(ctx context.Context, userName, userPwd string) (User, error)
	GetPwd(ctx context.Context, userName string) (User, error)
}

func NewUserModel(db *gorm.DB) UserModel {
	return &userModel{db: db}
}

type userModel struct {
	db *gorm.DB
}

func (m userModel) Exist(ctx context.Context, userName string) (bool, error) {
	u := User{}
	if db := m.db.Where("user_name = ?", userName).Find(&u); db.Error != nil {
		logc.Error(ctx, "get user name %v err %v", userName, db.Error)
		return false, db.Error
	}
	if u.ID == 0 {
		return false, nil
	}
	return true, nil
}

func (m userModel) GetPwd(ctx context.Context, userName string) (User, error) {
	u := User{}
	if db := m.db.Where("user_name = ?", userName).First(&u); db.Error != nil {
		logc.Error(ctx, "get user pwd %v err %v", userName, db.Error)
		return u, db.Error
	}
	return u, nil
}

func (m userModel) Get(ctx context.Context, userName, userPwd string) (User, error) {
	u := User{}
	if db := m.db.Where("user_name = ? AND user_pwd = ?", userName, userPwd).First(&u); db.Error != nil {
		logc.Error(ctx, "get user %v err %v", userName, db.Error)
		return u, db.Error
	}
	return u, nil
}

func (m userModel) Create(ctx context.Context, user User) (uint, error) {
	if db := m.db.Create(&user); db.Error != nil {
		logc.Error(ctx, "create user %v, err %v", user, db.Error)
		return 0, db.Error
	}
	return user.ID, nil
}
