package jobpostingdescriptionfetcher

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type Scraper struct {
	browser *rod.Browser
}

func NewScraper() (*Scraper, error) {
	controlURL, err := launcher.New().
		Bin(os.Getenv("CHROME_BIN")).
		Headless(true).
		NoSandbox(true).
		Launch()
	if err != nil {
		return nil, fmt.Errorf("error when launching browser")
	}

	browser := rod.New().ControlURL(controlURL)
	if err := browser.Connect(); err != nil {
		return nil, fmt.Errorf("error when attempting to control browser: %w", err) 
	}

	return &Scraper{browser: browser}, nil
}

func (s *Scraper) GetPageText(ctx context.Context, url string) (string, error) {

	url, err := s.NormalizeURL(ctx, url)
	if err != nil {
		return "", err
	}

	// Create browser page at given url
	page, err := s.browser.Page(proto.TargetCreateTarget{URL: url})
	if err != nil {
		return "", fmt.Errorf("error when opening browser page: %w", err)
	}
	defer page.Close()

	// Create a shallow copy of page to control page stablization timeout
	stablePage := page.Context(ctx).Timeout(time.Second * 30)

	err = stablePage.WaitLoad()
	if err == nil {
		err = stablePage.WaitIdle(time.Second * 15)
	}
	if err == nil {
		err = stablePage.WaitDOMStable(time.Second*3, 0.05)
	}

	stablePage.CancelTimeout()
	if errors.Is(err, context.DeadlineExceeded) {
		return "", fmt.Errorf("stable page timed out: %w", err)
	} else if err != nil {
		return "", fmt.Errorf("stable page failed: %w", err)
	}

	// Create a shallow copy of stable page to control page scraping timeout
	scrapePage := page.Context(ctx).Timeout(time.Second * 10)

	pageText := ""
	paragraphs, err := scrapePage.Elements("h1, h2, h3, h4, p, li")

	if err == nil {
		for _, p := range paragraphs {

			var text string

			text, err = p.Text()
			if err != nil {
				break
			}

			pageText += "\n" + text
		}

	}

	scrapePage.CancelTimeout()
	if errors.Is(err, context.DeadlineExceeded) {
		return "", fmt.Errorf("scrape page timed out: %w", err)
	} else if err != nil {
		return "", fmt.Errorf("scrape page failed: %w", err)
	}

	return pageText, nil
}

func (s *Scraper) NormalizeURL(ctx context.Context, originalURL string) (string, error) {
	normalizedURL, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}

	if strings.Contains(originalURL, "ashbyhq") {
		normalizedURL.Path = strings.TrimSuffix(normalizedURL.Path, "/application")
	}

	return normalizedURL.String(), nil
}

func NormalizeText(s string) string {
	return strings.ReplaceAll(strings.Join(strings.Fields(s), " "), "\n\n", "\n")
}
