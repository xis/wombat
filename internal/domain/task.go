package domain

import (
	"context"
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
	CreateTask(ctx context.Context, params CreateTaskParams) (Task, error)
}

type GetPendingTaskParams struct {
	WorkerID string
	QueueID  string
}

type CreateTaskParams struct {
	QueueID  string
	Priority int
	Payload  map[string]any
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
	Priority  int
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

func (s *taskService) CreateTask(ctx context.Context, params CreateTaskParams) (Task, error) {
	task := Task{
		ID:        xid.New().String(),
		QueueID:   params.QueueID,
		Status:    TaskStatusPending,
		CreatedAt: time.Now(),
		Priority:  params.Priority,
		Payload:   params.Payload,
	}

	_, err := s.repository.CreateTask(ctx, task)
	if err != nil {
		return Task{}, err
	}

	return task, nil
}
