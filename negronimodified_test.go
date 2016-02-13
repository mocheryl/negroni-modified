// Copyright 2016 Igor "Mocheryl" Zornik. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package negronimodified

import (
	`bufio`
	`net/http`
	`net/http/httptest`
	`testing`

	`github.com/codegangsta/negroni`
)

type respWriteTest struct {
	h http.Header
}

func (w *respWriteTest) Header() http.Header {
	return w.h
}

func (w *respWriteTest) WriteHeader(code int) {}

func (w *respWriteTest) Write(data []byte) (n int, err error) {
	return
}

func TestModifiedResponseWriter_Write(t *testing.T) {
	rw := &respWriteTest{
		h: make(http.Header),
	}
	nrw := negroni.NewResponseWriter(rw)
	crw := &modifiedResponseWriter{
		bufio.NewWriter(rw),
		nrw,
	}

	if n, err := crw.Write([]byte(`test`)); n != 4 || err != nil {
		t.Errorf(`negronimodified.modifiedResponseWriter.Write(%s) = %d, %v; want %d, nil`, []byte(`test`), n, err, 4)
	}
	// TODO: Check if correct data has been written.
}

func TestNewModified(t *testing.T) {
	handler := NewModified()
	if handler == nil {
		t.Fatal(`negronimodified.NewModified() cannot return nil`)
	}

	if handler.cacheControl != `public, must-revalidate` {
		t.Errorf(`negronimodified.NewModified().cacheControl = %q, want %q`, handler.cacheControl, `public, must-revalidate`)
	}
}

func TestNewModifiedWithCacheControl(t *testing.T) {
	handler := NewModifiedWithCacheControl([]string{`must-revalidate`, `public`})
	if handler == nil {
		t.Fatal(`negronimodified.NewModifiedWithCacheControl() cannot return nil`)
	}

	if handler.cacheControl != `must-revalidate, public` {
		t.Errorf(`negronimodified.NewModifiedWithCacheControl(%q, %q).cacheControl = %q, want %q`, `must-revalidate`, `public`, handler.cacheControl, `must-revalidate, public`)
	}
}

func TestModified_ServeHTTP(t *testing.T) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest(`GET`, `http://localhost/foo`, nil)
	if err != nil {
		t.Fatalf(`http.NewRequest(%q, %q, nil) = _, %v; want _, nil`, `GET`, `http://localhost/foo`, err)
	}

	// Test empty cache control headers.
	handler := NewModifiedWithCacheControl([]string{})

	handler.ServeHTTP(w, req, func(w http.ResponseWriter, r *http.Request) {

	})
	if h := w.Header().Get(headerCacheControl); h != `` {
		t.Errorf(`httputil.ResponseRecorder.Header().Get(%q) = %q, want %q`, headerCacheControl, h, ``)
	}

	// Test with cache control headers set.
	handler = NewModified()
	handler.ServeHTTP(w, req, func(w http.ResponseWriter, r *http.Request) {

	})
	if h := w.Header().Get(headerCacheControl); h != `public, must-revalidate` {
		t.Errorf(`httputil.ResponseRecorder.Header().Get(%q) = %q, want %q`, headerCacheControl, h, `public, must-revalidate`)
	}

	// Test again with cache control headers set but override them with custom
	// headers set by handler.
	handler.ServeHTTP(w, req, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerCacheControl, `must-revalidate, public`)
	})
	if h := w.Header().Get(headerCacheControl); h != `must-revalidate, public` {
		t.Errorf(`httputil.ResponseRecorder.Header().Get(%q) = %q, want %q`, headerCacheControl, h, `must-revalidate, public`)
	}

	// Check returned content with status not set.
	handler.ServeHTTP(w, req, func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`FooBar`))
	})
	if b := w.Body.String(); b != `FooBar` {
		t.Errorf(`httputil.ResponseRecorder.Body.String() = %q, want %q`, b, `FooBar`)
	}
	if c := w.Code; c != http.StatusOK {
		t.Errorf(`httputil.ResponseRecorder.Code = %d, want %d`, c, http.StatusOK)
	}
	// w.Body.Reset()
	// w.Flushed = false

	// FIXME: Could this test be performed without initializing a new recorder?
	w = httptest.NewRecorder()
	// When request and response headers for modified content match, we must
	// not receive any body content at all.
	req.Header.Set(headerIfModifiedSince, http.TimeFormat)
	handler.ServeHTTP(w, req, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerLastModified, http.TimeFormat)
		w.Write([]byte(`BarFoo`))
	})
	if b := w.Body.String(); b != `` {
		t.Errorf(`httputil.ResponseRecorder.Body.String() = %q, want %q`, b, ``)
	}
	if c := w.Code; c != http.StatusNotModified {
		t.Errorf(`httputil.ResponseRecorder.Code = %d, want %d`, c, http.StatusNotModified)
	}
}
