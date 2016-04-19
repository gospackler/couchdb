//File contains all the db related functions of couch
package couchdb

import (
	"encoding/json"
	"errors"
	"fmt"
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
	log.Info("couch : Create : Result, JsonResp", result, body)
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

func (db *Database) GetView(docName string, viewName string, key string) (error, []byte) {
	type ViewResponse struct {
		Error  string `json:"error"`
		Reason string `json:"reason"`
	}

	var body string
	var errs []error
	if key == "" {
		prefix := docName + "/_view/" + viewName
		log.Info("Getting view name " + prefix)
		_, body, errs = db.Req.Get(prefix).End()
	} else {
		prefix := docName + "/_view/" + viewName
		_, body, errs = db.Req.Get(prefix).Query("key=" + key).End()
	}

	if len(errs) > 0 {
		return errors.New("Database : " + fmt.Sprint("%v", errs)), nil
	}
	viewResp := &ViewResponse{}
	err := json.Unmarshal([]byte(body), viewResp)

	if err != nil {
		log.Error(body)
		return err, nil
	}

	if viewResp.Error != "" {
		err = errors.New(viewResp.Error + " " + viewResp.Reason)
		return err, nil
	}

	return nil, []byte(body)
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
	log.Info(result)
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
