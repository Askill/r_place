package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}
var img = GetImage(10, 10)
var tmpImage = GetImage(img.width, img.height)

func serve(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("error while upgrading", err)
		return
	}
	defer c.Close()
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for range ticker.C {
			diff := tmpImage.GetDiff(&img)
			msg, err := json.Marshal(diff)
			fmt.Println(diff)
			// TODO only write chenaged pixels to channel instead of entire image
			err = c.WriteMessage(1, msg)
			if err != nil {
				log.Println("error while writing image", err)
				break
			}
			copy(img.pixels, tmpImage.pixels)
		}
	}()
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		message := Message{}
		message.JsonToStruct(msg)
		img.SetPixel(message)
	}

}

func main() {
	var addr = flag.String("addr", "localhost:8080", "http service address")

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", serve)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
