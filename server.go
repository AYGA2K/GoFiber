package main

import (
	"log"

	"example.com/api/middleware"
	"example.com/api/routes"

	"example.com/api/database"
	"github.com/gofiber/fiber/v2"
)

func welcome(c *fiber.Ctx) error {
	return c.SendString("Welcome to an Awesome API")
}

func setupRoutes(app *fiber.App) {
	// Welcome endpoint
	app.Get("/api", welcome)

	// User endpoints
	app.Post("/api/users", routes.SignUp)
	app.Post("/api/login", routes.Login)
	app.Get("/api/users", middleware.AuthMiddleware, routes.GetUsers)
	app.Get("/api/users/:id", routes.GetUser)
	app.Delete("/api/users/:id", routes.DeleteUser)

}

func main() {
	database.ConnectDb()

	app := fiber.New()
	tokenApp := fiber.New()
	setupRoutes(app)

	log.Fatal(app.Listen(":3000"))
	log.Fatal(tokenApp.Listen(":4000"))
}
