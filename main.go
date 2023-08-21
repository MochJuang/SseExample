package main

import (
	"context"
	"fmt"
	"net/http"
)

// type Message struct {
// 	StreamId string
// 	Message  string
// }

// var messageCh chan Message
// var listConnection []string
// var counter int

// func main() {
// 	counter = 1
// 	messageCh = make(chan Message)

// 	fmt.Println("start")
// 	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {

// 		flusher, ok := w.(http.Flusher)
// 		if !ok {
// 			http.Error(w, "error stream not support", http.StatusInternalServerError)
// 			return
// 		}

// 		streamId := r.URL.Query().Get("id")
// 		listConnection = append(listConnection, streamId)
// 		fmt.Println(listConnection)

// 		w.Header().Set("Content-Type", "text/event-stream")
// 		w.Header().Set("Cache-Control", "no-cache")
// 		w.Header().Set("Connection", "keep-alive")
// 		w.Header().Set("Access-Control-Allow-Origin", "*")

// 		for {
// 			select {
// 			case mess := <-messageCh:
// 				fmt.Println(string(mess.Message))
// 				fmt.Fprint(w, mess.Message)
// 				flusher.Flush()
// 			}
// 		}

// 	})

// 	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {

// 		mode := r.URL.Query().Get("mode")
// 		counter += 1
// 		messageCh <- Message{
// 			StreamId: mode,
// 			Message:  fmt.Sprintf("data: counter ke %v\n\n", (counter)),
// 		}

// 		fmt.Fprintln(w, "success publish")
// 	})

// 	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
// 		http.ServeFile(w, r, "index.html")
// 		// w,r IS WRITE AND DELETE YOU INDEX.HTML

// 	})

// 	http.ListenAndServe(":8080", nil)
// }

func main() {
	counter := 1
	manager := NewManager()
	go manager.Run()

	fmt.Println("start")
	// block := make(chan struct{})

	http.HandleFunc("/stream", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		clientId := r.URL.Query().Get("id")
		var client = &SSEClient{
			ClientId: clientId,
			manager:  manager,
			Writer:   w,
			Ctx:      ctx,
		}

		client.NewClient()

		manager.register <- client

		client.write()

		// block <- struct{}{}
	})

	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {

		// mode := r.URL.Query().Get("mode")
		counter += 1

		message := fmt.Sprintf("data: counter ke %v\n\n", (counter))
		manager.broadcast <- message

		fmt.Fprintln(w, "success publish")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
		// w,r IS WRITE AND DELETE YOU INDEX.HTML

	})

	http.ListenAndServe(":8080", nil)
}

func Sum(a int, b int) int {
	return a + b
}
