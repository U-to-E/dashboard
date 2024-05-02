package routes

import (
	"github.com/U-to-E/dashboard/controller"
	"github.com/gofiber/fiber/v3"
)

func SetupStudentRoutes(app *fiber.App) {
	app.Get("/", controller.RenderLogin)
}
