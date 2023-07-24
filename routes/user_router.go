package routes

import (
	"github.com/georgmoe/repetition-tracker/controllers"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(app *fiber.App) {
	app.Get("/user", controllers.GetAllUsers)
	app.Get("/user/:userId", controllers.GetAUser)
	app.Post("/user", controllers.CreateUser)
	app.Put("/user/:userId", controllers.EditAUser)
	app.Delete("/user/:userId", controllers.DeleteAUser)
}
