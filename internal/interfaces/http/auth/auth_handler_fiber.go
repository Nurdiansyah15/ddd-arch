package auth

import (
	"github.com/Nurdiansyah15/ddd-arch/internal/usecase/auth"
	"github.com/gofiber/fiber/v2"
)

type AuthHandlerFiber struct {
	LoginUC *auth.LoginUsecase
}

func (h *AuthHandlerFiber) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"error": "invalid request"})
	}

	resp, err := h.LoginUC.Execute(auth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"error": "invalid credentials"})
	}

	return c.JSON(resp)
}
