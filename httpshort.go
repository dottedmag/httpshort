package httpshort

import (
	"context"
	"net/http"
	"net/http/httptest"
	"path"
)

// Transport is an http.RoundTripper that directly calls the HTTP handler
//
// This transport is useful for creating http.Client instances for testing:
// either short-cutting clients to existing HTTP server handlers, or
// mocking external HTTP endpoints.
type Transport struct {
	// Context, when specified, limits the lifetime of handler
	// invocation. It ensures that any request passed to the handler
	// gets cancelled if this context is cancelled.
	//
	// Deadlines and values from this context are disregarded.
	//
	// This is particularly useful for timely and clean shutdown of tests.
	Context context.Context //nolint:containedctx

	// Handler is the handler to be invoked
	Handler http.Handler
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Amend the request context if needed
	switch {
	case t.Context == nil:
		// keep the original context (nil or non-nil)
		if req.Context() == nil {
			req = req.Clone(context.Background())
		} else {
			req = req.Clone(req.Context())
		}
	case t.Context != nil && req.Context() == nil:
		// assign the context to the request that doesn't have one
		req = req.Clone(t.Context)
	default:
		// merge context cancellations
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		context.AfterFunc(t.Context, cancel)
		req = req.Clone(ctx)
	}

	// Clean the path like DefaltServeMux mux does
	req.URL.Path = path.Clean(req.URL.Path)

	recorder := httptest.NewRecorder()
	t.Handler.ServeHTTP(recorder, req)
	result := recorder.Result()
	result.Request = req
	return result, nil
}

// Client is a helper function that returns an http.Client.
// The client is configured with the Transport from this package.
func Client(ctx context.Context, handler http.Handler) *http.Client {
	return &http.Client{
		Transport: &Transport{
			Context: ctx,
			Handler: handler,
		},
	}
}
