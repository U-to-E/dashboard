package routes

import (
	"github.com/U-to-E/dashboard/controller"
	"github.com/U-to-E/dashboard/middleware"
	"github.com/gofiber/fiber/v3"
)

func SetupStudentRoutes(app *fiber.App) {

	//GET
	app.Get("/", controller.RenderLogin)
	app.Get("/adminpanel", controller.RenderAdmin, middleware.Protected)
	app.Get("/dashboard", controller.RenderDashboard, middleware.Protected)

	//POST
	app.Post("/signup", controller.Register)
	app.Post("/login", controller.Handlelogin)
	app.Post("/logout", controller.Logout, middleware.Protected)
}
