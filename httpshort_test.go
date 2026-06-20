package httpshort

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/dottedmag/must"
)

func ExampleTransport() {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /size", func(w http.ResponseWriter, r *http.Request) {
		// Check passed header
		if r.Header.Get("Content-Type") != "text/plain" {
			http.Error(w, "", http.StatusUnsupportedMediaType)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "x-text/digits")
		fmt.Fprintf(w, "%d", len(body))
	})

	// Create a HTTP client that uses mux

	client := &http.Client{
		Transport: Transport{
			Handler: mux,
		},
	}

	respPost := must.OK1(client.Post("/size", "text/plain", strings.NewReader("1234")))
	defer respPost.Body.Close()

	fmt.Println("POST /size response status:", respPost.Status)
	fmt.Println("POST /size response Content-Type:", respPost.Header.Get("Content-Type"))
	fmt.Println("POST /size response body:", string(must.OK1(io.ReadAll(respPost.Body))))

	// Output:
	// POST /size response status: 200 OK
	// POST /size response Content-Type: x-text/digits
	// POST /size response body: 4
}

func ExampleClient() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "index")
	})

	client := Client(nil, mux)

	resp := must.OK1(client.Get("/"))
	defer resp.Body.Close()

	fmt.Println("GET / response status:", resp.Status)
	fmt.Println("GET / response body:", string(must.OK1(io.ReadAll(resp.Body))))

	// Output:
	// GET / response status: 200 OK
	// GET / response body: index
}

func TestPathCleaning(t *testing.T) {
	var actualPath string

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		actualPath = r.URL.Path
	})

	client := Client(nil, mux)

	resp := must.OK1(client.Get("///foobar/././foo///"))
	defer resp.Body.Close()

	if actualPath != "/foobar/foo/" {
		t.Errorf("Expected ///foobar/././foo/// to clean to /foobar/foo/, got %q", actualPath)
	}
}

// TestPathCleaningPreservesTrailingSlashSubtree exercises the bug
// fixed by mirroring net/http.cleanPath: a request to "/admin/" must
// reach the registered "/admin/" subtree handler directly, not be
// stripped to "/admin" and bounce off ServeMux's add-trailing-slash
// redirect. Following redirects on the default http.Client would
// loop until "stopped after 10 redirects" without the fix.
func TestPathCleaningPreservesTrailingSlashSubtree(t *testing.T) {
	var hit bool

	mux := http.NewServeMux()
	mux.HandleFunc("/admin/", func(w http.ResponseWriter, r *http.Request) {
		hit = true
	})

	client := Client(nil, mux)

	resp := must.OK1(client.Get("/admin/"))
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("got status %d, want 200", resp.StatusCode)
	}
	if !hit {
		t.Fatal("/admin/ handler not invoked; request was stripped to /admin")
	}
}

func TestRequest(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {})

	client := Client(nil, mux)

	resp := must.OK1(client.Get("/"))
	defer resp.Body.Close()

	if resp.Request == nil {
		t.Fatal("resp.Request fails")
	}
}
