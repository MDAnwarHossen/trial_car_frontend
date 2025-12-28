package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

// Data Structures
type Specification struct {
	Engine       string `json:"engine"`
	Horsepower   int    `json:"horsepower"`
	Transmission string `json:"transmission"`
	Drivetrain   string `json:"drivetrain"`
}
type CarModel struct {
	ID             int           `json:"id"`
	Name           string        `json:"name"`
	ManufacturerID int           `json:"manufacturerId"`
	CategoryID     int           `json:"categoryId"`
	Year           int           `json:"year"`
	Specifications Specification `json:"specifications"`
	Image          string        `json:"image"`
}

func fetchModels() ([]CarModel, error) {
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	apiBase := os.Getenv("API_BASE_URL")
	if apiBase == "" {
		apiBase = "http://localhost:3000" // safe default
	}
	resp, err := client.Get(apiBase + "/api/models")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var models []CarModel
	err = json.NewDecoder(resp.Body).Decode(&models)
	return models, err
}

func main() {
	// Initialize HTML template engine
	engine := html.New("./templates", ".html")

	// Create Fiber app
	app := fiber.New(fiber.Config{
		Views: engine,
	})
	app.Static("/", "./public")

	// Home Page
	app.Get("/", func(c *fiber.Ctx) error {
		cars, err := fetchModels()
		if err != nil {
			return c.Status(500).SendString("Failed to load models")
		}

		return c.Render("index", fiber.Map{
			"Title": "Home",
			"cars":  cars,
		}, "layout")
	})

	// Form action
	app.Post("/action", func(c *fiber.Ctx) error {
		log.Println("Action triggered")
		return c.Redirect("/")
	})

	log.Println("Server running at http://localhost:8080")
	log.Fatal(app.Listen(":8080"))
}
