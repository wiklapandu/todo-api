package routes

import (
	auth_middleware "github.com/wiklapandu/todo-api/middlewares/auth"
	todo_service "github.com/wiklapandu/todo-api/services/todo"

	"github.com/gofiber/fiber/v2"
)

func RoutesTodos(route fiber.Router) {
	todoRoutes := route.Use(auth_middleware.AuthCheck).Group("/todo")

	todoRoutes.Post("/", todo_service.Store)
	todoRoutes.Put("/:id", todo_service.Update)
	todoRoutes.Delete("/:id", todo_service.Delete)
	todoRoutes.Get("/:id", todo_service.Show)
	todoRoutes.Get("/", todo_service.Index)
}
