package main

import (
	"context"
	"os"
	"wombat/internal/controllers"
	"wombat/internal/domain"
	mongorepo "wombat/internal/repositories/mongo"
	"wombat/internal/servers"

	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	DefaultDatabaseName = "wombat"
	DefaultHTTPAddr     = ":8080"
)

type EnvironmentVariables struct {
	MongoURI string `env:"MONGO_URI"`
	HTTPAddr string `env:"HTTP_ADDR"`
}

func main() {
	ctx := context.Background()

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	variables := getEnvironmentVariables()

	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(variables.MongoURI))
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to connect to MongoDB")
	}

	mongoTaskRepository := mongorepo.NewMongoTaskRepository(mongorepo.MongoTaskRepositoryDependencies{
		Database: mongoClient.Database(DefaultDatabaseName),
	})

	taskService := domain.NewTaskService(domain.TaskServiceDependencies{
		TaskRepository: mongoTaskRepository,
	})

	taskController := controllers.NewTaskController(controllers.TaskControllerDependencies{
		TaskService: taskService,
	})

	httpServer := servers.NewHTTPServer(servers.HTTPServerDependencies{
		TaskController: taskController,
	})

	httpServer.RegisterRoutes()

	if variables.HTTPAddr == "" {
		variables.HTTPAddr = DefaultHTTPAddr
	}

	err = httpServer.Start(variables.HTTPAddr)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to start HTTP server")
	}
}

func getEnvironmentVariables() EnvironmentVariables {
	return EnvironmentVariables{
		MongoURI: os.Getenv("MONGO_URI"),
		HTTPAddr: os.Getenv("HTTP_ADDR"),
	}
}
