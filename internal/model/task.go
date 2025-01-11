package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
)

type TaskState int

const (
	WaitTaskState    TaskState = 0
	DoingTaskState   TaskState = 1
	SuccessTaskState TaskState = 2
	FailedTaskState  TaskState = 3
)

type Task struct {
	gorm.Model
	State          TaskState `json:"state"`
	UserFileName   string    `json:"user_filename"`
	TaskFileName   string    `json:"task_filename"`
	ResultFileName string    `json:"result_filename"`
}

type TaskModel interface {
	Create(ctx context.Context, t Task) (uint, error)
	GetTask(ctx context.Context, id uint) (Task, error)
	UpdateState(ctx context.Context, id uint, sourceState, destState TaskState) (bool, error)
}

func NewTaskModel(db *gorm.DB) TaskModel {
	return taskModel{db: db}
}

type taskModel struct {
	db *gorm.DB
}

func (m taskModel) Create(ctx context.Context, t Task) (uint, error) {
	if db := m.db.Create(&t); db.Error != nil {
		logc.Error(ctx, "create task err %v", db.Error)
		return 0, db.Error
	}
	return t.ID, nil
}

func (m taskModel) GetTask(ctx context.Context, id uint) (Task, error) {
	t := Task{}
	if db := m.db.Where("id = ?", id).First(&t); db.Error != nil {
		logc.Error(ctx, "get task id err %v", db.Error)
		return t, db.Error
	}
	return t, nil
}

func (m taskModel) UpdateState(ctx context.Context, id uint, sourceState, destState TaskState) (bool, error) {
	db := m.db.Model(&Task{}).Where("id = ? AND state = ?", id, int(sourceState)).Update("state", int(destState))
	if db.Error != nil {
		logc.Error(ctx, "update task state err %v", db.Error)
		return false, db.Error
	}
	if db.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}
