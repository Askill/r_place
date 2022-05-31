package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var img = GetImage(1000, 1000)

func get(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("error while upgrading", err)
		return
	}

	c.SetReadLimit(maxMessageSize)
	c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	defer c.Close()
	ticker := time.NewTicker(1 * time.Second)
	var tmpImage = GetImage(img.width, img.height)

	for range ticker.C {

		diff := tmpImage.GetDiff(&img)
		for i := 0; i < int(diff.width*diff.height); i++ {
			pix := diff.pixels[i]
			if pix.pixel.UserID != 0 {
				x := i / int(diff.width)
				y := i % int(diff.height)
				msg := Message{X: uint32(x), Y: uint32(y), Timestamp: pix.pixel.Timestamp, UserID: pix.pixel.UserID, Color: pix.pixel.Color}
				marshalMsg, err := json.Marshal(msg)
				if err != nil {
					log.Println("error while writing image", err)
					break
				}
				err = c.WriteMessage(1, marshalMsg)
				_, msg2, _ := c.ReadMessage()
				_ = msg2
			}
		}
		if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			return
		}
		copy(tmpImage.pixels, img.pixels)
	}
}

func set(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("error while upgrading", err)
		return
	}
	c.SetReadLimit(maxMessageSize)

	c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	defer c.Close()
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		if mt == websocket.PingMessage {
			continue
		}
		message := Message{}
		json.Unmarshal(msg, &message)

		status := img.SetPixel(message)
		err = c.WriteMessage(1, []byte(strconv.Itoa(status)))
	}
}

func main() {
	var addr = flag.String("addr", "localhost:8080", "http service address")

	flag.Parse()
	log.SetFlags(0)
	log.Println("starting server on", *addr)
	http.HandleFunc("/get", get)
	http.HandleFunc("/set", set)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
