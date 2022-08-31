package rest

import "net/http"

type RequestOption struct {
	Body           any
	Header         http.Header
	MarshalFunc    MarshallFunc
	UnmarshallFunc UnmarshallFunc
}

type Option func(r *RequestOption)

// WithBody option set Body for request
func WithBody(body any) Option {
	return func(d *RequestOption) {
		d.Body = body
	}
}

// WithHeaders option set headers for request
func WithHeaders(header http.Header) Option {
	return func(d *RequestOption) {
		d.Header = header
	}
}

// WithMarshallFunc option set function for encoding for some contentType
func WithMarshallFunc(f MarshallFunc) Option {
	return func(d *RequestOption) {
		d.MarshalFunc = f
	}
}

// WithUnmarshalFunc option set the function f for decoding for some contentType
func WithUnmarshalFunc(f UnmarshallFunc) Option {
	return func(d *RequestOption) {
		d.UnmarshallFunc = f
	}
}
