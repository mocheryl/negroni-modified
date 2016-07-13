// Copyright 2016 Igor "Mocheryl" Zornik. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package negronimodified

import (
	"net/http"
	"time"
)

const (
	headerIfModifiedSince string = `If-Modified-Since` // Request HTTP header.
	headerLastModified    string = `Last-Modified`     // Response HTTP header.
	timeLocationGMT       string = `GMT`               // Timezone.
)

// TODO: Check for error.
var gmt, _ = time.LoadLocation(timeLocationGMT)

// IfModifiedSince checks HTTP request against the provided time string. It
// returns true if request requires a response based on the header and input
// time.
// Use this function to prematurely determine if any further data processing is
// necessary in case if you have large quantity of data to process and output.
// Input string must be in format as defined by RFC 1123.
func IfModifiedSince(r *http.Request, lastModified string) bool {
	ifMod := r.Header.Get(headerIfModifiedSince)
	if ifMod == `` || lastModified == `` || ifMod != lastModified {
		return true
	}

	return false
}

// IfModifiedSinceTime checks HTTP request against the provided time.Time. It
// returns true if request requires a response based on the header and input
// time.
// Use this function to prematurely determine if any further data processing is
// necessary in case if you have large quantity of data to process and output.
func IfModifiedSinceTime(r *http.Request, lastModified time.Time) bool {
	return IfModifiedSince(r, lastModified.In(gmt).Format(time.RFC1123))
}

// SetLastModified sets the provided time string into the response writer.
// Input string must be in format as defined by RFC 1123.
func SetLastModified(rw http.ResponseWriter, lastModified string) http.ResponseWriter {
	rw.Header().Set(headerLastModified, lastModified)
	return rw
}

// SetLastModifiedTime sets the time from the provided time object into the
// response writer.
func SetLastModifiedTime(rw http.ResponseWriter, lastModified time.Time) http.ResponseWriter {
	return SetLastModified(rw, lastModified.In(gmt).Format(time.RFC1123))
}
