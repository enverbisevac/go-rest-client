package rest

import "net/http"

type RequestData struct {
	Body   any
	Header http.Header
}

type Option func(r *RequestData)

func WithBody(body any) Option {
	return func(r *RequestData) {
		r.Body = body
	}
}

func WithHeaders(header http.Header) Option {
	return func(r *RequestData) {
		r.Header = header
	}
}
