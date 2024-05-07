package routes

import (
	"github.com/U-to-E/dashboard/controller"
	"github.com/U-to-E/dashboard/middleware"
	"github.com/gofiber/fiber/v3"
)

func SetupStudentRoutes(app *fiber.App) {

	student := app.Group("/student")
	mentor := app.Group("/mentor")
	admin := app.Group("/admin")

	//GET
	app.Get("/", controller.RenderLogin)
	admin.Get("/panel", controller.RenderAdmin, middleware.ProtectedAdmin)
	admin.Get("/panel/collage/:id", controller.GetStudentList, middleware.ProtectedAdmin)
	student.Get("/dashboard", controller.RenderDashboard, middleware.Protected)
	mentor.Get("/dashboard", controller.RenderMentorDash, middleware.Protected)
	admin.Get("/panel/mentormapping", controller.BulkMentorMapping, middleware.ProtectedAdmin)

	//POST
	app.Post("/login", controller.Handlelogin)
	app.Post("/logout", controller.Logout, middleware.Protected)
	admin.Post("/panel/register/student", controller.AddStudent, middleware.ProtectedAdmin)
	admin.Post("/panel/register/mentor", controller.AddMentor, middleware.ProtectedAdmin)
	admin.Post("/panel/register/singlestudent", controller.AddSingleStudent, middleware.ProtectedAdmin)
	admin.Post("/panel/register/singlementor", controller.AddSingleMentor, middleware.ProtectedAdmin)
	admin.Post("/panel/password/gen", controller.PasswordGenrator, middleware.ProtectedAdmin)
	admin.Post("/panel/collage/id", controller.PostCID, middleware.ProtectedAdmin)
	admin.Get("/panel/mentormapping", controller.UpdateMentorMapping, middleware.ProtectedAdmin)

}
