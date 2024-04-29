package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

func sendRequests(url string, numRequests int) {
	var wg sync.WaitGroup
	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer func() {
				if r := recover(); r != nil {
					fmt.Println("Recovered from panic:", r)
				}
			}()
			err := performRequest(url)
			if err != nil {
				fmt.Println("Error:", err)
			}
		}()
	}
	wg.Wait()
	fmt.Println("All requests completed")
}

func main() {
	url := "" //"http://101.34.70.9:8080/api/v1/accounts" // Update the URL with server address
	numRequests := 5
	sendRequests(url, numRequests)
}

func performRequest(url string) error {
	if url == "" {
		panic("URL cannot be empty")
	}

	data := map[string]interface{}{
		"id":   1,
		"name": "John Doe",
		"age":  30,
		// Add other required fields
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Print the response status
	fmt.Println("Response Status:", resp.Status)

	return nil
}
