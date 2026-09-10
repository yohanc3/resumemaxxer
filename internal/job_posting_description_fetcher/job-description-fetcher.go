package jobpostingdescriptionfetcher

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

 	observer "github.com/yohanc3/resumemaxxer/internal/observer"
	"github.com/lib/pq"
)

type JobDescriptionFetcher struct {
	db  *sql.DB
	sem chan struct{}
	wg  sync.WaitGroup
}

type MissingJobDescriptionURL struct {
	url            string
	job_posting_id string 
}

func (j *JobDescriptionFetcher) Start(ctx context.Context) error {

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(time.Second * 5)

	for {
		select {
		case <-ticker.C:
			if err := j.processMissingJobDescriptions(ctx); err != nil {
				slog.ErrorContext(ctx, "error when processing missing job descriptions", slog.String("error", err.Error()))
			}

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (j *JobDescriptionFetcher) processMissingJobDescriptions(ctx context.Context) error {

	return nil
}

func (j *JobDescriptionFetcher) getMissingJobDescriptionsURLs(ctx context.Context) ([]*string, error) {

	return nil, nil
}

// Pushes listings into job queue and into listings table
func (j *JobDescriptionFetcher) pushResumeCreationJobs(ctx context.Context, listings []*.Listing) error {

	// Slightly adapted from - https://stackoverflow.com/a/48070387

	// Holds list of values placeholders ($1, $2, $3 ...). 15 values per row.
	valueStrings := make([]string, 0, len(listings))

	// Holds list of arguments to be inserted in placeholders. 15 values per row.
	valueArgs := make([]interface{}, 0, len(listings)*15)

	i := 0
	for _, listing := range listings[:3] {

		// For each listing, store a new set of value placeholders
		// It allocates placeholder numbers based on the listing number
		// i.e., first iteration adds placeholders 1-15, second iteration 16-30, etc,
		// so that the appended arguments match the correct listing
		valueStrings = append(valueStrings, fmt.Sprintf(`
		($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d,d $%d, $%d, $%d, $%d, $%d, $%d)
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
		valueArgs = append(valueArgs, time.Unix(listing.DateUpdated, 0).UTC())
		valueArgs = append(valueArgs, time.Unix(listing.DatePosted, 0).UTC())
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
			job_posting_id, job_posting_url, company_name, user_id 
		)
		SELECT id, url, company_name, 'testid' 
		FROM upserted_resources
		ON CONFLICT (user_id, job_posting_id)
		DO NOTHING
		;
		`, strings.Join(valueStrings, ","))

	// Apply statement, and exclude the result.
	_, err := j.db.ExecContext(ctx, stmt, valueArgs...)

	if err != nil {
		// Error out for now. Should notify dev later.
		return fmt.Errorf("error when inserting batch of jobs. %w", err)
	}

	return nil

}
