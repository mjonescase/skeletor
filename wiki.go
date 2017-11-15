package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
)

var (
	clients            = make(map[*websocket.Conn]bool) // connected clients
	broadcast          = make(chan Message)             // broadcast channel
	upgrader           = websocket.Upgrader{}           // configure the upgrader
	config             = map[string]string{}
	configFile *string = flag.String("config", "config", "Path to config file")
	session            = &sql.DB{}
)

// define our message object
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

type Profile struct {
	Id           string `json:"id"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Username     string `json:"username"`
	Title        string `json:"title"`
	Password     string // TODO don't send this back
	MobileNumber string `json:"mobile"`
}

func handleConnections(writer http.ResponseWriter, request *http.Request) {
	ws, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}

		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
func initDb() {
	var err error
	session, err = sql.Open(
		"postgres", "host="+config["dbhost"]+
			"user="+config["dbuser"]+
			"dbname="+config["dbname"]+
			"sslmode="+config["dbsslmode"])

	if err != nil {
		panic(err)
	}
}

func initConfig() {
	conf, err := ioutil.ReadFile(*configFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(conf, &config)
	if err != nil {
		panic(err)
	}
}

func main() {
	// Create a simple file server
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)

	// Start listening for incoming chat messages
	go handleMessages()

	// Start the server on localhost port 5000 and log any errors
	log.Println("http server started on :5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
