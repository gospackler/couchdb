//File contains all the db related functions of couch
package couchdb

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

// Database defintes a database client
// It can be used to get Documents and Views
type Database struct {
	Name    string
	BaseURL string
	Client  *Client
	Req     *Request
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

	httpClient := new(http.Client)
	if c.Timeout != 0 {
		httpClient.Timeout = c.GetTimeoutDuration()
	}

	req := &Request{
		httpClient,
		url,
		c.Username,
		c.Password,
	}

	return Database{
		Name:    name,
		Client:  c,
		BaseURL: url,
		Req:     req,
	}
}

// Exists check to see if database exists
func (db *Database) Exists() error {
	resp := &struct {
		Error  string `json:"error"`
		DBName string `json:"db_name"`
	}{}
	body, err := db.Req.Get("", nil)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return err
	}
	if resp.Error != "" {
		return errors.New(resp.Error)
	}

	return nil
}

// Create creates a new database
func (db *Database) Create() error {
	resp := &struct {
		Error  string `json:"error"`
		Reason string `json:"reason"`
		Ok     bool   `json:"ok"`
	}{}

	body, err := db.Req.Put("")
	err = json.Unmarshal(body, resp)
	if err != nil {
		return err
	}

	if resp.Error != "" {
		return errors.New(resp.Error + " " + resp.Reason)
	}
	if !resp.Ok {
		return errors.New("Couch returned failure when creating [" + db.Name + "]")
	}

	return nil
}

func (db *Database) GetView(docName string, viewName string, args map[string]string) ([]byte, error) {
	resp := &struct {
		Error  string `json:"error"`
		Reason string `json:"reason"`
	}{}

	prefix := docName + "/_view/" + viewName
	body, err := db.Req.Get(prefix, args)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(body), resp)

	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		err = errors.New(resp.Error + " " + resp.Reason + "\n req :" + prefix)
		return nil, err
	}
	return body, nil
}

// Delete deletes database
func (db *Database) Delete() error {
	result := &struct {
		Error string `json:"error"`
		Ok    bool   `json:"ok"`
	}{}

	body, err := db.Req.Delete("")
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	if result.Error != "" {
		return errors.New(result.Error)
	}
	if !result.Ok {
		return errors.New("Couch returned failure when creating [" + db.Name + "]")
	}

	return nil
}
