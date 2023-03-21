package status

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

var urls string = "https://fly.io/jobs/"

func CheckStatus() {
	resp, err := http.Get(urls)
	if err != nil {
		fmt.Println("error happend: ", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Handle error
		fmt.Println("Error:", err)
		return
	}

	fmt.Println(string(body))
}
