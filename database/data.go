package database

import (
	"database/sql"
	"fmt"
	"log"
	m "github.com/SArtemJ/GoStreamControlAPI/model"
	_ "github.com/lib/pq"
)

var (
	DB *sql.DB
)

const (
	DBUser = "testu"
	DBpassword = "testup"
	DBName = "stream"
)

func init() {
	var err error
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DBUser, DBpassword, DBName)
	DB, err = sql.Open("postgres", dbinfo)
	if err != nil {
		panic(err)
	}
	//defer DB.Close()
}

func SelectAll() []m.Stream {

	var allStreams []m.Stream
	rows, err := DB.Query("SELECT * FROM stream")
	if err != nil {
		log.Println(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var s m.Stream
		err := rows.Scan(&s.ID, &s.Status)
		if err != nil {
			log.Println(err.Error())
		}
		allStreams = append(allStreams, s)
	}

	return allStreams
}

func InsertToDB(t *m.Stream) bool {
	var stringQ = "INSERT INTO stream (id, status) VALUES ($1, $2)"
	_, err := DB.Exec(stringQ, t.ID, t.Status)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
