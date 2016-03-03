package couchdb

import (
	"testing"
)

func TestMultipleView(t *testing.T) {

	view := NewView("test_view", "doc", "doc.age < 22", "doc.name, doc.age")

	view2 := NewView("raw_view", "", "", "")
	view2.RawStatus = true
	view2.RawJson = "function(rawDoc) {console.log(1234)}"

	fView := NewView("fred_view", "newVar", "newVar.age > 22", "newVar.name, newVar.age")

	doc := NewDesignDoc("test_design", &DBObject)

	doc.AddView(view)
	doc.AddView(fView)
	doc.AddView(view2)

	err := doc.SaveDoc()
	if err != nil {
		t.Error(err)
	}
}

func TestGetView(t *testing.T) {
	err, data := DBObject.GetView("test_design", "test_view")
	if err != nil {
		t.Error("Error :", err)
	} else {
		t.Log(string(data))
	}
}

func TestRetreiveUpdateDesignDoc(t *testing.T) {

	err, desDoc := RetreiveDocFromDb("test_design", &DBObject)
	if err == nil {
		desDoc.Views[0].Name = "test_view_updated"
		desDoc.RevStatus = true
		err := desDoc.SaveDoc()
		if err != nil {
			t.Error(err)
		}
	}
}
