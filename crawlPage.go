package main

import (
    "net/url"
    "fmt"
)

func crawlPage(rawBaseURL, rawCurrentURL string, pages map[string]int) (map[string]int, error) {
    fmt.Printf("Crawling: %s\n", rawCurrentURL)

    parsedBaseURL, err := url.Parse(rawBaseURL)
    if err != nil {
        return pages, fmt.Errorf("error parsing base URL: %v", err)
	}
    
	if err != nil {
		return pages, fmt.Errorf("error parsing base URL: %v", err)
	}
    parsedCurrentURL, err := url.Parse(rawCurrentURL)
    if err != nil {
        return pages, fmt.Errorf("error parsing current URL: %v", err)
	}
	

    if parsedBaseURL.Host != parsedCurrentURL.Host {
        fmt.Printf("Skipping external URL: %s (Base host: %s, Current host: %s)\n", 
           rawCurrentURL, parsedBaseURL.Host, parsedCurrentURL.Host)
        return pages, nil
	}
	fmt.Println(rawCurrentURL)
    normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		fmt.Println("error normalizing the url")
		return pages, err
	}
	fmt.Println(normalizedCurrentURL)
    if _, exists := pages[normalizedCurrentURL]; exists {
        fmt.Printf("Already crawled: %s\n", normalizedCurrentURL)
        pages[normalizedCurrentURL]++
		return pages, nil
	}

    pages[normalizedCurrentURL] = 1
    htmlContent, err := getHTML(normalizedCurrentURL)
    if err != nil {
        fmt.Println("err getting the html content", err)
        return pages, err
	}
	fmt.Printf("Got HTML content for: %s\n", normalizedCurrentURL) 

    urls, err := getURLsFromHTML(htmlContent, normalizedCurrentURL)
    if err != nil {
        fmt.Println("err getting the urls", err)
        return pages, err
	}

    fmt.Printf("Found %d URLs on page: %s\n", len(urls), normalizedCurrentURL)

    for _, url := range urls {
        pages, err = crawlPage(rawBaseURL, url, pages)
        if err != nil {
                fmt.Printf("Error crawling %s: %v\n", url, err)
                // Continue crawling other URLs even if one fails
		}
	}

	return pages, nil
}