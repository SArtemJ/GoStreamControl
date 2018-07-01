package libstream

import (
	"database/sql"
	"log"
)

var sd = &StreamData{}

func SelectAll(pn int, ps int) ([]Stream, bool) {
	var allStreams []Stream

	rows, err := DB.Query("SELECT * FROM stream")
	if err != nil {
		log.Println(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var s Stream
		err := rows.Scan(&s.ID, &s.Status)
		if err != nil {
			log.Println(err.Error())
		}
		allStreams = append(allStreams, s)
	}

	if validData, logic := validationPageSize(pn, ps, allStreams); logic {
		return validData, true
	}
	return nil, false
}

func InsertToDB(s *Stream) bool {
	stringQ := "INSERT INTO stream (id, status) VALUES ($1, $2)"
	_, err := DB.Exec(stringQ, s.ID, s.Status)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func DeleteFromDB(s *Stream) bool {
	stringQ := "DELETE FROM stream WHERE id = $1"
	_, err := DB.Exec(stringQ, s.ID)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}

func validationPageSize(number int, size int, sliceData []Stream) ([]Stream, bool) {
	startFromSlice := number * size
	endFromSlice := size * (number + 1)
	if startFromSlice < len(sliceData) && (endFromSlice < len(sliceData) && endFromSlice > startFromSlice) {
		return sliceData[startFromSlice:endFromSlice], true
	}
	return nil, false
}

func CheckFromDB(streamID string) (*StreamData, bool) {
	sd.mux.Lock()
	defer sd.mux.Unlock()
	stringQ := "SELECT * FROM stream WHERE id = $1"

	log.Println(stringQ)

	row := DB.QueryRow(stringQ, streamID)
	switch err := row.Scan(&sd.stream.ID, &sd.stream.Status); err {
	case sql.ErrNoRows:
		return sd, false
	case nil:
		return sd, true
	default:
		log.Println(err.Error())
		return sd, false
	}

}

func UpdateRow(sd *StreamData) bool {
	sd.mux.Lock()
	defer  sd.mux.Unlock()

	stringQ := "UPDATE stream SET status = $2 WHERE id = $1"
	_, err := DB.Exec(stringQ, sd.stream.ID, sd.stream.Status)
	if err != nil {
		log.Println(err.Error())
		return false
	}
	return true
}
