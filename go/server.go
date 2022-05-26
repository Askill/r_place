// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}

type Pixel struct {
	X         uint16 `json:"x"`
	Y         uint16 `json:"y"`
	Color     uint8  `json:"color"`
	Timestamp int64  `json:"timestamp"`
	UserID    uint64 `json:"userid"`
}

func JsonToStruct(input []byte) Pixel {
	pixel := Pixel{}
	json.Unmarshal(input, &pixel)
	return pixel
}

func serve(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("error while upgrading", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			//log.Println("read:", err)
			break
		}
		pixel := JsonToStruct(message)

		tm := time.Unix(pixel.Timestamp, 0)
		log.Println(tm)
		err = c.WriteMessage(mt, message)
		if err != nil {
			//log.Println("write:", err)
			break
		}
	}
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", serve)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
