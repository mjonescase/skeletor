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
	"net/http/httputil"
	_ "net/http/pprof"
	"net/url"
	"regexp"
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

type PublishedContent struct {
	Type     int         `json:"type"` // either PUBTYPE_MESSAGE or PUBTYPE_CONTACTS
	Contents interface{} `json:"contents"`
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
		sessionInfo := fmt.Sprintf("{'username': '%s', 'id': '%s', 'email': '%s'}", request.Username, request.Id, request.Email)
		expire := time.Now().AddDate(0, 0, 1)
		cookie := http.Cookie{"SessionInfo", sessionInfo, "/", config["hostname"], expire, expire.Format(time.UnixDate), 86400, true, false, "", []string{""}}
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

type Prox struct {
	target        *url.URL
	proxy         *httputil.ReverseProxy
	routePatterns []*regexp.Regexp // add some route patterns with regexp
}

func New(target string) *Prox {
	url, err := url.Parse(target)

	if err != nil {
		panic(err)
	}
	return &Prox{target: url, proxy: httputil.NewSingleHostReverseProxy(url)}
}

func (p *Prox) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-GoProxy", "GoProxy")

	if p.routePatterns == nil || p.parseWhiteList(r) {
		p.proxy.ServeHTTP(w, r)
	}
}

func (p *Prox) parseWhiteList(r *http.Request) bool {
	for _, regexp := range p.routePatterns {
		fmt.Println(r.URL.Path)
		if regexp.MatchString(r.URL.Path) {
			// let's forward it
			return true
		}
	}
	fmt.Println("Not accepted routes %x", r.URL.Path)
	return false
}

func main() {
	initConfig()
	initDb()

	// Create a simple file server
	proxy := New("http://frontend:8080")
	//fs := http.FileServer(http.Dir("./public"))
	http.HandleFunc("/", proxy.handle)
	http.HandleFunc("/ws", handleConnections)
	http.HandleFunc("/register/", handleRegistration)
	http.HandleFunc("/login/", handleLogin)

	// Start listening for incoming chat messages
	go handleMessages()

	// Start the server on localhost port 5000 and log any errors
	log.Println("http server started on :5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
