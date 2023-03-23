package cpm

import (
	"fmt"
	"os"
	"time"
)

func Add(args []string) {
	for _, url := range args {
		spinner := NewSpinner()
		spinner.Start(fmt.Sprintf("fetching info from %s", url))

		crawler := NewCrawler(url)
		_, err := crawler.Crawl()
		if err != nil {
			spinner.Failure()
			fmt.Fprintf(os.Stderr, "error fetching job postings from %s", url)
			continue
		}

		spinner.Success()
	}
}

type JobBoard struct {
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	LastFetch time.Time `json:"lastFetched"`
	JobsCount int       `json:"jobsCount"`
}
