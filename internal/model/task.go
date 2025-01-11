package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/logc"
	"gorm.io/gorm"
)

type TaskState int

const (
	WaitTaskState        TaskState = 0
	DoingTaskState       TaskState = 1
	SuccessTaskState     TaskState = 2
	FailedTaskState      TaskState = 3
	FinalFailedTaskState TaskState = 4
)

type Task struct {
	gorm.Model
	UserID         uint      `json:"user_id"`
	State          TaskState `json:"state"`
	UserFileName   string    `json:"user_filename"`
	TaskFileName   string    `json:"task_filename"`
	ResultFileName string    `json:"result_filename"`
	LLMCallCount   int       `json:"llm_call_count"`
}

type TaskModel interface {
	Create(ctx context.Context, t Task) (uint, error)
	GetTask(ctx context.Context, uid, id uint) (Task, error)
	UpdateState(ctx context.Context, id uint, sourceState, destState TaskState, count int) (bool, error)
	// for monitor
	GetTasks(ctx context.Context, begin uint, state []TaskState) ([]*Task, error)
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

func (m taskModel) GetTask(ctx context.Context, uid, id uint) (Task, error) {
	t := Task{}
	if db := m.db.Where("id = ? AND user_id = ?", id, uid).First(&t); db.Error != nil {
		logc.Error(ctx, "get task id err %v", db.Error)
		return t, db.Error
	}
	return t, nil
}

func (m taskModel) GetTasks(ctx context.Context, begin uint, states []TaskState) ([]*Task, error) {
	ts := []*Task{}
	if db := m.db.Where("id >= ? AND state IN ?", begin, states).Find(&ts); db.Error != nil {
		logc.Error(ctx, "get task id err %v", db.Error)
		return ts, db.Error
	}
	return ts, nil
}

func (m taskModel) UpdateState(ctx context.Context, id uint, sourceState, destState TaskState, count int) (bool, error) {
	updateColumns := map[string]interface{}{
		"state":          destState,
		"llm_call_count": count,
	}
	db := m.db.Model(&Task{}).Where("id = ? AND state = ?", id, int(sourceState)).Updates(updateColumns)
	if db.Error != nil {
		logc.Error(ctx, "update task state err %v", db.Error)
		return false, db.Error
	}
	if db.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}
