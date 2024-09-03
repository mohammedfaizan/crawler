package main

import (
	"net/url"
	"fmt"
	"strings"
)

func normalizeURL(inputURL string) (string, error) {

	output, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("couldn't parse URL")
	}

	fullPath := fmt.Sprintf("%v%v", output.Hostname(), output.Path)

	fullPath = strings.ToLower(fullPath)
	fullPath = strings.TrimSuffix(fullPath, "/")

	return fullPath, nil
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {
	return make([]string,0), nil
}