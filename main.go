package main

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	db "github.com/SArtemJ/GoStreamControlAPI/Database"
	m "github.com/SArtemJ/GoStreamControlAPI/Model"
)

//show all
func ShowAllStreams(w http.ResponseWriter, r *http.Request) {
	allStream := db.SelectAll()
	allStreamJSON, _ := json.Marshal(allStream)
	w.Write(allStreamJSON)
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

}

//set interrupted
func InterruptStream(w http.ResponseWriter, r *http.Request) {

}

//set finished
func FinishStream(w http.ResponseWriter, r *http.Request) {

}

func main() {

	router := mux.NewRouter()
	sub := router.PathPrefix("/api/v1").Subrouter()
	sub.HandleFunc("/s", ShowAllStreams).Methods("GET")
	sub.HandleFunc("/run", StartNewStream).Methods("GET")
	sub.HandleFunc("/activate", ActivateStream).Methods("POST")
	sub.HandleFunc("/interrupt", InterruptStream).Methods("POST")
	sub.HandleFunc("/finish", FinishStream).Methods("POST")
	http.ListenAndServe(":8000", router)
}
