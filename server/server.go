package main

import (
	"encoding/json"
	"flag"
	"fmt"
	go_image "image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
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
var diff = GetImage(img.Width, img.Height)

func calcDiff() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		diff = tmpImage.GetDiff(&img)
		copy(tmpImage.Pixels, img.Pixels)
	}
}

func get(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("error while upgrading", err)
		return
	}
	fmt.Println("incoming connection")
	c.SetReadLimit(maxMessageSize)
	c.SetPongHandler(func(string) error { c.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	defer c.Close()
	ticker := time.NewTicker(200 * time.Millisecond)
	for range ticker.C {

		for i := 0; i < int(diff.Width*diff.Height); i++ {

			pix := diff.Pixels[i]

			if pix.Pixel.UserID != 0 {
				x := i / int(diff.Width)
				y := i % int(diff.Height)
				msg := Message{X: uint32(x), Y: uint32(y), Timestamp: pix.Pixel.Timestamp, UserID: pix.Pixel.UserID, Color: pix.Pixel.Color}
				marshalMsg, err := json.Marshal(msg)
				if err != nil {
					log.Println("error while marshalling image", err)
					break
				}
				err = c.WriteMessage(1, marshalMsg)
				if err != nil {
					log.Println("error while writing image", err)
					break
				}
			}
		}
		if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			return
		}

		copy(tmpImage.Pixels, img.Pixels)
	}
}
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}
func getAll(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)
	w.Header().Set("Content-Type", "image/png")
	colors := [16]color.Color{color.RGBA{255, 255, 255, 0xff}, color.RGBA{228, 228, 228, 0xff}, color.RGBA{136, 136, 136, 0xff}, color.RGBA{34, 34, 34, 0xff}, color.RGBA{255, 167, 209, 0xff}, color.RGBA{229, 0, 0, 0xff}, color.RGBA{229, 149, 0, 0xff}, color.RGBA{160, 106, 66, 0xff}, color.RGBA{229, 217, 0, 0xff}, color.RGBA{148, 224, 68, 0xff}, color.RGBA{2, 190, 1, 0xff}, color.RGBA{0, 211, 221, 0xff}, color.RGBA{0, 131, 199, 0xff}, color.RGBA{0, 0, 234, 0xff}, color.RGBA{207, 110, 228, 0xff}, color.RGBA{130, 0, 128, 0xff}}
	upLeft := go_image.Point{0, 0}
	lowRight := go_image.Point{int(img.Width), int(img.Height)}

	png_img := go_image.NewRGBA(go_image.Rectangle{upLeft, lowRight})

	for x := uint32(0); x < img.Width; x++ {
		for y := uint32(0); y < img.Height; y++ {
			png_img.Set(int(y), int(x), colors[img.Pixels[x*img.Width+y].Pixel.Color])
		}
	}
	png.Encode(w, png_img)
}

func sendPing(ticker *time.Ticker, c *websocket.Conn, mutex *sync.Mutex) {
	for range ticker.C {
		mutex.Lock()
		defer mutex.Unlock()
		if err := c.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			fmt.Println("Error while sending ping")
			return
		}
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

	ticker := time.NewTicker(pingPeriod / 2)
	defer c.Close()
	defer ticker.Stop()
	mutex := sync.Mutex{}
	go sendPing(ticker, c, &mutex)

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
		mutex.Lock()

		err = c.WriteMessage(1, []byte(strconv.Itoa(status)))
		mutex.Unlock()
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
	var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

	flag.Parse()
	log.SetFlags(0)
	log.Println("starting server on", *addr)

	cachePath := "./state.json"

	loadState(&img, cachePath)
	go saveState(&img, cachePath, 10)
	go calcDiff()
	http.HandleFunc("/get", get)
	http.HandleFunc("/getAll", getAll)
	http.HandleFunc("/set", set)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
