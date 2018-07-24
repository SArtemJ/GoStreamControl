package libstream

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

type StreamServer struct {
	Address   string
	APIPrefix string
	RootToken string

	Timer   *time.Timer
	Router  *mux.Router
	Storage *MongoStorage
}

type ServerConfig struct {
	Address    string
	ApiPrefix  string
	RootToken  string
	Storage    *MongoStorage
	TimerValue int
}

func NewServer(config ServerConfig) *StreamServer {
	if config.ApiPrefix == "" {
		config.ApiPrefix = "/api/v1/"
	}
	if config.Address == "" {
		config.Address = "/"
	}
	if config.RootToken == "" {
		config.RootToken = "!csdf!25"
	}
	if config.TimerValue == 0 {
		config.TimerValue = 1
	}

	server := &StreamServer{
		Address:   config.Address,
		APIPrefix: config.ApiPrefix,
		RootToken: config.RootToken,
		Timer:     time.NewTimer(time.Minute * time.Duration(config.TimerValue)),
		Router:    mux.NewRouter(),
		Storage:   config.Storage,
	}
	server.SetupRouter()
	return server
}

func (s *StreamServer) SetupRouter() {
	s.Router = s.Router.PathPrefix(s.APIPrefix).Subrouter()
	Logger.Debugf(`API endpoint "%s"`, s.APIPrefix)
	Logger.Debugf(`Root token for delete "%s"`, s.RootToken)

	s.Router.HandleFunc("/s", s.ShowAllStreams).Methods("GET")
	s.Router.HandleFunc("/run", s.StartNewStream).Methods("GET")
	s.Router.HandleFunc("/activate/{id}", s.ActivateStream).Methods("PATCH")
	s.Router.HandleFunc("/interrupt/{id}", s.InterruptStream).Methods("PATCH")
	s.Router.HandleFunc("/finish/{id}", s.FinishStream).Methods("PATCH")
	s.Router.HandleFunc("/delete/{id}", s.DeleteStream).Methods("DELETE")

}

func (s *StreamServer) GetRouter() *mux.Router {
	return s.Router
}

func (s *StreamServer) Run() {
	Logger.Infof(`Stream server started on "%s"`, s.Address)
	http.ListenAndServe(s.Address, s.Router)
}

//show all
func (s *StreamServer) ShowAllStreams(w http.ResponseWriter, r *http.Request) {
	pn, _ := strconv.Atoi(r.URL.Query().Get("page[number]"))
	ps, _ := strconv.Atoi(r.URL.Query().Get("page[size]"))

	if allStream, success := s.Storage.SelectAll(pn, ps); success {
		allStreamJSON, _ := json.Marshal(allStream)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		w.Write(allStreamJSON)
	} else {
		w.WriteHeader(http.StatusNoContent)
		io.WriteString(w, "Parameters under the limit")
	}
}

//start new -- created
func (s *StreamServer) StartNewStream(w http.ResponseWriter, r *http.Request) {
	if stream, ok := s.Storage.NewStream(); ok == true {
		Logger.Debug(`New stream created with id `, stream.Stream.ID)
		streamJSON, _ := json.Marshal(stream.Stream)
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusCreated)
		w.Write(streamJSON)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Parameters under the limit")
	}
}

//set active
func (s *StreamServer) ActivateStream(w http.ResponseWriter, r *http.Request) {
	stream := mux.Vars(r)["id"]
	s.UpdateStream(w, stream, "a")
	s.Timer.Stop()
}

//set interrupted
func (s *StreamServer) InterruptStream(w http.ResponseWriter, r *http.Request) {
	stream := mux.Vars(r)["id"]
	s.UpdateStream(w, stream, "i")

	go s.finishByTimer(w, stream)
}

//set finished
func (s *StreamServer) FinishStream(w http.ResponseWriter, r *http.Request) {
	stream := mux.Vars(r)["id"]
	s.UpdateStream(w, stream, "f")
}

//set finished
func (s *StreamServer) DeleteStream(w http.ResponseWriter, r *http.Request) {
	stream := mux.Vars(r)["id"]
	token := r.URL.Query().Get("rt")
	if token == s.RootToken {
		if s.Storage.Remove(stream) == true {
			w.WriteHeader(http.StatusOK)
			resultString := `Stream  with id ` + stream + ` was deleted by Admin`
			io.WriteString(w, resultString)
		}
	}

}

func (s *StreamServer) UpdateStream(w http.ResponseWriter, streamID string, status string) {
	if sDB, check := s.Storage.CheckAndReturnStreamInDB(streamID); check {
		if name, success := sDB.UpdateStatus(status); success {
			if s.Storage.UpdateStream(bson.ObjectIdHex(streamID), "status", sDB.Stream.Status) == true {
				w.WriteHeader(http.StatusOK)
				resultString := "Stream status update on " + name
				io.WriteString(w, resultString)
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			io.WriteString(w, `Stream can't change status with current - `+sDB.Stream.Status)
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Stream update - error - check stream ID or Status name")
	}
}

func (s *StreamServer) finishByTimer(w http.ResponseWriter, streamID string) {
	<-s.Timer.C
	s.UpdateStream(w, streamID, "f")
}
