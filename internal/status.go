package internal

import (
	"fmt"
	"log"
	"os"
	"time"
)

func CheckStatus() {
	numOfUpdates := 0
	savedJobBoards, err := ReadConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	savedJobBoards = filterOutNewFetches(savedJobBoards)

	crawlersToCommit := make([]Crawler, 0)

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
		numOfUpdates += compareJobBoards(cr, jb)

		crawlersToCommit = append(crawlersToCommit, cr)
	}

	// Update the config file with the new LastFetch value and the new job info if any changes
	CommitToRepo(crawlersToCommit)

	if numOfUpdates == 0 {
		fmt.Println("No changes since last fetch")
	} else {
		fmt.Printf("\nTotal number of updates: %d\n", numOfUpdates)
	}
}

// return 0 if no changes, 1 if changes
func compareJobBoards(crawler Crawler, savedJobBoard JobBoard) int {

	if crawler.GetHash() == savedJobBoard.Hash {
		return 0 // No changes
	}

	fetchedjb, _ := crawler.GetJobLinks()

	fmt.Printf("\nThe job board %s has changed since last fetch\n", crawler.GetBoardName())
	fmt.Print("Possible Changes: ")
	if len(fetchedjb) > savedJobBoard.JobsCount {
		// If new jobs added, print message in green
		fmt.Printf("\033[32m new jobs added to: %s\033[0m\n\n", savedJobBoard.Url)
	} else if len(fetchedjb) < savedJobBoard.JobsCount {
		// If jobs removed, print message in red
		fmt.Printf("\033[31m jobs removed from: %s\033[0m\n\n", savedJobBoard.Url)
	} else {
		fmt.Printf("\033[33m jobs might have been updated in: %s\033[0m\n\n", savedJobBoard.Url)
	}
	return 1
}

// remove job boards that has LastFetch value less than 24 hours
func filterOutNewFetches(jobBoards []JobBoard) []JobBoard {
	var filteredJobBoards []JobBoard
	for _, jb := range jobBoards {
		diff := time.Since(jb.LastFetch)
		if diff.Hours() > 24 {
			filteredJobBoards = append(filteredJobBoards, jb)
		}
	}
	return filteredJobBoards
}
