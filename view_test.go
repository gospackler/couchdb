package couchdb

import (
	"testing"
)

var desDoc *DesignDoc

func TestMultipleView(t *testing.T) {

	view := NewView("test_view", "doc", "doc.age < 22", "\\\"name\\\", \\\"age\\\"")

	view2 := NewView("raw_view", "", "", "")
	view2.RawStatus = true
	view2.RawJson = "function(rawDoc) {console.log(1234)}"

	fView := NewView("fred_view", "newVar", "newVar.age > 22", "\\\"name\\\", \\\"age\\\"")

	doc := NewDesignDoc("test_design", &DBObject)

	doc.AddView(view)
	doc.AddView(fView)
	doc.AddView(view2)

	err := doc.SaveDoc()
	if err != nil {
		t.Error(err)
	}
	desDoc = doc
}

func getView(key string, t *testing.T) {

	data, err := DBObject.GetView(desDoc.Id, "test_view", key)
	if err != nil {
		t.Error("Error :", err)
	} else {
		t.Log(string(data))
	}
}

func TestGetView(t *testing.T) {
	getView("", t)
}

func TestRetreiveUpdateDesignDoc(t *testing.T) {

	err, desDoc := RetreiveDocFromDb("test_design", &DBObject)
	t.Log(desDoc)
	if err == nil {
		desDoc.Views[0].Name = "test_view_updated"
		desDoc.RevStatus = true
		err := desDoc.SaveDoc()
		if err != nil {
			t.Error(err)
		}
	} else {
		t.Error("Error while updating document")
	}
}

func TestGetKeyFromView(t *testing.T) {

	getView("\"8f752bab6b055d0702563c3672000979\"", t)
}
