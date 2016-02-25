// Representation of a document in couch DB
package couchdb

import (
	"time"
)

type CouchWrapperCreate struct {
	Subject  string
	Author   string
	PostDate time.Time
	Tags     []string
	Body     string //The JSON marshalled representation of the object.
}

type CouchWrapperUpdate struct {
	Id       string `json:"_id"`
	Rev      string `json:"_rev"`
	Subject  string
	Author   string
	PostDate time.Time
	Tags     []string
	Body     string
}

func NewCouchWrapperCreate(obj string) (wrap *CouchWrapperCreate) {
	wrap = &CouchWrapperCreate{
		Subject:  "",
		Author:   "get from user scope",
		PostDate: time.Now(),
		Body:     obj,
	}
	return
}

// Id and rev should be filled by the callee.
func NewCouchWrapperUpdate(obj string) (wrap *CouchWrapperUpdate) {
	wrap = &CouchWrapperUpdate{
		Id:       "",
		Rev:      "",
		Subject:  "",
		Author:   "get from user scope",
		PostDate: time.Now(),
		Body:     obj,
	}
	return
}

/**
Expecting this to last in on a higher level and data is abstract at this level.
func getJsonData(obj interface{}) (err error, jsonData string) {
	data, err := json.Marshal(obj)
	if err != nil {
		return err, string(data)
	}
	jsonData = string(data)
	return
}
*/
