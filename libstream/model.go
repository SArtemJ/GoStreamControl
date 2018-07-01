package libstream

import (
	"github.com/satori/go.uuid"
	"sync"
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
	stream Stream
	mux sync.Mutex
}

func NewStream() *Stream {
	idT, _ := uuid.NewV4()
	idTString := idT.String()

	stream := &Stream{
		ID:     idTString,
		Status: "Created",
	}
	return stream
}

func (sd *StreamData) UpdateStatus(status string) (string, bool) {
	if sd.stream.Status != "Finished" {
		switch status {
		case "a":
			sd.stream.Status = "Active"
			return "Active", true
		case "i":
			sd.stream.Status = "Interrupted"
			return "Interrupted", true
		case "f":
			sd.stream.Status = "Finished"
			return "Finished", true
		default:
			return "", false
		}
	} else {
		return "", false
	}
}
