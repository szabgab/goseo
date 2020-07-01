package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

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

	resp, err := http.Get(host)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	// fmt.Println(resp)
	//fmt.Println(resp.StatusCode)
	if resp.StatusCode != 200 {
		log.Printf("Status code is %v", resp.StatusCode)
	}
	// fmt.Println(resp.Header)
	//fmt.Printf("%T", resp)
}
