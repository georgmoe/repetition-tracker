package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/redis/v3"
)

var (
	store   *session.Store
	USER_ID string = "userId"
	// AUTH_KEY string = "authenticated"
)

func Setup() {
	app := fiber.New()
	storage := redis.New()
	store = session.New(session.Config{
		Storage:        storage,
		CookieSecure:   true,
		CookieHTTPOnly: true,
	})

	// registration route before auth middleware
	app.Post("/register", Register)
	app.Post("/login", Login)
	app.Post("/logout", Logout)

	app.Use(sessionAuth)

	// run database
	// configs.ConnectDB()

	// user routes
	app.Get("/user", GetAllUsers)
	// app.Post("/user", CreateUser)

	// workout routes
	app.Get("/workout", GetAllWorkouts)
	app.Get("/workout/:workoutId", GetWorkout)
	app.Post("/workout", PostWorkout)
	app.Put("/workout/:workoutId", PutWorkout)

	app.Post("/workout/:workoutId/exercise", PostExercise)
	app.Put("/workout/:workoutId/exercise/:exerciseIdx", PutExercise)
	app.Get("/workout/:workoutId/exercise/:exerciseIdx", GetExercise)

	// app.Post("/workout/:workoutId/exercise/:exerciseIdx", PostSet)
	// app.Put("/workout/:workoutId/exercise/:exerciseIdx/set/:setIdx", Put)

	app.Listen(":8080")
}
