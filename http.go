package rest

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type Doer interface {
	Do(request *http.Request) (*http.Response, error)
}

type Response struct {
	Body       []byte
	StatusCode int
	Header     http.Header
}

type Http struct {
	Client Doer
}

func NewHttp(client Doer) Http {
	return Http{
		Client: client,
	}
}

func (h *Http) Request(ctx context.Context, method string, requestURL string, body []byte,
	header http.Header) (*Response, error) {

	request, err := http.NewRequestWithContext(ctx, method, requestURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	request.Header = header

	response, err := h.Client.Do(request)
	if response != nil {
		defer func(rc io.ReadCloser) {
			_ = rc.Close()
		}(response.Body)
	}
	if err != nil {
		return nil, err
	}

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		Body:       data,
		StatusCode: response.StatusCode,
		Header:     response.Header,
	}, nil
}
