package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"time"
)

type Asset struct {
	BrowserDownloadURL string `json:"browser_download_url"`
	Name               string `json:"name"`
}

type Release struct {
	TagName   string    `json:"tag_name"`
	UpdatedAt time.Time `json:"updated_at"`
	Assets    []Asset   `json:"assets"`
}

func main() {
	// Make a GET request to the API
	resp, err := http.Get("https://api.github.com/repos/pallets/flask/releases")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Unmarshal the JSON response into a slice of Release structs
	var releases []Release
	err = json.Unmarshal(body, &releases)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return
	}

	// Sort the releases by updated_at in descending order
	sort.Slice(releases, func(i, j int) bool {
		return releases[i].UpdatedAt.After(releases[j].UpdatedAt)
	})

	// Print the latest release
	if len(releases) > 0 {
		latestRelease := releases[0]
		fmt.Printf("Latest release: %s (updated at %s)\n", latestRelease.TagName, latestRelease.UpdatedAt)

		// Find and print the download URL for flask-<version>.tar.gz
		for _, asset := range latestRelease.Assets {
			if strings.HasSuffix(asset.Name, ".tar.gz") {
				fmt.Printf(asset.BrowserDownloadURL)
				break
			}
		}
	} else {
		fmt.Println("No releases found")
	}
}
