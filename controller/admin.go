package controller

import "github.com/gofiber/fiber/v3"

func RenderAdmin(c fiber.Ctx) error {

	if c.IP() != "127.0.0.1" && c.IP() != "::1" {
		return c.Status(fiber.StatusForbidden).SendString("Forbidden")
	}

	return c.Render("admin", fiber.Map{})
}
