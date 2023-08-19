package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var (
	store   *session.Store
	USER_ID string = "userId"
	// AUTH_KEY string = "authenticated"
)

func Setup() {
	app := fiber.New()
	store = session.New(session.Config{
		CookieSecure:   true,
		CookieHTTPOnly: true,
	})

	// registration route before auth middleware
	app.Post("/register", Register)
	app.Post("/login", Login)
	app.Post("/logout", Logout)

	app.Use(returnSessionAuth())

	// run database
	// configs.ConnectDB()

	// user routes
	app.Get("/user", GetAllUsers)
	// app.Post("/user", CreateUser)

	// workout routes
	app.Get("/workout", GetAllWorkouts)
	app.Post("/workout", CreateWorkout)

	app.Listen(":8080")
}
