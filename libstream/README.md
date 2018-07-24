# Stream
API для управления stream

	s.Router.HandleFunc("/s", s.ShowAllStreams).Methods("GET")
	s.Router.HandleFunc("/run", s.StartNewStream).Methods("GET")
	s.Router.HandleFunc("/activate/{id}", s.ActivateStream).Methods("PATCH")
	s.Router.HandleFunc("/interrupt/{id}", s.InterruptStream).Methods("PATCH")
	s.Router.HandleFunc("/finish/{id}", s.FinishStream).Methods("PATCH")
	s.Router.HandleFunc("/delete/{id}", s.DeleteStream).Methods("DELETE")

Запросы:
- GET показать все stream из Mongo
    - http://localhost:8888/s
- GET создает новый stream
	- http://localhost:8888/run
- PATCH меняет статус stream на Active
	- http://localhost:8888/activate/{id}
- PATCH меняет статус stream на Interrupted
	- http://localhost:8888/interrupt/{id}
- PATCH меняет статус stream на Finished
	- http://localhost:8888/finish/{id}
- DELETE удаляет stream
	- http://localhost:8888/delete/{id}?rt=roottoken


# Сборка Docker
docker-compose up
