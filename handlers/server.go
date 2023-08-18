package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

var (
	store    *session.Store
	AUTH_KEY string = "authenticated"
	USER_ID  string = "user_id"
)

func Setup() {
	app := fiber.New()
	store = session.New(session.Config{
		CookieHTTPOnly: true,
		// CookieSecure: true, for https
		Expiration: 5 * time.Hour,
	})

	// registration route before auth middleware
	app.Post("/register", Register)
	app.Post("/login", Login)

	app.Use(ReturnSessionAuth())

	// run database
	// configs.ConnectDB()

	// user routes
	app.Get("/user", GetAllUsers)
	app.Post("/user", CreateUser)

	// workout routes
	app.Get("/workout", GetAllWorkouts)
	app.Post("/workout", CreateWorkout)

	app.Listen(":8080")
}
