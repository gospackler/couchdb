# couchdb
Minimalist CouchDB client

## Documentation

- API References
  - [gowalker.org](https://gowalker.org/github.com/bushwood/couchdb)
  - [godoc.com](https://godoc.org/github.com/bushwood/couchdb)

## Understanding

* **client.go** - represents the couchdb client for the others to use.
* **db.go** - This is the low level db interface for the object. (Couch works on rest.)
* **dbrequests** - Wrapper for all requests from couch.
* **view.go** - Deals with creation of views.
* **document.go** - Deals with the creation and updating of documents in the db.

## Running Tests.

Make sure there is a couch instane running on default port (5984) for it to work. Have disabled the DeleteDB for now so the DB created would persist.

```bash
$ export GOPATH=$PWD
$ go get -a
$ go test -v
```

## Example Usage

For any operation, there is a Database Object which represents a connection to a database in couchdb client.

```go
client := NewClient("127.0.0.1", 5984) // Creates the client conenction.
dbObj := client.DB("testdb") //Db object which represents the connection to db.
```

Exists checks if the database exists.
Create can create the database if it does not exist.

```go
status, err := dbObj.Exists() // Status contains the status of the check
err := dbObj.Create() //Creates a new database with the dbName passed to the object
```
## Documents

Creating a new document can be done as follows. Any json can be save and the example below shows how to Marshal an object.
```go
doc := NewDocument("", "", &DBObject) //args - ID and Revison of the docuemnt to pickup
byteObj, err := json.Marshal(obj)
err := doc.Create(byteObj) //Creates a document Id and Rev would be updated by now.

// For updating,
obj.Update()
byteObj, err = json.Marshal(obj)
err = doc.Update((byteObj)
```

To get the object back, the way to go about it would be to wrap the object with CouchWrapper which can take care of the extra information that comes back.

```go
doc := NewDocument(Id, Rev, &DBObject) //Rev can be empty and not used right now, have it there for if present case.
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
```
## Views
View has two parts- DesignDoc and View
DesignDoc is the representation of the design Document. Each design Document can have multiple Views assosciated with it.
The example below shows how to come up with a designDocument contating multiple views.

```go
view := NewView("test_view", "doc", "doc.age < 22", "doc.name, doc.age")

view2 := NewView("raw_view", "", "", "")
view2.RawStatus = true
view2.RawJson = "function(rawDoc) {console.log(1234)}"

fView := NewView("fred_view", "newVar", "newVar.age > 22", "newVar.name, newVar.age")

doc := NewDesignDoc("test_design", &DBObject) // Creating the doc

doc.AddView(view)  // Adding the views into it.
doc.AddView(fView)
doc.AddView(view2)

err := doc.SaveDoc() //Saving the document which creates the permanent view.
if err != nil {
	t.Error(err)
}
```

Requesting a view.

```go
//arguments are @designDoc - followed by test_view

err, data := DBObject.GetView("test_design", "test_view")
if err != nil {
	t.Error("Error :", err)
} else {
	t.Log(string(data))
}
```

The views in designDocuments can be updated as well.

```go
err, desDoc := RetreiveDocFromDb("test_design", &DBObject)
if err == nil {
	desDoc.Views[0].Name = "test_view_updated" //Demo of an update. Assuming a view exists.
	desDoc.RevStatus = true
	err := desDoc.SaveDoc()
	if err != nil {
		t.Error(err)
	}
}
```
