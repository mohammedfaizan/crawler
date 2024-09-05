package main

import (
	"net/url"
	"fmt"
	"strings"
	"golang.org/x/net/html"
	"io"
	"net/http"
)

func normalizeURL(inputURL string) (string, error) {

	output, err := url.Parse(inputURL)
	if err != nil {
		return "", fmt.Errorf("couldn't parse URL")
	}

	fullPath := fmt.Sprintf("%v://%v%v",output.Scheme, output.Hostname(), output.Path)

	fullPath = strings.ToLower(fullPath)
	fullPath = strings.TrimSuffix(fullPath, "/")
	return fullPath, nil
}

func getURLsFromHTML(htmlBody, rawBaseURL string) ([]string, error) {

	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return make([]string, 0), err
	}

	htmlReader := strings.NewReader(htmlBody)
	NodeTree, err := html.Parse(htmlReader)
	if err != nil {
		return make([]string, 0), err
	}

	var urls []string
	var traverseNode func(*html.Node)

	traverseNode = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
            for _, attr := range n.Attr {
                if attr.Key == "href" {
                    parsedURL, err := url.Parse(attr.Val)
                    if err != nil {
                        continue
                    }
					
                    absoluteURL := baseURL.ResolveReference(parsedURL)
					

                    
					urls = append(urls, absoluteURL.String())
					
                    break
                }
            }
	
		}

		for c:= n.FirstChild; c != nil; c = c.NextSibling {
			traverseNode(c)
		}

	}


	traverseNode(NodeTree)

	return urls, nil
}


func getHTML(rawURL string) (string, error) {
	resp, err := http.Get(rawURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return "", fmt.Errorf("error code")
	}
	fmt.Println(resp.Header.Get("Content-Type"))
	if !strings.Contains(resp.Header.Get("Content-Type"), "html") {
		return "", fmt.Errorf("not html")
	} 


	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading the body")
		return "", err
	}

	htmlContent := string(body)

	return htmlContent, nil
}

