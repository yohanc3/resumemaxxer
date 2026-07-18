package observer

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/lib/pq"
)

type Observer struct {
	Interval time.Duration
	URL      string
}

type Listing struct {
	Title       string   	`json:"title"`
	URL         string  	`json:"url"`
	CompanyName string  	`json:"company_name"`
	CompanyURL  string   	`json:"company_url"`
	Source      string   	`json:"source"`
	Category    string   	`json:"category"`
	ID          string   	`json:"id"`
	Active      bool        `json:"active"`
	Terms       []string 	`json:"terms"`
	Locations   []string 	`json:"locations"`
	Sponsorship string   	`json:"sponsorship"`
	IsVisible   bool     	`json:"is_visible"`
	Degrees     []string 	`json:"degrees"`
	DateUpdated time.Time	`json:"date_updated"`		
	DatePosted	time.Time	`json:"date_posted"`
}

func (o *Observer) Watch(ctx context.Context) error {

	return nil
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

// Pushes listings into job queue and into listings tables 
func (o *Observer) ProcessListings(ctx context.Context, db *sql.DB, listings []*Listing) error {

	// Slightly adapted from - https://stackoverflow.com/a/48070387

	// Holds list of values placeholders ($1, $2, $3 ...). 15 values per row.
    valueStrings := make([]string, 0, len(listings))

	// Holds list of arguments to be inserted in placeholders. 15 values per row.
    valueArgs := make([]interface{}, 0, len(listings) * 15)

    i := 0
    for _, listing := range listings {

		// For each listing, store a new set of value placeholders
		// It allocates placeholder numbers based on the listing number
		// i.e., first iteration adds placeholders 1-15, second iteration 16-30, etc,
		// so that the appended arguments match the correct listing
        valueStrings = append(valueStrings, fmt.Sprintf(`
		($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)
		`, i*15+1, i*15+2, i*15+3, i*15+4, i*15+5, i*15+6, i*15+7, i*15+8, i*15+9, i*15+10, i*15+11, i*15+12, i*15+13, i*15+14, i*15+15))

		// For each listing, store its 15 args. Add them in the same order than expected
		// in the sql statement.
        valueArgs = append(valueArgs, listing.ID)
        valueArgs = append(valueArgs, listing.Source)
        valueArgs = append(valueArgs, listing.Category)
        valueArgs = append(valueArgs, listing.CompanyName)
        valueArgs = append(valueArgs, listing.Title)
        valueArgs = append(valueArgs, listing.Active)
        valueArgs = append(valueArgs, pq.Array(listing.Terms))
        valueArgs = append(valueArgs, listing.DateUpdated)
        valueArgs = append(valueArgs, listing.DatePosted)
        valueArgs = append(valueArgs, listing.URL)
        valueArgs = append(valueArgs, pq.Array(listing.Locations))
        valueArgs = append(valueArgs, listing.CompanyURL)
        valueArgs = append(valueArgs, listing.IsVisible)
        valueArgs = append(valueArgs, listing.Sponsorship)
        valueArgs = append(valueArgs, pq.Array(listing.Degrees))
        i++
    }
	
	// Upserts into job_postings table, and inserts into the resume generation queue
	// job postings that either were inserted, or successfully edited (after id conflict).
    stmt := fmt.Sprintf(`
		WITH upserted_resources AS ( 
			INSERT INTO job_postings (id, source, category, company_name, title,
			active, terms, date_updated, date_posted, url, locations, company_url,
			is_visible, sponsorship, degrees)
			VALUES %s

			ON CONFLICT (id) 
			DO UPDATE SET
				source = EXCLUDED.source,
				category = EXCLUDED.category,
				company_name = EXCLUDED.company_name,
				title = EXCLUDED.title,
				active = EXCLUDED.active,
				terms = EXCLUDED.terms,
				date_updated = EXCLUDED.date_updated,
				date_posted = EXCLUDED.date_posted,
				url = EXCLUDED.url,
				locations = EXCLUDED.locations,
				company_url = EXCLUDED.company_url,
				is_visible = EXCLUDED.is_visible,
				sponsorship = EXCLUDED.sponsorship,
				degrees = EXCLUDED.degrees
			
			WHERE job_postings.active IS NOT FALSE 
			OR job_postings.date_updated IS DISTINCT FROM EXCLUDED.date_updated
			RETURNING id, url, company_name 
		)

		INSERT INTO resume_generation_queue (
			job_posting_id, job_posting_url, job_posting_company_name 
		)
		SELECT id, url, company_name 
		FROM upserted_resources
		;
		`, strings.Join(valueStrings, ","))
	
	// Apply statement, and exclude the result.
    _, err := db.ExecContext(ctx, stmt, valueArgs...)

	if err != nil {
		// Error out for now. Should notify dev later.
		return fmt.Errorf("error when inserting batch of jobs. %w", err)
	}

	return nil

}
