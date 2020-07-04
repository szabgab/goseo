package main

import (
	"container/list"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// TODO: collect which page(s) link to every given external page
// TODO: report which page(s) have broken links, make sure we report all of our pages, not just the first one with this link.
func checkExternalURLs(externalURLs *list.List, externalPages map[string]int) {
	log.Println("checkExternalURLs")
	for {
		if externalURLs.Len() == 0 {
			break
		}
		item := externalURLs.Front()
		extURL := item.Value.(string)
		_, ok := externalPages[extURL]
		if !ok {
			log.Printf("externalURL: %v", extURL)
			externalPages[extURL] = 0

			resp, err := http.Get(extURL)
			if err != nil {
				log.Println(err)
			} else {
				if resp.StatusCode != 200 {
					log.Printf("Status code is %v", resp.StatusCode)
				}
			}
		}
		externalPages[extURL]++

		externalURLs.Remove(item)
	}
}

func parseHTML(body io.ReadCloser, internalURLs *list.List, externalURLS *list.List) {
	tokenizer := html.NewTokenizer(body)
	for {
		tt := tokenizer.Next()
		if tt == html.ErrorToken {
			if tokenizer.Err() == io.EOF {
				break
			}
			log.Printf("Error: %v", tokenizer.Err())
			break // TODO: return error?
		}
		//fmt.Printf("Token: %v\n", tokenizer.Token())
		tag, hasAttr := tokenizer.TagName()
		//log.Printf("Tag: %v\n", string(tag))
		if !hasAttr {
			continue
		}

		isExternalLink := regexp.MustCompile(`https?://`)

		for {
			attrKey, attrValue, moreAttr := tokenizer.TagAttr()
			//log.Printf("Attr: %v\n", string(attrKey))
			//log.Printf("Attr: %v\n", string(attrValue))
			//log.Printf("Attr: %v\n", moreAttr)
			// TODO: What to do if the link goes to http://mysite while we are working on https://mysite - it probably should be noticed and reported as error
			// TODO: What to do with links to https://hu.code-maven.com/ when we are processing https://code-maven.com/ ? Should they be treated as external links?
			if string(tag) == "a" && string(attrKey) == "href" {
				href := string(attrValue)
				match := isExternalLink.MatchString(href)
				if match {
					//log.Printf("External: %v\n", href)
					externalURLS.PushBack(href)
				} else {
					//log.Printf("Internal: %v\n", href)
					internalURLs.PushBack(href)
				}
			}
			if !moreAttr {
				break
			}
		}
		// Process the current token.
	}
}

func processURL(currentURL string, externalPages map[string]int, internalURLs *list.List) {
	log.Printf("Processing page: %v", currentURL)
	externalURLs := list.New()

	resp, err := http.Get(currentURL)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	//fmt.Println(resp)
	//fmt.Println(resp.StatusCode)
	if resp.StatusCode != 200 {
		log.Printf("Status code is %v", resp.StatusCode)
		return
	}
	//fmt.Println(resp.Header)
	//fmt.Printf("%T", resp)
	defer resp.Body.Close()
	parseHTML(resp.Body, internalURLs, externalURLs)
	checkExternalURLs(externalURLs, externalPages)

	return
}

func main() {
	externalPages := make(map[string]int)

	var host string
	var limit int
	flag.StringVar(&host, "host", "", "The URL of the host")
	flag.IntVar(&limit, "limit", 0, "The max number of pages to visit")
	flag.Parse()
	if host == "" {
		fmt.Println("Need --host")
		os.Exit(1)
	}
	// TODO: check if the URL parses properly
	// TODO: parameter how many pages to check or how far to go from the initial page (how deep to go)
	// TODO: stop on first error?
	// TODO: stop on Nth error?
	log.Printf("Host: %v", host)

	internalURLs := list.New()
	count := 1
	processURL(host, externalPages, internalURLs)

	host = strings.TrimSuffix(host, "/")
	for {
		if internalURLs.Len() == 0 {
			break
		}
		if limit > 0 && count >= limit {
			log.Printf("Limit of %v pages was reached", limit)
			break
		}
		count++
		item := internalURLs.Front()
		thisPage := strings.TrimPrefix(item.Value.(string), "/")
		thisURL := host + "/" + thisPage
		//log.Printf("internalURL: %v", thisURL)
		processURL(thisURL, externalPages, internalURLs)
		internalURLs.Remove(item)
	}

	log.Println("------ Report ----")
	for key, value := range externalPages {
		fmt.Printf("%-4d  %s\n", value, key)
	}
}
