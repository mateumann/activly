package main

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/lib/pq"
	"github.com/mateumann/activly/adapter/handler"
	"github.com/mateumann/activly/adapter/repository"
	"github.com/mateumann/activly/core/services"
)

func main() {
	store := repository.NewUserPostgresRepository()
	service := services.NewUserService(store)

	app := initApp(service)

	panic(app.Listen(":8080"))
}

func initApp(userService *services.UserService) *fiber.App {
	app := fiber.New()
	h := handler.NewFiberHandler(userService)

	app.Use(logger.New(logger.Config{
		Format:     "${time} | ${status} | ${latency} | ${ip} | ${method} | ${path} | ${error}\n",
		TimeFormat: time.RFC3339Nano,
	}))

	app.Use(recover.New())

	app.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return userService.Ready() // && ...
		},
		ReadinessEndpoint: "/ready",
	}))

	app.Get("/users", h.ListUsers)

	return app
}
