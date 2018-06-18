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
	DBUser     = "testu"
	DBpassword = "testup"
	DBName     = "stream"
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

func SelectAll(pn int, ps int) ([]m.Stream, bool) {
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
		log.Println("Get All\n")
		log.Println(len(allStreams))
	}

	if validData, logic := validationPageSize(pn, ps, allStreams); logic {
		return validData, true
	}
	return nil, false
}

func InsertToDB(s *m.Stream) bool {
	stringQ := "INSERT INTO stream (id, status) VALUES ($1, $2)"
	_, err := DB.Exec(stringQ, s.ID, s.Status)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func DeleteFromDB(s *m.Stream) bool {
	stringQ := "DELETE FROM stream WHERE id = $1"
	_, err := DB.Exec(stringQ, s.ID)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func validationPageSize(number int, size int, sliceData []m.Stream) ([]m.Stream, bool) {
	startFromSlice := number * size
	endFromSlice := size * (number + 1)
	if startFromSlice < len(sliceData) && (endFromSlice < len(sliceData) && endFromSlice > startFromSlice) {
		return sliceData[startFromSlice:endFromSlice], true
	}
	return nil, false
}

func CheckFromDB(streamID string) (m.Stream, bool) {

	var stream m.Stream
	stringQ := "SELECT * FROM stream WHERE id = $1"

	log.Println(stringQ)

	row := DB.QueryRow(stringQ, streamID)
	switch err := row.Scan(&stream.ID, &stream.Status); err {
	case sql.ErrNoRows:
		return stream, false
	case nil:
		return stream, true
	default:
		log.Println(err.Error())
		return stream, false
	}
}

func UpdateRow(s m.Stream) bool {
	stringQ := "UPDATE stream SET status = $2 WHERE id = $1"
	_, err := DB.Exec(stringQ, s.ID, s.Status)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
