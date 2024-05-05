package routes

import (
	"github.com/U-to-E/dashboard/controller"
	"github.com/U-to-E/dashboard/middleware"
	"github.com/gofiber/fiber/v3"
)

func SetupStudentRoutes(app *fiber.App) {

	student := app.Group("/student")
	// mentor := app.Group("/mentor")
	admin := app.Group("/admin")

	//GET
	app.Get("/", controller.RenderLogin)
	admin.Get("/panel", controller.RenderAdmin, middleware.ProtectedAdmin)
	student.Get("/dashboard", controller.RenderDashboard, middleware.Protected)

	//POST
	app.Post("/login", controller.Handlelogin)
	app.Post("/logout", controller.Logout, middleware.Protected)
	admin.Post("/panel/register/student", controller.AddStudent, middleware.ProtectedAdmin)
	admin.Post("/panel/register/mentor", controller.AddMentor, middleware.ProtectedAdmin)
	admin.Post("/panel/password/gen", controller.PasswordGenrator, middleware.ProtectedAdmin)
}
