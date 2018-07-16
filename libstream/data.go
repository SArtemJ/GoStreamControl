package libstream

import (
	"database/sql"
)

var sd = &StreamData{}

func SelectAll(pn int, ps int) ([]Stream, bool) {
	var allStreams []Stream

	rows, err := DB.Query("SELECT * FROM stream")
	if err != nil {
		Logger.Debugw("Operation select from DB.Query", "err ", err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var s Stream
		err := rows.Scan(&s.ID, &s.Status)
		if err != nil {
			Logger.Infow("Operation select from rows.Scan", "err ", err.Error())
		}
		allStreams = append(allStreams, s)
	}

	if validData, logic := validationPageSize(pn, ps, allStreams); logic {
		return validData, true
	}
	return nil, false
}

func InsertToDB(sd *StreamData) bool {
	sd.M.Lock()
	defer sd.M.Unlock()
	stringQ := "INSERT INTO stream (id, status) VALUES ($1, $2)"
	_, err := DB.Exec(stringQ, &sd.S.ID, &sd.S.Status)
	if err != nil {
		Logger.Debugw("Operation insert to DB", "err ", err.Error())
		return false
	}
	return true
}

func DeleteFromDB(sd *StreamData) bool {
	sd.M.Lock()
	defer sd.M.Unlock()
	stringQ := "DELETE FROM stream WHERE id = $1"
	_, err := DB.Exec(stringQ, sd.S.ID)
	if err != nil {
		Logger.Debugw("Operation delete from DB", "err ", err.Error())
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

	sd := StreamData{}
	stringQ := "SELECT * FROM stream WHERE id = $1"
	row := DB.QueryRow(stringQ, streamID)
	switch err := row.Scan(&sd.S.ID, &sd.S.Status); err {
	case sql.ErrNoRows:
		Logger.Debugw("No stream in DB", "err ", err.Error())
		return &sd, false
	case nil:
		return &sd, true
	default:
		Logger.Debugw("No stream in DB", "err ", err.Error())
		return &sd, false
	}

}

func UpdateRow(sd *StreamData) bool {
	sd.M.Lock()
	defer sd.M.Unlock()

	stringQ := "UPDATE stream SET status = $2 WHERE id = $1"
	_, err := DB.Exec(stringQ, sd.S.ID, sd.S.Status)
	if err != nil {
		Logger.Debugw("Operation update stream in DB", "err ", err.Error())
		return false
	}
	return true
}
