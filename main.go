package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strings"
)

var (
	rooms                = map[string]Room{}
	upgrader             = websocket.Upgrader{}           // configure the upgrader
	config               = map[string]string{}
	DefaultError         = map[string]string{"ErrorReason": "You sent in a request with invalid json"}
	AuthError            = map[string]string{"ErrorReason": "You are not authorized to perform this action"}
	passcodesToRooms     = map[string]string{"blue123": "commBlue"," red456": "commRed"}
)

// room names
const (
	COMM_BLUE     = "commBlue"
	COMM_GREEN    = "commGreen"
	COMM_RED      = "commRed"
	LOCATION_BLUE = "locationBlue"
	LOCATION_RED  = "locationRed"
)

func handleConnections(writer http.ResponseWriter, request *http.Request) {
	roomName := passcodesToRooms[strings.Split(request.URL.RawQuery, "=")[1]]
	room := rooms[roomName]
	
	ws, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Fatal(err)
	}
	// call addConnection
	fmt.Sprintf("registering client to room")
	registerClient(ws, room)
	serveForever(ws, room)
}

func initRooms() {
	rooms [COMM_BLUE]     = buildRoom()
	rooms [COMM_GREEN]    = buildRoom()
	rooms [COMM_RED]      = buildRoom()
	rooms [LOCATION_BLUE] = buildRoom()
	rooms [LOCATION_RED]  = buildRoom()
	for  _, room := range rooms {
		go handleMessages(room)
	}
}

func main() {
	http.HandleFunc("/ws", handleConnections)

	// Start listening for incoming chat messages
	initRooms()

	// Start the server on localhost port 5000 and log any errors
	log.Println("http server started on :5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
