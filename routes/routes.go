package routes

import (
	auth_middleware "github.com/wiklapandu/todo-api/middlewares/auth"

	"github.com/gofiber/fiber/v2"
)

func Routes() {
	app := fiber.New()

	v1 := app.Group("/api/v1")
	RouteAuth(v1)

	RoutesTodos(v1)

	v1.Use(auth_middleware.AuthCheck).Get("/test", func(c *fiber.Ctx) error {
		return c.Send([]byte("Successfully message"))
	})

	app.Listen(":3000")
}
