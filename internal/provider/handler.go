package provider

import "github.com/gofiber/fiber/v2"

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Search(c *fiber.Ctx) error {
	results, err := h.service.Search(c.Context(), c.Query("userId"), c.Query("category"), c.Query("area"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(results)
}

func (h *Handler) GetProvider(c *fiber.Ctx) error {
	p, err := h.service.GetProvider(c.Context(), c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if p == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "provider not found"})
	}
	return c.JSON(p)
}

func (h *Handler) Recommendations(c *fiber.Ctx) error {
	recs, err := h.service.Recommendations(c.Context(), c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(recs)
}

type recommendRequest struct {
	UserID     string  `json:"userId"`
	ProviderID string  `json:"providerId"`
	Rating     float64 `json:"rating"`
}

func (h *Handler) Recommend(c *fiber.Ctx) error {
	var req recommendRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}
	if err := h.service.Recommend(c.Context(), req.UserID, req.ProviderID, req.Rating); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusCreated)
}