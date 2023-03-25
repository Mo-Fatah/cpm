package cpm

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

func Add(args []string) {
	for _, url := range args {
		spinner := NewSpinner()
		spinner.Start(fmt.Sprintf("fetching info from %s  ", url))

		crawler := NewCrawler(url)
		_, err := crawler.Crawl()
		if err != nil {
			spinner.Failure()
			fmt.Fprintf(os.Stderr, "error fetching job postings from %s\n%s", url, err)
			continue
		}

		spinner.Success()

		spinner.Start("commiting to the config repo ..  ")
		if err = commitToRepo(crawler); err != nil {
			spinner.Failure()
			fmt.Fprintf(os.Stderr, "error writing to the config repo\n%s\n", err)
			continue
		}
		spinner.Success()
	}
}

func commitToRepo(crawler Crawler) error {
	currUser, err := user.Current()
	if err != nil {
		return err
	}
	cpmConfigPath := filepath.Join(currUser.HomeDir, ".cpm", "cpmConfig.json")
	jobLinks, err := crawler.GetJobLinks()
	if err != nil {
		return err
	}
	currJobBoards, err := ReadConfigFile()
	if err != nil {
		return err
	}

	newJobBoard := JobBoard{
		Name:      crawler.GetBoardName(),
		Url:       crawler.Url,
		JobLinks:  jobLinks,
		JobsCount: len(jobLinks),
		Hash:      crawler.GetHash(),
		LastFetch: time.Now(),
	}

	newJobBoards := append(currJobBoards, newJobBoard)

	// open file in append mode
	file, err := os.OpenFile(cpmConfigPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// encode job board to json then write it to the config file
	b, err := json.MarshalIndent(newJobBoards, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(cpmConfigPath, b, 0644); err != nil {
		return err
	}
	//_, err = file.Write(b)
	//if err != nil {
	//	return err
	//}

	//// append a comma to the end of the file to separate the new job board from the previous one
	//_, err = file.Write([]byte(","))
	//if err != nil {
	//	return err
	//}
	return nil
}
