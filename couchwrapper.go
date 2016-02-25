// Representation of a document in couch DB
package couchdb

type CouchWrapperUpdate struct {
	Id  string `json:"_id"`
	Rev string `json:"_rev"`
}

/**
Expecting this to last in on a higher level and data is abstract at this level.
func getJsonData(obj interface{}) (err error, jsonData string) {
	data, err := json.Marshal(obj)
	if err != nil {
		return err, string(data)
	}
	jsonData = string(data)
	return
}
*/
