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
var img = GetImage(10, 10)
var tmpImage = GetImage(img.width, img.height)

func write(ticker time.Ticker, c websocket.Conn) {
	for range ticker.C {
		diff := tmpImage.GetDiff(&img)
		for i := 0; i < int(diff.Width*diff.Height); i++ {
			pix := diff.Pixels[i]
			if pix.UserID != 0 {
				x := i / int(diff.Width)
				y := i % int(diff.Height)
				msg := Message{X: uint16(x), Y: uint16(y), Timestamp: pix.Timestamp, UserID: pix.UserID, Color: pix.Color}
				marshalMsg, err := json.Marshal(msg)
				if err != nil {
					log.Println("error while writing image", err)
					break
				}
				err = c.WriteMessage(1, marshalMsg)
			}
		}
		copy(img.pixels, tmpImage.pixels)
	}
}

func serve(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("error while upgrading", err)
		return
	}
	defer c.Close()
	ticker := time.NewTicker(1 * time.Second)
	go write(*ticker, *c)
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
	log.Println("starting server on", *addr)
	http.HandleFunc("/", serve)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
