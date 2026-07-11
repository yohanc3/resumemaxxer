package observer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type Observer struct {
	Interval time.Duration
	URL      string
}

type Listing struct {
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	CompanyName string   `json:"comapny_name"`
	Source      string   `json:"source"`
	Category    string   `json:"category"`
	ID          string   `json:"id"`
	Active      bool     `json:"active"`
	Terms       []string `json:"terms"`
	Locations   []string `json:"locations"`
	Sponsorship string   `json:"sponsorship"`
	IsVisible   bool     `json:"is_visible"`
	Degrees     []string `json:"degrees"`
}

// Fetches all listings and returns them
func (o *Observer) FetchListings(ctx context.Context) ([]Listing, error) {
	url := o.URL

	res, err := http.Get(url)

	if err != nil {
		return nil, fmt.Errorf("error when fetching listings with url: %v. error: %w", o.URL, err)
	}

	defer res.Body.Close()

	// Parse all listings and return them
	var listings []Listing

	if err := json.NewDecoder(res.Body).Decode(&listings); err != nil {

		slog.Log(ctx,
			slog.LevelError,
			"error when decoding listings.json into listings object",
			slog.String("error", err.Error()))
	}

	return listings, nil
}

// Pushes lisitngs into job queue (SQLite)
func (o *Observer) PushJobsFromListings(listings []Listing) error {
	return nil
}

// Pushes listings into database (PostgreSQL)
func (o *Observer) SaveListings(listings []Listing) error {
	return nil
}
