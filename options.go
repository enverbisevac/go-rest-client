package rest

import "net/http"

type requestData struct {
	Body   any
	Header http.Header
}

type Option func(r *requestData)

// WithBody option set body for request
func WithBody(body any) Option {
	return func(r *requestData) {
		r.Body = body
	}
}

// WithHeaders option set headers for request
func WithHeaders(header http.Header) Option {
	return func(r *requestData) {
		r.Header = header
	}
}
