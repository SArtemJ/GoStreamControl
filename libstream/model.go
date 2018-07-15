package libstream

import (
	"sync"

	"github.com/satori/go.uuid"
)

// Created
// Active
// Interrupted
// Finished

type Stream struct {
	ID     string `json:"id,omitempty"`
	Status string `json:"status,omitempty"`
}

type StreamData struct {
	M sync.RWMutex
	S Stream
}

func NewStream() StreamData {
	ID, _ := uuid.NewV4()
	IDString := ID.String()

	stream := Stream{
		ID:     IDString,
		Status: "Created",
	}

	streamData := StreamData{}
	streamData.S = stream
	return streamData
}

func (sd *StreamData) UpdateStatus(status string) (string, bool) {
	if sd.S.Status != "Finished" {
		switch status {
		case "a":
			sd.S.Status = "Active"
			return "Active", true
		case "i":
			sd.S.Status = "Interrupted"
			return "Interrupted", true
		case "f":
			sd.S.Status = "Finished"
			return "Finished", true
		default:
			return "", false
		}
	} else {
		return "", false
	}
}
