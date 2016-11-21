//File contains all the db related functions of couch
package couchdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
)

// Database defintes a database client
// It can be used to get Documents and Views
type Database struct {
	Name    string
	BaseURL string
	Client  *Client
	Req     Request
}

// NewDB creates an instance of a Database
func NewDB(name string, c *Client) Database {
	protocol := "http"
	if c.Secure {
		protocol += "s"
	}

	port := ""
	if c.Port == 0 {
		port = "80"
	} else {
		port = strconv.Itoa(c.Port)
	}

	url := protocol + "://" + c.Host + ":" + port + "/" + name
	req := Request{gorequest.New(), url}
	if c.Timeout != 0 {
		req.Req.Timeout(c.GetTimeoutDuration())
	}
	if c.Username != "" && c.Password != "" {
		req.Req.SetBasicAuth(c.Username, c.Password)
	}
	return Database{
		Name:    name,
		Client:  c,
		BaseURL: url,
		Req:     req,
	}
}

// Exists check to see if database exists
func (db *Database) Exists() (bool, error) {
	type response struct {
		Error  string `json:"error"`
		DBName string `json:"db_name"`
	}
	_, body, _ := db.Req.Get("").End()
	result := response{}
	pErr := json.Unmarshal([]byte(body), &result)
	if pErr != nil {
		return false, pErr
	}
	if result.DBName != db.Name || result.Error != "" {
		return false, nil
	}

	return true, nil
}

// Create creates a new database
func (db *Database) Create() error {
	type response struct {
		Error  string `json:"error"`
		Reason string `json:"reason"`
		Ok     bool   `json:"ok"`
	}
	_, body, _ := db.Req.Put("").End()
	result := response{}
	pErr := json.Unmarshal([]byte(body), &result)
	log.Debug("couch : Create : Result, JsonResp", result, body)
	if pErr != nil {
		return pErr
	}
	if result.Error != "" {
		return errors.New(result.Error + " " + result.Reason)
	}
	if !result.Ok {
		return errors.New("Couch returned failure when creating [" + db.Name + "]")
	}

	return nil
}

func (db *Database) GetView(docName string, viewName string, query string) ([]byte, error) {

	log.Debugf("couch : GetView query %s in viewName %s of desDoc %s", query, viewName, docName)
	type ViewResponse struct {
		Error  string `json:"error"`
		Reason string `json:"reason"`
	}

	var body string
	var errs []error
	var prefix string
	var superAgent *gorequest.SuperAgent

	if query == "" {
		prefix = docName + "/_view/" + viewName
		log.Debug("Getting view name " + prefix)
		_, body, errs = db.Req.Get(prefix).End()
	} else {
		values, err := url.ParseQuery(query)
		if err != nil {
			return nil, errors.New("Unable to parse query string: " + query)
		}
		encodedKey := values.Encode()

		prefix = docName + "/_view/" + viewName
		superAgent = db.Req.Get(prefix).Query(encodedKey)
		_, body, errs = superAgent.End()
		log.Debug("Url" + superAgent.Url + encodedKey)
	}

	if len(errs) > 0 {
		return nil, errors.New("Database : Error making request " + fmt.Sprint("%v", errs))
	}
	viewResp := &ViewResponse{}
	err := json.Unmarshal([]byte(body), viewResp)

	if err != nil {
		log.Error(body)
		return nil, err
	}

	if viewResp.Error != "" {
		err = errors.New(viewResp.Error + " " + viewResp.Reason + "\n req " + superAgent.Url)
		return nil, err
	}

	log.Debug("Returning body", body)
	return []byte(body), nil
}

// Delete deletes database
func (db *Database) Delete() error {
	type response struct {
		Error string `json:"error"`
		Ok    bool   `json:"ok"`
	}
	_, body, _ := db.Req.Delete("").End()
	result := response{}
	pErr := json.Unmarshal([]byte(body), &result)
	log.Debug(result)
	if pErr != nil {
		return pErr
	}
	if result.Error != "" {
		return errors.New(result.Error)
	}
	if !result.Ok {
		return errors.New("Couch returned failure when creating [" + db.Name + "]")
	}

	return nil
}
