package couchdb

import (
	"encoding/json"
	"errors"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/parnurzeal/gorequest"
)

// Database defintes a database client
type Database struct {
	Name    string
	Client  *Client
	BaseURL string
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
		Error string `json:"error"`
		Ok    bool   `json:"ok"`
	}
	_, body, _ := db.Req.Put("").End()
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
