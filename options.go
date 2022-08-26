package rest

import "net/http"

type RequestData[T any] struct {
	Body           any
	Header         http.Header
	encoder        Encoder
	decoder        Decoder
	MarshalFunc    MarshallFunc
	UnmarshallFunc UnmarshallFunc
}

func (r *RequestData[T]) Encode() ([]byte, error) {
	ct := ContentType(r.Header.Get(Content))
	encoder := r.encoder
	if r.MarshalFunc != nil {
		encoder = r.encoder.Clone()
		encoder.Set(ct, r.MarshalFunc)
	}
	return encoder.Encode(r.Body, ct)
}

func (r *RequestData[T]) Decode(data []byte, v *T) error {
	ct := ContentType(r.Header.Get(Content))
	decoder := r.decoder
	if r.UnmarshallFunc != nil {
		decoder = r.decoder.Clone()
		decoder.Set(ct, r.UnmarshallFunc)
	}
	return r.decoder.Decode(data, v, ct)
}

type Option func(r *RequestData[any])

// WithBody option set Body for request
func WithBody(body any) Option {
	return func(r *RequestData[any]) {
		r.Body = body
	}
}

// WithHeaders option set headers for request
func WithHeaders(header http.Header) Option {
	return func(r *RequestData[any]) {
		r.Header = header
	}
}

// WithMarshallFunc option set function for encoding for some contentType
func WithMarshallFunc(f MarshallFunc) Option {
	return func(r *RequestData[any]) {
		r.MarshalFunc = f
	}
}

// WithUnmarshalFunc option set the function f for decoding for some contentType
func WithUnmarshalFunc(f UnmarshallFunc) Option {
	return func(r *RequestData[any]) {
		r.UnmarshallFunc = f
	}
}
