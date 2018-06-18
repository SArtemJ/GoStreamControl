package main

import (
	"encoding/json"
	db "github.com/SArtemJ/GoStreamControlAPI/database"
	m "github.com/SArtemJ/GoStreamControlAPI/model"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"strconv"
)

//show all
func ShowAllStreams(w http.ResponseWriter, r *http.Request) {

	pn, _ := strconv.Atoi(r.URL.Query().Get("page[number]"))
	ps, _ := strconv.Atoi(r.URL.Query().Get("page[size]"))

	if allStream, success := db.SelectAll(pn, ps); success {
		allStreamJSON, _ := json.Marshal(allStream)
		w.WriteHeader(http.StatusOK)
		w.Write(allStreamJSON)
	} else {
		w.WriteHeader(http.StatusNoContent)
		io.WriteString(w, "Parameters under the limit")
	}
}

//start new -- created
func StartNewStream(w http.ResponseWriter, r *http.Request) {
	s := m.NewStream()
	if db.InsertToDB(s) {
		streamJSON, _ := json.Marshal(s)
		w.WriteHeader(http.StatusCreated)
		w.Write(streamJSON)
	}
}

//set active
func ActivateStream(w http.ResponseWriter, r *http.Request) {
	stream := mux.Vars(r)["id"]
	UpdateStream(w, stream, "a")
}

//set interrupted
func InterruptStream(w http.ResponseWriter, r *http.Request) {
	stream := mux.Vars(r)["id"]
	UpdateStream(w, stream, "i")
}

//set finished
func FinishStream(w http.ResponseWriter, r *http.Request) {
	stream := mux.Vars(r)["id"]
	UpdateStream(w, stream, "f")
}

func main() {

	router := mux.NewRouter()
	sub := router.PathPrefix("/api/v1").Subrouter()
	sub.HandleFunc("/s", ShowAllStreams).Methods("GET")
	sub.HandleFunc("/run", StartNewStream).Methods("GET")
	sub.HandleFunc("/activate/{id}", ActivateStream).Methods("PATCH")
	sub.HandleFunc("/interrupt/{id}", InterruptStream).Methods("PATCH")
	sub.HandleFunc("/finish/{id}", FinishStream).Methods("PATCH")
	http.ListenAndServe(":8000", router)
}

func UpdateStream(w http.ResponseWriter, id string, status string) {

	if sDB, check := db.CheckFromDB(id); check {
		if name, success := sDB.UpdateStatus(status); success {
			w.WriteHeader(http.StatusOK)
			resultString := "Stream status update on " + name
			db.UpdateRow(sDB)
			io.WriteString(w, resultString)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Stream update - error - check stream ID or Status name")
	}
}
