package main

import (
	"fmt"
	"log"
)

type Manager struct {
	clients    map[*SSEClient]bool
	broadcast  chan string
	register   chan *SSEClient
	unregister chan *SSEClient
}

func NewManager() *Manager {
	return &Manager{
		clients:    make(map[*SSEClient]bool),
		register:   make(chan *SSEClient),
		unregister: make(chan *SSEClient),
		broadcast:  make(chan string),
	}
}

func (m *Manager) Run() {
	for {
		select {
		case client := <-m.register:
			if ok := m.clients[client]; ok {
				return
			}

			fmt.Println("Register Client ID : ", client.ClientId)
			m.clients[client] = true
		case client := <-m.unregister:
			if _, ok := m.clients[client]; ok {
				fmt.Println("Unregister Client ID : ", client.ClientId)
				delete(m.clients, client)
				close(client.send)
			}
		case message := <-m.broadcast:
			fmt.Println("--------------------------------------------")
			log.Println("broadcast message : ", string(message))
			log.Println("total subscriber : ", len(m.clients))
			fmt.Println("--------------------------------------------")
			for client := range m.clients {
				select {
				case client.send <- message:
				default:
					delete(m.clients, client)
					close(client.send)
				}
			}
		}
	}
}
