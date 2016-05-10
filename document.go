// This is where all the code with respect to the documents go in.
package couchdb

import (
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
)

// Document either has an Id, Rev and the DB it connects to.
type Document struct {
	Db  *Database
	Id  string `json:"_id"`
	Rev string `json:"_rev"`
}

func NewDocument(id string, rev string, Db *Database) *Document {
	return &Document{
		Db:  Db,
		Id:  id,
		Rev: rev,
	}
}

type DocCreateResoponse struct {
	Error string `json:"error"`
	Ok    bool   `json:"ok"`
	Id    string `json:"id"`
	Rev   string `json:"rev"`
}

//Function checks if the document exists and returns error if it does not
func (doc *Document) Exists() error {

	// Use the get operation to get it.
	_, body, errs := doc.Db.Req.Get(doc.Id).End()

	if len(errs) != 0 {
		//TODO Check other errors if any exists and make one error to return
		return errs[0]
	} else {

		result := &struct {
			Error string `json:"error"`
			Ok    bool   `json:"ok"`
			Id    string `json:"_id"`
			Rev   string `json:"_rev"`
		}{}
		pErr := json.Unmarshal([]byte(body), result)
		if pErr != nil {
			return pErr
		}

		if result.Error != "" {
			return errors.New(result.Error)
		}
		doc.Id = result.Id
		doc.Rev = result.Rev
	}
	return nil
}

// Does the document update in couch given a wrapped couch object with DB Exist error status
func (doc *Document) createOrUpdate(data []byte) (error, *DocCreateResoponse) {

	// TODO Fix the errs that are missed while making the request. Its dangerous to ignore.
	_, body, _ := doc.Db.Req.Post("").Send(string(data)).End()

	result := &DocCreateResoponse{}
	pErr := json.Unmarshal([]byte(body), result)
	log.Info("couch : createOrUpdate json resp:", body)
	log.Info(result)
	if pErr != nil {
		return pErr, result
	}
	if result.Error != "" {
		return errors.New("Failure while creating " + result.Error), result
	}
	if !result.Ok {
		return errors.New("Couch returned failure when creating [" + doc.Db.Name + "]"), result
	}
	return nil, result
}

// Creates a document if it does not already exist and generates an error if it already exists.
func (doc *Document) Create(data []byte) (err error) {
	err, docResp := doc.createOrUpdate(data)

	if err != nil {
		return
	}

	doc.Id = docResp.Id
	doc.Rev = docResp.Rev
	return
}

// Updates the document with the new Data.
// Data contains an encoded marshalled object that has the required fields, pre computed..
func (doc *Document) Update(data []byte) (err error) {

	err = doc.Exists()

	if err == nil {

		var val map[string]interface{}

		err = json.Unmarshal(data, &val)
		if err != nil {
			errors.New("Unmarshalling error " + err.Error())
		}

		val["_rev"] = doc.Rev
		// do the update operation.
		newData, err := json.Marshal(val)
		if err != nil {
			return errors.New("Marshalling error of new value in couch " + err.Error())
		}
		err, _ = doc.createOrUpdate(newData)
	}
	return
}

func (doc *Document) getDocFromId() ([]byte, error) {

	err := doc.Exists()
	if err == nil {
		// Use the get operation to get it.
		_, body, _ := doc.Db.Req.Get(doc.Id).End()
		return []byte(body), nil
	} else {
		return nil, err
	}
}

// Gets the document using the given id and error if it does not exist.
func (doc *Document) GetDocument() ([]byte, error) {

	if doc.Id != "" {
		return doc.getDocFromId()
	}

	return nil, errors.New("An id required to search for the document.")
}
