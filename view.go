// This is the space where all the view Related functions for couch go in

// View is part of a table. Views are made making use of the map reduce fuctions.
// The id of the view should start with _design/

// DesignDoc is needed with a name -> corresponds to the design document to query.
// Eash DesignDoc has a set of Views. -> Each view is a map reduce Function or raw javascript code.
package couchdb

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type View struct {
	Name         string
	VariableName string
	CondStatus   bool
	Condition    string
	EmitStr      string
	RawJson      string
	RawStatus    bool
}

func NewView(name string, varName string, condition string, emitStr string) (view *View) {

	view = &View{
		Name:         name,
		VariableName: varName,
		Condition:    condition,
		EmitStr:      emitStr,
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

	type ViewObj struct {
		Map    string `json:"map"`
		Reduce string `json:"reduce"`
	}

	type TempRetrieve struct {
		Id       string             `json:"_id"`
		Rev      string             `json:"_rev"`
		Language string             `json:"language"`
		Views    map[string]ViewObj `json:"views"`
	}

	tempRet := &TempRetrieve{}
	doc := NewDocument("_design/"+id, "", db)
	data, err := doc.GetDocument()
	if err == nil {
		fmt.Println("Data Read ", string(data))
		err = json.Unmarshal(data, tempRet)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Unmarshalled Json", tempRet)
			desDoc = &DesignDoc{} // allocating memory
			desDoc.Id = tempRet.Id
			desDoc.Db = db
			desDoc.Rev = tempRet.Rev
			for viewName, Data := range tempRet.Views {
				view := &View{}
				view.Name = viewName
				view.RawStatus = true
				view.RawJson = Data.Map
				desDoc.AddView(view)
			}
		}

	} else {
		fmt.Println(err)
	}
	return err, desDoc
}

// Test the DB for the revision of the document.
func (desDoc *DesignDoc) getRev(doc *Document) (error, string) {
	result := &DocCreateResoponse{}
	docBytes, err := doc.GetDocument()
	if err != nil {
		return err, ""
	}

	err = json.Unmarshal(docBytes, result)

	if err != nil {
		return err, ""
	}

	fmt.Println("Result", string(docBytes))
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

	fmt.Println("Checking the existance of ", viewName, " in ", doc.Id)
	index := 0
	for index, view := range doc.Views {
		if viewName == view.Name {
			return index, true
		}
	}
	if doc.LastView != nil {
		fmt.Println(doc.LastView)
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
		fmt.Println("Could not find revision ", err)
	}

	err, data := doc.CreateDoc()

	fmt.Println("Trying to create \n \n", string(data))

	if err == nil {
		err = dbDoc.Create(data)
	}
	return
}
