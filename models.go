package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// --- data structs ---
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

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Manufacturer struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func apiBase() string {
	if v := os.Getenv("API_BASE_URL"); v != "" {
		return v
	}
	return "http://localhost:3000"
}

func fetchModels() ([]CarModel, error) {
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(apiBase() + "/api/models")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var models []CarModel
	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
		return nil, err
	}
	return models, nil
}

func fetchCategories() ([]Category, error) {
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(apiBase() + "/api/categories")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cats []Category
	if err := json.NewDecoder(resp.Body).Decode(&cats); err != nil {
		return nil, err
	}
	return cats, nil
}

func fetchManufacturers() ([]Manufacturer, error) {
	client := http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(apiBase() + "/api/manufacturers")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var mans []Manufacturer
	if err := json.NewDecoder(resp.Body).Decode(&mans); err != nil {
		return nil, err
	}
	return mans, nil
}

// helper: fetch single model by id (keeps existing behavior)
func fetchModelByID(id int) (CarModel, error) {
	var cm CarModel
	client := http.Client{Timeout: 3 * time.Second}
	url := fmt.Sprintf("%s/api/models/%d", apiBase(), id)
	resp, err := client.Get(url)
	if err == nil && resp != nil {
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			if err := json.NewDecoder(resp.Body).Decode(&cm); err == nil {
				return cm, nil
			}
		}
	}
	// fallback
	models, err := fetchModels()
	if err != nil {
		return cm, err
	}
	for _, m := range models {
		if m.ID == id {
			return m, nil
		}
	}
	return cm, fmt.Errorf("model id %d not found", id)
}
