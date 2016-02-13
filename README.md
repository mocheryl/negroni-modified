# negroni-modified

Last-Modified middleware for [Negroni](https://github.com/codegangsta/negroni).

## Usage

~~~ go
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
~~~

See [godoc.org](http://godoc.org/github.com/mocheryl/negroni-modified) for more information.

## License

negroni-modified is released under the 3-Clause BSD license.
See [LICENSE](https://github.com/mocheryl/negroni-modified/blob/master/LICENSE).
