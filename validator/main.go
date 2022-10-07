package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"validator/api"
	"validator/database"
	"validator/notifications"
)

var config Config

func main() {
	config = LoadAppConfig()
	notifications.SetConfig(config.NotificationsToken, config.NotificationsURL)

	database.Migrate(config.DB)
	api.SetupDB(database.Connect(config.DB, config.DBMaxConnections, config.DBMaxIdleConnections))

	app := fiber.New(fiber.Config{
		Concurrency: config.HttpConcurrency,
		BodyLimit:   config.HttpBodyLimit,
		Prefork:     config.Prefork,
		ProxyHeader: config.ProxyHeader,
	})
	app.Post("/v1/request", api.ProcessRequest)
	app.Listen(fmt.Sprintf("%s:%d", config.HttpAddr, config.HttpPort))
}
