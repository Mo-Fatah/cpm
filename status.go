package cpm

import (
	"fmt"
	"log"
	"os"
	"time"
)

func CheckStatus() {
	savedJobBoards, err := ReadConfigFile()
	if err != nil {
		log.Fatal(err)
	}
	savedJobBoards = filterJobBoards(savedJobBoards)

	for _, jb := range savedJobBoards {
		cr := NewCrawler(jb.Url)
		if _, err := cr.Crawl(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}
		if err := cr.Parse(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		compareJobBoards(cr, jb)
	}
}

func compareJobBoards(crawler Crawler, savedJobBoard JobBoard) {

	if crawler.GetHash() == savedJobBoard.Hash {
		return // No changes
	}

	fetchedjb, _ := crawler.GetJobLinks()

	fmt.Printf("\nThe job board %s has changed since last fetch\n", crawler.GetBoardName())
	fmt.Print("Possible Changes: ")
	if len(fetchedjb) > savedJobBoard.JobsCount {
		// If new jobs added, print message in green
		fmt.Printf("\033[32m new jobs added to: %s\033[0m\n\n", savedJobBoard.Url)
	}
	if len(fetchedjb) < savedJobBoard.JobsCount {
		// If jobs removed, print message in red
		fmt.Printf("\033[31m jobs removed from: %s\033[0m\n\n", savedJobBoard.Url)
	}
}

// remove job boards that has LastFetch value less than 24 hours
func filterJobBoards(jobBoards []JobBoard) []JobBoard {
	var filteredJobBoards []JobBoard
	for _, jb := range jobBoards {
		diff := time.Since(jb.LastFetch)
		if diff.Hours() > 24 {
			filteredJobBoards = append(filteredJobBoards, jb)
		}
	}
	return filteredJobBoards
}
