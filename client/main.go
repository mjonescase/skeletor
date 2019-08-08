package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type PublishedContent struct {
	Username  string  `json:"username"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Message   string  `json: "message"`
}

func connect(u url.URL, username string) {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			var msg PublishedContent
			err := c.ReadJSON(&msg)
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", msg)
		}
	}()

	if username != "" {
		writeAtRegularIntervals(c, username, done)
	}
}

func writeAtRegularIntervals(c *websocket.Conn,	username string, done chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			msg := PublishedContent {
				Username: username,
				Latitude: float64(t.Hour()),
				Longitude: float64(t.Minute()),
			}
					
			err := c.WriteJSON(msg)
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func main() {
	var addr = flag.String("addr", "localhost:5000", "http service address")
	u := url.URL{
		Scheme: "ws",
		Host: *addr,
		Path: "/ws",
		RawQuery: "passcode=blue123",
	}
	args := os.Args[1:]
	if len(args) > 0 {
		connect(u, os.Args[len(os.Args) -1])
	} else {
		connect(u, "")
	}
}
