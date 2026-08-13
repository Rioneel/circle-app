package trust

import (

	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type trustRequest struct {
	FromID string `json:"fromId"`
	ToID string `json:"toId"`
	Weight float64 `json:"weight"`
}
func (h *Handler) CreateTrust(c *fiber.Ctx) error {
	var req trustRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"invalid request body"})
	}
	if err := h.service.CreateTrust(c.Context(), req.FromID, req.ToID, req.Weight); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":err.Error()})
	}
	return nil
}

func (h *Handler) UpdateTrust(c *fiber.Ctx) error {
	var req trustRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":"invalid request body"})
	}
	if err := h.service.UpdateTrust(c.Context(), req.FromID, req.ToID, req.Weight); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error":err.Error()})
	}
	return nil
}

func (h *Handler) RemoveTrust(c *fiber.Ctx) error {
	if err := h.service.RemoveTrust(c.Context(), c.Params("fromId"),c.Params("toId")); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	 }
	 return c.SendStatus(fiber.StatusNoContent)

}

func (h *Handler) ListTrusted(c *fiber.Ctx) (error) {
	users, err := h.service.ListTrusted(c.Context(), c.Params("userId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)

}

