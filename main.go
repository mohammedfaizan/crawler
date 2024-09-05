package main

import (
	"fmt"
	"os"

)

func main(){
	
	if len(os.Args) < 2 {
		fmt.Println("no website provided")
		os.Exit(1)
	} else if len(os.Args) > 2 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	} else {
		baseURL := os.Args[1]
		fmt.Println("starting crawl of:", baseURL)
		pages := make(map[string]int)
		updatedPages, err := crawlPage(baseURL, baseURL, pages)
		if err != nil {
			fmt.Printf("Error during crawl: %v\n", err)
		}

		fmt.Println("\nCrawl Results:")
		for url, count := range updatedPages {
			fmt.Printf("%s: %d\n", url, count)
		}
	}
	
	

}


