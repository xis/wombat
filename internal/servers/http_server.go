package servers

import (
	"wombat/internal/controllers"

	"github.com/gofiber/fiber/v2"
)

type HTTPServer struct {
	app *fiber.App

	taskController *controllers.TaskController
}

type HTTPServerDependencies struct {
	TaskController *controllers.TaskController
}

func NewHTTPServer(deps HTTPServerDependencies) *HTTPServer {
	app := fiber.New()

	return &HTTPServer{
		app:            app,
		taskController: deps.TaskController,
	}
}

func (s *HTTPServer) RegisterRoutes() {
	s.app.Get("/queues/:queueID/tasks", s.taskController.GetPendingTask)
	s.app.Post("/queues/:queueID/tasks", s.taskController.CreateTask)
	s.app.Put("/queues/:queueID/tasks/:taskID", s.taskController.UpdateTaskStatus)
}

func (s *HTTPServer) Start(addr string) error {
	return s.app.Listen(addr)
}

func (s *HTTPServer) Stop() error {
	return s.app.Shutdown()
}
