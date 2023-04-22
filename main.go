package main

import (
	"os"

	controller "github.com/barindebnath/gofiber/controllers"
	"github.com/barindebnath/gofiber/router"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/websocket/v2"
)

func main() {
	err := run()
	if err != nil {
		panic(err)
	}
}

func run() error {

	// create app
	app := fiber.New()

	// add basic middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	app.Use(func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	go controller.WSRunHub()

	// add routes
	router.Router(app)

	// start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	app.Listen(":" + port)

	return nil
}
