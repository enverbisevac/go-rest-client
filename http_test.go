package rest

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"reflect"
	"testing"
)

type mockRequestWrapper struct {
	NewRequestFunc func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error)
}

func (m *mockRequestWrapper) NewRequestWithContext(ctx context.Context, method, url string,
	body io.Reader) (*http.Request, error) {
	return m.NewRequestFunc(ctx, method, url, body)
}

type mockHttp struct {
	DoFunc func(request *http.Request) (*http.Response, error)
}

func (m *mockHttp) Do(request *http.Request) (*http.Response, error) {
	return m.DoFunc(request)
}

type mockReadCloserError struct {
}

func (r *mockReadCloserError) Read(p []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (r *mockReadCloserError) Close() error {
	return nil
}

func TestHttp_Request(t *testing.T) {
	type fields struct {
		RequestWrapper *mockRequestWrapper
		Client         *mockHttp
	}
	type args struct {
		ctx        context.Context
		method     string
		requestURL string
		body       []byte
		header     http.Header
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Response
		wantErr bool
	}{
		{
			name: "new request expect error",
			fields: fields{
				RequestWrapper: &mockRequestWrapper{
					NewRequestFunc: func(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
						return nil, errors.New("request error")
					},
				},
			},
			args: args{
				ctx:        nil,
				method:     "",
				requestURL: "",
				body:       nil,
				header:     nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "net call return error",
			fields: fields{
				Client: &mockHttp{
					DoFunc: func(request *http.Request) (*http.Response, error) {
						return nil, errors.New("network error")
					},
				},
				RequestWrapper: &mockRequestWrapper{
					NewRequestFunc: http.NewRequestWithContext,
				},
			},
			args: args{
				ctx:        context.Background(),
				method:     "",
				requestURL: "",
				body:       nil,
				header:     nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "empty method should work as GET",
			fields: fields{
				RequestWrapper: &mockRequestWrapper{
					NewRequestFunc: http.NewRequestWithContext,
				},
				Client: &mockHttp{
					DoFunc: func(request *http.Request) (*http.Response, error) {
						return &http.Response{
							StatusCode:    200,
							Body:          io.NopCloser(bytes.NewBufferString("")),
							ContentLength: 0,
							Header:        Header(WithContent(ApplicationJSON)),
						}, nil
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				method:     "",
				requestURL: "",
				body:       nil,
				header:     nil,
			},
			want: &Response{
				Body:       []byte(""),
				StatusCode: 200,
				Header: http.Header{
					Content: []string{string(ApplicationJSON)},
				},
			},
			wantErr: false,
		},
		{
			name: "response Body error",
			fields: fields{
				RequestWrapper: &mockRequestWrapper{
					NewRequestFunc: http.NewRequestWithContext,
				},
				Client: &mockHttp{
					DoFunc: func(request *http.Request) (*http.Response, error) {
						return &http.Response{
							Body: &mockReadCloserError{},
						}, nil
					},
				},
			},
			args: args{
				ctx:        context.Background(),
				method:     "",
				requestURL: "",
				body:       nil,
				header:     nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Http{
				Client:         tt.fields.Client,
				RequestWrapper: tt.fields.RequestWrapper,
			}
			got, err := h.Request(tt.args.ctx, tt.args.method, tt.args.requestURL, tt.args.body, tt.args.header)
			if (err != nil) != tt.wantErr {
				t.Errorf("Request() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Request() got = %v, want %v", got, tt.want)
			}
		})
	}
}
