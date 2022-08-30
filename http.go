package rest

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type Response struct {
	Body       []byte
	StatusCode int
	Header     http.Header
}

type RequestWrapper struct{}

func (h *RequestWrapper) NewRequestWithContext(ctx context.Context, method, url string,
	body io.Reader) (*http.Request, error) {
	return http.NewRequestWithContext(ctx, method, url, body)
}

type Http struct {
	Client interface {
		Do(request *http.Request) (*http.Response, error)
	}

	RequestWrapper interface {
		NewRequestWithContext(ctx context.Context, method, url string, body io.Reader) (*http.Request, error)
	}
}

func (h *Http) Request(ctx context.Context, method string, requestURL string, body []byte,
	header http.Header) (*Response, error) {

	req, err := h.RequestWrapper.NewRequestWithContext(ctx, method, requestURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header = header

	response, err := h.Client.Do(req)
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
