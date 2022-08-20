package rest

import "net/http"

const (
	Content       = "Content-Type"
	Authorization = "Authorization"
)

type ContentType string

func (t ContentType) String() string {
	return string(t)
}

const (
	TextPlain       ContentType = "text/plain"
	ApplicationJSON ContentType = "application/json"
	ApplicationXML  ContentType = "application/xml"
)

type AuthType string

const (
	Bearer AuthType = "Bearer"
)

type HeaderOption func(headers http.Header)

func Header(opts ...HeaderOption) http.Header {
	h := http.Header{}
	for _, f := range opts {
		f(h)
	}
	return h
}

func WithHeader(header http.Header) HeaderOption {
	return func(h http.Header) {
		for key, value := range header {
			h[key] = value
		}
	}
}

func WithAuth(auth AuthType, token string) HeaderOption {
	return func(headers http.Header) {
		headers[Authorization] = []string{string(auth) + " " + token}
	}
}

func WithContent(value ContentType) HeaderOption {
	return func(headers http.Header) {
		headers[Content] = []string{string(value)}
	}
}

func With(name string, value ...string) HeaderOption {
	return func(headers http.Header) {
		headers[name] = value
	}
}
