package main

import (
	"net/http"
	"fmt"
	"io"
)

func FetchConfig(repoOwner, repoName, filePath string) (string, error) {
	url := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/master/%s", repoOwner, repoName, filePath)
	
	// Making the HTTP request
	resp, err := http.Get(url)
		// Checking for errors and status code
	if resp.StatusCode != http.StatusOK {
    	return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	if err != nil {
		return "", fmt.Errorf("failed to fetch config: %w", err)
	}
	defer resp.Body.Close()

	// Reading the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read config: %w", err)
	}

	return string(body), nil
}





