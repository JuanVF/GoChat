package controller

import (
	"log"
	"net/http"

	"github.com/JuanVF/GoChat/model"
)

var SERVER_PORT = ":3000"

// This is basically a init function for the server
func Serve() {
	HandleView()

	http.HandleFunc(model.SOCKET_PATH, handlerSocket)

	go model.HandleMessages()

	log.Println("Server running at ", SERVER_PORT)

	if err := http.ListenAndServe(SERVER_PORT, nil); err != nil {
		log.Fatal("Server error: ", err)
	}
}

// Here will be serve our client
func HandleView() {
	fs := http.StripPrefix("/", http.FileServer(http.Dir("View")))

	http.Handle("/", fs)
}

// Socket Handler
func handlerSocket(w http.ResponseWriter, r *http.Request) {
	// We use the HTTP Conn as a socket one
	ws, err := model.Upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal(err)
	}

	// The conn will close when the app finishes
	defer ws.Close()

	// We add our client
	model.Clients[ws] = true

	for {
		var msg model.Message

		err := ws.ReadJSON(&msg)

		// In error case we close connection with our client
		if err != nil {
			log.Printf("Connection Err: %v", err)
			delete(model.Clients, ws)

			break
		}

		model.Broadcast <- msg
	}
}
