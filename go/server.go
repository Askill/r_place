package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}
var img = GetImage(1000, 1000)

func get(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("error while upgrading", err)
		return
	}
	defer c.Close()
	ticker := time.NewTicker(1 * time.Second)

	var refImage = GetImage(img.width, img.height)
	var tmpImage = GetImage(img.width, img.height)

	for range ticker.C {
		copy(refImage.pixels, img.pixels)
		diff := tmpImage.GetDiff(&refImage)
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
			}
		}
		copy(tmpImage.pixels, refImage.pixels)
	}
}

func set(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("error while upgrading", err)
		return
	}
	defer c.Close()
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		message := Message{}
		json.Unmarshal(msg, &message)

		img.SetPixel(message)
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
