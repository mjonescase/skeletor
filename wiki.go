package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"skeletor/utils"
	"time"
)

var (
	clients              = make(map[*websocket.Conn]bool) // connected clients
	broadcast            = make(chan PublishedContent)    // broadcast channel
	upgrader             = websocket.Upgrader{}           // configure the upgrader
	config               = map[string]string{}
	configFile   *string = flag.String("config", "config", "Path to config file")
	DefaultError         = map[string]string{"ErrorReason": "You sent in a request with invalid json"}
	AuthError            = map[string]string{"ErrorReason": "You are not authorized to perform this action"}
	session              = &sql.DB{}
	hashSalt             = ""
)

// published content type enum
const (
	PUBTYPE_MESSAGE  = iota // 0
	PUBTYPE_CONTACTS = iota // 1
)

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
		var msg PublishedContent
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

func validateLogin(req *Profile) bool {
	result := false
	req.Password = utils.HashPassword(req.Password)
	result = queryUserCredential(req)
	return result
}

func generateSessionFromProfile(request Profile) http.Cookie {
	sessionInfo := fmt.Sprintf("{'username': '%s', 'id': '%s', 'email': '%s'}", request.Username, request.Id, request.Email)
	expire := time.Now().AddDate(0, 0, 1)
	cookie := http.Cookie{"SessionInfo", sessionInfo, "/", config["hostname"], expire, expire.Format(time.UnixDate), 86400, true, false, "", []string{""}}
	return cookie
}

func handleLogin(rw http.ResponseWriter, req *http.Request) {
	request := Profile{}

	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&request)

	if err != nil {
		utils.MustEncode(rw, AuthError)
		return
	}

	authenticated := validateLogin(&request)
	if authenticated {
		cookie := generateSessionFromProfile(request)
		http.SetCookie(rw, &cookie)
		rw.WriteHeader(http.StatusOK)
	} else {
		rw.WriteHeader(http.StatusUnauthorized)
		utils.MustEncode(rw, AuthError)
		return
	}

	broadcast <- PublishedContent{Type: PUBTYPE_CONTACTS, Contents: getAllUsers()}
	utils.MustEncode(rw, request)
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
	// broadcast the new user
	broadcast <- PublishedContent{Type: PUBTYPE_CONTACTS, Contents: getAllUsers()}
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
	//proxy := New("http://0.0.0.0:8080")
	//http.HandleFunc("/", proxy.handle)

	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/register/", handleRegistration)
	http.HandleFunc("/login/", handleLogin)
	http.HandleFunc("/profiles/", handleProfilesRequest)

	// Start listening for incoming chat messages
	go handleMessages()

	// Start the server on localhost port 5000 and log any errors
	log.Println("http server started on :5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
