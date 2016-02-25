package couchdb

import (
	"encoding/json"
	"errors"
	//	"reflect"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
)

// Database defintes a database client
type Database struct {
	Name    string
	Client  *Client
	BaseURL string
	Req     Request
}

// NewDB creates an instance of a Database
func NewDB(name string, c *Client) Database {
	protocol := "http"
	if c.Secure {
		protocol += "s"
	}

	port := ""
	if c.Port == 0 {
		port = "80"
	} else {
		port = strconv.Itoa(c.Port)
	}

	url := protocol + "://" + c.Host + ":" + port + "/" + name
	req := Request{gorequest.New(), url}
	if c.Timeout != 0 {
		req.Req.Timeout(c.GetTimeoutDuration())
	}
	if c.Username != "" && c.Password != "" {
		req.Req.SetBasicAuth(c.Username, c.Password)
	}
	return Database{
		Name:    name,
		Client:  c,
		BaseURL: url,
		Req:     req,
	}
}

// Exists check to see if database exists
func (db *Database) Exists() (bool, error) {
	type response struct {
		Error  string `json:"error"`
		DBName string `json:"db_name"`
	}
	_, body, _ := db.Req.Get("").End()
	result := response{}
	pErr := json.Unmarshal([]byte(body), &result)
	if pErr != nil {
		return false, pErr
	}
	if result.DBName != db.Name || result.Error != "" {
		return false, nil
	}

	return true, nil
}

// Create creates a new database
func (db *Database) Create() error {
	type response struct {
		Error string `json:"error"`
		Ok    bool   `json:"ok"`
	}
	_, body, _ := db.Req.Put("").End()
	result := response{}
	pErr := json.Unmarshal([]byte(body), &result)
	log.Info(result)
	if pErr != nil {
		return pErr
	}
	if result.Error != "" {
		return errors.New(result.Error)
	}
	if !result.Ok {
		return errors.New("Couch returned failure when creating [" + db.Name + "]")
	}

	return nil
}

type DocCreateResoponse struct {
	Error string `json:"error"`
	Ok    bool   `json:"ok"`
	Id    string `json:"id"`
	Rev   string `json:"rev"`
}

// Does the document update in couch given a wrapped couch object with DB Exist error status
func (db *Database) updateDocument(err error, data []byte) (error, *DocCreateResoponse) {

	result := &DocCreateResoponse{}

	if err == nil {

		// TODO Fix the errs that are missed while making the request. Its dangerous to ignore.
		_, body, _ := db.Req.Post("").Send(string(data)).End()

		pErr := json.Unmarshal([]byte(body), result)
		log.Info(result)
		if pErr != nil {
			return pErr, result
		}
		if result.Error != "" {
			return errors.New(result.Error), result
		}
		if !result.Ok {
			return errors.New("Couch returned failure when creating [" + db.Name + "]"), result
		}
	} else {
		return err, nil
	}
	return nil, result
}

//Creates a new document to save the data.
func (db *Database) CreateDocument(obj []byte) (error, *DocCreateResoponse) {

	// Here is where the creation takes place.
	couchWrappedObj := NewCouchWrapperCreate(obj)
	data, err := json.Marshal(couchWrappedObj)
	return db.updateDocument(err, data)
}

//We always save the binary of the data.
func (db *Database) UpdateDocument(obj []byte, id string, rev string) (error, *DocCreateResoponse) {

	couchWrappedObj := NewCouchWrapperUpdate(obj)
	couchWrappedObj.Id = id
	couchWrappedObj.Rev = rev
	data, err := json.Marshal(couchWrappedObj)
	return db.updateDocument(err, data)
}

//Retrieve document from the database.
//Will deal with it when the use case comes up, higher up the tree.
func (db *Database) RetrieveDocument(id string) (error, []byte) {
	// Use the get operation to get it.

	// Response, string, error is what End() returns.
	// We need to get the Body out of it excluding all the junk and then unmarshal the data.
	_, body, _ := db.Req.Get(id).End()

	dummyRecv := &CouchWrapperUpdate{}
	err := json.Unmarshal([]byte(body), dummyRecv)

	if err != nil {
		return err, nil
	}

	return nil, dummyRecv.Body
}

// Create View function

//func (db *Database) GetView(viewName string) *View {
// Check if view already exists.

// Create it if it does not. (Map reduce fucnction)
//}

/*
func (db *Database) ReadView(viewName string, fn func()) {
	// Call GetView
	//Return List of objects
	// Otherwise panic

	// Read the view given a View object.
}
*/

// Delete deletes database
func (db *Database) Delete() error {
	type response struct {
		Error string `json:"error"`
		Ok    bool   `json:"ok"`
	}
	_, body, _ := db.Req.Delete("").End()
	result := response{}
	pErr := json.Unmarshal([]byte(body), &result)
	log.Info(result)
	if pErr != nil {
		return pErr
	}
	if result.Error != "" {
		return errors.New(result.Error)
	}
	if !result.Ok {
		return errors.New("Couch returned failure when creating [" + db.Name + "]")
	}

	return nil
}
