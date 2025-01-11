package model

import (
	"context"

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
	if tx := m.db.Create(&t); tx.Error != nil {
		return 0, tx.Error
	}
	return t.ID, nil
}

func (m taskModel) GetTask(ctx context.Context, id uint) (Task, error) {
	t := Task{}
	if db := m.db.Where("id = ?", id).First(&t); db.Error != nil {
		return t, db.Error
	}
	return t, nil
}

func (m taskModel) UpdateState(ctx context.Context, id uint, sourceState, destState TaskState) (bool, error) {
	db := m.db.Where("id = ? AND state = ?", id, sourceState).Update("state", destState)
	if db.Error != nil {
		return false, db.Error
	}
	if db.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}
