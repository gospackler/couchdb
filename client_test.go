package couchdb

import (
	"encoding/json"
	"fmt"
	"testing"
)

const TESTDBNAME = "testdb"
const HOST = "127.0.0.1"
const PORT = 5984

func TestNewDb(t *testing.T) {
	client := NewClient(HOST, PORT)
	dbObj := client.DB(TESTDBNAME)
	status, err := dbObj.Exists()
	if err == nil {
		t.Log("Test ran without errors: ", TESTDBNAME, " exists is ", status)
	} else {
		t.Error("Error checking Exists for DB")
	}
}

// DBObject representation of the database under consideration.
var DBObject Database

func TestCreateDb(t *testing.T) {
	client := NewClient(HOST, PORT)
	dbObj := client.DB(TESTDBNAME)
	status, err := dbObj.Exists()
	if err == nil {
		fmt.Printf("Test ran without errors ", TESTDBNAME, " --> ", status)
		if status == false {
			t.Log("Db does not exist, so let's create " + TESTDBNAME)
			err = dbObj.Create()
			if err != nil {
				t.Error("Error creating DB ", err)
			} else {
				t.Log("Successfully created db " + TESTDBNAME)
			}
		}
	} else {
		t.Error("Error running exists ", err)
	}
	DBObject = dbObj
}

type TestObj struct {
	Name string
	Age  int
}

var Id string
var Rev string

func TestCreateDocument(t *testing.T) {

	testObj := &TestObj{
		Name: "Fred",
		Age:  18,
	}

	strObj, err := json.Marshal(testObj)
	t.Log("Saving Object:", strObj)
	if err != nil {
		t.Error("Error Marshalling testObj")
	}

	err, status := DBObject.CreateDocument(strObj)
	if err != nil {
		t.Error("Error creating Document ", err)
	} else {
		t.Log("Successfully created document " + TESTDBNAME)
		t.Log("Document Id " + status.Id)
		t.Log("Document Revision " + status.Rev)
		Id = status.Id
		Rev = status.Rev
	}
}

func TestUpdateDocument(t *testing.T) {

	testObj := &TestObj{
		Name: "Gordon",
		Age:  28,
	}

	strObj, err := json.Marshal(testObj)
	if err != nil {
		t.Error("Error Marshalling testObj")
	}

	err, status := DBObject.UpdateDocument(strObj, Id, Rev)
	if err != nil {
		t.Error("Error Updating Document", err)
	} else {
		t.Log("Successfully updated document")
		t.Log("Document Id " + status.Id)
		t.Log("Revision " + status.Rev)
	}
}

func TestGetObject(t *testing.T) {

	err, jsonObj := DBObject.RetrieveDocument(Id)
	if err != nil {
		t.Error("Error ", err)
	}

	t.Log(string(jsonObj))
}
func TestDeleteDb(t *testing.T) {

	err := DBObject.Delete()
	if err == nil {
		t.Log("Deleted existing db " + TESTDBNAME + " Successful.")
	} else {
		t.Error("Error deleting "+TESTDBNAME, " ", err)
	}
}
