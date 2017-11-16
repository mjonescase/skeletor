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
	"skeletor/utils"
)

var (
	clients                  = make(map[*websocket.Conn]bool) // connected clients to messages
	broadcast                = make(chan Message)             // messages broadcast channel
	profileClients           = make(map[*websocket.Conn]bool) // clients connected to profile
	profileBroadcast         = make(chan []Profile)           // contacts list broadcast channel
	upgrader                 = websocket.Upgrader{}           // configure the upgrader
	config                   = map[string]string{}
	configFile       *string = flag.String("config", "config", "Path to config file")
	DefaultError             = map[string]string{"ErrorReason": "You sent in a request with invalid json"}
	session                  = &sql.DB{}
	hashSalt                 = ""
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
	Email        string `json:"email"`
	Title        string `json:"title"`
	Password     string `json:",omitempty"`
	MobileNumber string `json:"mobilenumber"`
}

func handleProfileConnection(writer http.ResponseWriter, request *http.Request) {
	ws, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()
	profileClients[ws] = true

	//send an update every time a new client registers here
	log.Printf("broadcasting all the users")
	profileBroadcast <- getAllUsers()

	// not sure if this is necessary. We don't need to read anything on this channel.
	for {
	}
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

func handleRegistration(rw http.ResponseWriter, req *http.Request) {
	request := Profile{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&request)

	if err != nil {
		utils.MustEncode(rw, DefaultError)
		return
	}

	saveUserProfile(&request)
	utils.MustEncode(rw, request)
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
			" user="+config["dbuser"]+
			" dbname="+config["dbname"]+
			" sslmode="+config["sslmode"])

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
	initConfig()
	initDb()

	// Create a simple file server
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/contacts/", handleProfileConnection)
	http.HandleFunc("/register/", handleRegistration)

	// Start listening for incoming chat messages
	go handleMessages()

	// Start the server on localhost port 5000 and log any errors
	log.Println("http server started on :5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
