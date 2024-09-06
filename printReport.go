package main

import (
	"fmt"
	"sort"
)

type Elem struct{
	Link string
	Count int
}

func printReport(pages map[string]int, baseURL string) {
	fmt.Println("REPORT for", baseURL)
	slice := make([]Elem, 0)

	for link, count := range pages {
		slice = append(slice, struct {
			Link string
			Count int
		}{link, count})
	}

	sort.Slice(slice, func(i, j int) bool {
		return slice[i].Count > slice[j].Count
	})

	for _, k := range slice {
		message := fmt.Sprintf("Found %d internal links to %v", k.Count, k.Link)
		fmt.Println(message)
	}
}