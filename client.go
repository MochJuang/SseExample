package main

import (
	"context"
	"fmt"
	"net/http"
)

type SSEClient struct {
	ClientId string
	Writer   http.ResponseWriter
	send     chan string
	Ctx      context.Context
	manager  *Manager
}

func (c *SSEClient) NewClient() {
	c.send = make(chan string)
}

func (c *SSEClient) write() {

	for {
		select {
		case message := <-c.send:
			flushMessage(c.Writer, message)

		case <-c.Ctx.Done():
			c.manager.unregister <- c
			return
		}
	}

}

func flushMessage(w http.ResponseWriter, message string) {
	if flusher, ok := w.(http.Flusher); ok && flusher != nil {
		fmt.Fprint(w, message)
		flusher.Flush()
	}
}
