package couchdb

import (
	"fmt"
	"testing"
)

var client Client

func TestNewClient(t *testing.T) {

	client = NewClient("127.0.0.1", 5984)
	fmt.Println(client.Host)
	fmt.Println(client.Port)
}

func TestNewDb(t *testing.T) {
	dbName := "testdb"
	dbObj := client.DB(dbName)
	status, err := dbObj.Exists()
	if err == nil {
		t.Log("Test ran without errors ", dbName, " --> ", status)
	} else {
		t.Error("Error checking Exists for DB")
	}
}

func TestCreateDb(t *testing.T) {
	dbName := "testdb"
	dbObj := client.DB(dbName)
	status, err := dbObj.Exists()
	if err == nil {
		fmt.Printf("Test ran without errors ", dbName, " --> ", status)
		if status == false {
			t.Log("Db does not exist, so let's create " + dbName)
			err = dbObj.Create()
			if err != nil {
				t.Error("Error creating DB ", err)
			} else {
				t.Log("Successfully created db " + dbName)
			}
		}
	} else {
		t.Error("Error running exists ", err)
	}
}

func TestDeleteDb(t *testing.T) {
	dbName := "testdb"
	dbObj := client.DB(dbName)
	status, err := dbObj.Exists()
	if err == nil {
		if status == true {
			t.Log("Db exists " + dbName)
			err = dbObj.Delete()
			if err == nil {
				t.Log("Deleted existing db " + dbName + " Successful.")
			} else {
				t.Error("Error deleting "+dbName, " ", err)
			}
		}
	}
}
