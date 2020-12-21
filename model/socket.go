package model

import (
	"log"

	"github.com/gorilla/websocket"
)

// This will be the struct we'll use to read the messages
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

var SOCKET_PATH = "/api/chat"

// These will be our clients who will be in a chat :)
var Clients = make(map[*websocket.Conn]bool)
var Broadcast = make(chan Message)

// This is how we'll use a HTTP conn as a socket :D
var Upgrader = websocket.Upgrader{}

// This is a function that will be running on a thread to
// Send the messages if there is one received
func HandleMessages() {
	for {
		msg := <-Broadcast

		// We send the message to each client registered
		for client := range Clients {
			err := client.WriteJSON(&msg)

			// In error case we close connection with our client
			if err != nil {
				log.Printf("Message error: %v", err)
				client.Close()
				delete(Clients, client)
			}
		}
	}
}
