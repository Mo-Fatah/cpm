package cpm

import (
	"fmt"
	"os"
)

// var url string =

var urls []string = []string{"https://www.hashicorp.com/careers/open-positions?department=Research+%26+Development&jobTypes=Software+Engineering&jobTypes=Infrastructure+Engineering",
	"https://fly.io/jobs/"}

func CheckStatus() {
	for _, url := range urls {
		cr := NewCrawler(url)
		if _, err := cr.Crawl(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		if err := cr.Parse(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		hashed := cr.GetHash()
		fmt.Println(hashed)
	}
}
