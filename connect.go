// Copyright (c) 2012-2013 Jason McVetta.  This is Free Software, released under
// the terms of the GPL v3.  See http://www.gnu.org/copyleft/gpl.html for details.
// Resist intellectual serfdom - the ownership of ideas is akin to slavery.

// +build !appengine

package neoism

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"gopkg.in/jmcvetta/napping.v3"
)

// option is a type alias for a function that takes a pointer to a Database.
// Used for functional configuration of our client.
type option func(*Database)

// WithClient is a configuration function that allows users of this library to
// supply their own http.Client. This is required if they want to use a self
// signed TLS certificate for example.
func WithClient(client *http.Client) option {
	return func(db *Database) {
		db.Session.Client = client
	}
}

// Connect sets up our client for connecting to the Neo4j server and calls
// ConnectWithRetry()
func Connect(uri string, options ...option) (*Database, error) {
	h := http.Header{}
	h.Add("User-Agent", fmt.Sprintf("neoism/%s (%s)", VERSION, runtime.GOOS))
	db := &Database{
		Session: &napping.Session{
			Header: &h,
		},
	}

	// apply our configuration functions
	for _, opt := range options {
		opt(db)
	}

	// trailing slash is important, check if it's not there and add it
	if !strings.HasSuffix(uri, "/") {
		uri += "/"
	}
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	if parsedURL.User != nil {
		db.Session.Userinfo = parsedURL.User
	}
	return connectWithRetry(db, parsedURL, 0)
}
