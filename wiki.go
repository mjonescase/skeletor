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
	"strings"
	"time"
)

var (
	rooms                = map[string]Room{}
	upgrader             = websocket.Upgrader{}           // configure the upgrader
	config               = map[string]string{}
	configFile   *string = flag.String("config", "config", "Path to config file")
	DefaultError         = map[string]string{"ErrorReason": "You sent in a request with invalid json"}
	AuthError            = map[string]string{"ErrorReason": "You are not authorized to perform this action"}
	session              = &sql.DB{}
	hashSalt             = ""
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
	fmt.Sprintf("Got a request to join room %s", request.URL.RawQuery)
	room := rooms[strings.Split(request.URL.RawQuery, "=")[1]]
	fmt.Sprintf("Got a request to join room %s", room)
	ws, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		log.Fatal(err)
	}
	// call addConnection
	fmt.Sprintf("registering client to room")
	registerClient(ws, room)
	serveForever(ws, room)
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

	utils.MustEncode(rw, request)
}

func handleRegistration(rw http.ResponseWriter, req *http.Request) {
	log.Println("got a registration request")
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
	initConfig()
	initDb()

	// Create a simple file server
	proxy := New("http://frontend:8080")
	http.HandleFunc("/", proxy.handle)

	//fs := http.FileServer(http.Dir("./public"))
	//http.Handle("/", fs)
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/register/", handleRegistration)
	http.HandleFunc("/login/", handleLogin)

	// Start listening for incoming chat messages
	initRooms()

	// Start the server on localhost port 5000 and log any errors
	log.Println("http server started on :5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
