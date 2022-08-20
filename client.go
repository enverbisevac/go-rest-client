package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"golang.org/x/exp/slices"
	"io"
	"net/http"
	"strings"
)

// Error is custom type for displaying information like status and body of the response message
type Error struct {
	StatusCode int
	Message    string
}

func (e Error) Error() string {
	return fmt.Sprintf("status %d, message %s", e.StatusCode, e.Message)
}

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
	Marshaller = map[ContentType]func(v any) ([]byte, error){
		ApplicationJSON: json.Marshal,
		ApplicationXML:  xml.Marshal,
	}

	// Unmarshaler maps some basic content types with top level decoding functions from stdlib
	Unmarshaler = map[ContentType]func(data []byte, v any) error{
		ApplicationJSON: json.Unmarshal,
		ApplicationXML:  xml.Unmarshal,
	}

	// DefaultContentType will be used if no content type specified in headers
	DefaultContentType = ApplicationJSON
)

// Modify resource on requestURL with options WithBody, WithHeaders and return T
// method can be POST, PUT, DELETE
// if error occurred T will be zero value
func Modify[T any](ctx context.Context, method string, requestURL string, options ...Option) (val T, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error in Modify: %w", err)
		}
	}()
	var rd requestData
	for _, f := range options {
		f(&rd)
	}

	allowedMethods := []string{http.MethodPost, http.MethodPut, http.MethodPatch}
	method = strings.ToUpper(method)

	if !slices.Contains(allowedMethods, method) {
		return val, fmt.Errorf("%w: %s, you can use one of %v", ErrMethodNotAllowed, method,
			allowedMethods)
	}

	data, err := marshal(rd.Body, ContentType(rd.Header.Get(Content)))
	if err != nil {
		return
	}

	request, err := http.NewRequestWithContext(ctx, method, requestURL, bytes.NewReader(data))
	if err != nil {
		return
	}

	if rd.Header != nil {
		request.Header = rd.Header
	}

	response, err := http.DefaultClient.Do(request)
	if response != nil {
		defer func(rc io.ReadCloser) {
			_ = rc.Close()
		}(response.Body)
	}
	if err != nil {
		return
	}

	data, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.StatusCode >= http.StatusBadRequest {
		return val, Error{
			StatusCode: response.StatusCode,
			Message:    string(data),
		}
	}

	if len(data) > 0 {
		err = unmarshal(data, &val, ContentType(response.Header.Get(Content)))
	}

	return
}

// Get resource T from requestedURL with options: WithBody, WithHeaders
// if error occurred T will be zero value
func Get[T any](ctx context.Context, requestURL string, options ...Option) (result T, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error in Get: %w", err)
		}
	}()

	var reader io.Reader
	var data []byte

	var rd requestData
	for _, f := range options {
		f(&rd)
	}

	if rd.Body != nil {
		data, err = marshal(rd.Body, ContentType(rd.Header.Get(Content)))
		if err != nil {
			return
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, reader)
	if err != nil {
		return
	}

	req.Header = rd.Header

	response, err := http.DefaultClient.Do(req)
	if response != nil {
		defer func(rc io.ReadCloser) {
			_ = rc.Close()
		}(response.Body)
	}

	if err != nil {
		return
	}

	data, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.StatusCode >= http.StatusBadRequest {
		return result, Error{
			StatusCode: response.StatusCode,
			Message:    string(data),
		}
	}

	err = unmarshal(data, &result, ContentType(response.Header.Get(Content)))
	return
}

// Delete resource from requestedURL with options WithBody, WithHeaders
// if error occurred T will be zero value
func Delete(ctx context.Context, requestURL string, options ...Option) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error in Delete: %w", err)
		}
	}()

	var rd requestData
	for _, f := range options {
		f(&rd)
	}

	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, http.MethodDelete, requestURL, nil)
	if err != nil {
		return
	}

	req.Header = rd.Header

	response, err := http.DefaultClient.Do(req)
	if response != nil {
		defer func(rc io.ReadCloser) {
			_ = rc.Close()
		}(response.Body)
	}

	if err != nil {
		return
	}

	var data []byte
	data, err = io.ReadAll(response.Body)
	if err != nil {
		return
	}

	if response.StatusCode >= http.StatusBadRequest {
		return Error{
			StatusCode: response.StatusCode,
			Message:    string(data),
		}
	}

	return
}

func marshal(object any, content ContentType) ([]byte, error) {
	if content == "" {
		content = DefaultContentType
	}
	content = getContentType(content)
	f, ok := Marshaller[content]
	if !ok {
		return []byte{}, ErrMarshallerFuncNotFound
	}
	return f(object)
}

func unmarshal(data []byte, object any, content ContentType) error {
	if content == "" {
		content = DefaultContentType
	}
	content = getContentType(content)
	f, ok := Unmarshaler[content]
	if !ok {
		return ErrUnmarshalerFuncNotFound
	}
	return f(data, object)
}

func getContentType(content ContentType) ContentType {
	parts := strings.Split(string(content), ";")
	return ContentType(strings.TrimSpace(parts[0]))
}
