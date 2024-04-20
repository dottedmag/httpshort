# Go library for short-cutting HTTP clients to HTTP handlers
[![Go Reference](https://pkg.go.dev/badge/github.com/dottedmag/httpshort.svg)](https://pkg.go.dev/github.com/dottedmag/httpshort)

Your client code uses an `http.Client`.
Your server code contains an `http.Handler`.

This package implements an `http.RoundTripper` that allows a HTTP client
to call a HTTP handler directly without starting a HTTP server.

    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello world")
    }

    client := &http.Client{
        Transport: Transport{
            Handler: mux,
        },
    }

    resp, err := client.Get("http://example.com")

There is also a helper function to create a `http.Client` from a handler:

    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello world")
    }

    client := Client(nil, mux)

    resp, err := client.Get("http://example.com")

## Legal

Copyright Mikhail Gusarov. Licensed under [Apache 2.0](LICENSE) license.
