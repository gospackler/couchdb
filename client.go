package couchdb

import "time"

// Client defines the base client for the orm
type Client struct {
	Host     string
	Port     int
	Username string
	Password string
	Secure   bool
	Timeout  int
}

// NewClient returns a new instance of a client with only host and port definted
func NewClient(h string, p int) Client {
	return Client{
		Host: h,
		Port: p,
	}
}

// SetPort is a thin utility for setting the port
func (c *Client) SetPort(p int) *Client {
	c.Port = p
	return c
}

// GetPort is a thin utility for getting the port
func (c *Client) GetPort() (p int) {
	return c.Port
}

// SetHost is a thin utility for setting the host
func (c *Client) SetHost(h string) *Client {
	c.Host = h
	return c
}

// GetHost is a thin utility for getting the host
func (c *Client) GetHost() string {
	return c.Host
}

// SetUser is a thin utility for setting the username
func (c *Client) SetUser(u string) *Client {
	c.Username = u
	return c
}

// GetUser is a thin utility for getting the username
func (c *Client) GetUser() string {
	return c.Username
}

// SetPwd is a thin utility for setting the password
func (c *Client) SetPwd(p string) *Client {
	c.Password = p
	return c
}

// GetPwd is a thin utility for getting the password
func (c *Client) GetPwd() string {
	return c.Password
}

// SetAuth is a thin utility for setting the password
func (c *Client) SetAuth(u string, p string) *Client {
	c.Username = u
	c.Password = p
	return c
}

// GetAuth is a thin utility for getting the password
func (c *Client) GetAuth() (string, string) {
	return c.Username, c.Password
}

// SetTimeout is a thin utility for setting the timeout
func (c *Client) SetTimeout(t int) *Client {
	c.Timeout = t
	return c
}

// GetTimeout is a thin utility for getting the timeout
func (c *Client) GetTimeout() int {
	return c.Timeout
}

// GetTimeoutDuration is a thin utility for getting the timeout in duration format
func (c *Client) GetTimeoutDuration() time.Duration {
	return time.Duration(c.Timeout) * time.Millisecond
}

// SetSecure is a thin utility for setting the client to https
func (c *Client) SetSecure() {
	c.Secure = true
}

// SetInsecure is a thin utility for setting the client to http
func (c *Client) SetInsecure() {
	c.Secure = false
}

// DB creates a Database client
func (c *Client) DB(name string) Database {
	return NewDB(name, c)
}
