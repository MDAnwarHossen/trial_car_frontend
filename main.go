package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/template/html/v2"
)

var (
	store *session.Store
	app   *fiber.App
)

func main() {
	// template engine
	engine := html.New("./templates", ".html")
	// small helper for dict usage in templates
	engine.AddFunc("dict", func(values ...interface{}) map[string]interface{} {
		m := make(map[string]interface{})
		for i := 0; i < len(values); i += 2 {
			key := values[i].(string)
			m[key] = values[i+1]
		}
		return m
	})

	// session store
	store = session.New()

	// fiber app with templates
	app = fiber.New(fiber.Config{Views: engine})

	// static files
	app.Static("/", "./public")

	// register routes from handlers.go
	registerRoutes()

	log.Println("Server running at http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
