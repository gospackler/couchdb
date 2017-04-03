package couchdb

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"

	log "github.com/Sirupsen/logrus"
)

// Request defines a base request used in DB connections
type Request struct {
	Req      *http.Client
	BaseURL  string
	username string
	password string
}

func (r *Request) makeRequest(method string, url string, body io.Reader, headers map[string]string) ([]byte, error) {
	httpReq, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	//	httpReq.Close = true
	if r.username != "" && r.password != "" {
		httpReq.SetBasicAuth(r.username, r.password)
	}

	log.Debugf("Method = %s, Url = %s", method, url)
	for key, value := range headers {
		httpReq.Header.Add(key, value)
	}
	resp, err := r.Req.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (r *Request) Get(p string, args map[string]string) ([]byte, error) {
	url, err := r.PathUrl(p)
	if err != nil {
		return nil, err
	}
	urlStr := url.String()
	if args != nil {
		values := url.Query()
		for key, value := range args {
			values.Add(key, "\""+value+"\"")
		}

		encVal := values.Encode()
		urlStr = urlStr + "?" + encVal
	}

	return r.makeRequest("GET", urlStr, nil, nil)
}

func (r *Request) Post(p string, body []byte, headers map[string]string) ([]byte, error) {
	url, err := r.PathUrl(p)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(body)
	return r.makeRequest("POST", url.String(), buffer, headers)
}

func (r *Request) Put(p string) ([]byte, error) {
	url, err := r.PathUrl(p)
	if err != nil {
		return nil, err
	}
	return r.makeRequest("PUT", url.String(), nil, nil)
}

func (r *Request) Delete(p string) ([]byte, error) {
	url, err := r.PathUrl(p)
	if err != nil {
		return nil, err
	}
	return r.makeRequest("DELETE", url.String(), nil, nil)
}

// PathUrl joins the provided path with the host
func (r *Request) PathUrl(p string) (*url.URL, error) {
	return url.Parse(r.BaseURL + path.Join("/", p))
}
