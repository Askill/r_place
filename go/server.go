package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		//origin := r.Header.Get("Origin")
		return true
	},
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
var tmpImage = GetImage(img.Width, img.Height)

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
	for range ticker.C {
		diff := tmpImage.GetDiff(&img)
		for i := 0; i < int(diff.Width*diff.Height); i++ {
			pix := diff.Pixels[i]
			if pix.Pixel.UserID != 0 {
				x := i / int(diff.Width)
				y := i % int(diff.Height)
				msg := Message{X: uint32(x), Y: uint32(y), Timestamp: pix.Pixel.Timestamp, UserID: pix.Pixel.UserID, Color: pix.Pixel.Color}
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

		copy(tmpImage.Pixels, img.Pixels)
	}
}

func getAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(img)
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

func loadState(img *image, path string) {
	stateJSON, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stateJSON.Close()
	byteValue, _ := ioutil.ReadAll(stateJSON)

	json.Unmarshal(byteValue, &img)
}

func saveState(img *image, path string, period time.Duration) {
	ticker := time.NewTicker(period * time.Second)

	for range ticker.C {
		imgJSON, _ := json.Marshal(img)
		file, err := os.Create(path)

		if err != nil {
			return
		}
		defer file.Close()
		file.WriteString(string(imgJSON))
		if err != nil {
			fmt.Println("Could not save state")
			fmt.Println(err)
		}
	}
}

func main() {
	var addr = flag.String("addr", "localhost:8080", "http service address")

	flag.Parse()
	log.SetFlags(0)
	log.Println("starting server on", *addr)

	cachePath := "./state.json"

	loadState(&img, cachePath)
	go saveState(&img, cachePath, 10)

	http.HandleFunc("/get", get)
	http.HandleFunc("/getAll", getAll)
	http.HandleFunc("/set", set)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
