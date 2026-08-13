package main

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"circle/internal/config"
	"circle/internal/db"
	"circle/internal/provider"
	"circle/internal/trust"
	"circle/internal/user"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	driver, err := db.NewDriver(ctx, cfg)
	if err != nil {
		log.Fatalf("could not connect to CognoDB: %v", err)
	}
	defer driver.Close(ctx)

	trustHandler := trust.NewHandler(trust.NewService(driver))
	providerHandler := provider.NewHandler(provider.NewService(driver))
	userHandler := user.NewHandler(user.NewService(driver))

	app := fiber.New(fiber.Config{AppName: "circle-api"})
	app.Use(logger.New())
	app.Use(cors.New())

	app.Get("/healthz", func(c *fiber.Ctx) error {
		if err := db.HealthCheck(c.Context(), driver); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"status": "unhealthy"})
		}
		return c.JSON(fiber.Map{"status": "healthy"})
	})

	RegisterRoutes(app, trustHandler, providerHandler, userHandler)

	app.Static("/", "./web/dist")

	log.Printf("circle-api listening on :%s", cfg.Port)
	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}

}