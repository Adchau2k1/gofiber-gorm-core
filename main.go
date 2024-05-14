package main

import (
	"backend/v1/database"
	"backend/v1/router"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func init() {
	database.Connect()
}

func main() {
	app := fiber.New()
	app.Use(cors.New())
	router.SetupRoutes(app)
	log.Fatal(app.Listen(":8000"))
}
