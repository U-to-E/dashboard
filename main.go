package main

import (
	"encoding/gob"
	"log"
	"time"

	"github.com/U-to-E/dashboard/controller"
	"github.com/U-to-E/dashboard/database"
	"github.com/U-to-E/dashboard/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/template/html/v2"
)

var store *session.Store

func main() {
	gob.Register(time.Time{})
	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	store = session.New(session.Config{
		CookieHTTPOnly: true,
		CookieSecure:   true,
		CookieSameSite: "Lax",
	})
	controller.Session(store)

	app.Static("/", "./materials")
	app.Static("/asserts", "./static")
	app.Static("/cert", "./certificates")
	// app.Get("/metrics", monitor.New())
	// app.Use(csrf.New())
	// app.Use(csrf.New(csrf.Config{
	// 	KeyLookup:      "header:X-Csrf-Token",
	// 	CookieName:     "csrf_",
	// 	CookieSameSite: "Lax",
	// 	CookieHTTPOnly: true,
	// 	CookieSecure:   true,
	// 	Expiration:     1 * time.Hour,
	// 	KeyGenerator:   utils.UUIDv4,
	// }))
	app.Use(cors.New(), logger.New())
	database.Connect()
	routes.SetupStudentRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
