package main

import (
	"github.com/gorilla/websocket"
	"log"
)

// define a websocket room
type Room struct {
	Clients map[*websocket.Conn]bool
	Broadcast chan PublishedContent
}

func buildRoom() Room {
	room := Room{}
	room.Clients = make(map[*websocket.Conn]bool)
	room.Broadcast = make(chan PublishedContent)
	return room
}

func registerClient(connection *websocket.Conn,
	room Room) {
	room.Clients[connection] = true
}

func serveForever(connection *websocket.Conn, room Room) {
	// close the connection when this function returns.
	defer connection.Close()
	for {
		var msg PublishedContent
		// Read in a new message as JSON and map it to a Message object
		err := connection.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(room.Clients, connection)
			break
		}

		// Send the newly received message to the broadcast channel
		log.Printf("Got a new message from username %s latitude: %f longitude: %f", msg.Username, msg.Latitude, msg.Longitude)
		room.Broadcast <- msg
	}
}

func handleMessages(room Room) {
	for {
		// Grab the next message from the broadcast channel
		msg := <- room.Broadcast
		log.Printf("got a message from: %s", msg.Username)
		// Send it out to every client that is currently connected
		for client := range room.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(room.Clients, client)
			}
		}
	}
}
