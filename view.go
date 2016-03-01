//This is the space where all the view Related functions for couch go in

// View is part of a table. Views are made making use of the map reduce fuctions.
// The id of the view should start with _design/

// DesignDoc is needed with a name -> corresponds to the design document to query.
// Eash DesignDoc has a set of Views. -> Each view is a map reduce Function or raw javascript code.
package couchdb

import (
	"bytes"
	"text/template"
)

type View struct {
	Name         string
	VariableName string
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
	return
}

type DesignDoc struct {
	Id       string //The id of the document
	Rev      string
	Views    []*View
	LastView *View
	Db       *Database
}

func NewDesignDoc(id string, db *Database) (doc *DesignDoc) {

	doc = &DesignDoc{
		Id: id,
		Db: db,
	}
	return
}

func (doc *DesignDoc) AddView(view *View) {

	if doc.LastView == nil {
		doc.LastView = view
	} else {
		doc.Views = append(doc.Views, doc.LastView)
		doc.LastView = view
	}
}

func (doc *DesignDoc) CreateDoc(tempFile string) (error, []byte) {

	tmpl, err := template.ParseFiles(tempFile)
	if err != nil {
		return err, nil
	}
	buffer := &bytes.Buffer{}
	err = tmpl.Execute(buffer, doc)
	return err, buffer.Bytes()
}

func (doc *DesignDoc) SaveDoc() (err error) {

	dbDoc := NewDocument("", "", doc.Db)
	err, data := doc.CreateDoc("templates/design.js")
	if err == nil {
		err = dbDoc.Create(data)
	}
	return
}
