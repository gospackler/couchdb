package couchdb

import (
	"encoding/json"
	"testing"
)

type TestObj struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

var Id string
var Rev string

func TestCreateDocument(t *testing.T) {

	testObj := &TestObj{
		Name: "Fred",
		Age:  18,
	}

	for i := 0; i < 10; i++ {
		strObj, err := json.Marshal(testObj)
		if err != nil {
			t.Error("Error Marshalling testObj")
		}
		testObj.Age++

		doc := NewDocument("", "", &DBObject)
		err = doc.Create(strObj)
		if err != nil {
			t.Error("Error creating Document ", err)
		} else {
			Id = doc.Id
			Rev = doc.Rev
		}
	}
}

func TestUpdateDocument(t *testing.T) {

	type UpdateObj struct {
		CouchWrapperUpdate
		TestObj
	}

	testObj := &UpdateObj{}

	testObj.Id = Id
	testObj.Rev = Rev

	testObj.Name = "Fred Updated"
	testObj.Age = 25

	objData, err := json.Marshal(testObj)
	if err != nil {
		t.Error("Error Marshalling testObj")
	}

	doc := NewDocument(Id, Rev, &DBObject)
	err = doc.Update(objData)
	if err != nil {
		t.Error("Error Updating Document", err)
	}
}

func TestGetObject(t *testing.T) {

	doc := NewDocument(Id, Rev, &DBObject)
	jsonObj, err := doc.GetDocument()
	if err != nil {
		t.Error("Error ", err)
	}

	type UpdateObj struct {
		CouchWrapperUpdate
		TestObj
	}

	obj := &UpdateObj{}
	json.Unmarshal(jsonObj, obj)
	if obj.Id != Id {
		t.Error("Id should be the same as requested")
	}
}
