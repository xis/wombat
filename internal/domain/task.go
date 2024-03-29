package domain

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/rs/xid"
)

var (
	ErrTaskNotFound   = errors.New("task not found")
	ErrNoPendingTasks = errors.New("no pending tasks")
)

type TaskService interface {
	GetPendingTask(ctx context.Context, params GetPendingTaskParams) (Task, error)
	UpdateTaskStatus(ctx context.Context, taskID string, status TaskStatus) error
	CreateTask(ctx context.Context, queueID string, payload []byte) (Task, error)
}

type GetPendingTaskParams struct {
	WorkerID string
	QueueID  string
}

type TaskRepository interface {
	GetPendingTask(ctx context.Context, params GetPendingTaskParams) (Task, error)
	UpdateTaskStatus(ctx context.Context, taskID string, status TaskStatus) error
	CreateTask(ctx context.Context, task Task) (Task, error)
}

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusProcessing TaskStatus = "processing"
	TaskStatusCompleted  TaskStatus = "completed"
)

type Task struct {
	ID        string
	QueueID   string
	Status    TaskStatus
	CreatedAt time.Time
	Payload   map[string]any
}

type TaskServiceDependencies struct {
	TaskRepository TaskRepository
}

type taskService struct {
	repository TaskRepository
}

func NewTaskService(dependencies TaskServiceDependencies) TaskService {
	return &taskService{
		repository: dependencies.TaskRepository,
	}
}

func (s *taskService) GetPendingTask(ctx context.Context, params GetPendingTaskParams) (Task, error) {
	task, err := s.repository.GetPendingTask(ctx, params)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}

func (s *taskService) UpdateTaskStatus(ctx context.Context, taskID string, status TaskStatus) error {
	err := s.repository.UpdateTaskStatus(ctx, taskID, status)
	if err != nil {
		return err
	}

	return nil
}

func (s *taskService) CreateTask(ctx context.Context, queueID string, payload []byte) (Task, error) {
	p := map[string]any{}

	err := json.Unmarshal(payload, &p)
	if err != nil {
		return Task{}, err
	}

	task := Task{
		ID:        xid.New().String(),
		QueueID:   queueID,
		Status:    TaskStatusPending,
		CreatedAt: time.Now(),
		Payload:   p,
	}

	_, err = s.repository.CreateTask(ctx, task)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}
