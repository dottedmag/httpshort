package httpshort

import (
	"context"
	"net/http"
	"net/http/httptest"
)

// Transport is http.RoundTripper that directly calls the HTTP handler
//
// This transport useful for creating http.Client instances in tests:
// - short-cutting clients to existing HTTP server handlers
// - mocking external HTTP endpoints
type Transport struct {
	// Context, if specified, limits the lifetime of handler
	// invocations: request passed to handler will be cancelled
	// whenever this context is cancelled.
	//
	// Deadlines and values from this context are ignored.
	//
	// This is useful to shut down tests timely and cleanly.
	Context context.Context //nolint:containedctx

	// Handler is the handler to be called
	Handler http.Handler
}

func (t Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Amend the request context if needed
	switch {
	case t.Context == nil:
		// keep the original context (nil or non-nil)
	case t.Context != nil && req.Context() == nil:
		// set the context in request
		req = req.WithContext(t.Context)
	default:
		// merge context cancellations
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		context.AfterFunc(t.Context, cancel)
		req = req.WithContext(ctx)
	}

	recorder := httptest.NewRecorder()
	t.Handler.ServeHTTP(recorder, req)
	return recorder.Result(), nil
}
