package cpm

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

type Spinner struct {
	chars  []rune
	delay  time.Duration
	active bool
}

func NewSpinner() *Spinner {
	return &Spinner{
		chars:  []rune{'-', '\\', '|', '/'},
		delay:  100 * time.Millisecond,
		active: true,
	}
}

func (s *Spinner) Start(msg string) {
	fmt.Print(msg)
	s.active = true
	go func() {
		for s.active {
			for _, char := range s.chars {
				fmt.Printf("\b%c", char)
				time.Sleep(s.delay)
			}
		}
	}()
}

func (s *Spinner) Success() {
	s.active = false
	fmt.Print("\b")
	fmt.Print("\033[32m\u2714\033[0m\n")
}

func (s *Spinner) Failure() {
	s.active = false
	fmt.Print("\b")
	fmt.Print("\033[31m\u2717\033[0m\n")
}

func (s *Spinner) Stop() {
	s.active = false
	fmt.Print("\b")
}

type JobBoard struct {
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	JobsCount int       `json:"jobsCount"`
	Hash      string    `json:"hash"`
	JobLinks  []string  `json:"jobPosts"`
	LastFetch time.Time `json:"lastFetched"`
}

func ReadConfigFile() ([]JobBoard, error) {
	var jobBoards []JobBoard
	currUser, err := user.Current()
	if err != nil {
		return nil, err
	}
	cpmConfigPath := filepath.Join(currUser.HomeDir, ".cpm", "cpmConfig.json")
	data, err := os.ReadFile(cpmConfigPath)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return []JobBoard{}, nil
	}

	err = json.Unmarshal(data, &jobBoards)
	if err != nil {
		return nil, err
	}
	return jobBoards, nil
}
