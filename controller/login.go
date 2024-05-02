package controller

import "github.com/gofiber/fiber/v3"

func RenderLogin(c fiber.Ctx) error {
	return c.Render("login", fiber.Map{})
}
