package couchdb

import (
	"github.com/parnurzeal/gorequest"
	"path"
)

// Request defines a base request used in DB connections
type Request struct {
	Req     *gorequest.SuperAgent
	BaseURL string
}

// Get build the GET superagent with the path
func (r *Request) Get(p string) *gorequest.SuperAgent {
	r.Req = r.Req.Get(r.PathUrl(p))
	return r.Req
}

// Post build the POST superagent with the path
func (r *Request) Post(p string) *gorequest.SuperAgent {
	r.Req = r.Req.Post(r.PathUrl(p))
	return r.Req
}

// Put build the PUT superagent with the path
func (r *Request) Put(p string) *gorequest.SuperAgent {
	r.Req = r.Req.Put(r.PathUrl(p))
	return r.Req
}

// Delete build the DELETE superagent with the path
func (r *Request) Delete(p string) *gorequest.SuperAgent {
	r.Req = r.Req.Delete(r.PathUrl(p))
	return r.Req
}

// PathUrl joins the provided path with the host
func (r *Request) PathUrl(p string) string {
	return r.BaseURL + path.Join("/", p)
}
