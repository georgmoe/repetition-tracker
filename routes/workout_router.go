package routes

import (
	"github.com/georgmoe/repetition-tracker/controllers"
	"github.com/gofiber/fiber/v2"
)

func WorkoutRouter(app *fiber.App) {
	app.Get("/workout", controllers.GetAllWorkouts)
	// app.Get("/user/:userId", controllers.GetAUser)
	app.Post("/workout", controllers.CreateWorkout)
	// app.Put("/user/:userId", controllers.EditAUser)
	// app.Delete("/user/:userId", controllers.DeleteAUser)
}
