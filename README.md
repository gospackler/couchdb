# couchdb
Minimalist CouchDB client

## Understanding.

client.go --> represents the couchdb client for the others to use.
db.go --> This is the low level db interface for the object. (Couch works on rest.)
dbrequests --> Wrapper for all requests from couch.
view.go --> Deals with creation of views.
document.go --> Deals with the creation and updating of documents in the db.

## Running Tests. 

$ go test -v

The *_test fils contains the tests written.

		

