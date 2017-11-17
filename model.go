package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
)

// define our message object
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

type Profile struct {
	Id           string `json:"id"`
	Firstname    string `json:"firstname"`
	Lastname     string `json:"lastname"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Title        string `json:"title"`
	Password     string `json:",omitempty"`
	MobileNumber string `json:"mobilenumber"`
}

type PublishedContent struct {
	Type     int         `json:"type"` // either PUBTYPE_MESSAGE or PUBTYPE_CONTACTS
	Contents interface{} `json:"contents"`
}

type Prox struct {
	target        *url.URL
	proxy         *httputil.ReverseProxy
	routePatterns []*regexp.Regexp // add some route patterns with regexp
}

func New(target string) *Prox {
	url, err := url.Parse(target)

	if err != nil {
		panic(err)
	}
	return &Prox{target: url, proxy: httputil.NewSingleHostReverseProxy(url)}
}

func (p *Prox) handle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("X-GoProxy", "GoProxy")

	if p.routePatterns == nil || p.parseWhiteList(r) {
		p.proxy.ServeHTTP(w, r)
	}
}

func (p *Prox) parseWhiteList(r *http.Request) bool {
	for _, regexp := range p.routePatterns {
		fmt.Println(r.URL.Path)
		if regexp.MatchString(r.URL.Path) {
			// let's forward it
			return true
		}
	}
	return false
}
