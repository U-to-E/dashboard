package routes

import (
	"github.com/U-to-E/dashboard/controller"
	"github.com/U-to-E/dashboard/middleware"
	"github.com/gofiber/fiber/v3"
)

func SetupStudentRoutes(app *fiber.App) {
	app.Get("/", controller.RenderLogin)
	app.Get("/adminpanel", middleware.Protected, controller.RenderAdmin)
	app.Post("/signup", controller.Register)
	app.Post("/login", controller.Handlelogin)
	app.Get("/dashboard", middleware.Protected, controller.RenderDashboard)
	app.Post("/logout", middleware.Protected, controller.Logout)
}
