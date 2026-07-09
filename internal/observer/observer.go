package observer

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type Observer struct {
	Interval time.Time
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
	Degrees     string   `json:"degrees"`
}

// Fetches all listings and returns them
func (o *Observer) fetchListings(ctx context.Context) ([]Listing, error) {
	url := o.URL

	res, err := http.Get(url)

	if err != nil {
		slog.Error("error: "+err.Error(), slog.String("url fetched", url))
		return nil, err
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
func (o *Observer) pushJobsFromListings(listings []Listing) error {
	return nil
}

// Pushes listings into database (PostgreSQL)
func (o *Observer) saveListings(listings []Listing) error {
	return nil
}
