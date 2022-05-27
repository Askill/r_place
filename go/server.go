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

func set(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("error while upgrading", err)
		return
	}
	defer c.Close()
	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		message := Message{}
		message.JsonToStruct(msg)
		img.SetPixel(message)

		err = c.WriteMessage(mt, msg)
		if err != nil {
			//log.Println("write:", err)
			break
		}
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("error while upgrading", err)
		return
	}
	defer c.Close()
	log.Println("Client Connected")
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			tmpImage := GetImage(img.width, img.height)
			diff := tmpImage.GetDiff(&img)
			msg, err := json.Marshal(diff)
			err = c.WriteMessage(1, msg)
			if err != nil {
				log.Print("error while writing image", err)
			}
			copy(tmpImage.pixels, img.pixels)
		}
	}()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	var addr = flag.String("addr", "localhost:8080", "http service address")

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/set", set)
	http.HandleFunc("/get", get)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
