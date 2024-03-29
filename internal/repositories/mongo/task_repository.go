package mongo

import (
	"context"
	"fmt"
	"wombat/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	tasksCollectionName = "tasks"
)

type MongoTaskRepository struct {
	db    *mongo.Database
	tasks *mongo.Collection
}

type MongoTaskRepositoryDependencies struct {
	Database *mongo.Database
}

func NewMongoTaskRepository(deps MongoTaskRepositoryDependencies) *MongoTaskRepository {
	return &MongoTaskRepository{
		db:    deps.Database,
		tasks: deps.Database.Collection(tasksCollectionName),
	}
}

func (r *MongoTaskRepository) GetPendingTask(ctx context.Context, params domain.GetPendingTaskParams) (domain.Task, error) {
	var t task

	filter := bson.M{
		"queue_id": params.QueueID,
		"status":   domain.TaskStatusPending,
	}

	opts := options.FindOneAndUpdate().
		SetReturnDocument(options.After).
		SetSort(bson.M{
			"created_at": 1,
		})

	update := bson.M{
		"$set": bson.M{
			"status": domain.TaskStatusProcessing,
		},
	}

	err := r.tasks.FindOneAndUpdate(ctx, filter, update, opts).Decode(&t)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Task{}, domain.ErrNoPendingTasks
		}

		return domain.Task{}, err
	}

	convertedTask, err := t.toDomain()
	if err != nil {
		return domain.Task{}, err
	}

	return convertedTask, nil
}

func (r *MongoTaskRepository) UpdateTaskStatus(ctx context.Context, taskID string, status domain.TaskStatus) error {
	result, err := r.tasks.UpdateOne(ctx, bson.M{
		"_id": taskID,
	}, bson.M{
		"$set": bson.M{
			"status": status,
		},
	})
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("%w, task id: %s", domain.ErrTaskNotFound, taskID)
	}

	return nil
}

func (r *MongoTaskRepository) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	converted := newTaskFromDomain(task)

	_, err := r.tasks.InsertOne(ctx, converted)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}
