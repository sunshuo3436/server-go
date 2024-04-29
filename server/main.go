package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	URL   string `json:"url"`
	Error string `json:"error,omitempty"`
} //存储http请求的相应：url、error，并且转为json

// add fsm
// Defining types of state and state constants
// Defining state machine names
type State int

const (
	WaitingForResponse State = iota
	Timeout
) //Two states of a state machine:WaitingForResponse, timeout

// saving fsm
// State machine constructs: state name, channel, response
type StateMachine struct {
	state     State
	urlChan   chan Response
	responses []Response
}

func (sm *StateMachine) transition() {
	switch sm.state {
	case WaitingForResponse:
		select {
		case res := <-sm.urlChan:
			sm.responses = append(sm.responses, res)
		case <-time.After(2 * time.Second):
			sm.responses = append(sm.responses, Response{Error: "Timed out"})
			sm.state = Timeout
		}
	}
}

// receives the URL, sends the response channel, sends the request to http
func fetchURL(url string, urlChan chan<- Response) {
	resp, err := http.Get(url)
	if err != nil {
		urlChan <- Response{URL: url, Error: fmt.Sprintf("Error: %v", err)}
		return
	}
	defer resp.Body.Close()
	urlChan <- Response{URL: url}
}

func handler(w http.ResponseWriter, r *http.Request) {
	urls := []string{
		"https://www.google.com",
		"https://www.youtube.com",
		"https://www.amazon.com",
		"https://www.github.com",
	}

	urlChan := make(chan Response, len(urls)) //创建一个response通道

	sm := StateMachine{
		state:     WaitingForResponse,
		urlChan:   urlChan,
		responses: make([]Response, 0, len(urls)),
	}

	for _, url := range urls {
		go fetchURL(url, urlChan)
	}

	for sm.state != Timeout && len(sm.responses) < len(urls) {
		sm.transition()
	}

	jsonData, err := json.Marshal(sm.responses)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func main() {
	http.HandleFunc("/rpc", handler)
	fmt.Println("Server is running at :8080")
	http.ListenAndServe(":8080", nil)
}
