package main

import (
	"log"

	"github.com/U-to-E/dashboard/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/template/html/v2"
)

func main() {
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	routes.SetupStudentRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
