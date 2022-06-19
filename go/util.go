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
	Color     uint8  `json:"c"`
	Timestamp int64  `json:"ts"`
	UserID    uint64 `json:"uid"`
}

type pixelContainer struct {
	Pixel pixel `json:"p"`
	mutex sync.Mutex
}

type image struct {
	Width  uint32           `json:"Width"`
	Height uint32           `json:"Height"`
	Pixels []pixelContainer `json:"Pixels"`
	Mutex  sync.Mutex
}

func GetImage(w uint32, h uint32) image {
	Pixels := make([]pixelContainer, w*h)
	for i := 0; i < int(w*h); i++ {
		Pixels[i] = pixelContainer{Pixel: pixel{Color: 0, Timestamp: 0, UserID: 0}, mutex: sync.Mutex{}}
	}
	return image{Width: w, Height: h, Pixels: Pixels, Mutex: sync.Mutex{}}
}

func (p *pixelContainer) setColor(color uint8, timestamp int64, userid uint64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if timestamp > p.Pixel.Timestamp {
		p.Pixel.Color = color
		p.Pixel.Timestamp = timestamp
		p.Pixel.UserID = userid
	}
}

func (img *image) SetPixel(message Message) int {
	if message.X >= img.Width || message.Y >= img.Height || message.X < 0 || message.Y < 0 {
		fmt.Printf("User %d tried accessing out of bounds \n", message.UserID)
		return 1
	}
	if message.Color > 15 || message.Color < 0 {
		fmt.Printf("User %d tried setting non existent color %d \n", message.UserID, message.Color)
		return 1
	}
	pos := uint32(message.X)*uint32(img.Width) + uint32(message.Y)
	img.Pixels[pos].setColor(message.Color, message.Timestamp, message.UserID)
	return 0
}

func comparePixels(pixel1 *pixelContainer, pixel2 *pixelContainer) bool {
	return pixel1.Pixel.Color == pixel2.Pixel.Color &&
		pixel1.Pixel.Timestamp == pixel2.Pixel.Timestamp &&
		pixel1.Pixel.UserID == pixel2.Pixel.UserID
}

func (img *image) GetDiff(img2 *image) image {
	diff := GetImage(img.Width, img.Height)
	for i := 0; i < int(img.Width*img.Height); i++ {
		if !comparePixels(&img.Pixels[i], &img2.Pixels[i]) {
			diff.Pixels[i].Pixel.Color = img2.Pixels[i].Pixel.Color
			diff.Pixels[i].Pixel.UserID = img2.Pixels[i].Pixel.UserID
			diff.Pixels[i].Pixel.Timestamp = img2.Pixels[i].Pixel.Timestamp
		}
	}
	return diff
}
