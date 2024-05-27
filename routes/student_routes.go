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
	app.Get("/forgot-password", controller.RenderForgotPass)
	app.Get("/reset/password", controller.RenderResetPass)
	admin.Get("/panel", controller.RenderAdmin, middleware.ProtectedAdmin)
	admin.Get("/panel/collage/:id", controller.GetStudentList, middleware.ProtectedAdmin)
	admin.Get("/panel/mentors", controller.GetMentorList, middleware.ProtectedAdmin)

	admin.Get("/panel/edit-student/:CID-:SID", controller.EditStudent, middleware.ProtectedAdmin)

	student.Get("/dashboard", controller.RenderDashboard, middleware.Protected)
	student.Get("/dashboard/quiz", controller.QuizPage, middleware.Protected)
	student.Get("/dashboard/marks", controller.RenderMarks, middleware.Protected)
	mentor.Get("/dashboard/:id", controller.GetStudentPage, middleware.MentorProtected)
	mentor.Get("/dashboard", controller.RenderMentorDash, middleware.MentorProtected)
	student.Get("/dashboard/changepass", controller.ChangePassPage, middleware.Protected)
	mentor.Get("/dashboard/chanagepass", controller.ChangePassPage, middleware.MentorProtected)

	//POST
	app.Post("/login", controller.Handlelogin)
	app.Post("/logout", controller.Logout, middleware.Protected)
	app.Post("/forgot-password", controller.ForgotPassword)
	app.Post("/reset/password", controller.ResetPassword)
	app.Post("/reset-password", controller.ChangePassword, middleware.Protected)
	admin.Post("/panel/register/student", controller.AddStudent, middleware.ProtectedAdmin)
	admin.Post("/panel/register/mentor", controller.AddMentor, middleware.ProtectedAdmin)
	admin.Post("/panel/register/singlestudent", controller.AddSingleStudent, middleware.ProtectedAdmin)
	admin.Post("/panel/register/singlementor", controller.AddSingleMentor, middleware.ProtectedAdmin)
	admin.Post("/panel/password/gen", controller.PasswordGenrator, middleware.ProtectedAdmin)
	admin.Post("/panel/collage/id", controller.PostCID, middleware.ProtectedAdmin)
	admin.Post("/panel/mentormapping", controller.MapMentorToCollage, middleware.ProtectedAdmin)
	admin.Post("/panel/update-student", controller.UpdateStudent, middleware.ProtectedAdmin)
	admin.Post("/panel/createcollage", controller.CreateCollageID, middleware.ProtectedAdmin)
	admin.Delete("/panel/delete-student/:CID-:SID", controller.DeleteStudent, middleware.ProtectedAdmin)
	admin.Delete("/panel/delete/mentor/:MID", controller.DeleteMentor, middleware.ProtectedAdmin)
	mentor.Post("/dashboard/material/add", controller.PostMaterial, middleware.MentorProtected)
	mentor.Post("/dashboard/material/delete", controller.DeleteMaterial, middleware.MentorProtected)
	mentor.Post("/dashboard/quiz/add", controller.CreateQuiz, middleware.MentorProtected)
	student.Post("/dashboard/submit/quiz", controller.SubmitQuiz, middleware.Protected)
	mentor.Post("/dashboard/level/set/:id", controller.SetLevel, middleware.MentorProtected)
	mentor.Post("/dashboard/cert/upload/:id", controller.UploadCert, middleware.MentorProtected)
	mentor.Post("/dashboard/cert/delete/:id", controller.DeleteCert, middleware.MentorProtected)

}
