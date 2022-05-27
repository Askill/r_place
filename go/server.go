// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build ignore
// +build ignore

package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
}
var img = GetImage(10, 10)

func serve(w http.ResponseWriter, r *http.Request) {
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

func main() {
	var addr = flag.String("addr", "localhost:8080", "http service address")

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", serve)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
