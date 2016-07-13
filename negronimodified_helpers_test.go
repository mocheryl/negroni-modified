// Copyright 2016 Igor "Mocheryl" Zornik. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package negronimodified

import (
	"net/http"
	"testing"
	"time"
)

func TestIfModifiedSince(t *testing.T) {
	r := &http.Request{
		Header: make(http.Header),
	}

	// When request header is not set it should return positive
	if c := IfModifiedSince(r, http.TimeFormat); c != true {
		t.Errorf(`negronimodified.IfModifiedSince(%v, %q) = %t, want %t`, r, http.TimeFormat, c, true)
	}

	// When response time is not set it should return positive.
	if c := IfModifiedSince(r, ``); c != true {
		t.Errorf(`negronimodified.IfModifiedSince(%v, %q) = %t, want %t`, r, ``, c, true)
	}

	// When either times are set, but don't match, it should return positive.
	r.Header.Set(headerIfModifiedSince, http.TimeFormat)
	if c := IfModifiedSince(r, `Mon, 02 Jan 2006 16:04:05 GMT`); c != true {
		t.Errorf(`negronimodified.IfModifiedSince(%v, %q) = %t, want %t`, r, http.TimeFormat, c, true)
	}

	// If both times match completely, it should return negative.
	if c := IfModifiedSince(r, http.TimeFormat); c != false {
		t.Errorf(`negronimodified.IfModifiedSince(%v, %q) = %t, want %t`, r, http.TimeFormat, c, false)
	}
}

func TestIfModifiedSinceTime(t *testing.T) {
	r := &http.Request{
		Header: make(http.Header),
	}
	r.Header.Set(headerIfModifiedSince, http.TimeFormat)
	l, err := time.LoadLocation(`CET`)
	if err != nil {
		t.Fatalf(`time.LoadLocation(%q) = _, %v; want _, nil`, `CET`, err)
	}
	m := time.Date(2006, 1, 2, 16, 4, 5, 0, l)

	if c := IfModifiedSinceTime(r, m); c != false {
		t.Errorf(`negronimodified.IfModifiedSinceTime(%v, %v) = %t, want %t`, r, m, c, false)
	}
}

func TestSetLastModified(t *testing.T) {
	rw := &respWriteTest{
		h: make(http.Header),
	}

	if h := SetLastModified(rw, http.TimeFormat).Header().Get(headerLastModified); h != http.TimeFormat {
		t.Errorf(`negronimodified.SetLastModified(%v, %q) = %q, want %q`, rw, http.TimeFormat, h, http.TimeFormat)
	}
}

func TestSetLastModifiedTime(t *testing.T) {
	rw := &respWriteTest{
		h: make(http.Header),
	}
	l, err := time.LoadLocation(`CET`)
	if err != nil {
		t.Fatalf(`time.LoadLocation(%q) = _, %v; want _, nil`, `CET`, err)
	}
	m := time.Date(2006, 1, 2, 16, 4, 5, 0, l)

	if h := SetLastModifiedTime(rw, m).Header().Get(headerLastModified); h != http.TimeFormat {
		t.Errorf(`negronimodified.SetLastModifiedTime(%v, %v) = %q, want %q`, rw, m, h, http.TimeFormat)
	}
}
