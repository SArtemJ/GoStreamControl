package libstream

import (
	"sync"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoStorage struct {
	Session    *mgo.Session
	Db         *mgo.Database
	DbName     string
	Collection string
}

type StreamStruct struct {
	ID     bson.ObjectId `bson:"_id" json:"id"`
	Status string        `bson:"status" json:"status,omitempty"`
}

type StreamWitMutex struct {
	Stream StreamStruct
	M      sync.Mutex
}

func NewMongoStorage(uri string, databaseName string) *MongoStorage {
	session, err := mgo.Dial(uri)
	if err != nil {
		Logger.Fatalw("Failer connect to MongoDB server",
			"uri", uri,
			"error", err,
		)
	}
	session.SetPoolLimit(200)
	s := &MongoStorage{
		Session:    session,
		DbName:     databaseName,
		Db:         session.DB(databaseName),
		Collection: "Stream",
	}
	Logger.Debugw("Connected to MongoDB server",
		"uri", uri,
		"database", databaseName,
	)
	return s
}

func (s *MongoStorage) Close() {
	s.Session.Close()
	s.Session = nil
	s.Db = nil
}

func (s MongoStorage) Reset() {
	s.Db.C(s.Collection).RemoveAll(nil)
}

func (s MongoStorage) NewStream() (*StreamWitMutex, bool) {
	var stream StreamWitMutex
	logic := false

	stream.Stream.ID = bson.NewObjectId()
	stream.Stream.Status = "Created"

	if err := s.Db.C(s.Collection).Insert(stream.Stream); err != nil {
		Logger.Debugw("Can't save stream in MongoDB", " - ", err)
		return nil, logic
	} else {
		logic = true
		return &stream, logic
	}
}

func (s MongoStorage) CheckAndReturnStreamInDB(streamID string) (*StreamWitMutex, bool) {
	var stream StreamWitMutex

	objID := bson.ObjectIdHex(streamID)
	if err := s.Db.C(s.Collection).Find(bson.M{"_id": objID}).One(&stream.Stream); err != nil {
		Logger.Debugw("Can't find stream in Mongo with", " id - ", streamID)
		return &stream, false
	}
	return &stream, true
}

func (s MongoStorage) UpdateStream(streamID bson.ObjectId, field string, value interface{}) bool {
	err := s.Db.C(s.Collection).Update(bson.M{"_id": streamID}, bson.M{"$set": bson.M{field: value}})
	if err != nil {
		Logger.Debugw("Can't save Stream in Mongo", err)
		return false
	}
	Logger.Debugw("Status Stream from DB update success", " streamID - ", streamID)
	return true
}

func (sd *StreamWitMutex) UpdateStatus(status string) (string, bool) {
	if sd.Stream.Status != "Finished" {
		sd.M.Lock()
		defer sd.M.Unlock()
		switch status {
		case "a":
			sd.Stream.Status = "Active"
			return "Active", true
		case "i":
			sd.Stream.Status = "Interrupted"
			return "Interrupted", true
		case "f":
			sd.Stream.Status = "Finished"
			return "Finished", true
		default:
			return "", false
		}
	} else {
		return "", false
	}
}

func (s MongoStorage) Remove(id string) bool {
	err := s.Db.C(s.Collection).RemoveId(bson.ObjectIdHex(id))
	if err != nil {
		Logger.Debugw("Can't remove Stream from Mongo", err)
		return false
	}
	return true
}

func (s MongoStorage) SelectAll(pn int, ps int) ([]StreamStruct, bool) {
	var allStreams []StreamStruct
	err := s.Db.C(s.Collection).Find(nil).All(&allStreams)
	if err != nil {
		Logger.Debugw("Can't get all Stream from Mongo", err)
		return nil, false
	}

	if validData, logic := validationPageSize(pn, ps, allStreams); logic {
		return validData, true
	}
	return allStreams, true
}

func validationPageSize(number int, size int, sliceData []StreamStruct) ([]StreamStruct, bool) {
	startFromSlice := number * size
	endFromSlice := size * (number + 1)
	if startFromSlice < len(sliceData) && (endFromSlice < len(sliceData) && endFromSlice > startFromSlice) {
		return sliceData[startFromSlice:endFromSlice], true
	}
	return nil, false
}
