package main

import (
    "net/url"
    "fmt"
	"sync"
)

type config struct {
	pages			   map[string]int
	baseURL 		   *url.URL
	mu				   *sync.Mutex
	concurrencyControl chan struct{}
	wg				   *sync.WaitGroup
	maxPages		   int
}

func (cfg *config) crawlPage(rawCurrentURL string) error {
    defer cfg.wg.Done()

	cfg.mu.Lock()
	if len(cfg.pages) >= cfg.maxPages {
		cfg.mu.Unlock()
		return fmt.Errorf("max pages reached")
	}
	cfg.mu.Unlock()

	select {
    case cfg.concurrencyControl <- struct{}{}:
        defer func() { <-cfg.concurrencyControl }()
    default:
        return nil
    }

	


    parsedBaseURL, err := url.Parse(cfg.baseURL.String())
    if err != nil {
        return fmt.Errorf("error parsing base URL: %v", err)
	}
    
	if err != nil {
		return fmt.Errorf("error parsing base URL: %v", err)
	}
    parsedCurrentURL, err := url.Parse(rawCurrentURL)
    if err != nil {
        return fmt.Errorf("error parsing current URL: %v", err)
	}
	

    if parsedBaseURL.Host != parsedCurrentURL.Host {
        fmt.Printf("Skipping external URL: %s (Base host: %s, Current host: %s)\n", 
           rawCurrentURL, parsedBaseURL.Host, parsedCurrentURL.Host)
        return nil
	}

    normalizedCurrentURL, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return err
	}
	fmt.Println(normalizedCurrentURL)
    if _, exists := cfg.pages[normalizedCurrentURL]; exists {
        fmt.Printf("Already crawled: %s\n", normalizedCurrentURL)
        cfg.pages[normalizedCurrentURL]++
		return nil
	}

    cfg.pages[normalizedCurrentURL] = 1
    htmlContent, err := getHTML(normalizedCurrentURL)
    if err != nil {
        
        return err
	}
	

    urls, err := getURLsFromHTML(htmlContent, normalizedCurrentURL)
    if err != nil {
        fmt.Println("err getting the urls", err)
        return err
	}


    for _, url := range urls {
		
		cfg.mu.Lock()
        if len(cfg.pages) >= cfg.maxPages {
            cfg.mu.Unlock()
            return nil
        }
        cfg.mu.Unlock()


        if cfg.addPageVisit(url) {
			cfg.wg.Add(1)
			cfg.crawlPage(url)
		}
	}

	return nil
}

func (cfg *config) addPageVisit(normalizedCurrentURL string) (isFirst bool) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()

	if len(cfg.pages) >= cfg.maxPages {
        return false
    }

	if _, exists := cfg.pages[normalizedCurrentURL]; exists {
		return false
	}
	return true
}