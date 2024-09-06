package main

import (
	"fmt"
	"os"
	"sync"
	"net/url"
	"strconv"
	"time"
)

func main(){
	
	if len(os.Args) != 4 {
		fmt.Println("Usage: ./crawler <URL> <maxConcurrency> <maxPages>")
		os.Exit(1)
	}


	baseURL := os.Args[1]
	maxConcurrency, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Errror: maxconcurrency must be an integer")
		os.Exit(1)
	}

	maxPages, err := strconv.Atoi(os.Args[3])
	if err != nil {
		fmt.Println("error: maxpages must be an integer")
		os.Exit(1)
	}

	parsedURL, err := url.Parse(baseURL)
    if err != nil {
        fmt.Println("Error parsing URL:", err)
        os.Exit(1)
    }

		
	fmt.Println("starting crawl of:", baseURL)


	cfg := &config {
		pages:				make(map[string]int),
		baseURL: 			parsedURL,
		mu:					&sync.Mutex{},
		concurrencyControl: make(chan struct{}, maxConcurrency),
		wg:					&sync.WaitGroup{},
		maxPages:			maxPages,
	}

	cfg.wg.Add(1)
	go cfg.crawlPage(baseURL)
	

	done := make(chan bool)
    go func() {
        cfg.wg.Wait()
        close(cfg.concurrencyControl)
        done <- true
    }()

    select {
    case <-done:
        fmt.Println("Crawling completed")
    case <-time.After(5 * time.Minute):
        fmt.Println("Crawling timed out after 5 minutes")
    }

		
	printReport(cfg.pages, cfg.baseURL.String())
	

}


