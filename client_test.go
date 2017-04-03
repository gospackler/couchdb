package couchdb

import (
	"testing"
)

const TESTDBNAME = "testdb"
const HOST = "127.0.0.1"
const PORT = 5984

func TestNewDb(t *testing.T) {
	client := NewClient(HOST, PORT)
	dbObj := client.DB(TESTDBNAME)
	err := dbObj.Exists()
	if err == nil {
		t.Log("Test ran without errors: ", TESTDBNAME)
	} else {
		t.Error("Error checking Exists for DB ", err)
	}
}

// DBObject representation of the database under consideration.
var DBObject Database

func TestCreateDb(t *testing.T) {
	client := NewClient(HOST, PORT)
	dbObj := client.DB(TESTDBNAME)
	err := dbObj.Exists()
	if err != nil {
		err = dbObj.Create()
		if err != nil {
			t.Error("Error creating DB ", err)
		}
	}
	DBObject = dbObj
}

/*
func TestDeleteDb(t *testing.T) {

	err := DBObject.Delete()
	if err == nil {
		t.Log("Deleted existing db " + TESTDBNAME + " Successful.")
	} else {
		t.Error("Error deleting "+TESTDBNAME, " ", err)
	}
}
*/
