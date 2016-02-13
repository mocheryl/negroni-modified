// Copyright 2016 Igor "Mocheryl" Zornik. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package negronimodified

import (
	`bufio`
	`io/ioutil`
	`net/http`
	`strings`

	`github.com/codegangsta/negroni`
)

const headerCacheControl string = `Cache-Control` // Cache mechanism

// modifiedResponseWriter is the ResponseWriter that negroni.ResponseWriter is
// wrapped in.
type modifiedResponseWriter struct {
	c *bufio.Writer
	negroni.ResponseWriter
}

// Write appends any data to writers buffer.
func (m *modifiedResponseWriter) Write(b []byte) (int, error) {
	return m.c.Write(b)
}

// NewModified returns a new modified middleware instance with default cache
// control headers set.
func NewModified() *modified {
	return NewModifiedWithCacheControl([]string{`public`, `must-revalidate`})
}

// NewModified returns a new modified middleware instance.
func NewModifiedWithCacheControl(cacheControl []string) *modified {
	return &modified{strings.Join(cacheControl, `, `)}
}

// modified handles the necessity of a response in case of client and server
// cache options.
type modified struct {
	cacheControl string
}

func (m *modified) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	// TODO: Somehow find out if response is already set to not modified, so we
	// we exit immediately.

	// Set any possible cache control headers if none have been set yet.
	if m.cacheControl != `` && rw.Header().Get(headerCacheControl) == `` {
		rw.Header().Set(headerCacheControl, m.cacheControl)
	}

	// Override default writer.
	nrw := negroni.NewResponseWriter(rw)
	crw := &modifiedResponseWriter{
		bufio.NewWriter(rw),
		nrw,
	}
	next(crw, r)

	// Check if both client and server are using cache capabilities.
	if IfModifiedSince(r, rw.Header().Get(headerLastModified)) {
		// Either of the required headers are not set. Write whatever we have
		// to regular writer and return.
		// TODO: Error check.
		crw.c.Flush()
		return
	}

	// Discard any data because user agent already has the latest version.
	crw.c.Reset(ioutil.Discard)
	rw.WriteHeader(http.StatusNotModified)
}
