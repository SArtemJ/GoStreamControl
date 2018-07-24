# Stream
API для управления stream
##### Для запуска параметры в stream.yaml

```
#запуск на локальной машине localhost:8888 
#запуск в docker 0.0.0.0:8888

server.addr: "0.0.0.0:8888"
server.apiPrefix: "/api/v1/"

#запуск на локальной машине localhost
#запуск в docker gostreamcontrolapi_db_1

storage.address: "gostreamcontrolapi_db_1"
storage.name: "Stream"

#время ожидания
#если stream переключили в состояние interruped после истечения timer.value в минутах
#stream переводится в состояние finished

timer.value: 1

#root token для удаления stream

r.t: "!)asd!567z"
```


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
