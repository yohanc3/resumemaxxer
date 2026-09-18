package jobpostingdescriptionfetcher

import (
	"context"
	"os"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

type Scraper struct {
	browser *rod.Browser
}

func NewScraper() *Scraper {
	return &Scraper{}
}

func (s *Scraper) GetPageText(ctx context.Context, url string) (string, error){

	controlURL := launcher.New().
		Bin(os.Getenv("CHROME_BIN")).
		Headless(true).
		NoSandbox(true).
		MustLaunch()
	
	browser := rod.New().ControlURL(controlURL).MustConnect()
	defer browser.MustClose()

	page := browser.MustPage(url)

	page.MustWaitLoad().MustWaitIdle()

	page.MustWaitDOMStable()

	pageText := ""

	paragraphs := page.MustElements("p")
	for _, p := range paragraphs {
		text := p.MustText()
		pageText += "\n" + text 
	}

	return pageText, nil
}
