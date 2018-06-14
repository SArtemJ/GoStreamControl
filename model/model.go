package model

import "github.com/satori/go.uuid"

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

func (mS *Stream) Active() {
	mS.Status = "Active"
}

func (mS *Stream) Interrupt() {
	mS.Status = "Interrupted"
}

func (mS *Stream) Finish() {
	mS.Status = "Finished"
}

func (mS *Stream) Delete(uuidT uuid.UUID) {
	//
}