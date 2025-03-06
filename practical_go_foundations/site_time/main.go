package main

import (
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

func main() {
	urls := []string{"https://golang.org", "https://google.com", "https://apple.com"}

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func() {
			defer wg.Done()
			siteTime(url)
		}()
	}

	wg.Wait()
}

func siteTime(url string) {
	start := time.Now()

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("fetching %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(io.Discard, resp.Body)
	if err != nil {
		log.Printf("reading %s: %v", url, err)
		return
	}

	duration := time.Since(start)

	log.Printf("%s: %s\n", url, duration)
}
