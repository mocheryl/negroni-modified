// Copyright 2016 Igor "Mocheryl" Zornik. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package negronimodified implements checking and setting of last modified
headers handler middleware for Negroni.

Basics

HTTP servers or content generators in the background can supress sending the
HTTP body content along with the response based on some of the headers set in
client request, most notably the "Last-Modified" on the server side and
"If-Modified-Since" on the user-agent side. There are some other headers that
impact this functionality, such as "Cache-Control" (more on this later), but
these two are the main ones we deal with on average basis.

When the server (or whatever is responsible for setting the header content)
decides that the content of interest has a certain generated time attached to
it, it can inform the client about in the response it sends back along side the
actual content. The client can then leverage this data for its cacheing
purposes. It does this in a way that in the next identical request it sends this
same time to the server where it can check it against the requested content time
and if the times match, none of it needs to be sent to the client. This way we
save on processing resources, bandwidth and, consequentially, time.

Usage

	package main

	import (
		`fmt`
		`net/http`

		`github.com/codegangsta/negroni`
		`github.com/mocheryl/negroni-modified`
	)

	func main() {
		mux := http.NewServeMux()
		mux.HandleFunc(`/`, func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set(`Last-Modified`, `Mon, 02 Jan 2006 16:04:05 GMT`)
			fmt.Fprintf(w, `This content will not be re-served if you ask about the modification time in the HTTP header!`)
		})

		n := negroni.Classic()
		n.Use(negronimodified.NewModified())
		n.UseHandler(mux)
		n.Run(`:3000`)
	}

The above code initializes the middleware with default settings. These include
the setting that will add "Cache-Control" header with a value of "public,
must-revalidate", which means that everything should be cached and reevaluated.

You can set your own value (or remove it all-together) by initializing the
middleware like this:

	NewModifiedWithCacheControl([]string{`private`})

Tips

To avoid having to study on how to manipulate the headers correctly, you can use
some of the included convenience function.

To set the modification time use:

	SetLastModified(responseWriter, `Mon, 02 Jan 2006 16:04:05 GMT`)

If you worry about setting the value in a correct date and time format, a
version of a function whose second argument is a time.Date type has also been
included:

	SetLastModifiedTime(responseWriter, time.Date(2006, 1, 2, 16, 4, 5, 0, tz))

Because the header checkup happens further down the line, it means that a lot of
content processing and generation goes on at first only to find out at the end
that it will be discarded due to the fact that the client already has the
latest version. We can prevent all this unnecessary data processing by executing
and checking for the returned value of the following included convenience
function:

	IfModifiedSince(request, `Mon, 02 Jan 2006 16:04:05 GMT`)

If times are off, if should return true and we can proceed as usual, otherwise
halt any further operations on data generating. A practical example would look
like this:

	if !IfModifiedSince(request, `Mon, 02 Jan 2006 16:04:05 GMT`) {
		return  // Content hasn't changed, i.e. client already has the latest version.
	}

Naturally, a time.Date variant of this function also exists:

	IfModifiedSince(request, time.Date(2006, 1, 2, 16, 4, 5, 0, tz))

Further Reading

For additional information, please check the following link on official
definition:
http://www.w3.org/Protocols/rfc2616/rfc2616-sec10.html#sec10.3.5

*/
package negronimodified
