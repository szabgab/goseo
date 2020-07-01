package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func processURL(currentURL string) {
	resp, err := http.Get(currentURL)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// fmt.Println(resp)
	//fmt.Println(resp.StatusCode)
	if resp.StatusCode != 200 {
		log.Printf("Status code is %v", resp.StatusCode)
		return
	}
	// fmt.Println(resp.Header)
	//fmt.Printf("%T", resp)
	defer resp.Body.Close()
	//fmt.Printf("%T\n", resp.Body)
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Printf("Could not read body: %v", err)
	// 	return
	// }
	// body = body
	// fmt.Println(string(body))
	// txt := "<html></html>"
	zz := html.NewTokenizer(resp.Body)
	for {
		tt := zz.Next()
		if tt == html.ErrorToken {
			if zz.Err() == io.EOF {
				return
			}
			log.Printf("Error: %v", zz.Err())
			return
		}
		//fmt.Printf("Token: %v\n", zz.Token())
		tag, _ := zz.TagName()
		//fmt.Print(x)
		// if err != nil {
		// 	log.Printf("Error: %v", err)
		// 	continue
		// }
		fmt.Printf("Tag: %v\n", string(tag))
		for {
			attrKey, attrValue, moreAttr := zz.TagAttr()
			if string(attrKey) == "" {
				break
			}
			fmt.Printf("Attr: %v\n", string(attrKey))
			fmt.Printf("Attr: %v\n", string(attrValue))
			fmt.Printf("Attr: %v\n", moreAttr)
			if !moreAttr {
				break
			}
		}
		// Process the current token.
	}

	return
}

func main() {
	var host string
	flag.StringVar(&host, "host", "", "The URL of the host")
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

	processURL(host)
}
