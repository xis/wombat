package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"wombat/internal/domain"

	"github.com/gofiber/fiber/v2"
)

var (
	TaskStatusPending    = "pending"
	TaskStatusProcessing = "processing"
	TaskStatusCompleted  = "completed"
)

type TaskController struct {
	taskService domain.TaskService
}

type TaskControllerDependencies struct {
	TaskService domain.TaskService
}

func NewTaskController(deps TaskControllerDependencies) *TaskController {
	return &TaskController{
		taskService: deps.TaskService,
	}
}

type task struct {
	ID        string         `json:"id"`
	QueueID   string         `json:"queue_id"`
	Status    string         `json:"status"`
	CreatedAt int64          `json:"created_at"`
	Priority  int            `json:"priority"`
	Payload   map[string]any `json:"payload"`
}

func newTaskFromDomain(t domain.Task) task {
	return task{
		ID:        t.ID,
		QueueID:   t.QueueID,
		Status:    string(t.Status),
		CreatedAt: t.CreatedAt.Unix(),
		Priority:  t.Priority,
		Payload:   t.Payload,
	}
}

func (c *TaskController) GetPendingTask(ctx *fiber.Ctx) error {
	queueID := ctx.Params("queueID")

	tasks, err := c.taskService.GetPendingTask(ctx.Context(), domain.GetPendingTaskParams{
		QueueID: queueID,
	})
	if err != nil {
		if errors.Is(err, domain.ErrNoPendingTasks) {
			return ctx.SendStatus(fiber.StatusNoContent)
		}

		return err
	}

	converted := newTaskFromDomain(tasks)

	return ctx.JSON(converted)
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status"`
}

func (c *TaskController) UpdateTaskStatus(ctx *fiber.Ctx) error {
	taskID := ctx.Params("taskID")

	var req UpdateTaskStatusRequest

	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	var status domain.TaskStatus

	switch req.Status {
	case TaskStatusPending:
		status = domain.TaskStatusPending
	case TaskStatusProcessing:
		status = domain.TaskStatusProcessing
	case TaskStatusCompleted:
		status = domain.TaskStatusCompleted
	default:
		return fmt.Errorf("invalid status: %s", req.Status)
	}

	err := c.taskService.UpdateTaskStatus(ctx.Context(), taskID, status)
	if err != nil {
		return err
	}

	return nil
}

func (c *TaskController) CreateTask(ctx *fiber.Ctx) error {
	payload := []byte(ctx.Body())

	queueID := ctx.Params("queueID")
	priority := ctx.QueryInt("priority")

	payloadMap := make(map[string]any)

	err := json.Unmarshal(payload, &payloadMap)
	if err != nil {
		return err
	}

	task, err := c.taskService.CreateTask(ctx.Context(), domain.CreateTaskParams{
		QueueID:  queueID,
		Priority: priority,
		Payload:  payloadMap,
	})
	if err != nil {
		return err
	}

	converted := newTaskFromDomain(task)

	return ctx.JSON(converted)
}
