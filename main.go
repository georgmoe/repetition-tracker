package main

import (
	"github.com/georgmoe/repetition-tracker/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	//run database
	// configs.ConnectDB()

	//routes
	routes.UserRouter(app)

	app.Listen(":8080")
}
