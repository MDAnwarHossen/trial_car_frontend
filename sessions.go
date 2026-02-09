package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2/middleware/session"
)

// comparisons helpers (store []CarModel JSON in session)
func readComparisonsFromSession(s *session.Session) ([]CarModel, error) {
	raw := s.Get("comparisons")
	if raw == nil {
		return nil, nil
	}
	// try common forms
	if models, ok := raw.([]CarModel); ok {
		return models, nil
	}
	if ia, ok := raw.([]interface{}); ok {
		out := make([]CarModel, 0, len(ia))
		for _, v := range ia {
			b, err := json.Marshal(v)
			if err != nil {
				continue
			}
			var m CarModel
			if err := json.Unmarshal(b, &m); err == nil {
				out = append(out, m)
			}
		}
		return out, nil
	}
	switch v := raw.(type) {
	case []byte:
		var out []CarModel
		if err := json.Unmarshal(v, &out); err == nil {
			return out, nil
		}
	case string:
		var out []CarModel
		if err := json.Unmarshal([]byte(v), &out); err == nil {
			return out, nil
		}
	}
	return nil, errors.New("unsupported comparisons session type")
}

func saveComparisonsToSession(s *session.Session, models []CarModel) error {
	b, err := json.Marshal(models)
	if err != nil {
		return fmt.Errorf("json marshal error: %w", err)
	}
	s.Set("comparisons", string(b))
	if err := s.Save(); err != nil {
		return fmt.Errorf("session save error: %w", err)
	}
	return nil
}

func addComparisonToSession(s *session.Session, model CarModel, max int) (bool, error) {
	if model.ID <= 0 {
		return false, errors.New("invalid model id")
	}
	models, err := readComparisonsFromSession(s)
	if err != nil {
		log.Println("readComparisonsFromSession warning:", err)
		models = nil
	}
	for _, m := range models {
		if m.ID == model.ID {
			return false, nil
		}
	}
	models = append(models, model)
	if max > 0 && len(models) > max {
		models = models[len(models)-max:]
	}
	if err := saveComparisonsToSession(s, models); err != nil {
		return false, err
	}
	return true, nil
}

func removeComparisonFromSession(s *session.Session, id int) (bool, error) {
	models, err := readComparisonsFromSession(s)
	if err != nil {
		return false, err
	}
	if models == nil || len(models) == 0 {
		return false, nil
	}
	out := make([]CarModel, 0, len(models))
	removed := false
	for _, m := range models {
		if m.ID == id {
			removed = true
			continue
		}
		out = append(out, m)
	}
	if !removed {
		return false, nil
	}
	if err := saveComparisonsToSession(s, out); err != nil {
		return false, err
	}
	return true, nil
}

// favorites helpers (store []int)
func readFavoritesFromSession(s *session.Session) ([]int, error) {
	raw := s.Get("favorites")
	if raw == nil {
		return nil, nil
	}
	switch v := raw.(type) {
	case []byte:
		var out []int
		if err := json.Unmarshal(v, &out); err == nil {
			return out, nil
		}
	case string:
		var out []int
		if err := json.Unmarshal([]byte(v), &out); err == nil {
			return out, nil
		}
	case []int:
		return v, nil
	case []interface{}:
		out := make([]int, 0, len(v))
		for _, x := range v {
			switch n := x.(type) {
			case float64:
				out = append(out, int(n))
			case int:
				out = append(out, n)
			case string:
				if iv, err := strconv.Atoi(n); err == nil {
					out = append(out, iv)
				}
			}
		}
		return out, nil
	}
	return nil, errors.New("unsupported favorites session type")
}

func saveFavoritesToSession(s *session.Session, ids []int) error {
	b, err := json.Marshal(ids)
	if err != nil {
		return fmt.Errorf("json marshal favorites: %w", err)
	}
	s.Set("favorites", string(b))
	if err := s.Save(); err != nil {
		return fmt.Errorf("session save error: %w", err)
	}
	return nil
}

// toggleFavoriteInSession returns true if added, false if removed.
func toggleFavoriteInSession(s *session.Session, id int) (bool, error) {
	if id <= 0 {
		return false, errors.New("invalid id")
	}
	ids, err := readFavoritesFromSession(s)
	if err != nil {
		log.Println("readFavoritesFromSession warning:", err)
		ids = nil
	}
	found := false
	out := make([]int, 0, len(ids))
	for _, v := range ids {
		if v == id {
			found = true
			continue
		}
		out = append(out, v)
	}
	if found {
		if err := saveFavoritesToSession(s, out); err != nil {
			return false, err
		}
		return false, nil
	}
	out = append(out, id)
	if err := saveFavoritesToSession(s, out); err != nil {
		return false, err
	}
	return true, nil
}
