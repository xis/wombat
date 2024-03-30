package mongo

import (
	"time"
	"wombat/internal/domain"
)

type task struct {
	ID        string            `bson:"_id"`
	QueueID   string            `bson:"queue_id"`
	Status    domain.TaskStatus `bson:"status"`
	CreatedAt time.Time         `bson:"created_at"`
	Priority  int               `bson:"priority"`
	Payload   map[string]any    `bson:"payload"`
}

func (t *task) toDomain() (domain.Task, error) {
	return domain.Task{
		ID:        t.ID,
		QueueID:   t.QueueID,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		Priority:  t.Priority,
		Payload:   t.Payload,
	}, nil
}

func newTaskFromDomain(t domain.Task) task {
	return task{
		ID:        t.ID,
		QueueID:   t.QueueID,
		Status:    t.Status,
		CreatedAt: t.CreatedAt,
		Priority:  t.Priority,
		Payload:   t.Payload,
	}
}
