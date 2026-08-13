package main

import (
	"github.com/gofiber/fiber/v2"

	"circle/internal/provider"
	"circle/internal/trust"
	"circle/internal/user"
)

func RegisterRoutes(app *fiber.App, trustHandler *trust.Handler, providerHandler *provider.Handler, userHandler *user.Handler) {
	app.Post("/api/trust", trustHandler.CreateTrust)
	app.Patch("/api/trust", trustHandler.UpdateTrust)
	app.Delete("/api/trust/:fromId/:toId", trustHandler.RemoveTrust)
	app.Get("/api/trust/:userId", trustHandler.ListTrusted)

	app.Get("/api/search", providerHandler.Search)
	app.Get("/api/providers/:id", providerHandler.GetProvider)
	app.Get("/api/providers/:id/recommendations", providerHandler.Recommendations)
	app.Post("/api/recommend", providerHandler.Recommend)

	app.Get("/api/users", userHandler.ListUsers)
	app.Get("/api/path", userHandler.TrustPath)
}