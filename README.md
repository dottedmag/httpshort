# Go library for short-cutting HTTP clients to HTTP handlers

Your client code uses a `http.Client`.
Your server code contains a `http.Handler`.

This package implements a `http.RoundTripper` that allows a HTTP client
to call a HTTP handler directly without starting a HTTP server.

    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "Hello world")
    }

    client := &http.Client{Transport: Transport{Handler: mux}}
    resp, err := client.Get("http://example.com")

## Legal

Copyright Mikhail Gusarov. Licensed under [Apache 2.0](LICENSE) license.
