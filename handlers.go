// route handlers
package main

import (
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func registerRoutes() {
	app.Get("/", homeHandler)
	app.Post("/action", addComparisonHandler)
	app.Get("/comparisons", comparisonsHandler)
	app.Post("/comparisons/remove", removeComparisonHandler)
	app.Post("/favorite/toggle", favoriteToggleHandler)
}

// --- filtering helper ---
func filterModels(models []CarModel, catID, manID int) []CarModel {
	out := make([]CarModel, 0, len(models))
	for _, m := range models {
		if catID != 0 && m.CategoryID != catID {
			continue
		}
		if manID != 0 && m.ManufacturerID != manID {
			continue
		}
		out = append(out, m)
	}
	return out
}

// GET /
func homeHandler(c *fiber.Ctx) error {
	s, _ := store.Get(c)
	catQ := c.Query("category")
	manQ := c.Query("manufacturer")
	catID := 0
	manID := 0
	if catQ != "" {
		if v, err := strconv.Atoi(catQ); err == nil {
			catID = v
		}
	}
	if manQ != "" {
		if v, err := strconv.Atoi(manQ); err == nil {
			manID = v
		}
	}

	models, err := fetchModels()
	if err != nil {
		log.Println("fetchModels error:", err)
		return c.Status(500).SendString("Failed to load models")
	}
	categories, _ := fetchCategories()
	manufacturers, _ := fetchManufacturers()
	filtered := filterModels(models, catID, manID)

	// favorites -> Recommendations
	favIDs, _ := readFavoritesFromSession(s)
	favSet := map[int]bool{}
	recs := make([]CarModel, 0, len(favIDs))
	if len(favIDs) > 0 {
		lookup := map[int]CarModel{}
		for _, m := range models {
			lookup[m.ID] = m
		}
		for _, id := range favIDs {
			if mm, ok := lookup[id]; ok {
				recs = append(recs, mm)
				favSet[id] = true
			}
		}
	}

	return c.Render("index", fiber.Map{
		"Title":           "Home",
		"cars":            filtered,
		"Categories":      categories,
		"Manufacturers":   manufacturers,
		"Recommendations": recs,
		"FavoriteSet":     favSet,
	}, "layout")
}

// POST /action
func addComparisonHandler(c *fiber.Ctx) error {
	s, err := store.Get(c)
	if err != nil {
		log.Println("session get error:", err)
		return c.Redirect("/")
	}
	idStr := c.FormValue("compare")
	if idStr == "" {
		return c.Redirect("/")
	}
	id, _ := strconv.Atoi(idStr)
	m, err := fetchModelByID(id)
	if err != nil {
		log.Println("fetchModelByID error:", err)
		return c.Redirect("/")
	}
	added, err := addComparisonToSession(s, m, 4)
	if err != nil {
		log.Println("addComparisonToSession error:", err)
	} else if !added {
		log.Println("already in comparisons:", id)
	}
	return c.Redirect("/")
}

// GET /comparisons
func comparisonsHandler(c *fiber.Ctx) error {
	s, err := store.Get(c)
	if err != nil {
		return c.Status(500).SendString("session error")
	}
	models, err := readComparisonsFromSession(s)
	if err != nil {
		log.Println("readComparisonsFromSession error:", err)
		models = nil
	}
	return c.Render("comparisons", fiber.Map{"Title": "Car Comparison", "Models": models}, "layout")
}

// POST /comparisons/remove
func removeComparisonHandler(c *fiber.Ctx) error {
	s, err := store.Get(c)
	if err != nil {
		return c.Redirect("/comparisons")
	}
	idStr := c.FormValue("id")
	id, _ := strconv.Atoi(idStr)
	_, _ = removeComparisonFromSession(s, id)
	return c.Redirect("/comparisons")
}

// POST /favorite/toggle
func favoriteToggleHandler(c *fiber.Ctx) error {
	s, err := store.Get(c)
	if err != nil {
		if c.Get("X-Requested-With") == "XMLHttpRequest" {
			return c.Status(500).JSON(fiber.Map{"error": "session error"})
		}
		return c.Redirect("/")
	}
	idStr := c.FormValue("id")
	if idStr == "" {
		if c.Get("X-Requested-With") == "XMLHttpRequest" {
			return c.Status(400).JSON(fiber.Map{"error": "missing id"})
		}
		return c.Redirect("/")
	}
	id, _ := strconv.Atoi(idStr)
	added, err := toggleFavoriteInSession(s, id)
	if err != nil {
		if c.Get("X-Requested-With") == "XMLHttpRequest" {
			return c.Status(500).JSON(fiber.Map{"error": "toggle error"})
		}
		return c.Redirect("/")
	}
	if c.Get("X-Requested-With") == "XMLHttpRequest" {
		return c.JSON(fiber.Map{"added": added, "id": id})
	}
	ref := c.Get("Referer")
	if ref == "" {
		ref = "/"
	}
	return c.Redirect(ref)
}
