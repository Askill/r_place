package main

import (
	"fmt"
	"sync"
)

type Message struct {
	X         uint32 `json:"x"`
	Y         uint32 `json:"y"`
	Color     uint8  `json:"color"`
	Timestamp int64  `json:"timestamp"`
	UserID    uint64 `json:"userid"`
}

type pixel struct {
	Color     uint8  `json:"color"`
	Timestamp int64  `json:"timestamp"`
	UserID    uint64 `json:"userid"`
}

type pixelContainer struct {
	pixel pixel
	Mutex sync.Mutex
}

type image struct {
	width  uint32
	height uint32
	pixels []pixelContainer
}

func GetImage(w uint32, h uint32) image {
	pixels := make([]pixelContainer, w*h)
	for i := 0; i < int(w*h); i++ {
		pixels[i] = pixelContainer{pixel: pixel{Color: 0, Timestamp: 0, UserID: 0}, Mutex: sync.Mutex{}}
	}
	return image{width: w, height: h, pixels: pixels}
}

func (p *pixelContainer) setColor(color uint8, timestamp int64, userid uint64) {
	p.Mutex.Lock()
	defer p.Mutex.Unlock()
	if timestamp > p.pixel.Timestamp {
		p.pixel.Color = color
		p.pixel.Timestamp = timestamp
		p.pixel.UserID = userid
	}
}

func (img *image) SetPixel(message Message) int {
	if message.X >= img.width || message.Y >= img.height || message.X < 0 || message.Y < 0 {
		fmt.Printf("User %d tried accessing out of bounds \n", message.UserID)
		return 1
	}
	if message.Color >= 255 || message.Color < 0 {
		fmt.Printf("User %d tried setting non existent color \n", message.UserID)
		return 1
	}
	pos := uint32(message.X)*uint32(img.width) + uint32(message.Y)
	img.pixels[pos].setColor(message.Color, message.Timestamp, message.UserID)
	return 0
}

func comparePixels(pixel1 *pixelContainer, pixel2 *pixelContainer) bool {
	return pixel1.pixel.Color == pixel2.pixel.Color &&
		pixel1.pixel.Timestamp == pixel2.pixel.Timestamp &&
		pixel1.pixel.UserID == pixel2.pixel.UserID
}

func (img *image) GetDiff(img2 *image) image {
	diff := GetImage(img.width, img.height)
	for i := 0; i < int(img.width*img.height); i++ {
		if !comparePixels(&img.pixels[i], &img2.pixels[i]) {
			diff.pixels[i].pixel.Color = img2.pixels[i].pixel.Color
			diff.pixels[i].pixel.UserID = img2.pixels[i].pixel.UserID
			diff.pixels[i].pixel.Timestamp = img2.pixels[i].pixel.Timestamp
		}
	}
	return diff
}
