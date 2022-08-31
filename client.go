package rest

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"net/http"
	"strings"
)

// Error is custom type for displaying information like status and Body of the response message
type Error struct {
	StatusCode int
	Message    string
}

func (e Error) Error() string {
	return fmt.Sprintf("status %d, message %s", e.StatusCode, e.Message)
}

type MarshallFunc func(v any) ([]byte, error)
type UnmarshallFunc func(data []byte, v any) error

var (
	// ErrMethodNotAllowed ...
	ErrMethodNotAllowed = errors.New("method not allowed")
	// ErrMarshallerFuncNotFound ...
	ErrMarshallerFuncNotFound = errors.New("marshaller function not found in map")
	// ErrUnmarshalerFuncNotFound ...
	ErrUnmarshalerFuncNotFound = errors.New("unmarshaler function not found in map")
)

var (
	// Marshaller maps some basic content types with top level encoding functions from stdlib
	Marshaller = EncoderRegistry{
		ApplicationJSON: json.Marshal,
		ApplicationXML:  xml.Marshal,
	}

	// Unmarshaler maps some basic content types with top level decoding functions from stdlib
	Unmarshaler = DecoderRegistry{
		ApplicationJSON: json.Unmarshal,
		ApplicationXML:  xml.Unmarshal,
	}

	// DefaultContentType will be used if no content type specified in headers
	DefaultContentType = ApplicationJSON

	DefaultHttp = &Http{
		Client:         http.DefaultClient,
		RequestWrapper: &RequestWrapper{},
	}
)

// Modify resource on requestURL with options WithBody, WithHeaders and return T
// method can be POST, PUT, DELETE
// if error occurred T will be zero value
func Modify[T any](ctx context.Context, method string, requestURL string, options ...Option) (T, error) {
	var result T
	allowedMethods := []string{http.MethodPost, http.MethodPut, http.MethodPatch}
	method = strings.ToUpper(method)

	if !slices.Contains(allowedMethods, method) {
		return result, fmt.Errorf("%w: %s, you can use one of %v", ErrMethodNotAllowed, method,
			allowedMethods)
	}

	return request[T](ctx, &config[T]{
		method:     method,
		requestURL: requestURL,
		requester:  DefaultHttp,
		encoder: Encoder{
			Registry: Marshaller,
		},
		decoder: Decoder[T]{
			Registry: Unmarshaler,
		},
	}, options...)
}

// Get resource T from requestedURL with options: WithBody, WithHeaders
// if error occurred T will be zero value
func Get[T any](ctx context.Context, requestURL string, options ...Option) (T, error) {
	return request[T](ctx, &config[T]{
		method:     http.MethodGet,
		requestURL: requestURL,
		requester:  DefaultHttp,
		encoder: Encoder{
			Registry: Marshaller,
		},
		decoder: Decoder[T]{
			Registry: Unmarshaler,
		},
	}, options...)
}

// Delete resource from requestedURL with options WithBody, WithHeaders
// if error occurred T will be zero value
func Delete[T any](ctx context.Context, requestURL string, options ...Option) (T, error) {
	return request[T](ctx, &config[T]{
		method:     http.MethodDelete,
		requestURL: requestURL,
		requester:  DefaultHttp,
		encoder: Encoder{
			Registry: Marshaller,
		},
		decoder: Decoder[T]{
			Registry: Unmarshaler,
		},
	}, options...)
}
