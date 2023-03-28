package internal

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

func Add(args []string) {
	crawlersToCommit := make([]Crawler, 0)
	spinner := NewSpinner()
	for _, url := range args {
		spinner.Start(fmt.Sprintf("fetching info from %s  ", url))

		crawler := NewCrawler(url)
		_, err := crawler.Crawl()
		if err != nil {
			spinner.Failure()
			fmt.Fprintf(os.Stderr, "error fetching job postings from %s\n%s\n", url, err)
			continue
		}
		crawlersToCommit = append(crawlersToCommit, crawler)
		spinner.Success()
	}

	spinner.Start("writing to config file  ")
	if err := CommitToRepo(crawlersToCommit); err != nil {
		spinner.Failure()
		fmt.Fprintf(os.Stderr, "error writing to config file: %s\n", err)
		return
	}
	spinner.Success()

}

// Takes a list of crawlers, write the new job boards to the config file
// If the job board already exists, update the LastFetch value and the job links
func CommitToRepo(crawlers []Crawler) error {
	newJobBoards := make([]JobBoard, 0)
	for _, crawler := range crawlers {
		jobLinks, _ := crawler.GetJobLinks()
		newJobBoard := JobBoard{
			Name:      crawler.GetBoardName(),
			Url:       crawler.Url,
			JobsCount: len(jobLinks),
			Hash:      crawler.GetHash(),
			JobLinks:  jobLinks,
			LastFetch: time.Now(),
		}
		newJobBoards = append(newJobBoards, newJobBoard)
	}
	oldJobBoards, err := ReadConfigFile()
	if err != nil {
		return err
	}

	newJobBoards = mergeJobBoards(oldJobBoards, newJobBoards)

	currUser, err := user.Current()
	if err != nil {
		return err
	}
	cpmConfigPath := filepath.Join(currUser.HomeDir, ".cpm", "cpmConfig.json")
	data, err := json.MarshalIndent(newJobBoards, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(cpmConfigPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func mergeJobBoards(oldJobBoards []JobBoard, newJobBoards []JobBoard) []JobBoard {
	JobBoardMap := make(map[string]JobBoard)
	for _, jobBoard := range oldJobBoards {
		JobBoardMap[jobBoard.Url] = jobBoard
	}
	for _, jobBoard := range newJobBoards {
		JobBoardMap[jobBoard.Url] = jobBoard
	}

	mergedJobBoards := make([]JobBoard, 0)
	for _, jobBoard := range JobBoardMap {
		mergedJobBoards = append(mergedJobBoards, jobBoard)
	}

	return mergedJobBoards
}
