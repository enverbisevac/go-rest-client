package rest

import (
	"context"
	"net/http"
)

type Encoder interface {
	Encode(object any, content Parser) ([]byte, error)
	Set(contentType ContentType, marshallFunc MarshallFunc)
	Clone() EncodeRegistry
}

type Decoder interface {
	Decode(data []byte, object any, content Parser) error
	Set(contentType ContentType, unmarshaler UnmarshallFunc)
	Clone() DecodeRegistry
}

type Requester interface {
	Request(ctx context.Context, method string, requestURL string, body []byte, header http.Header) (*Response, error)
}

type config struct {
	method     string
	requestURL string
	requester  Requester
	encoder    Encoder
	decoder    Decoder
}

func request[T any](ctx context.Context, cfg *config, options ...Option) (T, error) {
	var val T
	rd := RequestData[T]{
		encoder: cfg.encoder,
		decoder: cfg.decoder,
	}
	for _, f := range options {
		f((*RequestData[any])(&rd))
	}

	data, err := rd.Encode()
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

	if len(data) > 0 {
		err = rd.Decode(response.Body, &val)
	}

	return val, err
}
