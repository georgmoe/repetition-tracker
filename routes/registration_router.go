package routes

import (
	"github.com/georgmoe/repetition-tracker/controllers"
	"github.com/gofiber/fiber/v2"
)

func RegistrationRouter(app *fiber.App) {
	app.Post("/register", controllers.Register)
}
