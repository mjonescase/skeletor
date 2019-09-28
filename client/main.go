package main

import (
	"flag"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type PublishedContent struct {
	Username  string  `json:"username"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Message   string  `json: "message"`
}

func connect(u url.URL, username string, interval int) {
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
			log.Printf("recv: %s", msg.Message)
		}
	}()

	writeAtRegularIntervals(c, username, done, interval)
}

func writeAtRegularIntervals(c *websocket.Conn,	username string, done chan struct{}, interval int) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			seconds := t.Second()
			if interval > 0 && seconds % interval == 0 {
				msg := PublishedContent {
					Username: username,
						Latitude: float64(t.Minute()),
						Longitude: float64(seconds),
				}

				err := c.WriteJSON(msg)
				log.Printf("sent: %s", msg.Message)
				if err != nil {
					log.Println("write:", err)
					return
				}
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
	if len(args) > 1 {
		username := args[len(args) - 2]
		interval, err := strconv.Atoi(args[len(args) - 1])
		if err != nil {
			log.Println("Bad interval, must be int:", args[len(args) -1])
			os.Exit(1)
		}
		connect(u, username, interval)
	} else {
		connect(u, "", -1)
	}
}
