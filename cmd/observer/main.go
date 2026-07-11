package main 

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"flag"

	"github.com/yohanc3/resumemaxxer/internal/observer"
)

var (
	obs      *observer.Observer
)

func main() {
	
	fallbackURL := "https://raw.githubusercontent.com/SimplifyJobs/Summer2026-Internships/refs/heads/dev/.github/scripts/listings.json"

	url := flag.String("url", fallbackURL, "URL where listings will be fetched from")
	interval := flag.Int("interval", 3, "Fetch interval in seconds")

	flag.Parse()

	obs := &observer.Observer{Interval: time.Duration(*interval), URL: *url}

	for {

		fmt.Println("fetching...")

		res, err := obs.FetchListings(context.TODO())

		if err != nil {
			panic(err.Error())
		}
		
		content, err := json.Marshal(res)

		if err != nil {
			panic(err)
		}

		os.WriteFile("../../listings.json", content, 0644)
		fmt.Printf("succesfully saved %d jobs\n", len(res))

		fmt.Printf("sleeping %d seconds\n", int(*interval))
		time.Sleep(time.Second * time.Duration(*interval))

	}

}
