// Copyright 2016 Igor "Mocheryl" Zornik. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package negronimodified

import (
	`fmt`
	`net/http`

	`github.com/codegangsta/negroni`
)

// NewModified basic usage.
func ExampleNewModified() {
	mux := http.NewServeMux()
	mux.HandleFunc(`/`, func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set(`Last-Modified`, `Mon, 02 Jan 2006 16:04:05 GMT`)
		fmt.Fprintf(w, `This content will not be re-served if you ask about the modification time in the HTTP header!`)
	})

	n := negroni.Classic()
	n.Use(NewModified())
	n.UseHandler(mux)
	n.Run(`:3000`)
}
