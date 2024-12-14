package routes

import (
	"github.com/wiklapandu/todo-api/services"

	"github.com/gofiber/fiber/v2"
)

func RouteAuth(route fiber.Router) {
	auth := route.Group("/auth")

	auth.Post("login", services.Login)
	auth.Post("register", services.Register)
}
