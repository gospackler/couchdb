package couchdb

import (
	"encoding/json"
	"errors"
	"strconv"

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
	_, body, errs := db.Req.Get("").End()
	if len(errs) > 0 {
		return false, errs[0]
	}
	result := response{}
	pErr := json.Unmarshal([]byte(body), &result)
	if pErr != nil {
		return false, pErr
	}
	if result.Error != "" {
		return false, errors.New(result.Error)
	}
	if result.DBName != db.Name {
		return false, errors.New("Couch returned database [" + result.DBName + "]")
	}

	return true, nil
}

// Create creates a new database
func (db *Database) Create(name string) error {
	return nil
}
