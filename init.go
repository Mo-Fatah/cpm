package cpm

import (
	"fmt"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

func Initialize() {
	cpmDirPath := createCpmDir()
	createCpmConfig(cpmDirPath)
}

func createCpmDir() string {
	spinner := NewSpinner()
	spinner.Start("Initializing .cpm directory in your home page ..  ")

	currUser, err := user.Current()
	if err != nil {
		spinner.Failure()
		err = fmt.Errorf("error fetching the current user: %v", err)
		log.Fatal(err)
	}

	cpmDir := filepath.Join(currUser.HomeDir, ".cpm")
	err = os.Mkdir(cpmDir, 0755)
	if err != nil {
		spinner.Failure()
		log.Fatal(err)
	}
	spinner.Success()
	return cpmDir
}

func createCpmConfig(cpmDir string) {
	spinner := NewSpinner()
	spinner.Start("Creating cpmConfig.json ..  ")

	cpmConfigPath := filepath.Join(cpmDir, "cpmConfig.json")
	_, err := os.Create(cpmConfigPath)
	if err != nil {
		spinner.Failure()
		err = fmt.Errorf("error creating cpmConfig.json : %v", err)
		log.Fatal(err)
	}
	spinner.Success()
}
