package main

import (
	"encoding/json"
	"fmt"
	"sync"
)

type Message struct {
	X         uint16 `json:"x"`
	Y         uint16 `json:"y"`
	Color     uint8  `json:"color"`
	Timestamp int64  `json:"timestamp"`
	UserID    uint64 `json:"userid"`
}

func (message *Message) JsonToStruct(input []byte) *Message {
	json.Unmarshal(input, message)
	return message
}

type pixel struct {
	Color     uint8  `json:"color"`
	Timestamp int64  `json:"timestamp"`
	UserID    uint64 `json:"userid"`
	Mutex     sync.Mutex
}

type image struct {
	width  uint16
	height uint16
	pixels []pixel
}

type messagePixel struct {
	Color     uint8  `json:"color"`
	Timestamp int64  `json:"timestamp"`
	UserID    uint64 `json:"userid"`
}
type messageImage struct {
	Width  uint16         `json:"width"`
	Height uint16         `json:"height"`
	Pixels []messagePixel `json:"pixel"`
}

func GetMessageImage(w uint16, h uint16) messageImage {
	pixels := make([]messagePixel, w*h)
	for i := 0; i < int(w*h); i++ {
		pixels[i] = messagePixel{Color: 0, Timestamp: 0, UserID: 0}
	}
	return messageImage{Width: w, Height: h, Pixels: pixels}
}

func GetImage(w uint16, h uint16) image {
	pixels := make([]pixel, w*h)
	for i := 0; i < int(w*h); i++ {
		pixels[i] = pixel{Color: 0, Timestamp: 0, UserID: 0, Mutex: sync.Mutex{}}
	}
	return image{width: w, height: h, pixels: pixels}
}

func (p *pixel) setColor(color uint8, timestamp int64, userid uint64) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	if timestamp > p.Timestamp {
		p.Color = color
		p.Timestamp = timestamp
		p.UserID = userid
	}
}

func (img *image) SetPixel(message Message) *image {
	if message.X >= img.width || message.Y >= img.height || message.X < 0 || message.Y < 0 {
		fmt.Printf("User %d tried accessing out of bounds \n", message.UserID)
		return img
	}
	if message.Color >= 16 || message.Color < 0 {
		fmt.Printf("User %d tried setting non existent color \n", message.UserID)
		return img
	}
	pos := uint32(message.X)*uint32(img.width) + uint32(message.Y)
	img.pixels[pos].setColor(message.Color, message.Timestamp, message.UserID)
	return img
}

func comparePixels(pixel1 *pixel, pixel2 *pixel) bool {
	return pixel1.Color == pixel2.Color && pixel1.Timestamp == pixel2.Timestamp && pixel1.UserID == pixel2.UserID
}

func (img *image) GetDiff(img2 *image) messageImage {
	diff := GetMessageImage(img.width, img.height)
	for i := 0; i < int(img.width*img.height); i++ {
		if comparePixels(&img.pixels[i], &img2.pixels[i]) {
			diff.Pixels[i].Color = img2.pixels[i].Color
			diff.Pixels[i].UserID = img2.pixels[i].UserID
			diff.Pixels[i].Timestamp = img2.pixels[i].Timestamp
		}
	}
	return diff
}
