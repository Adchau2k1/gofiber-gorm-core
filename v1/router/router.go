package router

import (
	"backend/v1/config"
	"backend/v1/handler"
	"backend/v1/middleware"
	"backend/v1/response"

	"github.com/gofiber/fiber/v2"
)

// SetupRoutes setup router api
func SetupRoutes(app *fiber.App) {
	secretKey := config.GetConfigByKey("SECRET_KEY")
	auth := middleware.Auth(secretKey)

	// No auth
	app.Post("/api/v1/login", handler.Auth.Login)
	app.Post("/api/v1/refresh", handler.Auth.Refresh)
	app.Post("/api/v1/users", handler.User.CreateUser)

	// User
	user := app.Group("/api/v1/users", auth)
	user.Get("/", handler.User.GetUser)
	user.Patch("/:id", handler.User.UpdateUser)
	user.Delete("/:id", handler.User.DeleteUser)
	user.Put("/reset", handler.User.ResetUser)

	// Notfound router
	app.Get("/", func(c *fiber.Ctx) error {
		return response.Success(c, "Welcome to Go Fiber + GORM api ^-^", nil)
	})
	app.All("*", func(c *fiber.Ctx) error {
		return response.Custom(c, fiber.StatusNotFound, true, "404 Not Found", nil)
	})
}
