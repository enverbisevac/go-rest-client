package rest

import (
	"context"
	"net/http"
)

type config[T any] struct {
	method     string
	requestURL string
	requester  interface {
		Request(
			ctx context.Context,
			method string,
			requestURL string,
			body []byte,
			header http.Header) (*Response, error)
	}
	encoder interface {
		Encode(
			value any,
			contentType ContentType,
			marshalFunc MarshallFunc) ([]byte, error)
	}
	decoder interface {
		Decode(
			data []byte,
			val *T,
			contentType ContentType,
			unmarshallFunc UnmarshallFunc) error
	}
}

func request[T any](ctx context.Context, cfg *config[T], options ...Option) (T, error) {
	rd := RequestOption{}
	for _, f := range options {
		f(&rd)
	}

	var val T

	data, err := cfg.encoder.Encode(rd.Body, ContentType(rd.Header.Get(Content)), rd.MarshalFunc)
	if err != nil {
		return val, err
	}

	response, err := cfg.requester.Request(ctx, cfg.method, cfg.requestURL, data, rd.Header)
	if err != nil {
		return val, err
	}

	if response.StatusCode >= http.StatusBadRequest {
		return val, Error{
			StatusCode: response.StatusCode,
			Message:    string(data),
		}
	}

	if len(response.Body) > 0 {
		err = cfg.decoder.Decode(response.Body, &val, ContentType(response.Header.Get(Content)), rd.UnmarshallFunc)
	}

	return val, err
}
