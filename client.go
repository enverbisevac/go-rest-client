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

type Error struct {
	StatusCode int
	Message    string
}

func (e Error) Error() string {
	return fmt.Sprintf("status %d, message %s", e.StatusCode, e.Message)
}

var (
	ErrMethodNotAllowed         = errors.New("method not allowed")
	ErrMarshallerFuncNotFound   = errors.New("marshaller function not found in map")
	ErrUnmarshallerFuncNotFound = errors.New("unmarshaler function not found in map")
)

var (
	Marshaller = map[ContentType]func(v any) ([]byte, error){
		ApplicationJSON: json.Marshal,
		ApplicationXML:  xml.Marshal,
	}

	Unmarshaler = map[ContentType]func(data []byte, v any) error{
		ApplicationJSON: json.Unmarshal,
		ApplicationXML:  xml.Unmarshal,
	}

	DefaultContentType = ApplicationJSON
)

func Modify[T any](ctx context.Context, method string, requestURL string, body any,
	header http.Header) (val T, err error) {

	defer func() {
		if err != nil {
			err = fmt.Errorf("error in Modify: %w", err)
		}
	}()

	allowedMethods := []string{http.MethodPost, http.MethodPut, http.MethodPatch}
	method = strings.ToUpper(method)

	if !slices.Contains(allowedMethods, method) {
		return val, fmt.Errorf("%w: %s, you can use one of %v", ErrMethodNotAllowed, method,
			allowedMethods)
	}

	data, err := marshal(body, ContentType(header.Get(Content)))
	if err != nil {
		return
	}

	request, err := http.NewRequestWithContext(ctx, method, requestURL, bytes.NewReader(data))
	if err != nil {
		return
	}

	if header != nil {
		request.Header = header
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

func Get[T any](ctx context.Context, requestURL string, body any, header http.Header) (result T, err error) {
	var reader io.Reader
	var data []byte

	defer func() {
		if err != nil {
			err = fmt.Errorf("error in Get: %w", err)
		}
	}()

	if body != nil {
		data, err = marshal(body, ContentType(header.Get(Content)))
		if err != nil {
			return
		}
		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, reader)
	if err != nil {
		return
	}

	req.Header = header

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

func Delete(ctx context.Context, requestURL string, header http.Header) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("error in Delete: %w", err)
		}
	}()

	var req *http.Request
	req, err = http.NewRequestWithContext(ctx, http.MethodDelete, requestURL, nil)
	if err != nil {
		return
	}

	req.Header = header

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
	f, ok := Unmarshaler[content]
	if !ok {
		return ErrUnmarshallerFuncNotFound
	}
	return f(data, object)
}
