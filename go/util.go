package main

import (
	"encoding/json"
)

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
