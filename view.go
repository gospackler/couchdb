// This is the space where all the view Related functions for couch go in

// View is part of a table. Views are made making use of the map reduce fuctions.
// The id of the view should start with _design/

// DesignDoc is needed with a name -> corresponds to the design document to query.
// Eash DesignDoc has a set of Views. -> Each view is a map reduce Function or raw javascript code.
package couchdb

import (
	"bytes"
	"encoding/json"

	log "github.com/Sirupsen/logrus"
)

type View struct {
	Name         string
	VariableName string
	KeyName      string
	CondStatus   bool
	Condition    string
	EmitStr      string
	RawJson      string
}

func NewView(name string, varName string, condition string, emitStr string) (view *View) {

	view = &View{
		Name:         name,
		VariableName: varName,
		Condition:    condition,
		EmitStr:      emitStr,
		CondStatus:   true,
		KeyName:      "",
	}

	if condition == "" {
		view.CondStatus = false
	}
	return
}

type DesignDoc struct {
	Id        string //The id of the document
	Rev       string
	Views     []*View
	LastView  *View
	Db        *Database
	RevStatus bool
}

func NewDesignDoc(id string, db *Database) (doc *DesignDoc) {

	doc = &DesignDoc{
		Id: "_design/" + id,
		Db: db,
	}
	return
}

func RetreiveDocFromDb(id string, db *Database) (err error, desDoc *DesignDoc) {

	type TempRetrieve struct {
		Id       string                     `json:"_id"`
		Rev      string                     `json:"_rev"`
		Language string                     `json:"language"`
		Views    map[string]json.RawMessage `json:"views"`
	}

	tempRet := &TempRetrieve{}
	doc := NewDocument("_design/"+id, "", db)
	data, err := doc.GetDocument()
	if err == nil {
		err = json.Unmarshal(data, tempRet)
		if err != nil {
			log.Warn(err)
		} else {
			desDoc = &DesignDoc{} // allocating memory
			desDoc.Id = tempRet.Id
			desDoc.Db = db
			desDoc.Rev = tempRet.Rev
			for viewName, Data := range tempRet.Views {
				view := &View{}
				view.Name = viewName
				// Removing the prefix and suffix.
				view.RawJson = string(Data)[1 : len(string(Data))-1]
				desDoc.AddView(view)
			}
		}

	} else {
		log.Warn(err)
	}
	return err, desDoc
}

// Test the DB for the revision of the document.
func (desDoc *DesignDoc) getRev(doc *Document) (error, string) {

	type GetDocResp struct {
		Error string `json:"error"`
		Ok    bool   `json:"ok"`
		Id    string `json:"_id"`
		Rev   string `json:"_rev"`
	}

	result := new(GetDocResp)
	docBytes, err := doc.GetDocument()
	if err != nil {
		return err, ""
	}

	err = json.Unmarshal(docBytes, result)

	if err != nil {
		return err, ""
	}

	return nil, result.Rev
}

func (doc *DesignDoc) AddView(view *View) {

	if doc.LastView == nil {
		doc.LastView = view
	} else {
		doc.Views = append(doc.Views, doc.LastView)
		doc.LastView = view
	}
}

// Returning index as -1 => LastView has it.
// Otherwise the index returned is in the view.
// status returns true or false indicating presence.
func (doc *DesignDoc) CheckExists(viewName string) (int, bool) {

	index := 0
	for index, view := range doc.Views {
		if viewName == view.Name {
			return index, true
		}
	}
	if doc.LastView != nil {
		if viewName == doc.LastView.Name {
			return -1, true
		}
	}
	return index, false
}

// Works on the default Global value and not on config files.
func (doc *DesignDoc) CreateDoc() (error, []byte) {

	buffer := &bytes.Buffer{}
	err := DESIGNTMPL.Execute(buffer, doc)
	return err, buffer.Bytes()
}

func (doc *DesignDoc) SaveDoc() (err error) {

	dbDoc := NewDocument(doc.Id, "", doc.Db)
	err, rev := doc.getRev(dbDoc)
	if err == nil {
		// The document already exist
		doc.Rev = rev
		doc.RevStatus = true
	} else {
		log.Warn("Could not find revision ", err)
	}

	err, data := doc.CreateDoc()

	if err == nil {
		err = dbDoc.Create(data)
	}
	return
}
