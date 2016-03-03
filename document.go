// This is where all the code with respect to the documents go in.
package couchdb

import (
	"encoding/json"
	"errors"
	log "github.com/Sirupsen/logrus"
)

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

		result := &DocCreateResoponse{}
		pErr := json.Unmarshal([]byte(body), result)
		if pErr != nil {
			return pErr
		}

		if result.Error != "" {
			return errors.New(result.Error)
		}

	}
	return nil
}

// Does the document update in couch given a wrapped couch object with DB Exist error status
func (doc *Document) createOrUpdate(data []byte) (error, *DocCreateResoponse) {

	// TODO Fix the errs that are missed while making the request. Its dangerous to ignore.
	_, body, _ := doc.Db.Req.Post("").Send(string(data)).End()

	result := &DocCreateResoponse{}
	pErr := json.Unmarshal([]byte(body), result)
	log.Info(result)
	if pErr != nil {
		return pErr, result
	}
	if result.Error != "" {
		return errors.New(result.Error), result
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
		// do the update operation.
		err, _ = doc.createOrUpdate(data)
	}
	return
}

// Gets the document using the given id and error if it does not exist.
func (doc *Document) GetDocument() ([]byte, error) {

	err := doc.Exists()
	if err == nil {
		// Use the get operation to get it.
		_, body, _ := doc.Db.Req.Get(doc.Id).End()
		return []byte(body), nil
	} else {
		return nil, err
	}
}
