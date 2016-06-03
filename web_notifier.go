package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type WebNotifier struct {
	url string
}

func NewWebNotifier(url string) *WebNotifier {
	return &WebNotifier{url: url}
}

func (wn *WebNotifier) Notify(n *Notification) {
	buf := bytes.NewBuffer(nil)

	if err := json.NewEncoder(buf).Encode(n); err != nil {
		fmt.Println("failed to encode notification:", err)
		return
	}
	resp, err := http.Post(wn.url, "application/json", buf)

	if err != nil {
		fmt.Printf("failed notify webhook: %s: %s\n", wn.url, err)
		return
	}

	if resp.StatusCode < 200 || 299 < resp.StatusCode {
		fmt.Printf("notification webhook responded with unexpected status:", resp.Status)
	}
}
