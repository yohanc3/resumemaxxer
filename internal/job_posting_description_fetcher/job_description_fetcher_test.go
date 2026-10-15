package jobpostingdescriptionfetcher

import (
	"context"
	"os"
	"strings"
	"testing"
)

func TestFetchURL(t *testing.T){

	scraper, err := NewScraper()
	if err != nil {
		t.Fatalf("Scraper return error: %v", err.Error())
	}

	t.Run("fetching complex html job description", func(t *testing.T) {
		autozoneURL := "https://egud.fa.us2.oraclecloud.com/hcmUI/CandidateExperience/en/sites/CX_1/job/155451?utm_source=Simplify&ref=Simplify"

		got, err := scraper.GetPageText(context.Background(), autozoneURL)
		if err != nil {
			t.Fatalf("Scraper error when getting page text: %v", err.Error())
		}

		want := []string{
			"Must be able to reside in the Memphis  area for the duration of the internship",
			"Python, JavaScript, Java and C++",
		}

		for _, w := range want {
			if !strings.Contains(Normalize(got), Normalize(w)){
				t.Fatalf("Fetched text doesn't contain: %v", w)
			}
		}
	
		os.WriteFile("jobdescription.txt", []byte(Normalize(got)), 0644)

	})

}

func Normalize(s string) string {
	return strings.ReplaceAll(strings.Join(strings.Fields(s), " "), "\n\n", "\n")
}
