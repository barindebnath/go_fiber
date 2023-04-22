package router

import (
	controller "github.com/barindebnath/gofiber/controllers"
	"github.com/gofiber/fiber/v2"
)

func Router(app *fiber.App) {
	userRoutes(app)
	websocketRoute(app)
}

func userRoutes(app *fiber.App) {
	userGroup := app.Group("/user")

	userGroup.Post("/signup", controller.SignUp)
	userGroup.Post("/signin", controller.SignIn)
	userGroup.Post("/logout", controller.LogOut)
}

func websocketRoute(app *fiber.App) {
	app.Get("/ws", controller.WS)
}
