package model

import (
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

func NewStream() *Stream {
	idT, _ := uuid.NewV4()
	idTString := idT.String()

	stream := &Stream{
		ID:     idTString,
		Status: "Created",
	}
	return stream
}

func (mS *Stream) UpdateStatus(status string) (string, bool) {
	if mS.Status != "Finished" {
		switch status {
		case "a":
			mS.Status = "Active"
			return "Active", true
		case "i":
			mS.Status = "Interrupted"
			return "Interrupted", true
		case "f":
			mS.Status = "Finished"
			return "Finished", true
		default:
			return "", false
		}
	} else {
		return "", false
	}
}

func (mS *Stream) Delete(uuidT uuid.UUID) {
	//
}
