package main

import (
	"fmt"
	"net/http"
	"time"
)

func fetchURL(url string, ch chan<- string) {
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprintf("Error fetching %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- fmt.Sprintf("Error fetching %s: Status %d", url, resp.StatusCode)
		return
	}

	ch <- fmt.Sprintf("URL hit: %s", url)
}

func main() {
	urls := []string{
		"https://www.google.com",
		"https://www.youtube.com",
		"https://www.amazon.com",
		"https://www.github.com",
	}

	ch := make(chan string, len(urls))

	for _, url := range urls {
		go fetchURL(url, ch)
	}

	timeout := time.After(2 * time.Second)

	for i := 0; i < len(urls); i++ {
		select {
		case res := <-ch:
			fmt.Println(res)
		case <-timeout:
			fmt.Println("Timed out.")
			return
		}
	}
}
